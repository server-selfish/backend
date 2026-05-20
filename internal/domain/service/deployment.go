package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path/filepath"

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
	github_infra "github.com/server-selfish/backend/internal/infra/github"
	"github.com/server-selfish/backend/internal/pkg"
)

type (
	DeploymentService interface {
		GetDeploymentsByProjectId(ctx context.Context, userId, projectId pgtype.UUID) ([]schema.GetDeploymentData, error)
		GetDeploymentByDeploymentId(ctx context.Context, userId, deploymentId pgtype.UUID) (schema.GetSingleDeploymentData, error)
		GetActiveDeploymentByDeploymentId(ctx context.Context, userId, deploymentId pgtype.UUID) (schema.GetActiveDeploymentHistory, error)
		GetHistoryDeploymentByDeploymentId(ctx context.Context, userId, deploymentId pgtype.UUID) ([]schema.GetHistoryDeploymentHistory, error)
		GetTechstackName(ctx context.Context) (schema.GetTechstackList, error)
		GetTechstackVersionByName(ctx context.Context, techstackName string) ([]schema.GetTechstackVersion, error)
		CreateDeployment(ctx context.Context, params deployment_repository.CreateDeploymentParams) error
		CreateNewDeploymentVersionByDeploymentId(ctx context.Context, userID pgtype.UUID, installationID int64, params deployment_repository.CreateDeploymentHistoryParams) error
		buildAndRunContainer(ctx context.Context, p schema.BuildAndRunContainerParams) error
		DeleteDeploymentByDeploymentId(ctx context.Context, userId, deploymentId pgtype.UUID) error
	}
	deploymentService struct {
		dr *deployment_repository.Queries
		cr *container_repository.Queries
		tm pkg.TxManager
		dc *moby_client.Client
		gi github_infra.GithubInfra
	}
)

func NewDeploymentService(dr *deployment_repository.Queries, cr *container_repository.Queries, tm pkg.TxManager, dc *moby_client.Client, gi github_infra.GithubInfra) DeploymentService {
	return &deploymentService{
		dr: dr,
		tm: tm,
		dc: dc,
		gi: gi,
	}
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

// CreateNewDeploymentVersion implements [DeploymentService].
func (d *deploymentService) CreateNewDeploymentVersionByDeploymentId(ctx context.Context, userID pgtype.UUID, installationID int64, params deployment_repository.CreateDeploymentHistoryParams) error {
	activeContainer, err := d.cr.GetActiveDeploymentHistoryContainerByDeploymentId(ctx, container_repository.GetActiveDeploymentHistoryContainerByDeploymentIdParams{
		UserID:       pgtype.UUID{},
		DeploymentID: params.DeploymentID,
	})
	if err != nil {
		return err
	}
	deployment, err := d.dr.GetDeploymentByDeploymentId(ctx, deployment_repository.GetDeploymentByDeploymentIdParams{
		UserID: userID,
		ID:     params.DeploymentID,
	})
	if err != nil {
		return err
	}
	path := fmt.Sprintf("tmp/builds/%s/%s/%s", deployment.Name, params.Branch, params.Version)

	it, err := d.gi.CreateInstallationToken(ctx, installationID)
	if err != nil {
		return err
	}
	// Clone and extract metadata repository
	repo, err := git.PlainClone(path, &git.CloneOptions{
		URL: deployment.GitRemoteUrl,
		ClientOptions: []client.Option{
			client.WithHTTPAuth(&http.BasicAuth{
				Username: "x-access-token",
				Password: it.Token,
			}),
		},
		ReferenceName: plumbing.NewBranchReferenceName(params.Branch),
		SingleBranch:  true,
		Depth:         0,
	})
	if err != nil {
		return err
	}

	commitId, commitMsg, version, err := pkg.ExtractRepoMetaData(repo)
	if err != nil {
		return err
	}

	// container name
	cnUUID, err := uuid.NewV7()
	if err != nil {
		return err
	}
	cn := fmt.Sprintf("%s.%s", cnUUID.String(), deployment.Name)
	in := fmt.Sprintf("%s:%s", deployment.Name, version)

	// begin transaction
	if err := d.tm.WithTx(ctx, func(tx pgx.Tx) error {
		depQuery := d.dr.WithTx(tx)

		// deactivate + stop & remove container of active version
		if err := depQuery.SetActiveDeploymentHistoryNonActiveByDeploymentId(ctx, deployment_repository.SetActiveDeploymentHistoryNonActiveByDeploymentIdParams{
			UserID:       userID,
			DeploymentID: params.DeploymentID,
		}); err != nil {
			return err
		}
		if _, err := d.dc.ContainerStop(ctx, activeContainer.Name, moby_client.ContainerStopOptions{}); err != nil {
			return err
		}

		// build and run new container
		if err := d.buildAndRunContainer(ctx, schema.BuildAndRunContainerParams{
			DepQuery:              depQuery,
			Path:                  path,
			DeploymentId:          params.DeploymentID,
			ContainerName:         cn,
			ImageName:             in,
			BuildCommand:          params.BuildCommand.String,
			BuildFolder:           params.BuildFolder.String,
			RunCommand:            params.RunCommand.String,
			DeploymentTechstackID: params.DeploymentTechstackID,
			UserId:                userID,
		}); err != nil {
			return err
		}

		params.CommitID = commitId
		params.CommitMsg = commitMsg
		params.Version = version

		// QUERY: save deployment history
		historyId, err := d.dr.CreateDeploymentHistory(ctx, params)
		if err != nil {
			return err
		}

		// QUERY: save container
		if err := d.cr.CreateContainer(ctx, container_repository.CreateContainerParams{
			Name:                cn,
			DeploymentHistoryID: historyId,
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {

		//  remove folder in the path
		if err := pkg.DeleteDir(path); err != nil {
			return err
		}
		// stop and remove built container
		if err := pkg.StopAndRemoveContainer(ctx, d.dc, cn); err != nil {
			return err
		}
		// remove image
		if err := pkg.RemoveImage(ctx, d.dc, in); err != nil {
			return err
		}
		// start activeContainer
		_, err = d.dc.ContainerStart(ctx, activeContainer.Name, moby_client.ContainerStartOptions{})
		if err != nil {
			return err
		}
		return err
	}
	go func() {
		if _, err := d.dc.ContainerRemove(ctx, activeContainer.Name, moby_client.ContainerRemoveOptions{}); err != nil {
			log.Printf("failed to remove container: %v", err)
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
	cmdJson := pkg.ShellToExecForm(p.RunCommand)
	tpt := schema.DockerFileTemplate{
		DockerBaseImage:    techstack.DockerBaseImage,
		DockerRuntimeImage: techstack.DockerRuntimeImage,
		BuildCommand:       p.BuildCommand,
		BuildFolder:        p.BuildFolder,
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

	project, err := d.dr.GetProjectByDeploymentId(ctx, deployment_repository.GetProjectByDeploymentIdParams{
		UserID: p.UserId,
		ID:     p.DeploymentId,
	})
	if err != nil {
		return err
	}

	// network name
	nn := fmt.Sprintf("%s-network", project.Name)
	if err := pkg.EnsureDockerNetwork(ctx, d.dc, nn); err != nil {
		return err
	}

	// multiple port need to be parsed from input user
	parsedPort, err := network.ParsePort("8080")
	if err != nil {
		return err
	}

	if _, err := d.dc.ContainerCreate(ctx, moby_client.ContainerCreateOptions{
		Name: p.ContainerName,
		Config: &container.Config{
			ExposedPorts: network.PortSet{
				parsedPort: struct{}{},
			},
			Env: []string{
				"APP_ENV=production",
			},
			Image: p.ImageName,
		},
		HostConfig: &container.HostConfig{
			PortBindings: network.PortMap{
				parsedPort: []network.PortBinding{{
					HostIP:   constant.ALL_ADDR,
					HostPort: "8080",
				}},
			},
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

	// run container
	_, err = d.dc.ContainerStart(ctx, p.ContainerName, moby_client.ContainerStartOptions{})
	if err != nil {
		return err
	}
	return nil
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
func (d *deploymentService) CreateDeployment(ctx context.Context, params deployment_repository.CreateDeploymentParams) error {
	if err := d.dr.CreateDeployment(ctx, params); err != nil {
		return err
	}
	return nil
}

// GetHistoryDeploymentByDeploymentId implements [DeploymentService].
func (d *deploymentService) GetHistoryDeploymentByDeploymentId(ctx context.Context, userId, deploymentId pgtype.UUID) ([]schema.GetHistoryDeploymentHistory, error) {
	hd, err := d.dr.GetDeploymentHistoryByDeploymentId(ctx, deployment_repository.GetDeploymentHistoryByDeploymentIdParams{
		UserID:       userId,
		DeploymentID: deploymentId,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pkg.ErrNotFound
		}
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
			Port:              h.Port,
			CreatedAt:         h.CreatedAt.Time.String(),
			UpdatedAt:         h.UpdatedAt.Time.String(),
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
			return schema.GetActiveDeploymentHistory{}, pkg.ErrNotFound
		}
		return schema.GetActiveDeploymentHistory{}, err
	}
	res := schema.GetActiveDeploymentHistory{
		DeploymentHistoryID: ad.DeploymentHistoryID,
		Branch:              ad.Branch,
		CommitId:            ad.CommitID,
		CommitMessage:       ad.CommitMessage,
		DeploymentVersion:   ad.DeploymentVersion,
		Port:                ad.Port,
		BuildCommand:        ad.BuildCommand.String,
		TechstackID:         ad.TechstackID,
		TechstackName:       ad.TechstackName,
		TechstackVersion:    ad.TechstackVersion,
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
			return schema.GetSingleDeploymentData{}, pkg.ErrNotFound
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
	if len(deployments) == 0 {
		return nil, pkg.ErrNotFound
	}
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
			Port:              dep.Port,
			TechstackName:     dep.TechstackName,
			TechstackVersion:  dep.TechstackVersion,
			ContainerID:       dep.ContainerID,
			CreatedAt:         dep.CreatedAt.Time.String(),
			UpdatedAt:         dep.UpdatedAt.Time.String(),
		})
	}
	return res, nil
}
