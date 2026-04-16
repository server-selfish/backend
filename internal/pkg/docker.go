package pkg

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	build "github.com/docker/cli/cli/command/image/build"
	"github.com/moby/go-archive"
	"github.com/moby/moby/client"
)

// EnsureDockerNetwork checks if a Docker network exists, and creates it if not.
func EnsureDockerNetwork(ctx context.Context, dc *client.Client, networkName string) error {
	// List existing networks with the given name
	networks, err := dc.NetworkList(ctx, client.NetworkListOptions{})
	if err != nil {
		return err
	}
	for _, net := range networks.Items {
		if net.Name == networkName {
			return nil
		}
	}
	// Network does not exist, create it
	if _, err := dc.NetworkCreate(ctx, networkName, client.NetworkCreateOptions{
		Driver: "bridge",
	}); err != nil {
		return err
	}
	return nil
}

func BuildDockerImage(ctx context.Context, dc *client.Client, pPath string, imageTag string, buildArgs map[string]string) error {
	buildCtx, err := tarDirectoryWithDockerignore(pPath, "Dockerfile")
	if err != nil {
		return err
	}
	defer buildCtx.Close()

	apiBuildArgs := map[string]*string{}
	for k, v := range buildArgs {
		vv := v
		apiBuildArgs[k] = &vv
	}

	resp, err := dc.ImageBuild(ctx, buildCtx, client.ImageBuildOptions{
		Tags:       []string{imageTag}, // equivalent to: docker build -t imageTag
		Dockerfile: "Dockerfile",       // IMPORTANT: relative to build context tar root
		PullParent: true,               // like --pull
		Remove:     true,               // like --rm
		BuildArgs:  apiBuildArgs,       // --build-arg
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Optional: stream build output to logs/stdout
	_, _ = io.Copy(os.Stdout, resp.Body)
	return nil
}

func tarDirectoryWithDockerignore(dir string, dockerfileRel string) (io.ReadCloser, error) {
	excludes, err := build.ReadDockerignore(dir)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	excludes = append(excludes, "!"+filepath.ToSlash(dockerfileRel))

	return archive.TarWithOptions(dir, &archive.TarOptions{
		ExcludePatterns: excludes,
	})
}

func StopAndRemoveContainer(ctx context.Context, dc *client.Client, containerIDorName string) error {
	timeout := 10 * time.Second
	timeoutInt := int(timeout)
	if _, err := dc.ContainerStop(ctx, containerIDorName, client.ContainerStopOptions{
		Signal:  "",
		Timeout: &timeoutInt,
	}); err != nil {
		if strings.Contains(err.Error(), "is not running") {
		} else {
			return err
		}
	}
	_, err := dc.ContainerRemove(ctx, containerIDorName, client.ContainerRemoveOptions{})
	return err
}

func RemoveImage(ctx context.Context, dc *client.Client, imageIDorRef string) error {
	_, err := dc.ImageRemove(ctx, imageIDorRef, client.ImageRemoveOptions{
		Force:         true, // remove even if it has multiple tags (still won’t remove if in-use by a running container)
		PruneChildren: true, // also remove untagged parent images (where possible)
	})
	return err
}
