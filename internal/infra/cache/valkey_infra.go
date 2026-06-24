package cache_infra

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/valkey-io/valkey-go"
	"github.com/valkey-io/valkey-go/valkeycompat"
)

type (
	ValkeyInfra interface {
		// Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
		// SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error
		// Get(ctx context.Context, key string) (string, error)
		// GetJSON(ctx context.Context, key string, dest interface{}) error
		// Delete(ctx context.Context, key string) error
		valkeycompat.CoreCmdable
	}
	valkeyInfra struct {
		valkeyClient valkeycompat.Cmdable
		logger       zerolog.Logger
	}
)

func NewValkeyCache(valkeyClient valkey.Client, logger zerolog.Logger) ValkeyInfra {
	compat := valkeycompat.NewAdapter(valkeyClient)
	if err := compat.Ping(context.Background()).Err(); err != nil {
		logger.Fatal().Err(err).Msg("failed to ping valkey client")
	}
	return &valkeyInfra{
		valkeyClient: compat,
		logger:       logger,
	}
}

// ACLCat implements [ValkeyInfra].
func (v *valkeyInfra) ACLCat(ctx context.Context) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// ACLCatArgs implements [ValkeyInfra].
func (v *valkeyInfra) ACLCatArgs(ctx context.Context, options *valkeycompat.ACLCatArgs) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// ACLDelUser implements [ValkeyInfra].
func (v *valkeyInfra) ACLDelUser(ctx context.Context, username string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// ACLDryRun implements [ValkeyInfra].
func (v *valkeyInfra) ACLDryRun(ctx context.Context, username string, command ...any) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// ACLList implements [ValkeyInfra].
func (v *valkeyInfra) ACLList(ctx context.Context) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// ACLLog implements [ValkeyInfra].
func (v *valkeyInfra) ACLLog(ctx context.Context, count int64) *valkeycompat.ACLLogCmd {
	panic("unimplemented")
}

// ACLLogReset implements [ValkeyInfra].
func (v *valkeyInfra) ACLLogReset(ctx context.Context) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// ACLSetUser implements [ValkeyInfra].
func (v *valkeyInfra) ACLSetUser(ctx context.Context, username string, rules ...string) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// Append implements [ValkeyInfra].
func (v *valkeyInfra) Append(ctx context.Context, key string, value string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// BFAdd implements [ValkeyInfra].
func (v *valkeyInfra) BFAdd(ctx context.Context, key string, element interface{}) *valkeycompat.BoolCmd {
	panic("unimplemented")
}

// BFCard implements [ValkeyInfra].
func (v *valkeyInfra) BFCard(ctx context.Context, key string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// BFExists implements [ValkeyInfra].
func (v *valkeyInfra) BFExists(ctx context.Context, key string, element interface{}) *valkeycompat.BoolCmd {
	panic("unimplemented")
}

// BFInfo implements [ValkeyInfra].
func (v *valkeyInfra) BFInfo(ctx context.Context, key string) *valkeycompat.BFInfoCmd {
	panic("unimplemented")
}

// BFInfoArg implements [ValkeyInfra].
func (v *valkeyInfra) BFInfoArg(ctx context.Context, key string, option string) *valkeycompat.BFInfoCmd {
	panic("unimplemented")
}

// BFInfoCapacity implements [ValkeyInfra].
func (v *valkeyInfra) BFInfoCapacity(ctx context.Context, key string) *valkeycompat.BFInfoCmd {
	panic("unimplemented")
}

// BFInfoExpansion implements [ValkeyInfra].
func (v *valkeyInfra) BFInfoExpansion(ctx context.Context, key string) *valkeycompat.BFInfoCmd {
	panic("unimplemented")
}

// BFInfoFilters implements [ValkeyInfra].
func (v *valkeyInfra) BFInfoFilters(ctx context.Context, key string) *valkeycompat.BFInfoCmd {
	panic("unimplemented")
}

// BFInfoItems implements [ValkeyInfra].
func (v *valkeyInfra) BFInfoItems(ctx context.Context, key string) *valkeycompat.BFInfoCmd {
	panic("unimplemented")
}

// BFInfoSize implements [ValkeyInfra].
func (v *valkeyInfra) BFInfoSize(ctx context.Context, key string) *valkeycompat.BFInfoCmd {
	panic("unimplemented")
}

// BFInsert implements [ValkeyInfra].
func (v *valkeyInfra) BFInsert(ctx context.Context, key string, options *valkeycompat.BFInsertOptions, elements ...interface{}) *valkeycompat.BoolSliceCmd {
	panic("unimplemented")
}

// BFLoadChunk implements [ValkeyInfra].
func (v *valkeyInfra) BFLoadChunk(ctx context.Context, key string, iterator int64, data interface{}) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// BFMAdd implements [ValkeyInfra].
func (v *valkeyInfra) BFMAdd(ctx context.Context, key string, elements ...interface{}) *valkeycompat.BoolSliceCmd {
	panic("unimplemented")
}

// BFMExists implements [ValkeyInfra].
func (v *valkeyInfra) BFMExists(ctx context.Context, key string, elements ...interface{}) *valkeycompat.BoolSliceCmd {
	panic("unimplemented")
}

// BFReserve implements [ValkeyInfra].
func (v *valkeyInfra) BFReserve(ctx context.Context, key string, errorRate float64, capacity int64) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// BFReserveExpansion implements [ValkeyInfra].
func (v *valkeyInfra) BFReserveExpansion(ctx context.Context, key string, errorRate float64, capacity int64, expansion int64) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// BFReserveNonScaling implements [ValkeyInfra].
func (v *valkeyInfra) BFReserveNonScaling(ctx context.Context, key string, errorRate float64, capacity int64) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// BFReserveWithArgs implements [ValkeyInfra].
func (v *valkeyInfra) BFReserveWithArgs(ctx context.Context, key string, options *valkeycompat.BFReserveOptions) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// BFScanDump implements [ValkeyInfra].
func (v *valkeyInfra) BFScanDump(ctx context.Context, key string, iterator int64) *valkeycompat.ScanDumpCmd {
	panic("unimplemented")
}

// BLMPop implements [ValkeyInfra].
func (v *valkeyInfra) BLMPop(ctx context.Context, timeout time.Duration, direction string, count int64, keys ...string) *valkeycompat.KeyValuesCmd {
	panic("unimplemented")
}

// BLMove implements [ValkeyInfra].
func (v *valkeyInfra) BLMove(ctx context.Context, source string, destination string, srcpos string, destpos string, timeout time.Duration) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// BLPop implements [ValkeyInfra].
func (v *valkeyInfra) BLPop(ctx context.Context, timeout time.Duration, keys ...string) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// BRPop implements [ValkeyInfra].
func (v *valkeyInfra) BRPop(ctx context.Context, timeout time.Duration, keys ...string) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// BRPopLPush implements [ValkeyInfra].
func (v *valkeyInfra) BRPopLPush(ctx context.Context, source string, destination string, timeout time.Duration) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// BZMPop implements [ValkeyInfra].
func (v *valkeyInfra) BZMPop(ctx context.Context, timeout time.Duration, order string, count int64, keys ...string) *valkeycompat.ZSliceWithKeyCmd {
	panic("unimplemented")
}

// BZPopMax implements [ValkeyInfra].
func (v *valkeyInfra) BZPopMax(ctx context.Context, timeout time.Duration, keys ...string) *valkeycompat.ZWithKeyCmd {
	panic("unimplemented")
}

// BZPopMin implements [ValkeyInfra].
func (v *valkeyInfra) BZPopMin(ctx context.Context, timeout time.Duration, keys ...string) *valkeycompat.ZWithKeyCmd {
	panic("unimplemented")
}

// BgRewriteAOF implements [ValkeyInfra].
func (v *valkeyInfra) BgRewriteAOF(ctx context.Context) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// BgSave implements [ValkeyInfra].
func (v *valkeyInfra) BgSave(ctx context.Context) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// BitCount implements [ValkeyInfra].
func (v *valkeyInfra) BitCount(ctx context.Context, key string, bitCount *valkeycompat.BitCount) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// BitField implements [ValkeyInfra].
func (v *valkeyInfra) BitField(ctx context.Context, key string, args ...any) *valkeycompat.IntSliceCmd {
	panic("unimplemented")
}

// BitFieldRO implements [ValkeyInfra].
func (v *valkeyInfra) BitFieldRO(ctx context.Context, key string, values ...any) *valkeycompat.IntSliceCmd {
	panic("unimplemented")
}

// BitOpAnd implements [ValkeyInfra].
func (v *valkeyInfra) BitOpAnd(ctx context.Context, destKey string, keys ...string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// BitOpNot implements [ValkeyInfra].
func (v *valkeyInfra) BitOpNot(ctx context.Context, destKey string, key string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// BitOpOr implements [ValkeyInfra].
func (v *valkeyInfra) BitOpOr(ctx context.Context, destKey string, keys ...string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// BitOpXor implements [ValkeyInfra].
func (v *valkeyInfra) BitOpXor(ctx context.Context, destKey string, keys ...string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// BitPos implements [ValkeyInfra].
func (v *valkeyInfra) BitPos(ctx context.Context, key string, bit int64, pos ...int64) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// BitPosSpan implements [ValkeyInfra].
func (v *valkeyInfra) BitPosSpan(ctx context.Context, key string, bit int64, start int64, end int64, span string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// CFAdd implements [ValkeyInfra].
func (v *valkeyInfra) CFAdd(ctx context.Context, key string, element interface{}) *valkeycompat.BoolCmd {
	panic("unimplemented")
}

// CFAddNX implements [ValkeyInfra].
func (v *valkeyInfra) CFAddNX(ctx context.Context, key string, element interface{}) *valkeycompat.BoolCmd {
	panic("unimplemented")
}

// CFCount implements [ValkeyInfra].
func (v *valkeyInfra) CFCount(ctx context.Context, key string, element interface{}) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// CFDel implements [ValkeyInfra].
func (v *valkeyInfra) CFDel(ctx context.Context, key string, element interface{}) *valkeycompat.BoolCmd {
	panic("unimplemented")
}

// CFExists implements [ValkeyInfra].
func (v *valkeyInfra) CFExists(ctx context.Context, key string, element interface{}) *valkeycompat.BoolCmd {
	panic("unimplemented")
}

// CFInfo implements [ValkeyInfra].
func (v *valkeyInfra) CFInfo(ctx context.Context, key string) *valkeycompat.CFInfoCmd {
	panic("unimplemented")
}

// CFInsert implements [ValkeyInfra].
func (v *valkeyInfra) CFInsert(ctx context.Context, key string, options *valkeycompat.CFInsertOptions, elements ...interface{}) *valkeycompat.BoolSliceCmd {
	panic("unimplemented")
}

// CFInsertNX implements [ValkeyInfra].
func (v *valkeyInfra) CFInsertNX(ctx context.Context, key string, options *valkeycompat.CFInsertOptions, elements ...interface{}) *valkeycompat.IntSliceCmd {
	panic("unimplemented")
}

// CFLoadChunk implements [ValkeyInfra].
func (v *valkeyInfra) CFLoadChunk(ctx context.Context, key string, iterator int64, data interface{}) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// CFMExists implements [ValkeyInfra].
func (v *valkeyInfra) CFMExists(ctx context.Context, key string, elements ...interface{}) *valkeycompat.BoolSliceCmd {
	panic("unimplemented")
}

// CFReserve implements [ValkeyInfra].
func (v *valkeyInfra) CFReserve(ctx context.Context, key string, capacity int64) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// CFReserveBucketSize implements [ValkeyInfra].
func (v *valkeyInfra) CFReserveBucketSize(ctx context.Context, key string, capacity int64, bucketsize int64) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// CFReserveExpansion implements [ValkeyInfra].
func (v *valkeyInfra) CFReserveExpansion(ctx context.Context, key string, capacity int64, expansion int64) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// CFReserveMaxIterations implements [ValkeyInfra].
func (v *valkeyInfra) CFReserveMaxIterations(ctx context.Context, key string, capacity int64, maxiterations int64) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// CFReserveWithArgs implements [ValkeyInfra].
func (v *valkeyInfra) CFReserveWithArgs(ctx context.Context, key string, options *valkeycompat.CFReserveOptions) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// CFScanDump implements [ValkeyInfra].
func (v *valkeyInfra) CFScanDump(ctx context.Context, key string, iterator int64) *valkeycompat.ScanDumpCmd {
	panic("unimplemented")
}

// CMSIncrBy implements [ValkeyInfra].
func (v *valkeyInfra) CMSIncrBy(ctx context.Context, key string, elements ...interface{}) *valkeycompat.IntSliceCmd {
	panic("unimplemented")
}

// CMSInfo implements [ValkeyInfra].
func (v *valkeyInfra) CMSInfo(ctx context.Context, key string) *valkeycompat.CMSInfoCmd {
	panic("unimplemented")
}

// CMSInitByDim implements [ValkeyInfra].
func (v *valkeyInfra) CMSInitByDim(ctx context.Context, key string, width int64, height int64) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// CMSInitByProb implements [ValkeyInfra].
func (v *valkeyInfra) CMSInitByProb(ctx context.Context, key string, errorRate float64, probability float64) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// CMSMerge implements [ValkeyInfra].
func (v *valkeyInfra) CMSMerge(ctx context.Context, destKey string, sourceKeys ...string) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// CMSMergeWithWeight implements [ValkeyInfra].
func (v *valkeyInfra) CMSMergeWithWeight(ctx context.Context, destKey string, sourceKeys map[string]int64) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// CMSQuery implements [ValkeyInfra].
func (v *valkeyInfra) CMSQuery(ctx context.Context, key string, elements ...interface{}) *valkeycompat.IntSliceCmd {
	panic("unimplemented")
}

// ClientGetName implements [ValkeyInfra].
func (v *valkeyInfra) ClientGetName(ctx context.Context) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// ClientID implements [ValkeyInfra].
func (v *valkeyInfra) ClientID(ctx context.Context) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// ClientInfo implements [ValkeyInfra].
func (v *valkeyInfra) ClientInfo(ctx context.Context) *valkeycompat.ClientInfoCmd {
	panic("unimplemented")
}

// ClientKill implements [ValkeyInfra].
func (v *valkeyInfra) ClientKill(ctx context.Context, ipPort string) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// ClientKillByFilter implements [ValkeyInfra].
func (v *valkeyInfra) ClientKillByFilter(ctx context.Context, keys ...string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// ClientList implements [ValkeyInfra].
func (v *valkeyInfra) ClientList(ctx context.Context) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// ClientPause implements [ValkeyInfra].
func (v *valkeyInfra) ClientPause(ctx context.Context, dur time.Duration) *valkeycompat.BoolCmd {
	panic("unimplemented")
}

// ClientUnblock implements [ValkeyInfra].
func (v *valkeyInfra) ClientUnblock(ctx context.Context, id int64) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// ClientUnblockWithError implements [ValkeyInfra].
func (v *valkeyInfra) ClientUnblockWithError(ctx context.Context, id int64) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// ClientUnpause implements [ValkeyInfra].
func (v *valkeyInfra) ClientUnpause(ctx context.Context) *valkeycompat.BoolCmd {
	panic("unimplemented")
}

// ClusterAddSlots implements [ValkeyInfra].
func (v *valkeyInfra) ClusterAddSlots(ctx context.Context, slots ...int64) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// ClusterAddSlotsRange implements [ValkeyInfra].
func (v *valkeyInfra) ClusterAddSlotsRange(ctx context.Context, min int64, max int64) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// ClusterCountFailureReports implements [ValkeyInfra].
func (v *valkeyInfra) ClusterCountFailureReports(ctx context.Context, nodeID string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// ClusterCountKeysInSlot implements [ValkeyInfra].
func (v *valkeyInfra) ClusterCountKeysInSlot(ctx context.Context, slot int64) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// ClusterDelSlots implements [ValkeyInfra].
func (v *valkeyInfra) ClusterDelSlots(ctx context.Context, slots ...int64) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// ClusterDelSlotsRange implements [ValkeyInfra].
func (v *valkeyInfra) ClusterDelSlotsRange(ctx context.Context, min int64, max int64) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// ClusterFailover implements [ValkeyInfra].
func (v *valkeyInfra) ClusterFailover(ctx context.Context) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// ClusterForget implements [ValkeyInfra].
func (v *valkeyInfra) ClusterForget(ctx context.Context, nodeID string) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// ClusterGetKeysInSlot implements [ValkeyInfra].
func (v *valkeyInfra) ClusterGetKeysInSlot(ctx context.Context, slot int64, count int64) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// ClusterInfo implements [ValkeyInfra].
func (v *valkeyInfra) ClusterInfo(ctx context.Context) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// ClusterKeySlot implements [ValkeyInfra].
func (v *valkeyInfra) ClusterKeySlot(ctx context.Context, key string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// ClusterLinks implements [ValkeyInfra].
func (v *valkeyInfra) ClusterLinks(ctx context.Context) *valkeycompat.ClusterLinksCmd {
	panic("unimplemented")
}

// ClusterMeet implements [ValkeyInfra].
func (v *valkeyInfra) ClusterMeet(ctx context.Context, host string, port int64) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// ClusterMyShardID implements [ValkeyInfra].
func (v *valkeyInfra) ClusterMyShardID(ctx context.Context) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// ClusterNodes implements [ValkeyInfra].
func (v *valkeyInfra) ClusterNodes(ctx context.Context) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// ClusterReplicate implements [ValkeyInfra].
func (v *valkeyInfra) ClusterReplicate(ctx context.Context, nodeID string) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// ClusterResetHard implements [ValkeyInfra].
func (v *valkeyInfra) ClusterResetHard(ctx context.Context) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// ClusterResetSoft implements [ValkeyInfra].
func (v *valkeyInfra) ClusterResetSoft(ctx context.Context) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// ClusterSaveConfig implements [ValkeyInfra].
func (v *valkeyInfra) ClusterSaveConfig(ctx context.Context) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// ClusterShards implements [ValkeyInfra].
func (v *valkeyInfra) ClusterShards(ctx context.Context) *valkeycompat.ClusterShardsCmd {
	panic("unimplemented")
}

// ClusterSlaves implements [ValkeyInfra].
func (v *valkeyInfra) ClusterSlaves(ctx context.Context, nodeID string) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// ClusterSlots implements [ValkeyInfra].
func (v *valkeyInfra) ClusterSlots(ctx context.Context) *valkeycompat.ClusterSlotsCmd {
	panic("unimplemented")
}

// Command implements [ValkeyInfra].
func (v *valkeyInfra) Command(ctx context.Context) *valkeycompat.CommandsInfoCmd {
	panic("unimplemented")
}

// CommandGetKeys implements [ValkeyInfra].
func (v *valkeyInfra) CommandGetKeys(ctx context.Context, commands ...any) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// CommandGetKeysAndFlags implements [ValkeyInfra].
func (v *valkeyInfra) CommandGetKeysAndFlags(ctx context.Context, commands ...any) *valkeycompat.KeyFlagsCmd {
	panic("unimplemented")
}

// CommandList implements [ValkeyInfra].
func (v *valkeyInfra) CommandList(ctx context.Context, filter valkeycompat.FilterBy) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// ConfigGet implements [ValkeyInfra].
func (v *valkeyInfra) ConfigGet(ctx context.Context, parameter string) *valkeycompat.StringStringMapCmd {
	panic("unimplemented")
}

// ConfigResetStat implements [ValkeyInfra].
func (v *valkeyInfra) ConfigResetStat(ctx context.Context) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// ConfigRewrite implements [ValkeyInfra].
func (v *valkeyInfra) ConfigRewrite(ctx context.Context) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// ConfigSet implements [ValkeyInfra].
func (v *valkeyInfra) ConfigSet(ctx context.Context, parameter string, value string) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// Copy implements [ValkeyInfra].
func (v *valkeyInfra) Copy(ctx context.Context, sourceKey string, destKey string, db int64, replace bool) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// DBSize implements [ValkeyInfra].
func (v *valkeyInfra) DBSize(ctx context.Context) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// DebugObject implements [ValkeyInfra].
func (v *valkeyInfra) DebugObject(ctx context.Context, key string) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// Decr implements [ValkeyInfra].
func (v *valkeyInfra) Decr(ctx context.Context, key string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// DecrBy implements [ValkeyInfra].
func (v *valkeyInfra) DecrBy(ctx context.Context, key string, decrement int64) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// Del implements [ValkeyInfra].
func (v *valkeyInfra) Del(ctx context.Context, keys ...string) *valkeycompat.IntCmd {
	return v.valkeyClient.Del(ctx, keys...)
}

// Dump implements [ValkeyInfra].
func (v *valkeyInfra) Dump(ctx context.Context, key string) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// Echo implements [ValkeyInfra].
func (v *valkeyInfra) Echo(ctx context.Context, message any) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// Eval implements [ValkeyInfra].
func (v *valkeyInfra) Eval(ctx context.Context, script string, keys []string, args ...any) *valkeycompat.Cmd {
	panic("unimplemented")
}

// EvalRO implements [ValkeyInfra].
func (v *valkeyInfra) EvalRO(ctx context.Context, script string, keys []string, args ...any) *valkeycompat.Cmd {
	panic("unimplemented")
}

// EvalSha implements [ValkeyInfra].
func (v *valkeyInfra) EvalSha(ctx context.Context, sha1 string, keys []string, args ...any) *valkeycompat.Cmd {
	panic("unimplemented")
}

// EvalShaRO implements [ValkeyInfra].
func (v *valkeyInfra) EvalShaRO(ctx context.Context, sha1 string, keys []string, args ...any) *valkeycompat.Cmd {
	panic("unimplemented")
}

// Exists implements [ValkeyInfra].
func (v *valkeyInfra) Exists(ctx context.Context, keys ...string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// Expire implements [ValkeyInfra].
func (v *valkeyInfra) Expire(ctx context.Context, key string, expiration time.Duration) *valkeycompat.BoolCmd {
	panic("unimplemented")
}

// ExpireAt implements [ValkeyInfra].
func (v *valkeyInfra) ExpireAt(ctx context.Context, key string, tm time.Time) *valkeycompat.BoolCmd {
	panic("unimplemented")
}

// ExpireGT implements [ValkeyInfra].
func (v *valkeyInfra) ExpireGT(ctx context.Context, key string, expiration time.Duration) *valkeycompat.BoolCmd {
	panic("unimplemented")
}

// ExpireLT implements [ValkeyInfra].
func (v *valkeyInfra) ExpireLT(ctx context.Context, key string, expiration time.Duration) *valkeycompat.BoolCmd {
	panic("unimplemented")
}

// ExpireNX implements [ValkeyInfra].
func (v *valkeyInfra) ExpireNX(ctx context.Context, key string, expiration time.Duration) *valkeycompat.BoolCmd {
	panic("unimplemented")
}

// ExpireTime implements [ValkeyInfra].
func (v *valkeyInfra) ExpireTime(ctx context.Context, key string) *valkeycompat.DurationCmd {
	panic("unimplemented")
}

// ExpireXX implements [ValkeyInfra].
func (v *valkeyInfra) ExpireXX(ctx context.Context, key string, expiration time.Duration) *valkeycompat.BoolCmd {
	panic("unimplemented")
}

// FCall implements [ValkeyInfra].
func (v *valkeyInfra) FCall(ctx context.Context, function string, keys []string, args ...any) *valkeycompat.Cmd {
	panic("unimplemented")
}

// FCallRO implements [ValkeyInfra].
func (v *valkeyInfra) FCallRO(ctx context.Context, function string, keys []string, args ...any) *valkeycompat.Cmd {
	panic("unimplemented")
}

// FTAggregate implements [ValkeyInfra].
func (v *valkeyInfra) FTAggregate(ctx context.Context, index string, query string) *valkeycompat.MapStringInterfaceCmd {
	panic("unimplemented")
}

// FTAggregateWithArgs implements [ValkeyInfra].
func (v *valkeyInfra) FTAggregateWithArgs(ctx context.Context, index string, query string, options *valkeycompat.FTAggregateOptions) *valkeycompat.AggregateCmd {
	panic("unimplemented")
}

// FTAliasAdd implements [ValkeyInfra].
func (v *valkeyInfra) FTAliasAdd(ctx context.Context, index string, alias string) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// FTAliasDel implements [ValkeyInfra].
func (v *valkeyInfra) FTAliasDel(ctx context.Context, alias string) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// FTAliasUpdate implements [ValkeyInfra].
func (v *valkeyInfra) FTAliasUpdate(ctx context.Context, index string, alias string) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// FTAlter implements [ValkeyInfra].
func (v *valkeyInfra) FTAlter(ctx context.Context, index string, skipInitialScan bool, definition []interface{}) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// FTConfigGet implements [ValkeyInfra].
func (v *valkeyInfra) FTConfigGet(ctx context.Context, option string) *valkeycompat.MapMapStringInterfaceCmd {
	panic("unimplemented")
}

// FTConfigSet implements [ValkeyInfra].
func (v *valkeyInfra) FTConfigSet(ctx context.Context, option string, value interface{}) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// FTCreate implements [ValkeyInfra].
func (v *valkeyInfra) FTCreate(ctx context.Context, index string, options *valkeycompat.FTCreateOptions, schema ...*valkeycompat.FieldSchema) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// FTCursorDel implements [ValkeyInfra].
func (v *valkeyInfra) FTCursorDel(ctx context.Context, index string, cursorId int) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// FTCursorRead implements [ValkeyInfra].
func (v *valkeyInfra) FTCursorRead(ctx context.Context, index string, cursorId int, count int) *valkeycompat.MapStringInterfaceCmd {
	panic("unimplemented")
}

// FTDictAdd implements [ValkeyInfra].
func (v *valkeyInfra) FTDictAdd(ctx context.Context, dict string, term ...interface{}) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// FTDictDel implements [ValkeyInfra].
func (v *valkeyInfra) FTDictDel(ctx context.Context, dict string, term ...interface{}) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// FTDictDump implements [ValkeyInfra].
func (v *valkeyInfra) FTDictDump(ctx context.Context, dict string) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// FTDropIndex implements [ValkeyInfra].
func (v *valkeyInfra) FTDropIndex(ctx context.Context, index string) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// FTDropIndexWithArgs implements [ValkeyInfra].
func (v *valkeyInfra) FTDropIndexWithArgs(ctx context.Context, index string, options *valkeycompat.FTDropIndexOptions) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// FTExplain implements [ValkeyInfra].
func (v *valkeyInfra) FTExplain(ctx context.Context, index string, query string) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// FTExplainWithArgs implements [ValkeyInfra].
func (v *valkeyInfra) FTExplainWithArgs(ctx context.Context, index string, query string, options *valkeycompat.FTExplainOptions) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// FTInfo implements [ValkeyInfra].
func (v *valkeyInfra) FTInfo(ctx context.Context, index string) *valkeycompat.FTInfoCmd {
	panic("unimplemented")
}

// FTSearch implements [ValkeyInfra].
func (v *valkeyInfra) FTSearch(ctx context.Context, index string, query string) *valkeycompat.FTSearchCmd {
	panic("unimplemented")
}

// FTSearchWithArgs implements [ValkeyInfra].
func (v *valkeyInfra) FTSearchWithArgs(ctx context.Context, index string, query string, options *valkeycompat.FTSearchOptions) *valkeycompat.FTSearchCmd {
	panic("unimplemented")
}

// FTSpellCheck implements [ValkeyInfra].
func (v *valkeyInfra) FTSpellCheck(ctx context.Context, index string, query string) *valkeycompat.FTSpellCheckCmd {
	panic("unimplemented")
}

// FTSpellCheckWithArgs implements [ValkeyInfra].
func (v *valkeyInfra) FTSpellCheckWithArgs(ctx context.Context, index string, query string, options *valkeycompat.FTSpellCheckOptions) *valkeycompat.FTSpellCheckCmd {
	panic("unimplemented")
}

// FTSynDump implements [ValkeyInfra].
func (v *valkeyInfra) FTSynDump(ctx context.Context, index string) *valkeycompat.FTSynDumpCmd {
	panic("unimplemented")
}

// FTSynUpdate implements [ValkeyInfra].
func (v *valkeyInfra) FTSynUpdate(ctx context.Context, index string, synGroupId interface{}, terms []interface{}) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// FTSynUpdateWithArgs implements [ValkeyInfra].
func (v *valkeyInfra) FTSynUpdateWithArgs(ctx context.Context, index string, synGroupId interface{}, options *valkeycompat.FTSynUpdateOptions, terms []interface{}) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// FTTagVals implements [ValkeyInfra].
func (v *valkeyInfra) FTTagVals(ctx context.Context, index string, field string) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// FT_List implements [ValkeyInfra].
func (v *valkeyInfra) FT_List(ctx context.Context) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// FlushAll implements [ValkeyInfra].
func (v *valkeyInfra) FlushAll(ctx context.Context) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// FlushAllAsync implements [ValkeyInfra].
func (v *valkeyInfra) FlushAllAsync(ctx context.Context) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// FlushDB implements [ValkeyInfra].
func (v *valkeyInfra) FlushDB(ctx context.Context) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// FlushDBAsync implements [ValkeyInfra].
func (v *valkeyInfra) FlushDBAsync(ctx context.Context) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// FunctionDelete implements [ValkeyInfra].
func (v *valkeyInfra) FunctionDelete(ctx context.Context, libName string) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// FunctionDump implements [ValkeyInfra].
func (v *valkeyInfra) FunctionDump(ctx context.Context) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// FunctionFlush implements [ValkeyInfra].
func (v *valkeyInfra) FunctionFlush(ctx context.Context) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// FunctionFlushAsync implements [ValkeyInfra].
func (v *valkeyInfra) FunctionFlushAsync(ctx context.Context) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// FunctionKill implements [ValkeyInfra].
func (v *valkeyInfra) FunctionKill(ctx context.Context) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// FunctionList implements [ValkeyInfra].
func (v *valkeyInfra) FunctionList(ctx context.Context, q valkeycompat.FunctionListQuery) *valkeycompat.FunctionListCmd {
	panic("unimplemented")
}

// FunctionLoad implements [ValkeyInfra].
func (v *valkeyInfra) FunctionLoad(ctx context.Context, code string) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// FunctionLoadReplace implements [ValkeyInfra].
func (v *valkeyInfra) FunctionLoadReplace(ctx context.Context, code string) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// FunctionRestore implements [ValkeyInfra].
func (v *valkeyInfra) FunctionRestore(ctx context.Context, libDump string) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// FunctionStats implements [ValkeyInfra].
func (v *valkeyInfra) FunctionStats(ctx context.Context) *valkeycompat.FunctionStatsCmd {
	panic("unimplemented")
}

// GeoAdd implements [ValkeyInfra].
func (v *valkeyInfra) GeoAdd(ctx context.Context, key string, geoLocation ...valkeycompat.GeoLocation) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// GeoDist implements [ValkeyInfra].
func (v *valkeyInfra) GeoDist(ctx context.Context, key string, member1 string, member2 string, unit string) *valkeycompat.FloatCmd {
	panic("unimplemented")
}

// GeoHash implements [ValkeyInfra].
func (v *valkeyInfra) GeoHash(ctx context.Context, key string, members ...string) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// GeoPos implements [ValkeyInfra].
func (v *valkeyInfra) GeoPos(ctx context.Context, key string, members ...string) *valkeycompat.GeoPosCmd {
	panic("unimplemented")
}

// GeoRadius implements [ValkeyInfra].
func (v *valkeyInfra) GeoRadius(ctx context.Context, key string, longitude float64, latitude float64, query valkeycompat.GeoRadiusQuery) *valkeycompat.GeoLocationCmd {
	panic("unimplemented")
}

// GeoRadiusByMember implements [ValkeyInfra].
func (v *valkeyInfra) GeoRadiusByMember(ctx context.Context, key string, member string, query valkeycompat.GeoRadiusQuery) *valkeycompat.GeoLocationCmd {
	panic("unimplemented")
}

// GeoRadiusByMemberStore implements [ValkeyInfra].
func (v *valkeyInfra) GeoRadiusByMemberStore(ctx context.Context, key string, member string, query valkeycompat.GeoRadiusQuery) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// GeoRadiusStore implements [ValkeyInfra].
func (v *valkeyInfra) GeoRadiusStore(ctx context.Context, key string, longitude float64, latitude float64, query valkeycompat.GeoRadiusQuery) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// GeoSearch implements [ValkeyInfra].
func (v *valkeyInfra) GeoSearch(ctx context.Context, key string, q valkeycompat.GeoSearchQuery) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// GeoSearchLocation implements [ValkeyInfra].
func (v *valkeyInfra) GeoSearchLocation(ctx context.Context, key string, q valkeycompat.GeoSearchLocationQuery) *valkeycompat.GeoLocationCmd {
	panic("unimplemented")
}

// GeoSearchStore implements [ValkeyInfra].
func (v *valkeyInfra) GeoSearchStore(ctx context.Context, key string, store string, q valkeycompat.GeoSearchStoreQuery) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// Get implements [ValkeyInfra].
func (v *valkeyInfra) Get(ctx context.Context, key string) *valkeycompat.StringCmd {
	return v.valkeyClient.Get(ctx, key)
}

// GetBit implements [ValkeyInfra].
func (v *valkeyInfra) GetBit(ctx context.Context, key string, offset int64) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// GetDel implements [ValkeyInfra].
func (v *valkeyInfra) GetDel(ctx context.Context, key string) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// GetEx implements [ValkeyInfra].
func (v *valkeyInfra) GetEx(ctx context.Context, key string, expiration time.Duration) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// GetRange implements [ValkeyInfra].
func (v *valkeyInfra) GetRange(ctx context.Context, key string, start int64, end int64) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// GetSet implements [ValkeyInfra].
func (v *valkeyInfra) GetSet(ctx context.Context, key string, value any) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// HDel implements [ValkeyInfra].
func (v *valkeyInfra) HDel(ctx context.Context, key string, fields ...string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// HExists implements [ValkeyInfra].
func (v *valkeyInfra) HExists(ctx context.Context, key string, field string) *valkeycompat.BoolCmd {
	panic("unimplemented")
}

// HExpire implements [ValkeyInfra].
func (v *valkeyInfra) HExpire(ctx context.Context, key string, expiration time.Duration, fields ...string) *valkeycompat.IntSliceCmd {
	panic("unimplemented")
}

// HExpireAt implements [ValkeyInfra].
func (v *valkeyInfra) HExpireAt(ctx context.Context, key string, tm time.Time, fields ...string) *valkeycompat.IntSliceCmd {
	panic("unimplemented")
}

// HExpireAtWithArgs implements [ValkeyInfra].
func (v *valkeyInfra) HExpireAtWithArgs(ctx context.Context, key string, tm time.Time, expirationArgs valkeycompat.HExpireArgs, fields ...string) *valkeycompat.IntSliceCmd {
	panic("unimplemented")
}

// HExpireTime implements [ValkeyInfra].
func (v *valkeyInfra) HExpireTime(ctx context.Context, key string, fields ...string) *valkeycompat.IntSliceCmd {
	panic("unimplemented")
}

// HExpireWithArgs implements [ValkeyInfra].
func (v *valkeyInfra) HExpireWithArgs(ctx context.Context, key string, expiration time.Duration, expirationArgs valkeycompat.HExpireArgs, fields ...string) *valkeycompat.IntSliceCmd {
	panic("unimplemented")
}

// HGet implements [ValkeyInfra].
func (v *valkeyInfra) HGet(ctx context.Context, key string, field string) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// HGetAll implements [ValkeyInfra].
func (v *valkeyInfra) HGetAll(ctx context.Context, key string) *valkeycompat.StringStringMapCmd {
	panic("unimplemented")
}

// HGetDel implements [ValkeyInfra].
func (v *valkeyInfra) HGetDel(ctx context.Context, key string, fields ...string) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// HGetEX implements [ValkeyInfra].
func (v *valkeyInfra) HGetEX(ctx context.Context, key string, fields ...string) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// HGetEXWithArgs implements [ValkeyInfra].
func (v *valkeyInfra) HGetEXWithArgs(ctx context.Context, key string, options *valkeycompat.HGetEXOptions, fields ...string) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// HIncrBy implements [ValkeyInfra].
func (v *valkeyInfra) HIncrBy(ctx context.Context, key string, field string, incr int64) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// HIncrByFloat implements [ValkeyInfra].
func (v *valkeyInfra) HIncrByFloat(ctx context.Context, key string, field string, incr float64) *valkeycompat.FloatCmd {
	panic("unimplemented")
}

// HKeys implements [ValkeyInfra].
func (v *valkeyInfra) HKeys(ctx context.Context, key string) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// HLen implements [ValkeyInfra].
func (v *valkeyInfra) HLen(ctx context.Context, key string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// HMGet implements [ValkeyInfra].
func (v *valkeyInfra) HMGet(ctx context.Context, key string, fields ...string) *valkeycompat.SliceCmd {
	panic("unimplemented")
}

// HMSet implements [ValkeyInfra].
func (v *valkeyInfra) HMSet(ctx context.Context, key string, values ...any) *valkeycompat.BoolCmd {
	panic("unimplemented")
}

// HPExpire implements [ValkeyInfra].
func (v *valkeyInfra) HPExpire(ctx context.Context, key string, expiration time.Duration, fields ...string) *valkeycompat.IntSliceCmd {
	panic("unimplemented")
}

// HPExpireAt implements [ValkeyInfra].
func (v *valkeyInfra) HPExpireAt(ctx context.Context, key string, tm time.Time, fields ...string) *valkeycompat.IntSliceCmd {
	panic("unimplemented")
}

// HPExpireAtWithArgs implements [ValkeyInfra].
func (v *valkeyInfra) HPExpireAtWithArgs(ctx context.Context, key string, tm time.Time, expirationArgs valkeycompat.HExpireArgs, fields ...string) *valkeycompat.IntSliceCmd {
	panic("unimplemented")
}

// HPExpireTime implements [ValkeyInfra].
func (v *valkeyInfra) HPExpireTime(ctx context.Context, key string, fields ...string) *valkeycompat.IntSliceCmd {
	panic("unimplemented")
}

// HPExpireWithArgs implements [ValkeyInfra].
func (v *valkeyInfra) HPExpireWithArgs(ctx context.Context, key string, expiration time.Duration, expirationArgs valkeycompat.HExpireArgs, fields ...string) *valkeycompat.IntSliceCmd {
	panic("unimplemented")
}

// HPTTL implements [ValkeyInfra].
func (v *valkeyInfra) HPTTL(ctx context.Context, key string, fields ...string) *valkeycompat.IntSliceCmd {
	panic("unimplemented")
}

// HPersist implements [ValkeyInfra].
func (v *valkeyInfra) HPersist(ctx context.Context, key string, fields ...string) *valkeycompat.IntSliceCmd {
	panic("unimplemented")
}

// HRandField implements [ValkeyInfra].
func (v *valkeyInfra) HRandField(ctx context.Context, key string, count int64) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// HRandFieldWithValues implements [ValkeyInfra].
func (v *valkeyInfra) HRandFieldWithValues(ctx context.Context, key string, count int64) *valkeycompat.KeyValueSliceCmd {
	panic("unimplemented")
}

// HScan implements [ValkeyInfra].
func (v *valkeyInfra) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) *valkeycompat.ScanCmd {
	panic("unimplemented")
}

// HScanNoValues implements [ValkeyInfra].
func (v *valkeyInfra) HScanNoValues(ctx context.Context, key string, cursor uint64, match string, count int64) *valkeycompat.ScanCmd {
	panic("unimplemented")
}

// HSet implements [ValkeyInfra].
func (v *valkeyInfra) HSet(ctx context.Context, key string, values ...any) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// HSetEX implements [ValkeyInfra].
func (v *valkeyInfra) HSetEX(ctx context.Context, key string, fieldsAndValues ...string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// HSetEXWithArgs implements [ValkeyInfra].
func (v *valkeyInfra) HSetEXWithArgs(ctx context.Context, key string, options *valkeycompat.HSetEXOptions, fieldsAndValues ...string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// HSetNX implements [ValkeyInfra].
func (v *valkeyInfra) HSetNX(ctx context.Context, key string, field string, value any) *valkeycompat.BoolCmd {
	panic("unimplemented")
}

// HStrLen implements [ValkeyInfra].
func (v *valkeyInfra) HStrLen(ctx context.Context, key string, field string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// HTTL implements [ValkeyInfra].
func (v *valkeyInfra) HTTL(ctx context.Context, key string, fields ...string) *valkeycompat.IntSliceCmd {
	panic("unimplemented")
}

// HVals implements [ValkeyInfra].
func (v *valkeyInfra) HVals(ctx context.Context, key string) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// Incr implements [ValkeyInfra].
func (v *valkeyInfra) Incr(ctx context.Context, key string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// IncrBy implements [ValkeyInfra].
func (v *valkeyInfra) IncrBy(ctx context.Context, key string, value int64) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// IncrByFloat implements [ValkeyInfra].
func (v *valkeyInfra) IncrByFloat(ctx context.Context, key string, value float64) *valkeycompat.FloatCmd {
	panic("unimplemented")
}

// Info implements [ValkeyInfra].
func (v *valkeyInfra) Info(ctx context.Context, section ...string) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// JSONArrAppend implements [ValkeyInfra].
func (v *valkeyInfra) JSONArrAppend(ctx context.Context, key string, path string, values ...interface{}) *valkeycompat.IntSliceCmd {
	panic("unimplemented")
}

// JSONArrIndex implements [ValkeyInfra].
func (v *valkeyInfra) JSONArrIndex(ctx context.Context, key string, path string, value ...interface{}) *valkeycompat.IntSliceCmd {
	panic("unimplemented")
}

// JSONArrIndexWithArgs implements [ValkeyInfra].
func (v *valkeyInfra) JSONArrIndexWithArgs(ctx context.Context, key string, path string, options *valkeycompat.JSONArrIndexArgs, value ...interface{}) *valkeycompat.IntSliceCmd {
	panic("unimplemented")
}

// JSONArrInsert implements [ValkeyInfra].
func (v *valkeyInfra) JSONArrInsert(ctx context.Context, key string, path string, index int64, values ...interface{}) *valkeycompat.IntSliceCmd {
	panic("unimplemented")
}

// JSONArrLen implements [ValkeyInfra].
func (v *valkeyInfra) JSONArrLen(ctx context.Context, key string, path string) *valkeycompat.IntSliceCmd {
	panic("unimplemented")
}

// JSONArrPop implements [ValkeyInfra].
func (v *valkeyInfra) JSONArrPop(ctx context.Context, key string, path string, index int) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// JSONArrTrim implements [ValkeyInfra].
func (v *valkeyInfra) JSONArrTrim(ctx context.Context, key string, path string) *valkeycompat.IntSliceCmd {
	panic("unimplemented")
}

// JSONArrTrimWithArgs implements [ValkeyInfra].
func (v *valkeyInfra) JSONArrTrimWithArgs(ctx context.Context, key string, path string, options *valkeycompat.JSONArrTrimArgs) *valkeycompat.IntSliceCmd {
	panic("unimplemented")
}

// JSONClear implements [ValkeyInfra].
func (v *valkeyInfra) JSONClear(ctx context.Context, key string, path string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// JSONDebugMemory implements [ValkeyInfra].
func (v *valkeyInfra) JSONDebugMemory(ctx context.Context, key string, path string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// JSONDel implements [ValkeyInfra].
func (v *valkeyInfra) JSONDel(ctx context.Context, key string, path string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// JSONForget implements [ValkeyInfra].
func (v *valkeyInfra) JSONForget(ctx context.Context, key string, path string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// JSONGet implements [ValkeyInfra].
func (v *valkeyInfra) JSONGet(ctx context.Context, key string, paths ...string) *valkeycompat.JSONCmd {
	return v.valkeyClient.JSONGet(ctx, key, paths...)
}

// JSONGetWithArgs implements [ValkeyInfra].
func (v *valkeyInfra) JSONGetWithArgs(ctx context.Context, key string, options *valkeycompat.JSONGetArgs, paths ...string) *valkeycompat.JSONCmd {
	panic("unimplemented")
}

// JSONMGet implements [ValkeyInfra].
func (v *valkeyInfra) JSONMGet(ctx context.Context, path string, keys ...string) *valkeycompat.JSONSliceCmd {
	panic("unimplemented")
}

// JSONMSet implements [ValkeyInfra].
func (v *valkeyInfra) JSONMSet(ctx context.Context, params ...interface{}) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// JSONMSetArgs implements [ValkeyInfra].
func (v *valkeyInfra) JSONMSetArgs(ctx context.Context, docs []valkeycompat.JSONSetArgs) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// JSONMerge implements [ValkeyInfra].
func (v *valkeyInfra) JSONMerge(ctx context.Context, key string, path string, value string) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// JSONNumIncrBy implements [ValkeyInfra].
func (v *valkeyInfra) JSONNumIncrBy(ctx context.Context, key string, path string, value float64) *valkeycompat.JSONCmd {
	panic("unimplemented")
}

// JSONObjKeys implements [ValkeyInfra].
func (v *valkeyInfra) JSONObjKeys(ctx context.Context, key string, path string) *valkeycompat.SliceCmd {
	panic("unimplemented")
}

// JSONObjLen implements [ValkeyInfra].
func (v *valkeyInfra) JSONObjLen(ctx context.Context, key string, path string) *valkeycompat.IntPointerSliceCmd {
	panic("unimplemented")
}

// JSONSet implements [ValkeyInfra].
func (v *valkeyInfra) JSONSet(ctx context.Context, key string, path string, value interface{}) *valkeycompat.StatusCmd {
	return v.valkeyClient.JSONMSet(ctx, key, path, value)
}

// JSONSetMode implements [ValkeyInfra].
func (v *valkeyInfra) JSONSetMode(ctx context.Context, key string, path string, value interface{}, mode string) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// JSONStrAppend implements [ValkeyInfra].
func (v *valkeyInfra) JSONStrAppend(ctx context.Context, key string, path string, value string) *valkeycompat.IntPointerSliceCmd {
	panic("unimplemented")
}

// JSONStrLen implements [ValkeyInfra].
func (v *valkeyInfra) JSONStrLen(ctx context.Context, key string, path string) *valkeycompat.IntPointerSliceCmd {
	panic("unimplemented")
}

// JSONToggle implements [ValkeyInfra].
func (v *valkeyInfra) JSONToggle(ctx context.Context, key string, path string) *valkeycompat.IntPointerSliceCmd {
	panic("unimplemented")
}

// JSONType implements [ValkeyInfra].
func (v *valkeyInfra) JSONType(ctx context.Context, key string, path string) *valkeycompat.JSONSliceCmd {
	panic("unimplemented")
}

// Keys implements [ValkeyInfra].
func (v *valkeyInfra) Keys(ctx context.Context, pattern string) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// LCS implements [ValkeyInfra].
func (v *valkeyInfra) LCS(ctx context.Context, q *valkeycompat.LCSQuery) *valkeycompat.LCSCmd {
	panic("unimplemented")
}

// LIndex implements [ValkeyInfra].
func (v *valkeyInfra) LIndex(ctx context.Context, key string, index int64) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// LInsert implements [ValkeyInfra].
func (v *valkeyInfra) LInsert(ctx context.Context, key string, op string, pivot any, value any) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// LInsertAfter implements [ValkeyInfra].
func (v *valkeyInfra) LInsertAfter(ctx context.Context, key string, pivot any, value any) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// LInsertBefore implements [ValkeyInfra].
func (v *valkeyInfra) LInsertBefore(ctx context.Context, key string, pivot any, value any) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// LLen implements [ValkeyInfra].
func (v *valkeyInfra) LLen(ctx context.Context, key string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// LMPop implements [ValkeyInfra].
func (v *valkeyInfra) LMPop(ctx context.Context, direction string, count int64, keys ...string) *valkeycompat.KeyValuesCmd {
	panic("unimplemented")
}

// LMove implements [ValkeyInfra].
func (v *valkeyInfra) LMove(ctx context.Context, source string, destination string, srcpos string, destpos string) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// LPop implements [ValkeyInfra].
func (v *valkeyInfra) LPop(ctx context.Context, key string) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// LPopCount implements [ValkeyInfra].
func (v *valkeyInfra) LPopCount(ctx context.Context, key string, count int64) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// LPos implements [ValkeyInfra].
func (v *valkeyInfra) LPos(ctx context.Context, key string, value string, args valkeycompat.LPosArgs) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// LPosCount implements [ValkeyInfra].
func (v *valkeyInfra) LPosCount(ctx context.Context, key string, value string, count int64, args valkeycompat.LPosArgs) *valkeycompat.IntSliceCmd {
	panic("unimplemented")
}

// LPush implements [ValkeyInfra].
func (v *valkeyInfra) LPush(ctx context.Context, key string, values ...any) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// LPushX implements [ValkeyInfra].
func (v *valkeyInfra) LPushX(ctx context.Context, key string, values ...any) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// LRange implements [ValkeyInfra].
func (v *valkeyInfra) LRange(ctx context.Context, key string, start int64, stop int64) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// LRem implements [ValkeyInfra].
func (v *valkeyInfra) LRem(ctx context.Context, key string, count int64, value any) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// LSet implements [ValkeyInfra].
func (v *valkeyInfra) LSet(ctx context.Context, key string, index int64, value any) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// LTrim implements [ValkeyInfra].
func (v *valkeyInfra) LTrim(ctx context.Context, key string, start int64, stop int64) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// LastSave implements [ValkeyInfra].
func (v *valkeyInfra) LastSave(ctx context.Context) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// MGet implements [ValkeyInfra].
func (v *valkeyInfra) MGet(ctx context.Context, keys ...string) *valkeycompat.SliceCmd {
	panic("unimplemented")
}

// MSet implements [ValkeyInfra].
func (v *valkeyInfra) MSet(ctx context.Context, values ...any) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// MSetNX implements [ValkeyInfra].
func (v *valkeyInfra) MSetNX(ctx context.Context, values ...any) *valkeycompat.BoolCmd {
	panic("unimplemented")
}

// MemoryUsage implements [ValkeyInfra].
func (v *valkeyInfra) MemoryUsage(ctx context.Context, key string, samples ...int64) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// Migrate implements [ValkeyInfra].
func (v *valkeyInfra) Migrate(ctx context.Context, host string, port int64, key string, db int64, timeout time.Duration) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// ModuleLoadex implements [ValkeyInfra].
func (v *valkeyInfra) ModuleLoadex(ctx context.Context, conf *valkeycompat.ModuleLoadexConfig) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// Move implements [ValkeyInfra].
func (v *valkeyInfra) Move(ctx context.Context, key string, db int64) *valkeycompat.BoolCmd {
	panic("unimplemented")
}

// ObjectEncoding implements [ValkeyInfra].
func (v *valkeyInfra) ObjectEncoding(ctx context.Context, key string) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// ObjectIdleTime implements [ValkeyInfra].
func (v *valkeyInfra) ObjectIdleTime(ctx context.Context, key string) *valkeycompat.DurationCmd {
	panic("unimplemented")
}

// ObjectRefCount implements [ValkeyInfra].
func (v *valkeyInfra) ObjectRefCount(ctx context.Context, key string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// PExpire implements [ValkeyInfra].
func (v *valkeyInfra) PExpire(ctx context.Context, key string, expiration time.Duration) *valkeycompat.BoolCmd {
	panic("unimplemented")
}

// PExpireAt implements [ValkeyInfra].
func (v *valkeyInfra) PExpireAt(ctx context.Context, key string, tm time.Time) *valkeycompat.BoolCmd {
	panic("unimplemented")
}

// PExpireTime implements [ValkeyInfra].
func (v *valkeyInfra) PExpireTime(ctx context.Context, key string) *valkeycompat.DurationCmd {
	panic("unimplemented")
}

// PFAdd implements [ValkeyInfra].
func (v *valkeyInfra) PFAdd(ctx context.Context, key string, els ...any) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// PFCount implements [ValkeyInfra].
func (v *valkeyInfra) PFCount(ctx context.Context, keys ...string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// PFMerge implements [ValkeyInfra].
func (v *valkeyInfra) PFMerge(ctx context.Context, dest string, keys ...string) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// PTTL implements [ValkeyInfra].
func (v *valkeyInfra) PTTL(ctx context.Context, key string) *valkeycompat.DurationCmd {
	panic("unimplemented")
}

// Persist implements [ValkeyInfra].
func (v *valkeyInfra) Persist(ctx context.Context, key string) *valkeycompat.BoolCmd {
	panic("unimplemented")
}

// Ping implements [ValkeyInfra].
func (v *valkeyInfra) Ping(ctx context.Context) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// Pipeline implements [ValkeyInfra].
func (v *valkeyInfra) Pipeline() valkeycompat.Pipeliner {
	panic("unimplemented")
}

// Pipelined implements [ValkeyInfra].
func (v *valkeyInfra) Pipelined(ctx context.Context, fn func(valkeycompat.Pipeliner) error) ([]valkeycompat.Cmder, error) {
	panic("unimplemented")
}

// PubSubChannels implements [ValkeyInfra].
func (v *valkeyInfra) PubSubChannels(ctx context.Context, pattern string) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// PubSubNumPat implements [ValkeyInfra].
func (v *valkeyInfra) PubSubNumPat(ctx context.Context) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// PubSubNumSub implements [ValkeyInfra].
func (v *valkeyInfra) PubSubNumSub(ctx context.Context, channels ...string) *valkeycompat.StringIntMapCmd {
	panic("unimplemented")
}

// PubSubShardChannels implements [ValkeyInfra].
func (v *valkeyInfra) PubSubShardChannels(ctx context.Context, pattern string) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// PubSubShardNumSub implements [ValkeyInfra].
func (v *valkeyInfra) PubSubShardNumSub(ctx context.Context, channels ...string) *valkeycompat.StringIntMapCmd {
	panic("unimplemented")
}

// Publish implements [ValkeyInfra].
func (v *valkeyInfra) Publish(ctx context.Context, channel string, message any) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// Quit implements [ValkeyInfra].
func (v *valkeyInfra) Quit(ctx context.Context) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// RPop implements [ValkeyInfra].
func (v *valkeyInfra) RPop(ctx context.Context, key string) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// RPopCount implements [ValkeyInfra].
func (v *valkeyInfra) RPopCount(ctx context.Context, key string, count int64) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// RPopLPush implements [ValkeyInfra].
func (v *valkeyInfra) RPopLPush(ctx context.Context, source string, destination string) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// RPush implements [ValkeyInfra].
func (v *valkeyInfra) RPush(ctx context.Context, key string, values ...any) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// RPushX implements [ValkeyInfra].
func (v *valkeyInfra) RPushX(ctx context.Context, key string, values ...any) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// RandomKey implements [ValkeyInfra].
func (v *valkeyInfra) RandomKey(ctx context.Context) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// ReadOnly implements [ValkeyInfra].
func (v *valkeyInfra) ReadOnly(ctx context.Context) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// ReadWrite implements [ValkeyInfra].
func (v *valkeyInfra) ReadWrite(ctx context.Context) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// Rename implements [ValkeyInfra].
func (v *valkeyInfra) Rename(ctx context.Context, key string, newkey string) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// RenameNX implements [ValkeyInfra].
func (v *valkeyInfra) RenameNX(ctx context.Context, key string, newkey string) *valkeycompat.BoolCmd {
	panic("unimplemented")
}

// Restore implements [ValkeyInfra].
func (v *valkeyInfra) Restore(ctx context.Context, key string, ttl time.Duration, value string) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// RestoreReplace implements [ValkeyInfra].
func (v *valkeyInfra) RestoreReplace(ctx context.Context, key string, ttl time.Duration, value string) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// SAdd implements [ValkeyInfra].
func (v *valkeyInfra) SAdd(ctx context.Context, key string, members ...any) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// SCard implements [ValkeyInfra].
func (v *valkeyInfra) SCard(ctx context.Context, key string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// SDiff implements [ValkeyInfra].
func (v *valkeyInfra) SDiff(ctx context.Context, keys ...string) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// SDiffStore implements [ValkeyInfra].
func (v *valkeyInfra) SDiffStore(ctx context.Context, destination string, keys ...string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// SInter implements [ValkeyInfra].
func (v *valkeyInfra) SInter(ctx context.Context, keys ...string) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// SInterCard implements [ValkeyInfra].
func (v *valkeyInfra) SInterCard(ctx context.Context, limit int64, keys ...string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// SInterStore implements [ValkeyInfra].
func (v *valkeyInfra) SInterStore(ctx context.Context, destination string, keys ...string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// SIsMember implements [ValkeyInfra].
func (v *valkeyInfra) SIsMember(ctx context.Context, key string, member any) *valkeycompat.BoolCmd {
	panic("unimplemented")
}

// SMIsMember implements [ValkeyInfra].
func (v *valkeyInfra) SMIsMember(ctx context.Context, key string, members ...any) *valkeycompat.BoolSliceCmd {
	panic("unimplemented")
}

// SMembers implements [ValkeyInfra].
func (v *valkeyInfra) SMembers(ctx context.Context, key string) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// SMembersMap implements [ValkeyInfra].
func (v *valkeyInfra) SMembersMap(ctx context.Context, key string) *valkeycompat.StringStructMapCmd {
	panic("unimplemented")
}

// SMove implements [ValkeyInfra].
func (v *valkeyInfra) SMove(ctx context.Context, source string, destination string, member any) *valkeycompat.BoolCmd {
	panic("unimplemented")
}

// SPop implements [ValkeyInfra].
func (v *valkeyInfra) SPop(ctx context.Context, key string) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// SPopN implements [ValkeyInfra].
func (v *valkeyInfra) SPopN(ctx context.Context, key string, count int64) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// SPublish implements [ValkeyInfra].
func (v *valkeyInfra) SPublish(ctx context.Context, channel string, message any) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// SRandMember implements [ValkeyInfra].
func (v *valkeyInfra) SRandMember(ctx context.Context, key string) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// SRandMemberN implements [ValkeyInfra].
func (v *valkeyInfra) SRandMemberN(ctx context.Context, key string, count int64) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// SRem implements [ValkeyInfra].
func (v *valkeyInfra) SRem(ctx context.Context, key string, members ...any) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// SScan implements [ValkeyInfra].
func (v *valkeyInfra) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) *valkeycompat.ScanCmd {
	panic("unimplemented")
}

// SUnion implements [ValkeyInfra].
func (v *valkeyInfra) SUnion(ctx context.Context, keys ...string) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// SUnionStore implements [ValkeyInfra].
func (v *valkeyInfra) SUnionStore(ctx context.Context, destination string, keys ...string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// Save implements [ValkeyInfra].
func (v *valkeyInfra) Save(ctx context.Context) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// Scan implements [ValkeyInfra].
func (v *valkeyInfra) Scan(ctx context.Context, cursor uint64, match string, count int64) *valkeycompat.ScanCmd {
	panic("unimplemented")
}

// ScanType implements [ValkeyInfra].
func (v *valkeyInfra) ScanType(ctx context.Context, cursor uint64, match string, count int64, keyType string) *valkeycompat.ScanCmd {
	panic("unimplemented")
}

// ScriptExists implements [ValkeyInfra].
func (v *valkeyInfra) ScriptExists(ctx context.Context, hashes ...string) *valkeycompat.BoolSliceCmd {
	panic("unimplemented")
}

// ScriptFlush implements [ValkeyInfra].
func (v *valkeyInfra) ScriptFlush(ctx context.Context) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// ScriptKill implements [ValkeyInfra].
func (v *valkeyInfra) ScriptKill(ctx context.Context) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// ScriptLoad implements [ValkeyInfra].
func (v *valkeyInfra) ScriptLoad(ctx context.Context, script string) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// Set implements [ValkeyInfra].
func (v *valkeyInfra) Set(ctx context.Context, key string, value any, expiration time.Duration) *valkeycompat.StatusCmd {
	return v.valkeyClient.Set(ctx, key, value, expiration)
}

// SetArgs implements [ValkeyInfra].
func (v *valkeyInfra) SetArgs(ctx context.Context, key string, value any, a valkeycompat.SetArgs) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// SetBit implements [ValkeyInfra].
func (v *valkeyInfra) SetBit(ctx context.Context, key string, offset int64, value int64) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// SetEX implements [ValkeyInfra].
func (v *valkeyInfra) SetEX(ctx context.Context, key string, value any, expiration time.Duration) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// SetNX implements [ValkeyInfra].
func (v *valkeyInfra) SetNX(ctx context.Context, key string, value any, expiration time.Duration) *valkeycompat.BoolCmd {
	panic("unimplemented")
}

// SetRange implements [ValkeyInfra].
func (v *valkeyInfra) SetRange(ctx context.Context, key string, offset int64, value string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// SetXX implements [ValkeyInfra].
func (v *valkeyInfra) SetXX(ctx context.Context, key string, value any, expiration time.Duration) *valkeycompat.BoolCmd {
	panic("unimplemented")
}

// Shutdown implements [ValkeyInfra].
func (v *valkeyInfra) Shutdown(ctx context.Context) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// ShutdownNoSave implements [ValkeyInfra].
func (v *valkeyInfra) ShutdownNoSave(ctx context.Context) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// ShutdownSave implements [ValkeyInfra].
func (v *valkeyInfra) ShutdownSave(ctx context.Context) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// SlaveOf implements [ValkeyInfra].
func (v *valkeyInfra) SlaveOf(ctx context.Context, host string, port string) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// SlowLogGet implements [ValkeyInfra].
func (v *valkeyInfra) SlowLogGet(ctx context.Context, num int64) *valkeycompat.SlowLogCmd {
	panic("unimplemented")
}

// SlowLogReset implements [ValkeyInfra].
func (v *valkeyInfra) SlowLogReset(ctx context.Context) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// Sort implements [ValkeyInfra].
func (v *valkeyInfra) Sort(ctx context.Context, key string, sort valkeycompat.Sort) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// SortInterfaces implements [ValkeyInfra].
func (v *valkeyInfra) SortInterfaces(ctx context.Context, key string, sort valkeycompat.Sort) *valkeycompat.SliceCmd {
	panic("unimplemented")
}

// SortRO implements [ValkeyInfra].
func (v *valkeyInfra) SortRO(ctx context.Context, key string, sort valkeycompat.Sort) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// SortStore implements [ValkeyInfra].
func (v *valkeyInfra) SortStore(ctx context.Context, key string, store string, sort valkeycompat.Sort) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// StrLen implements [ValkeyInfra].
func (v *valkeyInfra) StrLen(ctx context.Context, key string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// TDigestAdd implements [ValkeyInfra].
func (v *valkeyInfra) TDigestAdd(ctx context.Context, key string, elements ...float64) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// TDigestByRank implements [ValkeyInfra].
func (v *valkeyInfra) TDigestByRank(ctx context.Context, key string, rank ...uint64) *valkeycompat.FloatSliceCmd {
	panic("unimplemented")
}

// TDigestByRevRank implements [ValkeyInfra].
func (v *valkeyInfra) TDigestByRevRank(ctx context.Context, key string, rank ...uint64) *valkeycompat.FloatSliceCmd {
	panic("unimplemented")
}

// TDigestCDF implements [ValkeyInfra].
func (v *valkeyInfra) TDigestCDF(ctx context.Context, key string, elements ...float64) *valkeycompat.FloatSliceCmd {
	panic("unimplemented")
}

// TDigestCreate implements [ValkeyInfra].
func (v *valkeyInfra) TDigestCreate(ctx context.Context, key string) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// TDigestCreateWithCompression implements [ValkeyInfra].
func (v *valkeyInfra) TDigestCreateWithCompression(ctx context.Context, key string, compression int64) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// TDigestInfo implements [ValkeyInfra].
func (v *valkeyInfra) TDigestInfo(ctx context.Context, key string) *valkeycompat.TDigestInfoCmd {
	panic("unimplemented")
}

// TDigestMax implements [ValkeyInfra].
func (v *valkeyInfra) TDigestMax(ctx context.Context, key string) *valkeycompat.FloatCmd {
	panic("unimplemented")
}

// TDigestMerge implements [ValkeyInfra].
func (v *valkeyInfra) TDigestMerge(ctx context.Context, destKey string, options *valkeycompat.TDigestMergeOptions, sourceKeys ...string) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// TDigestMin implements [ValkeyInfra].
func (v *valkeyInfra) TDigestMin(ctx context.Context, key string) *valkeycompat.FloatCmd {
	panic("unimplemented")
}

// TDigestQuantile implements [ValkeyInfra].
func (v *valkeyInfra) TDigestQuantile(ctx context.Context, key string, elements ...float64) *valkeycompat.FloatSliceCmd {
	panic("unimplemented")
}

// TDigestRank implements [ValkeyInfra].
func (v *valkeyInfra) TDigestRank(ctx context.Context, key string, values ...float64) *valkeycompat.IntSliceCmd {
	panic("unimplemented")
}

// TDigestReset implements [ValkeyInfra].
func (v *valkeyInfra) TDigestReset(ctx context.Context, key string) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// TDigestRevRank implements [ValkeyInfra].
func (v *valkeyInfra) TDigestRevRank(ctx context.Context, key string, values ...float64) *valkeycompat.IntSliceCmd {
	panic("unimplemented")
}

// TDigestTrimmedMean implements [ValkeyInfra].
func (v *valkeyInfra) TDigestTrimmedMean(ctx context.Context, key string, lowCutQuantile float64, highCutQuantile float64) *valkeycompat.FloatCmd {
	panic("unimplemented")
}

// TFCall implements [ValkeyInfra].
func (v *valkeyInfra) TFCall(ctx context.Context, libName string, funcName string, numKeys int) *valkeycompat.Cmd {
	panic("unimplemented")
}

// TFCallASYNC implements [ValkeyInfra].
func (v *valkeyInfra) TFCallASYNC(ctx context.Context, libName string, funcName string, numKeys int) *valkeycompat.Cmd {
	panic("unimplemented")
}

// TFCallASYNCArgs implements [ValkeyInfra].
func (v *valkeyInfra) TFCallASYNCArgs(ctx context.Context, libName string, funcName string, numKeys int, options *valkeycompat.TFCallOptions) *valkeycompat.Cmd {
	panic("unimplemented")
}

// TFCallArgs implements [ValkeyInfra].
func (v *valkeyInfra) TFCallArgs(ctx context.Context, libName string, funcName string, numKeys int, options *valkeycompat.TFCallOptions) *valkeycompat.Cmd {
	panic("unimplemented")
}

// TFunctionDelete implements [ValkeyInfra].
func (v *valkeyInfra) TFunctionDelete(ctx context.Context, libName string) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// TFunctionList implements [ValkeyInfra].
func (v *valkeyInfra) TFunctionList(ctx context.Context) *valkeycompat.MapStringInterfaceSliceCmd {
	panic("unimplemented")
}

// TFunctionListArgs implements [ValkeyInfra].
func (v *valkeyInfra) TFunctionListArgs(ctx context.Context, options *valkeycompat.TFunctionListOptions) *valkeycompat.MapStringInterfaceSliceCmd {
	panic("unimplemented")
}

// TFunctionLoad implements [ValkeyInfra].
func (v *valkeyInfra) TFunctionLoad(ctx context.Context, lib string) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// TFunctionLoadArgs implements [ValkeyInfra].
func (v *valkeyInfra) TFunctionLoadArgs(ctx context.Context, lib string, options *valkeycompat.TFunctionLoadOptions) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// TSAdd implements [ValkeyInfra].
func (v *valkeyInfra) TSAdd(ctx context.Context, key string, timestamp interface{}, value float64) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// TSAddWithArgs implements [ValkeyInfra].
func (v *valkeyInfra) TSAddWithArgs(ctx context.Context, key string, timestamp interface{}, value float64, options *valkeycompat.TSOptions) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// TSAlter implements [ValkeyInfra].
func (v *valkeyInfra) TSAlter(ctx context.Context, key string, options *valkeycompat.TSAlterOptions) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// TSCreate implements [ValkeyInfra].
func (v *valkeyInfra) TSCreate(ctx context.Context, key string) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// TSCreateRule implements [ValkeyInfra].
func (v *valkeyInfra) TSCreateRule(ctx context.Context, sourceKey string, destKey string, aggregator valkeycompat.Aggregator, bucketDuration int) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// TSCreateRuleWithArgs implements [ValkeyInfra].
func (v *valkeyInfra) TSCreateRuleWithArgs(ctx context.Context, sourceKey string, destKey string, aggregator valkeycompat.Aggregator, bucketDuration int, options *valkeycompat.TSCreateRuleOptions) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// TSCreateWithArgs implements [ValkeyInfra].
func (v *valkeyInfra) TSCreateWithArgs(ctx context.Context, key string, options *valkeycompat.TSOptions) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// TSDecrBy implements [ValkeyInfra].
func (v *valkeyInfra) TSDecrBy(ctx context.Context, Key string, timestamp float64) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// TSDecrByWithArgs implements [ValkeyInfra].
func (v *valkeyInfra) TSDecrByWithArgs(ctx context.Context, key string, timestamp float64, options *valkeycompat.TSIncrDecrOptions) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// TSDel implements [ValkeyInfra].
func (v *valkeyInfra) TSDel(ctx context.Context, Key string, fromTimestamp int, toTimestamp int) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// TSDeleteRule implements [ValkeyInfra].
func (v *valkeyInfra) TSDeleteRule(ctx context.Context, sourceKey string, destKey string) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// TSGet implements [ValkeyInfra].
func (v *valkeyInfra) TSGet(ctx context.Context, key string) *valkeycompat.TSTimestampValueCmd {
	panic("unimplemented")
}

// TSGetWithArgs implements [ValkeyInfra].
func (v *valkeyInfra) TSGetWithArgs(ctx context.Context, key string, options *valkeycompat.TSGetOptions) *valkeycompat.TSTimestampValueCmd {
	panic("unimplemented")
}

// TSIncrBy implements [ValkeyInfra].
func (v *valkeyInfra) TSIncrBy(ctx context.Context, Key string, timestamp float64) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// TSIncrByWithArgs implements [ValkeyInfra].
func (v *valkeyInfra) TSIncrByWithArgs(ctx context.Context, key string, timestamp float64, options *valkeycompat.TSIncrDecrOptions) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// TSInfo implements [ValkeyInfra].
func (v *valkeyInfra) TSInfo(ctx context.Context, key string) *valkeycompat.MapStringInterfaceCmd {
	panic("unimplemented")
}

// TSInfoWithArgs implements [ValkeyInfra].
func (v *valkeyInfra) TSInfoWithArgs(ctx context.Context, key string, options *valkeycompat.TSInfoOptions) *valkeycompat.MapStringInterfaceCmd {
	panic("unimplemented")
}

// TSMAdd implements [ValkeyInfra].
func (v *valkeyInfra) TSMAdd(ctx context.Context, ktvSlices [][]interface{}) *valkeycompat.IntSliceCmd {
	panic("unimplemented")
}

// TSMGet implements [ValkeyInfra].
func (v *valkeyInfra) TSMGet(ctx context.Context, filters []string) *valkeycompat.MapStringSliceInterfaceCmd {
	panic("unimplemented")
}

// TSMGetWithArgs implements [ValkeyInfra].
func (v *valkeyInfra) TSMGetWithArgs(ctx context.Context, filters []string, options *valkeycompat.TSMGetOptions) *valkeycompat.MapStringSliceInterfaceCmd {
	panic("unimplemented")
}

// TSMRange implements [ValkeyInfra].
func (v *valkeyInfra) TSMRange(ctx context.Context, fromTimestamp int, toTimestamp int, filterExpr []string) *valkeycompat.MapStringSliceInterfaceCmd {
	panic("unimplemented")
}

// TSMRangeWithArgs implements [ValkeyInfra].
func (v *valkeyInfra) TSMRangeWithArgs(ctx context.Context, fromTimestamp int, toTimestamp int, filterExpr []string, options *valkeycompat.TSMRangeOptions) *valkeycompat.MapStringSliceInterfaceCmd {
	panic("unimplemented")
}

// TSMRevRange implements [ValkeyInfra].
func (v *valkeyInfra) TSMRevRange(ctx context.Context, fromTimestamp int, toTimestamp int, filterExpr []string) *valkeycompat.MapStringSliceInterfaceCmd {
	panic("unimplemented")
}

// TSMRevRangeWithArgs implements [ValkeyInfra].
func (v *valkeyInfra) TSMRevRangeWithArgs(ctx context.Context, fromTimestamp int, toTimestamp int, filterExpr []string, options *valkeycompat.TSMRevRangeOptions) *valkeycompat.MapStringSliceInterfaceCmd {
	panic("unimplemented")
}

// TSQueryIndex implements [ValkeyInfra].
func (v *valkeyInfra) TSQueryIndex(ctx context.Context, filterExpr []string) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// TSRange implements [ValkeyInfra].
func (v *valkeyInfra) TSRange(ctx context.Context, key string, fromTimestamp int, toTimestamp int) *valkeycompat.TSTimestampValueSliceCmd {
	panic("unimplemented")
}

// TSRangeWithArgs implements [ValkeyInfra].
func (v *valkeyInfra) TSRangeWithArgs(ctx context.Context, key string, fromTimestamp int, toTimestamp int, options *valkeycompat.TSRangeOptions) *valkeycompat.TSTimestampValueSliceCmd {
	panic("unimplemented")
}

// TSRevRange implements [ValkeyInfra].
func (v *valkeyInfra) TSRevRange(ctx context.Context, key string, fromTimestamp int, toTimestamp int) *valkeycompat.TSTimestampValueSliceCmd {
	panic("unimplemented")
}

// TSRevRangeWithArgs implements [ValkeyInfra].
func (v *valkeyInfra) TSRevRangeWithArgs(ctx context.Context, key string, fromTimestamp int, toTimestamp int, options *valkeycompat.TSRevRangeOptions) *valkeycompat.TSTimestampValueSliceCmd {
	panic("unimplemented")
}

// TTL implements [ValkeyInfra].
func (v *valkeyInfra) TTL(ctx context.Context, key string) *valkeycompat.DurationCmd {
	panic("unimplemented")
}

// Time implements [ValkeyInfra].
func (v *valkeyInfra) Time(ctx context.Context) *valkeycompat.TimeCmd {
	panic("unimplemented")
}

// TopKAdd implements [ValkeyInfra].
func (v *valkeyInfra) TopKAdd(ctx context.Context, key string, elements ...interface{}) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// TopKCount implements [ValkeyInfra].
func (v *valkeyInfra) TopKCount(ctx context.Context, key string, elements ...interface{}) *valkeycompat.IntSliceCmd {
	panic("unimplemented")
}

// TopKIncrBy implements [ValkeyInfra].
func (v *valkeyInfra) TopKIncrBy(ctx context.Context, key string, elements ...interface{}) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// TopKInfo implements [ValkeyInfra].
func (v *valkeyInfra) TopKInfo(ctx context.Context, key string) *valkeycompat.TopKInfoCmd {
	panic("unimplemented")
}

// TopKList implements [ValkeyInfra].
func (v *valkeyInfra) TopKList(ctx context.Context, key string) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// TopKListWithCount implements [ValkeyInfra].
func (v *valkeyInfra) TopKListWithCount(ctx context.Context, key string) *valkeycompat.MapStringIntCmd {
	panic("unimplemented")
}

// TopKQuery implements [ValkeyInfra].
func (v *valkeyInfra) TopKQuery(ctx context.Context, key string, elements ...interface{}) *valkeycompat.BoolSliceCmd {
	panic("unimplemented")
}

// TopKReserve implements [ValkeyInfra].
func (v *valkeyInfra) TopKReserve(ctx context.Context, key string, k int64) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// TopKReserveWithOptions implements [ValkeyInfra].
func (v *valkeyInfra) TopKReserveWithOptions(ctx context.Context, key string, k int64, width int64, depth int64, decay float64) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// Touch implements [ValkeyInfra].
func (v *valkeyInfra) Touch(ctx context.Context, keys ...string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// TxPipeline implements [ValkeyInfra].
func (v *valkeyInfra) TxPipeline() valkeycompat.Pipeliner {
	panic("unimplemented")
}

// TxPipelined implements [ValkeyInfra].
func (v *valkeyInfra) TxPipelined(ctx context.Context, fn func(valkeycompat.Pipeliner) error) ([]valkeycompat.Cmder, error) {
	panic("unimplemented")
}

// Type implements [ValkeyInfra].
func (v *valkeyInfra) Type(ctx context.Context, key string) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// Unlink implements [ValkeyInfra].
func (v *valkeyInfra) Unlink(ctx context.Context, keys ...string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// XAck implements [ValkeyInfra].
func (v *valkeyInfra) XAck(ctx context.Context, stream string, group string, ids ...string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// XAdd implements [ValkeyInfra].
func (v *valkeyInfra) XAdd(ctx context.Context, a valkeycompat.XAddArgs) *valkeycompat.StringCmd {
	panic("unimplemented")
}

// XAutoClaim implements [ValkeyInfra].
func (v *valkeyInfra) XAutoClaim(ctx context.Context, a valkeycompat.XAutoClaimArgs) *valkeycompat.XAutoClaimCmd {
	panic("unimplemented")
}

// XAutoClaimJustID implements [ValkeyInfra].
func (v *valkeyInfra) XAutoClaimJustID(ctx context.Context, a valkeycompat.XAutoClaimArgs) *valkeycompat.XAutoClaimJustIDCmd {
	panic("unimplemented")
}

// XCfgSet implements [ValkeyInfra].
func (v *valkeyInfra) XCfgSet(ctx context.Context, a valkeycompat.XCfgSetArgs) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// XClaim implements [ValkeyInfra].
func (v *valkeyInfra) XClaim(ctx context.Context, a valkeycompat.XClaimArgs) *valkeycompat.XMessageSliceCmd {
	panic("unimplemented")
}

// XClaimJustID implements [ValkeyInfra].
func (v *valkeyInfra) XClaimJustID(ctx context.Context, a valkeycompat.XClaimArgs) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// XDel implements [ValkeyInfra].
func (v *valkeyInfra) XDel(ctx context.Context, stream string, ids ...string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// XGroupCreate implements [ValkeyInfra].
func (v *valkeyInfra) XGroupCreate(ctx context.Context, stream string, group string, start string) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// XGroupCreateConsumer implements [ValkeyInfra].
func (v *valkeyInfra) XGroupCreateConsumer(ctx context.Context, stream string, group string, consumer string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// XGroupCreateMkStream implements [ValkeyInfra].
func (v *valkeyInfra) XGroupCreateMkStream(ctx context.Context, stream string, group string, start string) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// XGroupDelConsumer implements [ValkeyInfra].
func (v *valkeyInfra) XGroupDelConsumer(ctx context.Context, stream string, group string, consumer string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// XGroupDestroy implements [ValkeyInfra].
func (v *valkeyInfra) XGroupDestroy(ctx context.Context, stream string, group string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// XGroupSetID implements [ValkeyInfra].
func (v *valkeyInfra) XGroupSetID(ctx context.Context, stream string, group string, start string) *valkeycompat.StatusCmd {
	panic("unimplemented")
}

// XInfoConsumers implements [ValkeyInfra].
func (v *valkeyInfra) XInfoConsumers(ctx context.Context, key string, group string) *valkeycompat.XInfoConsumersCmd {
	panic("unimplemented")
}

// XInfoGroups implements [ValkeyInfra].
func (v *valkeyInfra) XInfoGroups(ctx context.Context, key string) *valkeycompat.XInfoGroupsCmd {
	panic("unimplemented")
}

// XInfoStream implements [ValkeyInfra].
func (v *valkeyInfra) XInfoStream(ctx context.Context, key string) *valkeycompat.XInfoStreamCmd {
	panic("unimplemented")
}

// XInfoStreamFull implements [ValkeyInfra].
func (v *valkeyInfra) XInfoStreamFull(ctx context.Context, key string, count int64) *valkeycompat.XInfoStreamFullCmd {
	panic("unimplemented")
}

// XLen implements [ValkeyInfra].
func (v *valkeyInfra) XLen(ctx context.Context, stream string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// XPending implements [ValkeyInfra].
func (v *valkeyInfra) XPending(ctx context.Context, stream string, group string) *valkeycompat.XPendingCmd {
	panic("unimplemented")
}

// XPendingExt implements [ValkeyInfra].
func (v *valkeyInfra) XPendingExt(ctx context.Context, a valkeycompat.XPendingExtArgs) *valkeycompat.XPendingExtCmd {
	panic("unimplemented")
}

// XRange implements [ValkeyInfra].
func (v *valkeyInfra) XRange(ctx context.Context, stream string, start string, stop string) *valkeycompat.XMessageSliceCmd {
	panic("unimplemented")
}

// XRangeN implements [ValkeyInfra].
func (v *valkeyInfra) XRangeN(ctx context.Context, stream string, start string, stop string, count int64) *valkeycompat.XMessageSliceCmd {
	panic("unimplemented")
}

// XRead implements [ValkeyInfra].
func (v *valkeyInfra) XRead(ctx context.Context, a valkeycompat.XReadArgs) *valkeycompat.XStreamSliceCmd {
	panic("unimplemented")
}

// XReadGroup implements [ValkeyInfra].
func (v *valkeyInfra) XReadGroup(ctx context.Context, a valkeycompat.XReadGroupArgs) *valkeycompat.XStreamSliceCmd {
	panic("unimplemented")
}

// XReadStreams implements [ValkeyInfra].
func (v *valkeyInfra) XReadStreams(ctx context.Context, streams ...string) *valkeycompat.XStreamSliceCmd {
	panic("unimplemented")
}

// XRevRange implements [ValkeyInfra].
func (v *valkeyInfra) XRevRange(ctx context.Context, stream string, start string, stop string) *valkeycompat.XMessageSliceCmd {
	panic("unimplemented")
}

// XRevRangeN implements [ValkeyInfra].
func (v *valkeyInfra) XRevRangeN(ctx context.Context, stream string, start string, stop string, count int64) *valkeycompat.XMessageSliceCmd {
	panic("unimplemented")
}

// XTrimMaxLen implements [ValkeyInfra].
func (v *valkeyInfra) XTrimMaxLen(ctx context.Context, key string, maxLen int64) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// XTrimMaxLenApprox implements [ValkeyInfra].
func (v *valkeyInfra) XTrimMaxLenApprox(ctx context.Context, key string, maxLen int64, limit int64) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// XTrimMinID implements [ValkeyInfra].
func (v *valkeyInfra) XTrimMinID(ctx context.Context, key string, minID string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// XTrimMinIDApprox implements [ValkeyInfra].
func (v *valkeyInfra) XTrimMinIDApprox(ctx context.Context, key string, minID string, limit int64) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// ZAdd implements [ValkeyInfra].
func (v *valkeyInfra) ZAdd(ctx context.Context, key string, members ...valkeycompat.Z) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// ZAddArgs implements [ValkeyInfra].
func (v *valkeyInfra) ZAddArgs(ctx context.Context, key string, args valkeycompat.ZAddArgs) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// ZAddArgsIncr implements [ValkeyInfra].
func (v *valkeyInfra) ZAddArgsIncr(ctx context.Context, key string, args valkeycompat.ZAddArgs) *valkeycompat.FloatCmd {
	panic("unimplemented")
}

// ZAddGT implements [ValkeyInfra].
func (v *valkeyInfra) ZAddGT(ctx context.Context, key string, members ...valkeycompat.Z) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// ZAddLT implements [ValkeyInfra].
func (v *valkeyInfra) ZAddLT(ctx context.Context, key string, members ...valkeycompat.Z) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// ZAddNX implements [ValkeyInfra].
func (v *valkeyInfra) ZAddNX(ctx context.Context, key string, members ...valkeycompat.Z) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// ZAddXX implements [ValkeyInfra].
func (v *valkeyInfra) ZAddXX(ctx context.Context, key string, members ...valkeycompat.Z) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// ZCard implements [ValkeyInfra].
func (v *valkeyInfra) ZCard(ctx context.Context, key string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// ZCount implements [ValkeyInfra].
func (v *valkeyInfra) ZCount(ctx context.Context, key string, min string, max string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// ZDiff implements [ValkeyInfra].
func (v *valkeyInfra) ZDiff(ctx context.Context, keys ...string) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// ZDiffStore implements [ValkeyInfra].
func (v *valkeyInfra) ZDiffStore(ctx context.Context, destination string, keys ...string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// ZDiffWithScores implements [ValkeyInfra].
func (v *valkeyInfra) ZDiffWithScores(ctx context.Context, keys ...string) *valkeycompat.ZSliceCmd {
	panic("unimplemented")
}

// ZIncrBy implements [ValkeyInfra].
func (v *valkeyInfra) ZIncrBy(ctx context.Context, key string, increment float64, member string) *valkeycompat.FloatCmd {
	panic("unimplemented")
}

// ZInter implements [ValkeyInfra].
func (v *valkeyInfra) ZInter(ctx context.Context, store valkeycompat.ZStore) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// ZInterCard implements [ValkeyInfra].
func (v *valkeyInfra) ZInterCard(ctx context.Context, limit int64, keys ...string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// ZInterStore implements [ValkeyInfra].
func (v *valkeyInfra) ZInterStore(ctx context.Context, destination string, store valkeycompat.ZStore) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// ZInterWithScores implements [ValkeyInfra].
func (v *valkeyInfra) ZInterWithScores(ctx context.Context, store valkeycompat.ZStore) *valkeycompat.ZSliceCmd {
	panic("unimplemented")
}

// ZLexCount implements [ValkeyInfra].
func (v *valkeyInfra) ZLexCount(ctx context.Context, key string, min string, max string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// ZMPop implements [ValkeyInfra].
func (v *valkeyInfra) ZMPop(ctx context.Context, order string, count int64, keys ...string) *valkeycompat.ZSliceWithKeyCmd {
	panic("unimplemented")
}

// ZMScore implements [ValkeyInfra].
func (v *valkeyInfra) ZMScore(ctx context.Context, key string, members ...string) *valkeycompat.FloatSliceCmd {
	panic("unimplemented")
}

// ZPopMax implements [ValkeyInfra].
func (v *valkeyInfra) ZPopMax(ctx context.Context, key string, count ...int64) *valkeycompat.ZSliceCmd {
	panic("unimplemented")
}

// ZPopMin implements [ValkeyInfra].
func (v *valkeyInfra) ZPopMin(ctx context.Context, key string, count ...int64) *valkeycompat.ZSliceCmd {
	panic("unimplemented")
}

// ZRandMember implements [ValkeyInfra].
func (v *valkeyInfra) ZRandMember(ctx context.Context, key string, count int64) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// ZRandMemberWithScores implements [ValkeyInfra].
func (v *valkeyInfra) ZRandMemberWithScores(ctx context.Context, key string, count int64) *valkeycompat.ZSliceCmd {
	panic("unimplemented")
}

// ZRange implements [ValkeyInfra].
func (v *valkeyInfra) ZRange(ctx context.Context, key string, start int64, stop int64) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// ZRangeArgs implements [ValkeyInfra].
func (v *valkeyInfra) ZRangeArgs(ctx context.Context, z valkeycompat.ZRangeArgs) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// ZRangeArgsWithScores implements [ValkeyInfra].
func (v *valkeyInfra) ZRangeArgsWithScores(ctx context.Context, z valkeycompat.ZRangeArgs) *valkeycompat.ZSliceCmd {
	panic("unimplemented")
}

// ZRangeByLex implements [ValkeyInfra].
func (v *valkeyInfra) ZRangeByLex(ctx context.Context, key string, opt valkeycompat.ZRangeBy) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// ZRangeByScore implements [ValkeyInfra].
func (v *valkeyInfra) ZRangeByScore(ctx context.Context, key string, opt valkeycompat.ZRangeBy) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// ZRangeByScoreWithScores implements [ValkeyInfra].
func (v *valkeyInfra) ZRangeByScoreWithScores(ctx context.Context, key string, opt valkeycompat.ZRangeBy) *valkeycompat.ZSliceCmd {
	panic("unimplemented")
}

// ZRangeStore implements [ValkeyInfra].
func (v *valkeyInfra) ZRangeStore(ctx context.Context, dst string, z valkeycompat.ZRangeArgs) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// ZRangeWithScores implements [ValkeyInfra].
func (v *valkeyInfra) ZRangeWithScores(ctx context.Context, key string, start int64, stop int64) *valkeycompat.ZSliceCmd {
	panic("unimplemented")
}

// ZRank implements [ValkeyInfra].
func (v *valkeyInfra) ZRank(ctx context.Context, key string, member string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// ZRankWithScore implements [ValkeyInfra].
func (v *valkeyInfra) ZRankWithScore(ctx context.Context, key string, member string) *valkeycompat.RankWithScoreCmd {
	panic("unimplemented")
}

// ZRem implements [ValkeyInfra].
func (v *valkeyInfra) ZRem(ctx context.Context, key string, members ...any) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// ZRemRangeByLex implements [ValkeyInfra].
func (v *valkeyInfra) ZRemRangeByLex(ctx context.Context, key string, min string, max string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// ZRemRangeByRank implements [ValkeyInfra].
func (v *valkeyInfra) ZRemRangeByRank(ctx context.Context, key string, start int64, stop int64) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// ZRemRangeByScore implements [ValkeyInfra].
func (v *valkeyInfra) ZRemRangeByScore(ctx context.Context, key string, min string, max string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// ZRevRange implements [ValkeyInfra].
func (v *valkeyInfra) ZRevRange(ctx context.Context, key string, start int64, stop int64) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// ZRevRangeByLex implements [ValkeyInfra].
func (v *valkeyInfra) ZRevRangeByLex(ctx context.Context, key string, opt valkeycompat.ZRangeBy) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// ZRevRangeByScore implements [ValkeyInfra].
func (v *valkeyInfra) ZRevRangeByScore(ctx context.Context, key string, opt valkeycompat.ZRangeBy) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// ZRevRangeByScoreWithScores implements [ValkeyInfra].
func (v *valkeyInfra) ZRevRangeByScoreWithScores(ctx context.Context, key string, opt valkeycompat.ZRangeBy) *valkeycompat.ZSliceCmd {
	panic("unimplemented")
}

// ZRevRangeWithScores implements [ValkeyInfra].
func (v *valkeyInfra) ZRevRangeWithScores(ctx context.Context, key string, start int64, stop int64) *valkeycompat.ZSliceCmd {
	panic("unimplemented")
}

// ZRevRank implements [ValkeyInfra].
func (v *valkeyInfra) ZRevRank(ctx context.Context, key string, member string) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// ZRevRankWithScore implements [ValkeyInfra].
func (v *valkeyInfra) ZRevRankWithScore(ctx context.Context, key string, member string) *valkeycompat.RankWithScoreCmd {
	panic("unimplemented")
}

// ZScan implements [ValkeyInfra].
func (v *valkeyInfra) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) *valkeycompat.ScanCmd {
	panic("unimplemented")
}

// ZScore implements [ValkeyInfra].
func (v *valkeyInfra) ZScore(ctx context.Context, key string, member string) *valkeycompat.FloatCmd {
	panic("unimplemented")
}

// ZUnion implements [ValkeyInfra].
func (v *valkeyInfra) ZUnion(ctx context.Context, store valkeycompat.ZStore) *valkeycompat.StringSliceCmd {
	panic("unimplemented")
}

// ZUnionStore implements [ValkeyInfra].
func (v *valkeyInfra) ZUnionStore(ctx context.Context, dest string, store valkeycompat.ZStore) *valkeycompat.IntCmd {
	panic("unimplemented")
}

// ZUnionWithScores implements [ValkeyInfra].
func (v *valkeyInfra) ZUnionWithScores(ctx context.Context, store valkeycompat.ZStore) *valkeycompat.ZSliceCmd {
	panic("unimplemented")
}
