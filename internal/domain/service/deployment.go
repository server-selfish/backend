package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/transport/http"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/moby/moby/client"
	deployment_repository "github.com/server-selfish/backend/internal/domain/repository/deployment"
	"github.com/server-selfish/backend/internal/domain/schema"
	"github.com/server-selfish/backend/internal/pkg"
)

type (
	DeploymentService interface {
		GetDeploymentsByProjectId(ctx context.Context, projectId pgtype.UUID) ([]schema.GetDeploymentData, error)
		GetDeploymentByDeploymentId(ctx context.Context, deploymentId pgtype.UUID) (schema.GetSingleDeploymentData, error)
		GetActiveDeploymentByDeploymentId(ctx context.Context, deploymentId pgtype.UUID) (schema.GetActiveDeploymentHistory, error)
		GetHistoryDeploymentByDeploymentId(ctx context.Context, deploymentId pgtype.UUID) ([]schema.GetHistoryDeploymentHistory, error)
		CreateDeployment(ctx context.Context, params deployment_repository.CreateDeploymentParams) error
		CreateNewDeploymentVersionByDeploymentId(ctx context.Context, params deployment_repository.CreateDeploymentHistoryParams) error
		DeleteDeploymentByDeploymentId(ctx context.Context, deploymentId pgtype.UUID) error
	}
	deploymentService struct {
		dr *deployment_repository.Queries
		tm pkg.TxManager
		dc *client.Client
	}
)

func NewDeploymentService(dr *deployment_repository.Queries, tm pkg.TxManager, dc *client.Client) DeploymentService {
	return &deploymentService{
		dr: dr,
		tm: tm,
		dc: dc,
	}
}

// CreateNewDeploymentVersion implements [DeploymentService].
func (d *deploymentService) CreateNewDeploymentVersionByDeploymentId(ctx context.Context, params deployment_repository.CreateDeploymentHistoryParams) error {
	activeContainer, err := d.dr.GetActiveDeploymentHistoryContainerByDeploymentId(ctx, params.DeploymentID)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("tmp/builds/%s/%s/%s", params.DeploymentID, params.Branch, params.Version)
	// begin transaction
	// deactivate + stop & remove container of active version
	// check Port
	// build image newest version (git clone, generate docker file based on techstack, and build)
	// run the container of newest version
	// insert deployment history and the container of the new version
	//
	if err := d.tm.WithTx(ctx, func(tx pgx.Tx) error {
		depQuery := d.dr.WithTx(tx)
		if err := depQuery.SetActiveDeploymentHistoryNonActiveByDeploymentId(ctx, params.DeploymentID); err != nil {
			return err
		}
		if _, err := d.dc.ContainerStop(ctx, activeContainer.ID.String(), client.ContainerStopOptions{}); err != nil {
			return err
		}

		// get access token by installation id

		repo, err := git.PlainClone(path, &git.CloneOptions{
			URL: params.GitRemoteUrl,
			Auth: &http.BasicAuth{
				Username: "x-access-token",
				Password: "",
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

		techstack, err := depQuery.GetTechstackByTechstackId(ctx, params.DeploymentTechstackID)
		if err != nil {
			return err
		}
		cmdJson := pkg.ShellToExecForm(params.RunCommand.String)
		tpt := schema.DockerFileTemplate{
			DockerBaseImage:    techstack.DockerBaseImage,
			DockerRuntimeImage: techstack.DockerRuntimeImage,
			BuildCommand:       params.BuildCommand.String,
			BuildFolder:        params.BuildFolder.String,
			RunCommand:         cmdJson,
		}
		template, err := pkg.ParseTemplateFromEmbed(pkg.GetFileNameByTechstack(techstack.Name), tpt)
		if err != nil {
			return err
		}
		dockerfilePath := fmt.Sprintf("%s/Dockerfile", path)
		if err := pkg.WriteFile(dockerfilePath, []byte(template)); err != nil {
			return err
		}

		params.CommitID = commitId
		params.CommitMsg = commitMsg
		params.Version = version
		return nil
	}); err != nil {
		// remove the folder in path /tmp/build deployment_id and etc
		// remove
		// rollback: start the container of the active version previously
		return err
	}
	go func() {
		d.dc.ContainerRemove(ctx, activeContainer.ID.String(), client.ContainerRemoveOptions{})
	}()
	return nil
}

// DeleteDeploymentByDeploymentId implements [DeploymentService].
func (d *deploymentService) DeleteDeploymentByDeploymentId(ctx context.Context, deploymentId pgtype.UUID) error {
	if err := d.dr.DeleteDeploymentByDeploymentId(ctx, deploymentId); err != nil {
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
func (d *deploymentService) GetHistoryDeploymentByDeploymentId(ctx context.Context, deploymentId pgtype.UUID) ([]schema.GetHistoryDeploymentHistory, error) {
	hd, err := d.dr.GetDeploymentHistoryByDeploymentId(ctx, deploymentId)
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
func (d *deploymentService) GetActiveDeploymentByDeploymentId(ctx context.Context, deploymentId pgtype.UUID) (schema.GetActiveDeploymentHistory, error) {
	ad, err := d.dr.GetActiveDeploymentHistoryByDeploymentId(ctx, deploymentId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return schema.GetActiveDeploymentHistory{}, pkg.ErrNotFound
		}
		return schema.GetActiveDeploymentHistory{}, err
	}
	res := schema.GetActiveDeploymentHistory{
		DeploymentHistoryID: ad.DeploymentHistoryID,
		GitRemoteURL:        ad.GitRemoteUrl,
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
func (d *deploymentService) GetDeploymentByDeploymentId(ctx context.Context, deploymentId pgtype.UUID) (schema.GetSingleDeploymentData, error) {
	deployment, err := d.dr.GetDeploymentByDeploymentId(ctx, deploymentId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return schema.GetSingleDeploymentData{}, pkg.ErrNotFound
		}
		return schema.GetSingleDeploymentData{}, err
	}
	res := schema.GetSingleDeploymentData{
		ID:        deployment.ID.String(),
		Name:      deployment.Name,
		CreatedAt: deployment.CreatedAt.Time.String(),
		UpdatedAt: deployment.UpdatedAt.Time.String(),
	}
	return res, nil
}

// GetDeploymentsByProjectId implements [DeploymentService].
func (d *deploymentService) GetDeploymentsByProjectId(ctx context.Context, projectId pgtype.UUID) ([]schema.GetDeploymentData, error) {
	deployments, err := d.dr.GetDeploymentsByProjectId(ctx, projectId)
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
