package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/containerd/errdefs"
	"github.com/go-git/go-billy/v6"
	"github.com/go-git/go-billy/v6/memfs"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/client"
	"github.com/go-git/go-git/v6/plumbing/transport/http"
	"github.com/go-git/go-git/v6/storage/memory"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	moby_client "github.com/moby/moby/client"
	"github.com/rs/zerolog"
	cache_repository "github.com/server-selfish/backend/internal/domain/repository/cache"
	container_repository "github.com/server-selfish/backend/internal/domain/repository/container"
	deployment_repository "github.com/server-selfish/backend/internal/domain/repository/deployment"
	"github.com/server-selfish/backend/internal/domain/schema"
	parentservice "github.com/server-selfish/backend/internal/domain/service"
	docker_infra "github.com/server-selfish/backend/internal/infra/docker"
	git_infra "github.com/server-selfish/backend/internal/infra/git"
	github_infra "github.com/server-selfish/backend/internal/infra/github"
	"github.com/server-selfish/backend/internal/pkg"
	defined_error "github.com/server-selfish/backend/internal/pkg/error"
)

type (
	DeploymentService interface {
		GetDeploymentsByProjectId(ctx context.Context, userId, projectId pgtype.UUID) ([]schema.GetDeploymentData, error)
		GetDeploymentByDeploymentId(ctx context.Context, userId, deploymentId pgtype.UUID) (schema.GetSingleDeploymentData, error)
		GetActiveDeploymentByDeploymentName(ctx context.Context, userId pgtype.UUID, projectName, deploymentName string) (schema.GetActiveDeploymentHistory, error)
		GetHistoryDeploymentByDeploymentName(ctx context.Context, userId pgtype.UUID, projectName, deploymentName string) ([]schema.GetHistoryDeploymentHistory, error)
		GetTechstackName(ctx context.Context) (schema.GetTechstackList, error)
		GetTechstackVersionByName(ctx context.Context, techstackName string) ([]schema.GetTechstackVersion, error)
		GetDeploymentSettings(ctx context.Context, userId pgtype.UUID, projectName, deploymentName string) (schema.GetDeploymentSettings, error)
		CreateNewDeployment(ctx context.Context, userID pgtype.UUID, installationID int64, params schema.CreateDeploymentHistoryParams) error
		UpdateDeploymentVersionToLatest(ctx context.Context, userID pgtype.UUID, projectName, deploymentName string) error
		UpdateDeployment(ctx context.Context, userID pgtype.UUID, params schema.UpdateDeploymentParams) error
		DeleteDeploymentByDeploymentName(ctx context.Context, userId pgtype.UUID, projectName, deploymentName string) error
		DeleteDeploymentByDeploymentId(ctx context.Context, userId, deploymentId pgtype.UUID) error
		buildAndRunContainer(ctx context.Context, p schema.BuildAndRunContainerParams) error
		ensureDockerNetwork(ctx context.Context, networkName string) error
		buildDockerImage(ctx context.Context, fs billy.Filesystem, imageTag string, buildArgs map[string]string) error
	}
	deploymentService struct {
		gs       parentservice.GithubAppService
		dr       *deployment_repository.Queries
		cr       *container_repository.Queries
		tm       pkg.TxManager
		conRep   container_repository.ContainerRepository
		gi       github_infra.GithubInfra
		cache    cache_repository.CacheRepository
		gitInfra git_infra.GitInfra
		log      zerolog.Logger
	}
)

func NewDeploymentService(
	dr *deployment_repository.Queries,
	gs parentservice.GithubAppService,
	cache cache_repository.CacheRepository,
	cr *container_repository.Queries,
	tm pkg.TxManager,
	conRep container_repository.ContainerRepository,
	gi github_infra.GithubInfra,
	gitInfra git_infra.GitInfra,
	log zerolog.Logger,
) DeploymentService {
	return &deploymentService{
		gs:       gs,
		dr:       dr,
		cr:       cr,
		tm:       tm,
		cache:    cache,
		conRep:   conRep,
		gi:       gi,
		gitInfra: gitInfra,
		log:      log,
	}
}

// UpdateDeployment implements [DeploymentService].
func (d *deploymentService) UpdateDeployment(ctx context.Context, userID pgtype.UUID, params schema.UpdateDeploymentParams) error {
	prevActiveData, err :=
		d.dr.GetActiveDeploymentDetailByDeploymentName(ctx, deployment_repository.GetActiveDeploymentDetailByDeploymentNameParams{
			UserID: userID,
			Name:   params.ProjectName,
			Name_2: params.DeploymentName,
		})
	if err != nil {
		return err
	}

	// ponytail: description-only change skips rebuild, single UPDATE
	onlyDesc, noChange := isOnlyDescriptionChange(prevActiveData, params)
	if noChange {
		return nil
	}
	if onlyDesc {
		return d.dr.UpdateDeploymentDescriptionByDeploymentId(ctx, deployment_repository.UpdateDeploymentDescriptionByDeploymentIdParams{
			UserID:      userID,
			ID:          prevActiveData.ID,
			Description: pgtype.Text{String: params.DeploymentDescription, Valid: true},
		})
	}

	reps, err := d.gs.ListInstallationRepositories(ctx, userID, prevActiveData.InstallationID)
	if err != nil {
		return err
	}

	found := false
	for _, val := range reps {
		if val.ID == int64(prevActiveData.RepositoryID) {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("repository's permission is not sufficient")
	}

	activeContainerName := prevActiveData.ContainerName
	var cn, in string
	// transaction start
	if mainErr := d.tm.WithTx(ctx, func(tx pgx.Tx) error {
		depQuery := d.dr.WithTx(tx)
		conQuery := d.cr.WithTx(tx)

		// query the latest commit from repository
		it, err := d.gi.CreateInstallationToken(ctx, prevActiveData.InstallationID)
		if err != nil {
			return err
		}
		fs := memfs.New()
		repo, err := d.gitInfra.Clone(memory.NewStorage(), fs, &git.CloneOptions{
			URL: prevActiveData.GitRemoteUrl,
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
		commitId, commitMsg, version, err := extractRepoMetaData(repo)
		if err != nil {
			return err
		}

		// container name
		cnUUID, err := uuid.NewV7()
		if err != nil {
			return err
		}
		newVersion := version
		if commitId == prevActiveData.CommitID {
			count, err := depQuery.GetDeploymentHistoryCountByCommitID(ctx, deployment_repository.GetDeploymentHistoryCountByCommitIDParams{
				ID:       prevActiveData.ID,
				CommitID: commitId,
			})
			if err != nil {
				return err
			}
			newVersion = nextOverlapVersion(prevActiveData.Version, count)
		}
		cn = fmt.Sprintf("%s.%s", cnUUID.String(), docker_infra.NormalizeDockerName(prevActiveData.DeploymentName))
		in = fmt.Sprintf("%s:%s", docker_infra.NormalizeDockerName(prevActiveData.DeploymentName), newVersion)
		// deactivate + stop & remove container of active version
		if err := depQuery.SetActiveDeploymentHistoryNonActiveByDeploymentId(ctx, deployment_repository.SetActiveDeploymentHistoryNonActiveByDeploymentIdParams{
			UserID:       userID,
			DeploymentID: prevActiveData.ID,
		}); err != nil {
			return err
		}
		// container stop
		if _, err := d.conRep.ContainerStop(ctx, activeContainerName, moby_client.ContainerStopOptions{}); err != nil {
			return err
		}

		techstack, err := depQuery.GetTechstackByTechstackId(ctx, params.DeploymentTechstackID)
		if err != nil {
			return err
		}

		mainFilePath := fmt.Sprintf("%s/%s", params.BuildFolder, params.MainFileName)
		runCommand := getRunCommandByTechstack(techstack.Name, mainFilePath, techstack.DockerBaseImage)

		// build and run new container
		if err := d.buildAndRunContainer(ctx, schema.BuildAndRunContainerParams{
			DepQuery:              depQuery,
			DeploymentId:          prevActiveData.ID,
			FileSystem:            fs,
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
			RunCommandJSON:        shellToExecForm(runCommand),
			DockerBaseImage:       techstack.DockerBaseImage,
			DockerRuntimeImage:    techstack.DockerRuntimeImage,
			TechstackName:         techstack.Name,
		}); err != nil {
			return err
		}
		// QUERY: save deployment history
		historyId, err := depQuery.CreateDeploymentHistory(ctx, deployment_repository.CreateDeploymentHistoryParams{
			DeploymentID:          prevActiveData.ID,
			Branch:                params.Branch,
			CommitID:              commitId,
			CommitMsg:             commitMsg,
			Version:               newVersion,
			DeploymentTechstackID: params.DeploymentTechstackID,
			BuildCommand:          pgtype.Text{String: params.BuildCommand, Valid: true},
			BuildFolder:           pgtype.Text{String: params.BuildFolder, Valid: true},
			RunCommand:            pgtype.Text{String: runCommand, Valid: true},
			MainFilePath:          pgtype.Text{String: params.MainFileName, Valid: true},
		})
		if err != nil {
			return err
		}

		// QUERY: save container
		conId, err := conQuery.CreateContainer(ctx, container_repository.CreateContainerParams{
			Name:                cn,
			ImageName:           in,
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
	}); mainErr != nil {
		// stop and remove built container
		if cn != "" {
			timeout := 10 * time.Second
			timeoutInt := int(timeout)
			if _, err := d.conRep.ContainerStop(ctx, cn, moby_client.ContainerStopOptions{
				Timeout: &timeoutInt,
			}); err != nil {
				log.Printf("failed to stop container: %v", err)
			}
			if _, err := d.conRep.ContainerRemove(ctx, cn, moby_client.ContainerRemoveOptions{}); err != nil {
				log.Printf("failed to remove container: %v", err)
			}
		}
		// remove image
		if in != "" {
			if _, err := d.conRep.ImageRemove(ctx, in, moby_client.ImageRemoveOptions{
				Force:         true,
				PruneChildren: true,
			}); err != nil {
				log.Printf("failed to remove image: %v", err)
			}
		}
		// start activeContainer
		if activeContainerName != "" {
			_, err = d.conRep.ContainerStart(ctx, activeContainerName, moby_client.ContainerStartOptions{})
			if err != nil {
				log.Printf("failed to start container: %v", err)
			}
		}
		return mainErr
	}
	// remove prev active deployment container
	go func() {
		if activeContainerName != "" {
			if _, err := d.conRep.ContainerRemove(ctx, activeContainerName, moby_client.ContainerRemoveOptions{}); err != nil {
				log.Printf("failed to remove container: %v", err)
			}
		}
	}()
	return nil
}

// UpdateDeploymentVersionToLatest implements [DeploymentService].
func (d *deploymentService) UpdateDeploymentVersionToLatest(ctx context.Context, userID pgtype.UUID, projectName, deploymentName string) error {
	prevActiveData, err :=
		d.dr.GetActiveDeploymentDetailByDeploymentName(ctx, deployment_repository.GetActiveDeploymentDetailByDeploymentNameParams{
			UserID: userID,
			Name:   projectName,
			Name_2: deploymentName,
		})
	if err != nil {
		return err
	}
	var portList []schema.Port
	if err := json.Unmarshal(prevActiveData.Port, &portList); err != nil {
		return err
	}
	var envList []schema.ENV
	if err := json.Unmarshal(prevActiveData.Env, &envList); err != nil {
		return err
	}

	reps, err := d.gs.ListInstallationRepositories(ctx, userID, prevActiveData.InstallationID)
	if err != nil {
		return err
	}

	found := false

	for _, val := range reps {
		if val.ID == int64(prevActiveData.RepositoryID) {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("repository's permission is not sufficient")
	}
	activeContainerName := prevActiveData.ContainerName
	var cn, in string
	// transaction start
	if mainErr := d.tm.WithTx(ctx, func(tx pgx.Tx) error {
		depQuery := d.dr.WithTx(tx)
		conQuery := d.cr.WithTx(tx)

		// query the latest commit from repository
		it, err := d.gi.CreateInstallationToken(ctx, prevActiveData.InstallationID)
		if err != nil {
			return err
		}
		fs := memfs.New()
		repo, err := d.gitInfra.Clone(memory.NewStorage(), fs, &git.CloneOptions{
			URL: prevActiveData.GitRemoteUrl,
			ClientOptions: []client.Option{
				client.WithHTTPAuth(&http.BasicAuth{
					Username: "x-access-token",
					Password: it.Token,
				}),
			},
			ReferenceName: plumbing.NewBranchReferenceName(prevActiveData.Branch),
			SingleBranch:  true,
			Depth:         1,
		})
		if err != nil {
			return err
		}
		commitId, commitMsg, version, err := extractRepoMetaData(repo)
		if err != nil {
			return err
		}

		// container name
		cnUUID, err := uuid.NewV7()
		if err != nil {
			return err
		}
		newVersion := version
		if commitId == prevActiveData.CommitID {
			count, err := depQuery.GetDeploymentHistoryCountByCommitID(ctx, deployment_repository.GetDeploymentHistoryCountByCommitIDParams{
				ID:       prevActiveData.ID,
				CommitID: commitId,
			})
			if err != nil {
				return err
			}
			newVersion = nextOverlapVersion(prevActiveData.Version, count)
		}
		cn = fmt.Sprintf("%s.%s", cnUUID.String(), docker_infra.NormalizeDockerName(prevActiveData.DeploymentName))
		in = fmt.Sprintf("%s:%s", docker_infra.NormalizeDockerName(prevActiveData.DeploymentName), newVersion)
		// deactivate + stop & remove container of active version
		if err := depQuery.SetActiveDeploymentHistoryNonActiveByDeploymentId(ctx, deployment_repository.SetActiveDeploymentHistoryNonActiveByDeploymentIdParams{
			UserID:       userID,
			DeploymentID: prevActiveData.ID,
		}); err != nil {
			return err
		}
		// container stop
		if _, err := d.conRep.ContainerStop(ctx, activeContainerName, moby_client.ContainerStopOptions{}); err != nil {
			return err
		}

		mainFilePath := fmt.Sprintf("%s/%s", prevActiveData.BuildFolder.String, prevActiveData.MainFilePath.String)
		runCommand := getRunCommandByTechstack(prevActiveData.TechstackName, mainFilePath, prevActiveData.DockerBaseImage)

		// build and run new container
		if err := d.buildAndRunContainer(ctx, schema.BuildAndRunContainerParams{
			DepQuery:              depQuery,
			DeploymentId:          prevActiveData.ID,
			FileSystem:            fs,
			ContainerName:         cn,
			ImageName:             in,
			BuildCommand:          prevActiveData.BuildCommand.String,
			BuildFolder:           prevActiveData.BuildFolder.String,
			Env:                   envList,
			Port:                  portList,
			DeploymentTechstackID: prevActiveData.TechstackID,
			ProjectName:           prevActiveData.ProjectName,
			MainFileName:          prevActiveData.MainFilePath.String,
			UserId:                userID,
			RunCommandJSON:        shellToExecForm(runCommand),
			DockerBaseImage:       prevActiveData.DockerBaseImage,
			DockerRuntimeImage:    prevActiveData.DockerRuntimeImage,
			TechstackName:         prevActiveData.TechstackName,
		}); err != nil {
			return err
		}
		// QUERY: save deployment history
		historyId, err := depQuery.CreateDeploymentHistory(ctx, deployment_repository.CreateDeploymentHistoryParams{
			DeploymentID:          prevActiveData.ID,
			Branch:                prevActiveData.Branch,
			CommitID:              commitId,
			CommitMsg:             commitMsg,
			Version:               newVersion,
			DeploymentTechstackID: prevActiveData.TechstackID,
			BuildCommand:          prevActiveData.BuildCommand,
			BuildFolder:           prevActiveData.BuildFolder,
			RunCommand:            pgtype.Text{String: runCommand, Valid: true},
			MainFilePath:          prevActiveData.MainFilePath,
		})
		if err != nil {
			return err
		}

		// QUERY: save container
		conId, err := conQuery.CreateContainer(ctx, container_repository.CreateContainerParams{
			Name:                cn,
			ImageName:           in,
			DeploymentHistoryID: historyId,
		})
		if err != nil {
			return err
		}

		// QUERY: save port and env
		keys := make([]string, 0, len(envList))
		values := make([]string, 0, len(envList))
		for _, val := range envList {
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

		external := make([]int32, 0, len(portList))
		internal := make([]int32, 0, len(portList))
		protocol := make([]string, 0, len(portList))

		for _, val := range portList {
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
	}); mainErr != nil {
		// stop and remove built container
		if cn != "" {
			timeout := 10 * time.Second
			timeoutInt := int(timeout)
			if _, err := d.conRep.ContainerStop(ctx, cn, moby_client.ContainerStopOptions{
				Timeout: &timeoutInt,
			}); err != nil {
				log.Printf("failed to stop container: %v", err)
			}
			if _, err := d.conRep.ContainerRemove(ctx, cn, moby_client.ContainerRemoveOptions{}); err != nil {
				log.Printf("failed to remove container: %v", err)
			}
		}
		// remove image
		if in != "" {
			if _, err := d.conRep.ImageRemove(ctx, in, moby_client.ImageRemoveOptions{
				Force:         true,
				PruneChildren: true,
			}); err != nil {
				log.Printf("failed to remove image: %v", err)
			}
		}
		// start activeContainer
		if activeContainerName != "" {
			_, err = d.conRep.ContainerStart(ctx, activeContainerName, moby_client.ContainerStartOptions{})
			if err != nil {
				log.Printf("failed to start container: %v", err)
			}
		}
		return mainErr
	}
	// remove prev active deployment container
	go func() {
		if activeContainerName != "" {
			if _, err := d.conRep.ContainerRemove(ctx, activeContainerName, moby_client.ContainerRemoveOptions{}); err != nil {
				log.Printf("failed to remove container: %v", err)
			}
		}
	}()
	return nil
}

// CreateNewDeploymentVersion implements [DeploymentService].
func (d *deploymentService) CreateNewDeployment(ctx context.Context, userID pgtype.UUID, installationID int64, params schema.CreateDeploymentHistoryParams) error {
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

	var cn, in, activeContainerName string

	// begin transaction
	if mainErr := d.tm.WithTx(ctx, func(tx pgx.Tx) error {
		depQuery := d.dr.WithTx(tx)
		conQuery := d.cr.WithTx(tx)
		depId, err := depQuery.UpsertDeployment(ctx, deployment_repository.UpsertDeploymentParams{
			Name:           params.DeploymentName,
			GitRemoteUrl:   rep.HTMLURL,
			InstallationID: installationID,
			Name_2:         params.ProjectName,
			UserID:         userID,
			RepositoryID:   int32(repId),
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

		it, err := d.gi.CreateInstallationToken(ctx, installationID)
		if err != nil {
			return err
		}
		fs := memfs.New()
		repo, err := d.gitInfra.Clone(memory.NewStorage(), fs, &git.CloneOptions{
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

		commitId, commitMsg, version, err := extractRepoMetaData(repo)
		if err != nil {
			return err
		}

		// container name
		cnUUID, err := uuid.NewV7()
		if err != nil {
			return err
		}
		cn = fmt.Sprintf("%s.%s", cnUUID.String(), docker_infra.NormalizeDockerName(deployment.Name))
		in = fmt.Sprintf("%s:%s", docker_infra.NormalizeDockerName(deployment.Name), version)

		// deactivate + stop & remove container of active version
		if err := depQuery.SetActiveDeploymentHistoryNonActiveByDeploymentId(ctx, deployment_repository.SetActiveDeploymentHistoryNonActiveByDeploymentIdParams{
			UserID:       userID,
			DeploymentID: depId,
		}); err != nil {
			return err
		}
		if activeContainerName != "" {
			if _, err := d.conRep.ContainerStop(ctx, activeContainer.Name, moby_client.ContainerStopOptions{}); err != nil {
				return err
			}
		}

		techstack, err := depQuery.GetTechstackByTechstackId(ctx, params.DeploymentTechstackID)
		if err != nil {
			return err
		}

		mainFilePath := fmt.Sprintf("%s/%s", params.BuildFolder, params.MainFileName)
		runCommand := getRunCommandByTechstack(techstack.Name, mainFilePath, techstack.DockerBaseImage)

		// build and run new container
		if err := d.buildAndRunContainer(ctx, schema.BuildAndRunContainerParams{
			DepQuery:              depQuery,
			DeploymentId:          depId,
			FileSystem:            fs,
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
			RunCommandJSON:        shellToExecForm(runCommand),
			DockerBaseImage:       techstack.DockerBaseImage,
			DockerRuntimeImage:    techstack.DockerRuntimeImage,
			TechstackName:         techstack.Name,
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
			RunCommand:            pgtype.Text{String: runCommand, Valid: true},
			MainFilePath:          pgtype.Text{String: params.MainFileName, Valid: true},
		})
		if err != nil {
			return err
		}

		// QUERY: save container
		conId, err := conQuery.CreateContainer(ctx, container_repository.CreateContainerParams{
			Name:                cn,
			ImageName:           in,
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
	}); mainErr != nil {

		// stop and remove built container
		if cn != "" {
			timeout := 10 * time.Second
			timeoutInt := int(timeout)
			if _, err := d.conRep.ContainerStop(ctx, cn, moby_client.ContainerStopOptions{
				Timeout: &timeoutInt,
			}); err != nil {
				log.Printf("failed to stop container: %v", err)
			}
			if _, err := d.conRep.ContainerRemove(ctx, cn, moby_client.ContainerRemoveOptions{}); err != nil {
				log.Printf("failed to remove container: %v", err)
			}
		}
		// remove image
		if in != "" {
			if _, err := d.conRep.ImageRemove(ctx, in, moby_client.ImageRemoveOptions{
				Force:         true,
				PruneChildren: true,
			}); err != nil {
				log.Printf("failed to remove image: %v", err)
			}
		}
		// start activeContainer
		if activeContainerName != "" {
			_, err = d.conRep.ContainerStart(ctx, activeContainerName, moby_client.ContainerStartOptions{})
			if err != nil {
				log.Printf("failed to start container: %v", err)
			}
		}
		return mainErr
	}
	go func() {
		if activeContainerName != "" {
			if _, err := d.conRep.ContainerRemove(ctx, activeContainerName, moby_client.ContainerRemoveOptions{}); err != nil {
				log.Printf("failed to remove container: %v", err)
			}
		}
	}()
	return nil
}

// DeleteDeploymentByName implements [DeploymentService].
func (d *deploymentService) DeleteDeploymentByDeploymentName(ctx context.Context, userId pgtype.UUID, projectName string, deploymentName string) error {
	// get container and image
	cnt, err := d.cr.GetActiveContainerByDelploymentName(ctx, container_repository.GetActiveContainerByDelploymentNameParams{
		UserID: userId,
		Name:   projectName,
		Name_2: deploymentName,
	})
	if err != nil {

		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
	}
	// stop, container
	if _, err := d.conRep.ContainerStop(ctx, cnt.Name, moby_client.ContainerStopOptions{}); err != nil {
		if !errdefs.IsNotFound(err) {
			return err
		}
	}
	if err := d.dr.DeleteDeploymentByDeploymentName(ctx, deployment_repository.DeleteDeploymentByDeploymentNameParams{
		UserID: userId,
		Name:   projectName,
		Name_2: deploymentName,
	}); err != nil {
		if _, cErr := d.conRep.ContainerStart(ctx, cnt.Name, moby_client.ContainerStartOptions{}); cErr != nil {
			d.log.Err(cErr).Msg("error to start container")
		}
		return err
	}
	var wg sync.WaitGroup
	wg.Go(func() {
		// remove container
		if _, rmErr := d.conRep.ContainerRemove(ctx, cnt.Name, moby_client.ContainerRemoveOptions{}); rmErr != nil {
			d.log.Err(rmErr).Msg("error to remove container")
		}
		// remove image
		if _, rmiErr := d.conRep.ImageRemove(ctx, cnt.ImageName, moby_client.ImageRemoveOptions{}); rmiErr != nil {
			d.log.Err(rmiErr).Msg("error to remove container")
		}
	})
	wg.Wait()
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

// GetDeploymentSetting implements [DeploymentService].
func (d *deploymentService) GetDeploymentSettings(ctx context.Context, userId pgtype.UUID, projectName string, deploymentName string) (schema.GetDeploymentSettings, error) {
	ad, err := d.dr.GetActiveDeploymentHistoryByDeploymentName(ctx, deployment_repository.GetActiveDeploymentHistoryByDeploymentNameParams{
		UserID: userId,
		Name:   projectName,
		Name_2: deploymentName,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return schema.GetDeploymentSettings{}, defined_error.ErrActiveDeploymentNotFound
		}
		return schema.GetDeploymentSettings{}, err
	}

	var portList []schema.Port
	if err := json.Unmarshal(ad.Port, &portList); err != nil {
		return schema.GetDeploymentSettings{}, err
	}
	var envList []schema.ENV
	if err := json.Unmarshal(ad.Env, &envList); err != nil {
		return schema.GetDeploymentSettings{}, err
	}
	t1 := ad.UpdatedAt.Time
	t2 := ad.DeploymentUpdatedAt.Time
	latest := t1
	if t2.After(t1) {
		latest = t2
	}
	res := schema.GetDeploymentSettings{
		DeploymentName:        ad.DeploymentName,
		DeploymentDescription: ad.DeploymentDescription.String,
		GithubAccount:         ad.GithubAccount.String,
		RemoteUrl:             ad.RemoteUrl,
		Branch:                ad.Branch,
		Port:                  portList,
		Env:                   envList,
		BuildCommand:          ad.BuildCommand.String,
		BuildFolder:           ad.BuildFolder.String,
		MainFilePath:          ad.MainFilePath.String,
		TechstackID:           ad.TechstackID,
		TechstackName:         ad.TechstackName,
		UpdatedAt:             latest.String(),
	}
	return res, nil
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

// GetHistoryDeploymentByDeploymentId implements [DeploymentService].
func (d *deploymentService) GetHistoryDeploymentByDeploymentName(ctx context.Context, userId pgtype.UUID, projectName, deploymentName string) ([]schema.GetHistoryDeploymentHistory, error) {
	hd, err := d.dr.GetDeploymentHistoryByDeploymentName(ctx, deployment_repository.GetDeploymentHistoryByDeploymentNameParams{
		UserID: userId,
		Name:   projectName,
		Name_2: deploymentName,
	})
	if err != nil {
		return nil, err
	}
	res := make([]schema.GetHistoryDeploymentHistory, 0, len(hd))
	for _, h := range hd {
		res = append(res, schema.GetHistoryDeploymentHistory{
			ID:                h.ID,
			Branch:            h.Branch,
			CommitID:          h.CommitID,
			CommitMessage:     h.CommitMessage,
			DeploymentVersion: h.DeploymentVersion,
			BuildCommand:      h.BuildCommand.String,
			TechstackID:       h.TechstackID,
			TechstackName:     h.TechstackName,
			TechstackVersion:  h.TechstackVersion,
			CreatedAt:         h.CreatedAt.Time.String(),
			UpdatedAt:         h.UpdatedAt.Time.String(),
		})
	}
	return res, nil
}

// GetActiveDeploymentByDeploymentName implements [DeploymentService].
func (d *deploymentService) GetActiveDeploymentByDeploymentName(ctx context.Context, userId pgtype.UUID, projectName, deploymentName string) (schema.GetActiveDeploymentHistory, error) {
	ad, err := d.dr.GetActiveDeploymentHistoryByDeploymentName(ctx, deployment_repository.GetActiveDeploymentHistoryByDeploymentNameParams{
		UserID: userId,
		Name:   projectName,
		Name_2: deploymentName,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return schema.GetActiveDeploymentHistory{}, defined_error.ErrActiveDeploymentNotFound
		}
		return schema.GetActiveDeploymentHistory{}, err
	}
	var portList []schema.Port
	if err := json.Unmarshal(ad.Port, &portList); err != nil {
		return schema.GetActiveDeploymentHistory{}, err
	}
	res := schema.GetActiveDeploymentHistory{
		DeploymentHistoryID: ad.DeploymentHistoryID,
		RemoteUrl:           ad.RemoteUrl,
		Branch:              ad.Branch,
		CommitId:            ad.CommitID,
		CommitMessage:       ad.CommitMessage,
		DeploymentVersion:   ad.DeploymentVersion,
		Port:                portList,
		BuildCommand:        ad.BuildCommand.String,
		TechstackID:         ad.TechstackID,
		TechstackName:       ad.TechstackName,
		TechstackVersion:    ad.TechstackVersion,
		ContainerName:       ad.ContainerName,
		CreatedAt:           ad.CreatedAt.Time.String(),
		UpdatedAt:           ad.UpdatedAt.Time.String(),
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
