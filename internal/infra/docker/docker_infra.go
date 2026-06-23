package docker_infra

import (
	"context"
	"io"
	"net"

	"github.com/moby/moby/client"
)

type (
	DockerInfra interface {
		client.APIClient
	}
	dockerInfra struct {
		dc *client.Client
	}
)

func NewDockerInfra(
	dc *client.Client,
) DockerInfra {
	return dockerInfra{
		dc: dc,
	}
}

// BuildCachePrune implements [DockerInfra].
func (d dockerInfra) BuildCachePrune(ctx context.Context, opts client.BuildCachePruneOptions) (client.BuildCachePruneResult, error) {
	panic("unimplemented")
}

// BuildCancel implements [DockerInfra].
func (d dockerInfra) BuildCancel(ctx context.Context, id string, opts client.BuildCancelOptions) (client.BuildCancelResult, error) {
	panic("unimplemented")
}

// CheckpointCreate implements [DockerInfra].
func (d dockerInfra) CheckpointCreate(ctx context.Context, container string, options client.CheckpointCreateOptions) (client.CheckpointCreateResult, error) {
	panic("unimplemented")
}

// CheckpointList implements [DockerInfra].
func (d dockerInfra) CheckpointList(ctx context.Context, container string, options client.CheckpointListOptions) (client.CheckpointListResult, error) {
	panic("unimplemented")
}

// CheckpointRemove implements [DockerInfra].
func (d dockerInfra) CheckpointRemove(ctx context.Context, container string, options client.CheckpointRemoveOptions) (client.CheckpointRemoveResult, error) {
	panic("unimplemented")
}

// ClientVersion implements [DockerInfra].
func (d dockerInfra) ClientVersion() string {
	panic("unimplemented")
}

// Close implements [DockerInfra].
func (d dockerInfra) Close() error {
	panic("unimplemented")
}

// ConfigCreate implements [DockerInfra].
func (d dockerInfra) ConfigCreate(ctx context.Context, options client.ConfigCreateOptions) (client.ConfigCreateResult, error) {
	panic("unimplemented")
}

// ConfigInspect implements [DockerInfra].
func (d dockerInfra) ConfigInspect(ctx context.Context, id string, options client.ConfigInspectOptions) (client.ConfigInspectResult, error) {
	panic("unimplemented")
}

// ConfigList implements [DockerInfra].
func (d dockerInfra) ConfigList(ctx context.Context, options client.ConfigListOptions) (client.ConfigListResult, error) {
	panic("unimplemented")
}

// ConfigRemove implements [DockerInfra].
func (d dockerInfra) ConfigRemove(ctx context.Context, id string, options client.ConfigRemoveOptions) (client.ConfigRemoveResult, error) {
	panic("unimplemented")
}

// ConfigUpdate implements [DockerInfra].
func (d dockerInfra) ConfigUpdate(ctx context.Context, id string, options client.ConfigUpdateOptions) (client.ConfigUpdateResult, error) {
	panic("unimplemented")
}

// ContainerAttach implements [DockerInfra].
func (d dockerInfra) ContainerAttach(ctx context.Context, container string, options client.ContainerAttachOptions) (client.ContainerAttachResult, error) {
	panic("unimplemented")
}

// ContainerCommit implements [DockerInfra].
func (d dockerInfra) ContainerCommit(ctx context.Context, container string, options client.ContainerCommitOptions) (client.ContainerCommitResult, error) {
	panic("unimplemented")
}

// ContainerCreate implements [DockerInfra].
func (d dockerInfra) ContainerCreate(ctx context.Context, options client.ContainerCreateOptions) (client.ContainerCreateResult, error) {
	return d.dc.ContainerCreate(ctx, options)
}

// ContainerDiff implements [DockerInfra].
func (d dockerInfra) ContainerDiff(ctx context.Context, container string, options client.ContainerDiffOptions) (client.ContainerDiffResult, error) {
	panic("unimplemented")
}

// ContainerExport implements [DockerInfra].
func (d dockerInfra) ContainerExport(ctx context.Context, container string, options client.ContainerExportOptions) (client.ContainerExportResult, error) {
	panic("unimplemented")
}

// ContainerInspect implements [DockerInfra].
func (d dockerInfra) ContainerInspect(ctx context.Context, container string, options client.ContainerInspectOptions) (client.ContainerInspectResult, error) {
	return d.dc.ContainerInspect(ctx, container, options)
}

// ContainerKill implements [DockerInfra].
func (d dockerInfra) ContainerKill(ctx context.Context, container string, options client.ContainerKillOptions) (client.ContainerKillResult, error) {
	panic("unimplemented")
}

// ContainerList implements [DockerInfra].
func (d dockerInfra) ContainerList(ctx context.Context, options client.ContainerListOptions) (client.ContainerListResult, error) {
	panic("unimplemented")
}

// ContainerLogs implements [DockerInfra].
func (d dockerInfra) ContainerLogs(ctx context.Context, container string, options client.ContainerLogsOptions) (client.ContainerLogsResult, error) {
	return d.dc.ContainerLogs(ctx, container, options)
}

// ContainerPause implements [DockerInfra].
func (d dockerInfra) ContainerPause(ctx context.Context, container string, options client.ContainerPauseOptions) (client.ContainerPauseResult, error) {
	return d.dc.ContainerPause(ctx, container, options)
}

// ContainerPrune implements [DockerInfra].
func (d dockerInfra) ContainerPrune(ctx context.Context, opts client.ContainerPruneOptions) (client.ContainerPruneResult, error) {
	panic("unimplemented")
}

// ContainerRemove implements [DockerInfra].
func (d dockerInfra) ContainerRemove(ctx context.Context, container string, options client.ContainerRemoveOptions) (client.ContainerRemoveResult, error) {
	return d.dc.ContainerRemove(ctx, container, options)
}

// ContainerRename implements [DockerInfra].
func (d dockerInfra) ContainerRename(ctx context.Context, container string, options client.ContainerRenameOptions) (client.ContainerRenameResult, error) {
	panic("unimplemented")
}

// ContainerResize implements [DockerInfra].
func (d dockerInfra) ContainerResize(ctx context.Context, container string, options client.ContainerResizeOptions) (client.ContainerResizeResult, error) {
	panic("unimplemented")
}

// ContainerRestart implements [DockerInfra].
func (d dockerInfra) ContainerRestart(ctx context.Context, container string, options client.ContainerRestartOptions) (client.ContainerRestartResult, error) {
	return d.dc.ContainerRestart(ctx, container, options)
}

// ContainerStart implements [DockerInfra].
func (d dockerInfra) ContainerStart(ctx context.Context, container string, options client.ContainerStartOptions) (client.ContainerStartResult, error) {
	return d.dc.ContainerStart(ctx, container, options)
}

// ContainerStatPath implements [DockerInfra].
func (d dockerInfra) ContainerStatPath(ctx context.Context, container string, options client.ContainerStatPathOptions) (client.ContainerStatPathResult, error) {
	panic("unimplemented")
}

// ContainerStats implements [DockerInfra].
func (d dockerInfra) ContainerStats(ctx context.Context, container string, options client.ContainerStatsOptions) (client.ContainerStatsResult, error) {
	panic("unimplemented")
}

// ContainerStop implements [DockerInfra].
func (d dockerInfra) ContainerStop(ctx context.Context, container string, options client.ContainerStopOptions) (client.ContainerStopResult, error) {
	return d.dc.ContainerStop(ctx, container, options)
}

// ContainerTop implements [DockerInfra].
func (d dockerInfra) ContainerTop(ctx context.Context, container string, options client.ContainerTopOptions) (client.ContainerTopResult, error) {
	panic("unimplemented")
}

// ContainerUnpause implements [DockerInfra].
func (d dockerInfra) ContainerUnpause(ctx context.Context, container string, options client.ContainerUnpauseOptions) (client.ContainerUnpauseResult, error) {
	return d.dc.ContainerUnpause(ctx, container, options)
}

// ContainerUpdate implements [DockerInfra].
func (d dockerInfra) ContainerUpdate(ctx context.Context, container string, updateConfig client.ContainerUpdateOptions) (client.ContainerUpdateResult, error) {
	panic("unimplemented")
}

// ContainerWait implements [DockerInfra].
func (d dockerInfra) ContainerWait(ctx context.Context, container string, options client.ContainerWaitOptions) client.ContainerWaitResult {
	panic("unimplemented")
}

// CopyFromContainer implements [DockerInfra].
func (d dockerInfra) CopyFromContainer(ctx context.Context, container string, options client.CopyFromContainerOptions) (client.CopyFromContainerResult, error) {
	panic("unimplemented")
}

// CopyToContainer implements [DockerInfra].
func (d dockerInfra) CopyToContainer(ctx context.Context, container string, options client.CopyToContainerOptions) (client.CopyToContainerResult, error) {
	panic("unimplemented")
}

// DaemonHost implements [DockerInfra].
func (d dockerInfra) DaemonHost() string {
	panic("unimplemented")
}

// DialHijack implements [DockerInfra].
func (d dockerInfra) DialHijack(ctx context.Context, url string, proto string, meta map[string][]string) (net.Conn, error) {
	panic("unimplemented")
}

// Dialer implements [DockerInfra].
func (d dockerInfra) Dialer() func(context.Context) (net.Conn, error) {
	panic("unimplemented")
}

// DiskUsage implements [DockerInfra].
func (d dockerInfra) DiskUsage(ctx context.Context, options client.DiskUsageOptions) (client.DiskUsageResult, error) {
	panic("unimplemented")
}

// DistributionInspect implements [DockerInfra].
func (d dockerInfra) DistributionInspect(ctx context.Context, image string, options client.DistributionInspectOptions) (client.DistributionInspectResult, error) {
	panic("unimplemented")
}

// Events implements [DockerInfra].
func (d dockerInfra) Events(ctx context.Context, options client.EventsListOptions) client.EventsResult {
	panic("unimplemented")
}

// ExecAttach implements [DockerInfra].
func (d dockerInfra) ExecAttach(ctx context.Context, execID string, options client.ExecAttachOptions) (client.ExecAttachResult, error) {
	panic("unimplemented")
}

// ExecCreate implements [DockerInfra].
func (d dockerInfra) ExecCreate(ctx context.Context, container string, options client.ExecCreateOptions) (client.ExecCreateResult, error) {
	panic("unimplemented")
}

// ExecInspect implements [DockerInfra].
func (d dockerInfra) ExecInspect(ctx context.Context, execID string, options client.ExecInspectOptions) (client.ExecInspectResult, error) {
	panic("unimplemented")
}

// ExecResize implements [DockerInfra].
func (d dockerInfra) ExecResize(ctx context.Context, execID string, options client.ExecResizeOptions) (client.ExecResizeResult, error) {
	panic("unimplemented")
}

// ExecStart implements [DockerInfra].
func (d dockerInfra) ExecStart(ctx context.Context, execID string, options client.ExecStartOptions) (client.ExecStartResult, error) {
	panic("unimplemented")
}

// ImageAttestations implements [DockerInfra].
func (d dockerInfra) ImageAttestations(ctx context.Context, image string, _ ...client.ImageAttestationsOption) (client.ImageAttestationsResult, error) {
	panic("unimplemented")
}

// ImageBuild implements [DockerInfra].
func (d dockerInfra) ImageBuild(ctx context.Context, context io.Reader, options client.ImageBuildOptions) (client.ImageBuildResult, error) {
	return d.dc.ImageBuild(ctx, context, options)
}

// ImageHistory implements [DockerInfra].
func (d dockerInfra) ImageHistory(ctx context.Context, image string, _ ...client.ImageHistoryOption) (client.ImageHistoryResult, error) {
	panic("unimplemented")
}

// ImageImport implements [DockerInfra].
func (d dockerInfra) ImageImport(ctx context.Context, source client.ImageImportSource, ref string, options client.ImageImportOptions) (client.ImageImportResult, error) {
	panic("unimplemented")
}

// ImageInspect implements [DockerInfra].
func (d dockerInfra) ImageInspect(ctx context.Context, image string, _ ...client.ImageInspectOption) (client.ImageInspectResult, error) {
	panic("unimplemented")
}

// ImageList implements [DockerInfra].
func (d dockerInfra) ImageList(ctx context.Context, options client.ImageListOptions) (client.ImageListResult, error) {
	panic("unimplemented")
}

// ImageLoad implements [DockerInfra].
func (d dockerInfra) ImageLoad(ctx context.Context, input io.Reader, _ ...client.ImageLoadOption) (client.ImageLoadResult, error) {
	panic("unimplemented")
}

// ImagePrune implements [DockerInfra].
func (d dockerInfra) ImagePrune(ctx context.Context, opts client.ImagePruneOptions) (client.ImagePruneResult, error) {
	return d.dc.ImagePrune(ctx, opts)
}

// ImagePull implements [DockerInfra].
func (d dockerInfra) ImagePull(ctx context.Context, ref string, options client.ImagePullOptions) (client.ImagePullResponse, error) {
	panic("unimplemented")
}

// ImagePush implements [DockerInfra].
func (d dockerInfra) ImagePush(ctx context.Context, ref string, options client.ImagePushOptions) (client.ImagePushResponse, error) {
	panic("unimplemented")
}

// ImageRemove implements [DockerInfra].
func (d dockerInfra) ImageRemove(ctx context.Context, image string, options client.ImageRemoveOptions) (client.ImageRemoveResult, error) {
	return d.dc.ImageRemove(ctx, image, options)
}

// ImageSave implements [DockerInfra].
func (d dockerInfra) ImageSave(ctx context.Context, images []string, _ ...client.ImageSaveOption) (client.ImageSaveResult, error) {
	panic("unimplemented")
}

// ImageSearch implements [DockerInfra].
func (d dockerInfra) ImageSearch(ctx context.Context, term string, options client.ImageSearchOptions) (client.ImageSearchResult, error) {
	panic("unimplemented")
}

// ImageTag implements [DockerInfra].
func (d dockerInfra) ImageTag(ctx context.Context, options client.ImageTagOptions) (client.ImageTagResult, error) {
	panic("unimplemented")
}

// Info implements [DockerInfra].
func (d dockerInfra) Info(ctx context.Context, options client.InfoOptions) (client.SystemInfoResult, error) {
	panic("unimplemented")
}

// NetworkConnect implements [DockerInfra].
func (d dockerInfra) NetworkConnect(ctx context.Context, network string, options client.NetworkConnectOptions) (client.NetworkConnectResult, error) {
	panic("unimplemented")
}

// NetworkCreate implements [DockerInfra].
func (d dockerInfra) NetworkCreate(ctx context.Context, name string, options client.NetworkCreateOptions) (client.NetworkCreateResult, error) {
	return d.dc.NetworkCreate(ctx, name, options)
}

// NetworkDisconnect implements [DockerInfra].
func (d dockerInfra) NetworkDisconnect(ctx context.Context, network string, options client.NetworkDisconnectOptions) (client.NetworkDisconnectResult, error) {
	panic("unimplemented")
}

// NetworkInspect implements [DockerInfra].
func (d dockerInfra) NetworkInspect(ctx context.Context, network string, options client.NetworkInspectOptions) (client.NetworkInspectResult, error) {
	panic("unimplemented")
}

// NetworkList implements [DockerInfra].
func (d dockerInfra) NetworkList(ctx context.Context, options client.NetworkListOptions) (client.NetworkListResult, error) {
	return d.dc.NetworkList(ctx, options)
}

// NetworkPrune implements [DockerInfra].
func (d dockerInfra) NetworkPrune(ctx context.Context, opts client.NetworkPruneOptions) (client.NetworkPruneResult, error) {
	panic("unimplemented")
}

// NetworkRemove implements [DockerInfra].
func (d dockerInfra) NetworkRemove(ctx context.Context, network string, options client.NetworkRemoveOptions) (client.NetworkRemoveResult, error) {
	panic("unimplemented")
}

// NodeInspect implements [DockerInfra].
func (d dockerInfra) NodeInspect(ctx context.Context, nodeID string, options client.NodeInspectOptions) (client.NodeInspectResult, error) {
	panic("unimplemented")
}

// NodeList implements [DockerInfra].
func (d dockerInfra) NodeList(ctx context.Context, options client.NodeListOptions) (client.NodeListResult, error) {
	panic("unimplemented")
}

// NodeRemove implements [DockerInfra].
func (d dockerInfra) NodeRemove(ctx context.Context, nodeID string, options client.NodeRemoveOptions) (client.NodeRemoveResult, error) {
	panic("unimplemented")
}

// NodeUpdate implements [DockerInfra].
func (d dockerInfra) NodeUpdate(ctx context.Context, nodeID string, options client.NodeUpdateOptions) (client.NodeUpdateResult, error) {
	panic("unimplemented")
}

// Ping implements [DockerInfra].
func (d dockerInfra) Ping(ctx context.Context, options client.PingOptions) (client.PingResult, error) {
	panic("unimplemented")
}

// PluginCreate implements [DockerInfra].
func (d dockerInfra) PluginCreate(ctx context.Context, createContext io.Reader, options client.PluginCreateOptions) (client.PluginCreateResult, error) {
	panic("unimplemented")
}

// PluginDisable implements [DockerInfra].
func (d dockerInfra) PluginDisable(ctx context.Context, name string, options client.PluginDisableOptions) (client.PluginDisableResult, error) {
	panic("unimplemented")
}

// PluginEnable implements [DockerInfra].
func (d dockerInfra) PluginEnable(ctx context.Context, name string, options client.PluginEnableOptions) (client.PluginEnableResult, error) {
	panic("unimplemented")
}

// PluginInspect implements [DockerInfra].
func (d dockerInfra) PluginInspect(ctx context.Context, name string, options client.PluginInspectOptions) (client.PluginInspectResult, error) {
	panic("unimplemented")
}

// PluginInstall implements [DockerInfra].
func (d dockerInfra) PluginInstall(ctx context.Context, name string, options client.PluginInstallOptions) (client.PluginInstallResult, error) {
	panic("unimplemented")
}

// PluginList implements [DockerInfra].
func (d dockerInfra) PluginList(ctx context.Context, options client.PluginListOptions) (client.PluginListResult, error) {
	panic("unimplemented")
}

// PluginPush implements [DockerInfra].
func (d dockerInfra) PluginPush(ctx context.Context, name string, options client.PluginPushOptions) (client.PluginPushResult, error) {
	panic("unimplemented")
}

// PluginRemove implements [DockerInfra].
func (d dockerInfra) PluginRemove(ctx context.Context, name string, options client.PluginRemoveOptions) (client.PluginRemoveResult, error) {
	panic("unimplemented")
}

// PluginSet implements [DockerInfra].
func (d dockerInfra) PluginSet(ctx context.Context, name string, options client.PluginSetOptions) (client.PluginSetResult, error) {
	panic("unimplemented")
}

// PluginUpgrade implements [DockerInfra].
func (d dockerInfra) PluginUpgrade(ctx context.Context, name string, options client.PluginUpgradeOptions) (client.PluginUpgradeResult, error) {
	panic("unimplemented")
}

// RegistryLogin implements [DockerInfra].
func (d dockerInfra) RegistryLogin(ctx context.Context, auth client.RegistryLoginOptions) (client.RegistryLoginResult, error) {
	panic("unimplemented")
}

// SecretCreate implements [DockerInfra].
func (d dockerInfra) SecretCreate(ctx context.Context, options client.SecretCreateOptions) (client.SecretCreateResult, error) {
	panic("unimplemented")
}

// SecretInspect implements [DockerInfra].
func (d dockerInfra) SecretInspect(ctx context.Context, id string, options client.SecretInspectOptions) (client.SecretInspectResult, error) {
	panic("unimplemented")
}

// SecretList implements [DockerInfra].
func (d dockerInfra) SecretList(ctx context.Context, options client.SecretListOptions) (client.SecretListResult, error) {
	panic("unimplemented")
}

// SecretRemove implements [DockerInfra].
func (d dockerInfra) SecretRemove(ctx context.Context, id string, options client.SecretRemoveOptions) (client.SecretRemoveResult, error) {
	panic("unimplemented")
}

// SecretUpdate implements [DockerInfra].
func (d dockerInfra) SecretUpdate(ctx context.Context, id string, options client.SecretUpdateOptions) (client.SecretUpdateResult, error) {
	panic("unimplemented")
}

// ServerVersion implements [DockerInfra].
func (d dockerInfra) ServerVersion(ctx context.Context, options client.ServerVersionOptions) (client.ServerVersionResult, error) {
	panic("unimplemented")
}

// ServiceCreate implements [DockerInfra].
func (d dockerInfra) ServiceCreate(ctx context.Context, options client.ServiceCreateOptions) (client.ServiceCreateResult, error) {
	panic("unimplemented")
}

// ServiceInspect implements [DockerInfra].
func (d dockerInfra) ServiceInspect(ctx context.Context, serviceID string, options client.ServiceInspectOptions) (client.ServiceInspectResult, error) {
	panic("unimplemented")
}

// ServiceList implements [DockerInfra].
func (d dockerInfra) ServiceList(ctx context.Context, options client.ServiceListOptions) (client.ServiceListResult, error) {
	panic("unimplemented")
}

// ServiceLogs implements [DockerInfra].
func (d dockerInfra) ServiceLogs(ctx context.Context, serviceID string, options client.ServiceLogsOptions) (client.ServiceLogsResult, error) {
	panic("unimplemented")
}

// ServiceRemove implements [DockerInfra].
func (d dockerInfra) ServiceRemove(ctx context.Context, serviceID string, options client.ServiceRemoveOptions) (client.ServiceRemoveResult, error) {
	panic("unimplemented")
}

// ServiceUpdate implements [DockerInfra].
func (d dockerInfra) ServiceUpdate(ctx context.Context, serviceID string, options client.ServiceUpdateOptions) (client.ServiceUpdateResult, error) {
	panic("unimplemented")
}

// SwarmGetUnlockKey implements [DockerInfra].
func (d dockerInfra) SwarmGetUnlockKey(ctx context.Context) (client.SwarmGetUnlockKeyResult, error) {
	panic("unimplemented")
}

// SwarmInit implements [DockerInfra].
func (d dockerInfra) SwarmInit(ctx context.Context, options client.SwarmInitOptions) (client.SwarmInitResult, error) {
	panic("unimplemented")
}

// SwarmInspect implements [DockerInfra].
func (d dockerInfra) SwarmInspect(ctx context.Context, options client.SwarmInspectOptions) (client.SwarmInspectResult, error) {
	panic("unimplemented")
}

// SwarmJoin implements [DockerInfra].
func (d dockerInfra) SwarmJoin(ctx context.Context, options client.SwarmJoinOptions) (client.SwarmJoinResult, error) {
	panic("unimplemented")
}

// SwarmLeave implements [DockerInfra].
func (d dockerInfra) SwarmLeave(ctx context.Context, options client.SwarmLeaveOptions) (client.SwarmLeaveResult, error) {
	panic("unimplemented")
}

// SwarmUnlock implements [DockerInfra].
func (d dockerInfra) SwarmUnlock(ctx context.Context, options client.SwarmUnlockOptions) (client.SwarmUnlockResult, error) {
	panic("unimplemented")
}

// SwarmUpdate implements [DockerInfra].
func (d dockerInfra) SwarmUpdate(ctx context.Context, options client.SwarmUpdateOptions) (client.SwarmUpdateResult, error) {
	panic("unimplemented")
}

// TaskInspect implements [DockerInfra].
func (d dockerInfra) TaskInspect(ctx context.Context, taskID string, options client.TaskInspectOptions) (client.TaskInspectResult, error) {
	panic("unimplemented")
}

// TaskList implements [DockerInfra].
func (d dockerInfra) TaskList(ctx context.Context, options client.TaskListOptions) (client.TaskListResult, error) {
	panic("unimplemented")
}

// TaskLogs implements [DockerInfra].
func (d dockerInfra) TaskLogs(ctx context.Context, taskID string, options client.TaskLogsOptions) (client.TaskLogsResult, error) {
	panic("unimplemented")
}

// VolumeCreate implements [DockerInfra].
func (d dockerInfra) VolumeCreate(ctx context.Context, options client.VolumeCreateOptions) (client.VolumeCreateResult, error) {
	panic("unimplemented")
}

// VolumeInspect implements [DockerInfra].
func (d dockerInfra) VolumeInspect(ctx context.Context, volumeID string, options client.VolumeInspectOptions) (client.VolumeInspectResult, error) {
	panic("unimplemented")
}

// VolumeList implements [DockerInfra].
func (d dockerInfra) VolumeList(ctx context.Context, options client.VolumeListOptions) (client.VolumeListResult, error) {
	panic("unimplemented")
}

// VolumePrune implements [DockerInfra].
func (d dockerInfra) VolumePrune(ctx context.Context, options client.VolumePruneOptions) (client.VolumePruneResult, error) {
	panic("unimplemented")
}

// VolumeRemove implements [DockerInfra].
func (d dockerInfra) VolumeRemove(ctx context.Context, volumeID string, options client.VolumeRemoveOptions) (client.VolumeRemoveResult, error) {
	panic("unimplemented")
}

// VolumeUpdate implements [DockerInfra].
func (d dockerInfra) VolumeUpdate(ctx context.Context, volumeID string, options client.VolumeUpdateOptions) (client.VolumeUpdateResult, error) {
	panic("unimplemented")
}
