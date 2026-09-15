package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/go-git/go-billy/v6"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	moby_client "github.com/moby/moby/client"
	"github.com/moby/moby/client/pkg/jsonmessage"
	"github.com/server-selfish/backend/internal/constant"
	"github.com/server-selfish/backend/internal/domain/schema"
	docker_infra "github.com/server-selfish/backend/internal/infra/docker"
)

// buildAndRunContainer implements [DeploymentService].`
func (d *deploymentService) buildAndRunContainer(ctx context.Context, p schema.BuildAndRunContainerParams) error {
	tpt := schema.DockerFileTemplate{
		DockerBaseImage:    p.DockerBaseImage,
		DockerRuntimeImage: p.DockerRuntimeImage,
		BuildFolder:        p.BuildFolder,
		BuildCommand:       p.BuildCommand,
		MainFileName:       p.MainFileName,
		RunCommand:         p.RunCommandJSON,
	}
	template, err := parseTemplateFromEmbed(getFileNameByTechstack(p.TechstackName), tpt)
	if err != nil {
		return err
	}
	dockerignoreTemplate, err := parseTemplateFromEmbed("dockerignore", nil)
	if err != nil {
		return err
	}

	// dockerfilePath := filepath.Join(p.Path, "Dockerfile")
	if err := docker_infra.WriteFileToBillyFs(p.FileSystem, "Dockerfile", []byte(template)); err != nil {
		return err
	}

	// dockerignorePath := filepath.Join(p.Path, ".dockerignore")
	if err := docker_infra.WriteFileToBillyFs(p.FileSystem, ".dockerignore", []byte(dockerignoreTemplate)); err != nil {
		return err
	}

	// build docker image
	if err = d.buildDockerImage(ctx, p.FileSystem, p.ImageName, map[string]string{}); err != nil {
		return err
	}
	// network name
	nn := fmt.Sprintf("%s-network", docker_infra.NormalizeDockerName(p.ProjectName))
	if err := d.ensureDockerNetwork(ctx, nn); err != nil {
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

	if _, err := d.conRep.ContainerCreate(ctx, moby_client.ContainerCreateOptions{
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

	_, err = d.conRep.ContainerStart(ctx, p.ContainerName, moby_client.ContainerStartOptions{})
	if err != nil {
		return err
	}
	return nil
}

// buildDockerImage implements [DeploymentService].
func (d *deploymentService) buildDockerImage(ctx context.Context, fs billy.Filesystem, imageTag string, buildArgs map[string]string) error {
	excludes, err := docker_infra.ReadDockerignore(fs)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	// excludes = append(excludes, "!"+filepath.ToSlash("Dockerfile"))
	buildCtx, err := docker_infra.TarFilesystem(fs, excludes)
	if err != nil {
		return err
	}
	defer func() {
		if err := buildCtx.Close(); err != nil {
			log.Printf("failed to close build context: %v", err)
		}
	}()

	apiBuildArgs := map[string]*string{}
	for k, v := range buildArgs {
		vv := v
		apiBuildArgs[k] = &vv
	}

	resp, err := d.conRep.ImageBuild(ctx, buildCtx, moby_client.ImageBuildOptions{
		Tags:        []string{imageTag},
		Dockerfile:  "Dockerfile",
		PullParent:  true,
		Remove:      true,
		BuildArgs:   apiBuildArgs,
		ForceRemove: true,
	})
	if err != nil {
		return err
	}
	defer func() {
		if _, err := d.conRep.ImagePrune(ctx, moby_client.ImagePruneOptions{
			Filters: moby_client.Filters{},
		}); err != nil {
			log.Printf("failed remove dangling image: %v", err)
		}
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	err = jsonmessage.DisplayJSONMessagesStream(
		resp.Body,
		io.Discard,
		0,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

// ensureDockerNetwork implements [DeploymentService].
func (d *deploymentService) ensureDockerNetwork(ctx context.Context, networkName string) error {
	networks, err := d.conRep.NetworkList(ctx, moby_client.NetworkListOptions{})
	if err != nil {
		return err
	}
	for _, net := range networks.Items {
		if net.Name == networkName {
			return nil
		}
	}
	// Network does not exist, create it
	if _, err := d.conRep.NetworkCreate(ctx, networkName, moby_client.NetworkCreateOptions{
		Driver: "bridge",
	}); err != nil {
		return err
	}
	return nil
}
