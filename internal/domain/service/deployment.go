package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/client"
	"github.com/go-git/go-git/v6/plumbing/transport/http"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	moby_client "github.com/moby/moby/client"
	"github.com/server-selfish/backend/internal/constant"
	container_repository "github.com/server-selfish/backend/internal/domain/repository/container"
	deployment_repository "github.com/server-selfish/backend/internal/domain/repository/deployment"
	"github.com/server-selfish/backend/internal/domain/schema"
	cache_infra "github.com/server-selfish/backend/internal/infra/cache"
	github_infra "github.com/server-selfish/backend/internal/infra/github"
	"github.com/server-selfish/backend/internal/pkg"
	defined_error "github.com/server-selfish/backend/internal/pkg/error"
)

type (
	DeploymentService interface {
		GetDeploymentsByProjectId(ctx context.Context, userId, projectId pgtype.UUID) ([]schema.GetDeploymentData, error)
		GetDeploymentByDeploymentId(ctx context.Context, userId, deploymentId pgtype.UUID) (schema.GetSingleDeploymentData, error)
		GetActiveDeploymentByDeploymentId(ctx context.Context, userId, deploymentId pgtype.UUID) (schema.GetActiveDeploymentHistory, error)
		GetHistoryDeploymentByDeploymentId(ctx context.Context, userId, deploymentId pgtype.UUID) ([]schema.GetHistoryDeploymentHistory, error)
		GetTechstackName(ctx context.Context) (schema.GetTechstackList, error)
		GetTechstackVersionByName(ctx context.Context, techstackName string) ([]schema.GetTechstackVersion, error)
		// CreateDeployment(ctx context.Context, params deployment_repository.CreateDeploymentParams) error
		CreateNewDeploymentVersionByDeploymentId(ctx context.Context, userID pgtype.UUID, installationID int64, params schema.CreateDeploymentHistoryParams) error
		buildAndRunContainer(ctx context.Context, p schema.BuildAndRunContainerParams) error
		DeleteDeploymentByDeploymentId(ctx context.Context, userId, deploymentId pgtype.UUID) error
	}
	deploymentService struct {
		gs    GithubAppService
		dr    *deployment_repository.Queries
		cr    *container_repository.Queries
		tm    pkg.TxManager
		dc    *moby_client.Client
		gi    github_infra.GithubInfra
		cache cache_infra.ValkeyInfra
	}
)

func NewDeploymentService(
	dr *deployment_repository.Queries,
	gs GithubAppService,
	cache cache_infra.ValkeyInfra,
	cr *container_repository.Queries,
	tm pkg.TxManager,
	dc *moby_client.Client,
	gi github_infra.GithubInfra,
) DeploymentService {
	return &deploymentService{
		gs:    gs,
		dr:    dr,
		cr:    cr,
		tm:    tm,
		cache: cache,
		dc:    dc,
		gi:    gi,
	}
}

// CreateNewDeploymentVersion implements [DeploymentService].
func (d *deploymentService) CreateNewDeploymentVersionByDeploymentId(ctx context.Context, userID pgtype.UUID, installationID int64, params schema.CreateDeploymentHistoryParams) error {
	reps, err := d.gs.ListInstallationRepositories(ctx, userID, installationID)
	if err != nil {
		return err
	}
	var rep schema.GithubInstallationRepository

	repId, err := strconv.ParseInt(params.RepositoryID, 10, 64)
	for _, val := range reps {
		if val.ID == repId {
			rep = val
			break
		}
	}
	if err != nil {
		return err
	}

	var cn, in, activeContainerName, toBeDeletedPath string

	// begin transaction
	if err := d.tm.WithTx(ctx, func(tx pgx.Tx) error {
		depQuery := d.dr.WithTx(tx)
		conQuery := d.cr.WithTx(tx)
		depId, err := depQuery.UpsertDeployment(ctx, deployment_repository.UpsertDeploymentParams{
			Name:           params.DeploymentName,
			GitRemoteUrl:   rep.HTMLURL,
			InstallationID: installationID,
			Name_2:         params.ProjectName,
			UserID:         userID,
		})
		if err != nil {
			return err
		}

		activeContainer, err := conQuery.GetActiveDeploymentHistoryContainerByDeploymentId(ctx, container_repository.GetActiveDeploymentHistoryContainerByDeploymentIdParams{
			UserID:       userID,
			DeploymentID: depId,
		})
		if err == nil {
			activeContainerName = activeContainer.Name
		}

		deployment, err := depQuery.GetDeploymentByDeploymentId(ctx, deployment_repository.GetDeploymentByDeploymentIdParams{
			UserID: userID,
			ID:     depId,
		})
		if err != nil {
			return err
		}
		tmpPath := fmt.Sprintf(
			"tmp/builds/%s/%s/tmp",
			deployment.Name,
			params.Branch,
		)
		toBeDeletedPath = fmt.Sprintf(
			"tmp/builds/%s/%s",
			deployment.Name,
			params.Branch,
		)

		it, err := d.gi.CreateInstallationToken(ctx, installationID)
		if err != nil {
			return err
		}
		repo, err := git.PlainClone(tmpPath, &git.CloneOptions{
			URL: deployment.GitRemoteUrl,
			ClientOptions: []client.Option{
				client.WithHTTPAuth(&http.BasicAuth{
					Username: "x-access-token",
					Password: it.Token,
				}),
			},
			ReferenceName: plumbing.NewBranchReferenceName(params.Branch),
			SingleBranch:  true,
			Depth:         1,
		})
		if err != nil {
			return err
		}

		commitId, commitMsg, version, err := pkg.ExtractRepoMetaData(repo)
		if err != nil {
			return err
		}

		finalPath := fmt.Sprintf(
			"tmp/builds/%s/%s/%s",
			deployment.Name,
			params.Branch,
			version,
		)
		if err := os.Rename(tmpPath, finalPath); err != nil {
			return fmt.Errorf("failed to rename %s to %s: %w", tmpPath, finalPath, err)
		}

		// container name
		cnUUID, err := uuid.NewV7()
		if err != nil {
			return err
		}
		cn = fmt.Sprintf("%s.%s", cnUUID.String(), deployment.Name)
		in = fmt.Sprintf("%s:%s", deployment.Name, version)

		// deactivate + stop & remove container of active version
		if err := depQuery.SetActiveDeploymentHistoryNonActiveByDeploymentId(ctx, deployment_repository.SetActiveDeploymentHistoryNonActiveByDeploymentIdParams{
			UserID:       userID,
			DeploymentID: depId,
		}); err != nil {
			return err
		}
		if activeContainerName != "" {
			if _, err := d.dc.ContainerStop(ctx, activeContainer.Name, moby_client.ContainerStopOptions{}); err != nil {
				return err
			}
		}

		// build and run new container
		if err := d.buildAndRunContainer(ctx, schema.BuildAndRunContainerParams{
			DepQuery:              depQuery,
			Path:                  finalPath,
			DeploymentId:          depId,
			ContainerName:         cn,
			ImageName:             in,
			BuildCommand:          params.BuildCommand,
			BuildFolder:           params.BuildFolder,
			Env:                   params.Env,
			Port:                  params.Port,
			DeploymentTechstackID: params.DeploymentTechstackID,
			ProjectName:           params.ProjectName,
			MainFileName:          params.MainFileName,
			UserId:                userID,
		}); err != nil {
			return err
		}

		// QUERY: save deployment history
		historyId, err := depQuery.CreateDeploymentHistory(ctx, deployment_repository.CreateDeploymentHistoryParams{
			DeploymentID:          depId,
			Branch:                params.Branch,
			CommitID:              commitId,
			CommitMsg:             commitMsg,
			Version:               version,
			DeploymentTechstackID: params.DeploymentTechstackID,
			BuildCommand:          pgtype.Text{String: params.BuildCommand, Valid: true},
			BuildFolder:           pgtype.Text{String: params.BuildFolder, Valid: true},
		})
		if err != nil {
			return err
		}

		// QUERY: save container
		conId, err := conQuery.CreateContainer(ctx, container_repository.CreateContainerParams{
			Name:                cn,
			DeploymentHistoryID: historyId,
		})
		if err != nil {
			return err
		}

		// QUERY: save port and env
		keys := make([]string, 0, len(params.Env))
		values := make([]string, 0, len(params.Env))
		for _, val := range params.Env {
			keys = append(keys, val.Key)
			values = append(values, val.Value)
		}

		if err := conQuery.CreateContainerEnv(ctx, container_repository.CreateContainerEnvParams{
			ContainerIds: conId,
			Keys:         keys,
			Values:       values,
		}); err != nil {
			return err
		}

		external := make([]int32, 0, len(params.Port))
		internal := make([]int32, 0, len(params.Port))
		protocol := make([]string, 0, len(params.Port))

		for _, val := range params.Port {
			internal = append(internal, val.Internal)
			external = append(external, val.External)
			protocol = append(protocol, val.Protocol)
		}

		if err := conQuery.CreateContainerPort(ctx, container_repository.CreateContainerPortParams{
			ContainerIds: conId,
			External:     external,
			Internal:     internal,
			Protocol:     protocol,
		}); err != nil {
			return err
		}

		return nil
	}); err != nil {
		if toBeDeletedPath != "" {
			if err := pkg.DeleteDir(toBeDeletedPath); err != nil {
				return err
			}
		}
		// stop and remove built container
		if cn != "" {
			if err := pkg.StopAndRemoveContainer(ctx, d.dc, cn); err != nil {
				return err
			}
		}
		// remove image
		if in != "" {
			if err := pkg.RemoveImage(ctx, d.dc, in); err != nil {
				return err
			}
		}
		// start activeContainer
		if activeContainerName != "" {
			_, err = d.dc.ContainerStart(ctx, activeContainerName, moby_client.ContainerStartOptions{})
			if err != nil {
				return err
			}
		}
		return err
	}
	go func() {
		if activeContainerName != "" {
			if _, err := d.dc.ContainerRemove(ctx, activeContainerName, moby_client.ContainerRemoveOptions{}); err != nil {
				log.Printf("failed to remove container: %v", err)
			}
		}
	}()
	return nil
}

// buildAndRunContainer implements [DeploymentService].`
func (d *deploymentService) buildAndRunContainer(ctx context.Context, p schema.BuildAndRunContainerParams) error {
	techstack, err := p.DepQuery.GetTechstackByTechstackId(ctx, p.DeploymentTechstackID)
	if err != nil {
		return err
	}
	mainFilePath := fmt.Sprintf("%s/%s", p.BuildFolder, p.MainFileName)
	cmdJson := pkg.GetRunCommandByTechstack(techstack.Name, mainFilePath, techstack.DockerBaseImage)
	tpt := schema.DockerFileTemplate{
		DockerBaseImage:    techstack.DockerBaseImage,
		DockerRuntimeImage: techstack.DockerRuntimeImage,
		BuildFolder:        p.BuildFolder,
		BuildCommand:       p.BuildCommand,
		RunCommand:         cmdJson,
	}
	template, err := pkg.ParseTemplateFromEmbed(pkg.GetFileNameByTechstack(techstack.Name), tpt)
	if err != nil {
		return err
	}
	dockerignoreTemplate, err := pkg.ParseTemplateFromEmbed("dockerignore", nil)
	if err != nil {
		return err
	}

	dockerfilePath := filepath.Join(p.Path, "Dockerfile")
	if err := pkg.WriteFile(dockerfilePath, []byte(template)); err != nil {
		return err
	}

	dockerignorePath := filepath.Join(p.Path, ".dockerignore")
	if err := pkg.WriteFile(dockerignorePath, []byte(dockerignoreTemplate)); err != nil {
		return err
	}
	// build docker image
	if err = pkg.BuildDockerImage(ctx, d.dc, p.Path, p.ImageName, map[string]string{}); err != nil {
		return err
	}
	// network name
	nn := fmt.Sprintf("%s-network", p.ProjectName)
	if err := pkg.EnsureDockerNetwork(ctx, d.dc, nn); err != nil {
		return err
	}

	exposedPorts := network.PortSet{}
	portBindings := network.PortMap{}

	for _, p := range p.Port {
		port := network.Port(network.MustParsePort(fmt.Sprintf(
			"%d/%s",
			p.Internal,
			strings.ToLower(p.Protocol),
		)))

		exposedPorts[port] = struct{}{}
		portBindings[port] = []network.PortBinding{
			{
				HostIP:   constant.ALL_ADDR,
				HostPort: strconv.Itoa(int(p.External)),
			},
		}
	}

	containerEnv := make([]string, 0, len(p.Env))

	for _, env := range p.Env {
		containerEnv = append(
			containerEnv,
			fmt.Sprintf("%s=%s", env.Key, env.Value),
		)
	}

	if _, err := d.dc.ContainerCreate(ctx, moby_client.ContainerCreateOptions{
		Name: p.ContainerName,
		Config: &container.Config{
			ExposedPorts: exposedPorts,
			Env:          containerEnv,
			Image:        p.ImageName,
		},
		HostConfig: &container.HostConfig{
			PortBindings: portBindings,
			RestartPolicy: container.RestartPolicy{
				Name:              container.RestartPolicyOnFailure,
				MaximumRetryCount: constant.MAXIMUM_RESTART,
			},
		},
		NetworkingConfig: &network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				nn: {},
			},
		},
	}); err != nil {
		return err
	}

	_, err = d.dc.ContainerStart(ctx, p.ContainerName, moby_client.ContainerStartOptions{})
	if err != nil {
		return err
	}
	return nil
}

// GetTechstackName implements [DeploymentService].
func (d *deploymentService) GetTechstackName(ctx context.Context) (schema.GetTechstackList, error) {
	tn, err := d.dr.GetTechstackName(ctx)
	if err != nil {
		return schema.GetTechstackList{}, err
	}
	return schema.GetTechstackList{
		Name: tn,
	}, nil
}

// GetTechstackVersionByName implements [DeploymentService].
func (d *deploymentService) GetTechstackVersionByName(ctx context.Context, techstackName string) ([]schema.GetTechstackVersion, error) {
	vl, err := d.dr.GetTechstackVersionByName(ctx, techstackName)
	if err != nil {
		return []schema.GetTechstackVersion{}, err
	}

	resp := make([]schema.GetTechstackVersion, 0, len(vl))
	for _, v := range vl {
		resp = append(resp, schema.GetTechstackVersion{
			ID:      v.ID,
			Version: v.Version,
		})
	}
	return resp, nil
}

// DeleteDeploymentByDeploymentId implements [DeploymentService].
func (d *deploymentService) DeleteDeploymentByDeploymentId(ctx context.Context, userID, deploymentId pgtype.UUID) error {
	if err := d.dr.DeleteDeploymentByDeploymentId(ctx, deployment_repository.DeleteDeploymentByDeploymentIdParams{
		UserID: userID,
		ID:     deploymentId,
	}); err != nil {
		return err
	}
	return nil
}

// CreateDeployment implements [DeploymentService].
// func (d *deploymentService) CreateDeployment(ctx context.Context, params deployment_repository.CreateDeploymentParams) error {
// 	if err := d.dr.CreateDeployment(ctx, params); err != nil {
// 		return err
// 	}
// 	return nil
// }

// GetHistoryDeploymentByDeploymentId implements [DeploymentService].
func (d *deploymentService) GetHistoryDeploymentByDeploymentId(ctx context.Context, userId, deploymentId pgtype.UUID) ([]schema.GetHistoryDeploymentHistory, error) {
	hd, err := d.dr.GetDeploymentHistoryByDeploymentId(ctx, deployment_repository.GetDeploymentHistoryByDeploymentIdParams{
		UserID:       userId,
		DeploymentID: deploymentId,
	})
	if err != nil {
		return nil, err
	}
	var res []schema.GetHistoryDeploymentHistory
	for _, h := range hd {
		res = append(res, schema.GetHistoryDeploymentHistory{
			ID:                h.ID,
			Branch:            h.Branch,
			CommitID:          h.CommitID,
			CommitMessage:     h.CommitMessage,
			DeploymentVersion: h.DeploymentVersion,
			// Port:              h.Port,
			CreatedAt: h.CreatedAt.Time.String(),
			UpdatedAt: h.UpdatedAt.Time.String(),
		})
	}
	return res, nil
}

// GetActiveDeploymentByDeploymentId implements [DeploymentService].
func (d *deploymentService) GetActiveDeploymentByDeploymentId(ctx context.Context, userId, deploymentId pgtype.UUID) (schema.GetActiveDeploymentHistory, error) {
	ad, err := d.dr.GetActiveDeploymentHistoryByDeploymentId(ctx, deployment_repository.GetActiveDeploymentHistoryByDeploymentIdParams{
		UserID:       userId,
		DeploymentID: deploymentId,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return schema.GetActiveDeploymentHistory{}, defined_error.ErrActiveDeploymentNotFound
		}
		return schema.GetActiveDeploymentHistory{}, err
	}
	res := schema.GetActiveDeploymentHistory{
		DeploymentHistoryID: ad.DeploymentHistoryID,
		Branch:              ad.Branch,
		CommitId:            ad.CommitID,
		CommitMessage:       ad.CommitMessage,
		DeploymentVersion:   ad.DeploymentVersion,
		// Port:                ad.Port,
		BuildCommand:     ad.BuildCommand.String,
		TechstackID:      ad.TechstackID,
		TechstackName:    ad.TechstackName,
		TechstackVersion: ad.TechstackVersion,
	}
	return res, nil
}

// GetDeploymentByDeploymentId implements [DeploymentService].
func (d *deploymentService) GetDeploymentByDeploymentId(ctx context.Context, userId, deploymentId pgtype.UUID) (schema.GetSingleDeploymentData, error) {
	deployment, err := d.dr.GetDeploymentByDeploymentId(ctx, deployment_repository.GetDeploymentByDeploymentIdParams{
		UserID: userId,
		ID:     deploymentId,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return schema.GetSingleDeploymentData{}, defined_error.ErrDeploymentNotFound
		}
		return schema.GetSingleDeploymentData{}, err
	}
	res := schema.GetSingleDeploymentData{
		ID:           deployment.ID.String(),
		Name:         deployment.Name,
		GitRemoteURL: deployment.GitRemoteUrl,
		CreatedAt:    deployment.CreatedAt.Time.String(),
		UpdatedAt:    deployment.UpdatedAt.Time.String(),
	}
	return res, nil
}

// GetDeploymentsByProjectId implements [DeploymentService].
func (d *deploymentService) GetDeploymentsByProjectId(ctx context.Context, userId, projectId pgtype.UUID) ([]schema.GetDeploymentData, error) {
	deployments, err := d.dr.GetDeploymentsByProjectId(ctx, deployment_repository.GetDeploymentsByProjectIdParams{
		UserID:    userId,
		ProjectID: projectId,
	})
	if err != nil {
		return nil, err
	}
	// if len(deployments) == 0 {
	// 	return nil, pkg.ErrNotFound
	// }
	var res []schema.GetDeploymentData
	for _, dep := range deployments {
		res = append(res, schema.GetDeploymentData{
			DeploymentID:      dep.DeploymentID,
			DeploymentName:    dep.DeploymentName,
			GitRemoteUrl:      dep.GitRemoteUrl,
			Branch:            dep.Branch,
			CommitID:          dep.CommitID,
			CommitMessage:     dep.CommitMessage,
			DeploymentVersion: dep.DeploymentVersion,
			// Port:              dep.Port,
			TechstackName:    dep.TechstackName,
			TechstackVersion: dep.TechstackVersion,
			ContainerID:      dep.ContainerID,
			CreatedAt:        dep.CreatedAt.Time.String(),
			UpdatedAt:        dep.UpdatedAt.Time.String(),
		})
	}
	return res, nil
}
