package collector

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace        = "goat"
	scrapeRPCTimeout = 10 * time.Second
	syncSubsystem    = "sync_progress"
)

type GethRPCCollector struct {
	endpoint                       string
	client                         *GoatClient
	blockDesc                      *prometheus.Desc
	chainDesc                      *prometheus.Desc
	syncDoneDesc                   *prometheus.Desc
	syncStartingBlockDesc          *prometheus.Desc
	syncCurrentBlockDesc           *prometheus.Desc
	syncHighestBlockDesc           *prometheus.Desc
	syncPulledStatesDesc           *prometheus.Desc
	syncKnownStatesDesc            *prometheus.Desc
	syncSyncedAccountsDesc         *prometheus.Desc
	syncSyncedAccountBytesDesc     *prometheus.Desc
	syncSyncedBytecodesDesc        *prometheus.Desc
	syncSyncedBytecodeBytesDesc    *prometheus.Desc
	syncSyncedStorageDesc          *prometheus.Desc
	syncSyncedStorageBytesDesc     *prometheus.Desc
	syncHealedTrienodesDesc        *prometheus.Desc
	syncHealedTrienodeBytesDesc    *prometheus.Desc
	syncHealedBytecodesDesc        *prometheus.Desc
	syncHealedBytecodeBytesDesc    *prometheus.Desc
	syncHealingTrienodesDesc       *prometheus.Desc
	syncHealingBytecodeDesc        *prometheus.Desc
	syncTxIndexFinishedBlocksDesc  *prometheus.Desc
	syncTxIndexRemainingBlocksDesc *prometheus.Desc
	syncStateIndexRemainingDesc    *prometheus.Desc
}

func NewGethRPCCollector(endpoint string) (*GethRPCCollector, error) {
	client, err := NewGoatClient(endpoint)
	if err != nil {
		return nil, err
	}

	collector := &GethRPCCollector{
		endpoint:                       endpoint,
		client:                         client,
		blockDesc:                      prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "block_height"), "Current block height reported by the Goat network RPC node.", nil, nil),
		chainDesc:                      prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "chain_id"), "Chain ID reported by the Goat network RPC node.", nil, nil),
		syncDoneDesc:                   prometheus.NewDesc(prometheus.BuildFQName(namespace, syncSubsystem, "done"), "Whether the node has completed initial sync (1 when SyncProgress.Done() reports true).", nil, nil),
		syncStartingBlockDesc:          prometheus.NewDesc(prometheus.BuildFQName(namespace, syncSubsystem, "starting_block"), "Block number where sync began.", nil, nil),
		syncCurrentBlockDesc:           prometheus.NewDesc(prometheus.BuildFQName(namespace, syncSubsystem, "current_block"), "Current block number where sync is at.", nil, nil),
		syncHighestBlockDesc:           prometheus.NewDesc(prometheus.BuildFQName(namespace, syncSubsystem, "highest_block"), "Highest alleged block number in the chain.", nil, nil),
		syncPulledStatesDesc:           prometheus.NewDesc(prometheus.BuildFQName(namespace, syncSubsystem, "pulled_states"), "Number of state trie entries already downloaded.", nil, nil),
		syncKnownStatesDesc:            prometheus.NewDesc(prometheus.BuildFQName(namespace, syncSubsystem, "known_states"), "Total number of state trie entries known about.", nil, nil),
		syncSyncedAccountsDesc:         prometheus.NewDesc(prometheus.BuildFQName(namespace, syncSubsystem, "synced_accounts"), "Number of accounts downloaded.", nil, nil),
		syncSyncedAccountBytesDesc:     prometheus.NewDesc(prometheus.BuildFQName(namespace, syncSubsystem, "synced_account_bytes"), "Number of account trie bytes persisted to disk.", nil, nil),
		syncSyncedBytecodesDesc:        prometheus.NewDesc(prometheus.BuildFQName(namespace, syncSubsystem, "synced_bytecodes"), "Number of bytecodes downloaded.", nil, nil),
		syncSyncedBytecodeBytesDesc:    prometheus.NewDesc(prometheus.BuildFQName(namespace, syncSubsystem, "synced_bytecode_bytes"), "Number of bytecode bytes downloaded.", nil, nil),
		syncSyncedStorageDesc:          prometheus.NewDesc(prometheus.BuildFQName(namespace, syncSubsystem, "synced_storage"), "Number of storage slots downloaded.", nil, nil),
		syncSyncedStorageBytesDesc:     prometheus.NewDesc(prometheus.BuildFQName(namespace, syncSubsystem, "synced_storage_bytes"), "Number of storage trie bytes persisted to disk.", nil, nil),
		syncHealedTrienodesDesc:        prometheus.NewDesc(prometheus.BuildFQName(namespace, syncSubsystem, "healed_trienodes"), "Number of state trie nodes downloaded.", nil, nil),
		syncHealedTrienodeBytesDesc:    prometheus.NewDesc(prometheus.BuildFQName(namespace, syncSubsystem, "healed_trienode_bytes"), "Number of state trie bytes persisted to disk.", nil, nil),
		syncHealedBytecodesDesc:        prometheus.NewDesc(prometheus.BuildFQName(namespace, syncSubsystem, "healed_bytecodes"), "Number of bytecodes downloaded.", nil, nil),
		syncHealedBytecodeBytesDesc:    prometheus.NewDesc(prometheus.BuildFQName(namespace, syncSubsystem, "healed_bytecode_bytes"), "Number of bytecodes persisted to disk.", nil, nil),
		syncHealingTrienodesDesc:       prometheus.NewDesc(prometheus.BuildFQName(namespace, syncSubsystem, "healing_trienodes"), "Number of state trie nodes pending.", nil, nil),
		syncHealingBytecodeDesc:        prometheus.NewDesc(prometheus.BuildFQName(namespace, syncSubsystem, "healing_bytecode"), "Number of bytecodes pending.", nil, nil),
		syncTxIndexFinishedBlocksDesc:  prometheus.NewDesc(prometheus.BuildFQName(namespace, syncSubsystem, "txindex_finished_blocks"), "Number of blocks whose transactions are already indexed.", nil, nil),
		syncTxIndexRemainingBlocksDesc: prometheus.NewDesc(prometheus.BuildFQName(namespace, syncSubsystem, "txindex_remaining_blocks"), "Number of blocks whose transactions are not indexed yet.", nil, nil),
		syncStateIndexRemainingDesc:    prometheus.NewDesc(prometheus.BuildFQName(namespace, syncSubsystem, "state_index_remaining"), "Number of states remain unindexed.", nil, nil),
	}

	return collector, nil
}

// Close releases the underlying RPC client resources.
func (c *GethRPCCollector) Close() {
	c.client.Close()
}

// Describe announces the metric descriptors to Prometheus.
func (c *GethRPCCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.blockDesc
	ch <- c.chainDesc
	ch <- c.syncDoneDesc
	ch <- c.syncStartingBlockDesc
	ch <- c.syncCurrentBlockDesc
	ch <- c.syncHighestBlockDesc
	ch <- c.syncPulledStatesDesc
	ch <- c.syncKnownStatesDesc
	ch <- c.syncSyncedAccountsDesc
	ch <- c.syncSyncedAccountBytesDesc
	ch <- c.syncSyncedBytecodesDesc
	ch <- c.syncSyncedBytecodeBytesDesc
	ch <- c.syncSyncedStorageDesc
	ch <- c.syncSyncedStorageBytesDesc
	ch <- c.syncHealedTrienodesDesc
	ch <- c.syncHealedTrienodeBytesDesc
	ch <- c.syncHealedBytecodesDesc
	ch <- c.syncHealedBytecodeBytesDesc
	ch <- c.syncHealingTrienodesDesc
	ch <- c.syncHealingBytecodeDesc
	ch <- c.syncTxIndexFinishedBlocksDesc
	ch <- c.syncTxIndexRemainingBlocksDesc
	ch <- c.syncStateIndexRemainingDesc
}

// Collect executes RPC calls per scrape and publishes the resulting metrics.
func (c *GethRPCCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), scrapeRPCTimeout)
	defer cancel()

	func() {
		blockHeight, err := c.client.BlockNumber(ctx)
		if err != nil {
			slog.Error("fetch block number", slog.String("error", err.Error()), slog.String("endpoint", c.endpoint))
			return
		}
		ch <- prometheus.MustNewConstMetric(c.blockDesc, prometheus.GaugeValue, float64(blockHeight))
	}()

	func() {
		chainID, err := c.client.ChainID(ctx)
		if err != nil {
			slog.Error("fetch chain id", slog.String("error", err.Error()), slog.String("endpoint", c.endpoint))
			return
		}
		ch <- prometheus.MustNewConstMetric(c.chainDesc, prometheus.GaugeValue, float64(chainID))
	}()

	func() {
		syncProgress, err := c.client.SyncProgress(ctx)
		if err != nil {
			slog.Error("fetch sync progress", slog.String("error", err.Error()), slog.String("endpoint", c.endpoint))
			return
		}

		// syncProgress is nil when the node is fully synced (not actively syncing).
		// This is normal behavior - eth_syncing returns false when sync is complete.
		// See: https://github.com/ethereum/go-ethereum/blob/1e4b39ed122f475ac3f776ae66c8d065e845a84e/ethclient/ethclient.go#L353

		if syncProgress == nil {
			// Node is fully synced, report sync_done=1 and zero out other sync metrics
			ch <- prometheus.MustNewConstMetric(c.syncDoneDesc, prometheus.GaugeValue, 1)
			return
		}

		doneValue := 0.0
		if syncProgress.Done() {
			doneValue = 1
		}
		ch <- prometheus.MustNewConstMetric(c.syncDoneDesc, prometheus.GaugeValue, doneValue)

		ch <- prometheus.MustNewConstMetric(c.syncStartingBlockDesc, prometheus.GaugeValue, float64(syncProgress.StartingBlock))
		ch <- prometheus.MustNewConstMetric(c.syncCurrentBlockDesc, prometheus.GaugeValue, float64(syncProgress.CurrentBlock))
		ch <- prometheus.MustNewConstMetric(c.syncHighestBlockDesc, prometheus.GaugeValue, float64(syncProgress.HighestBlock))
		ch <- prometheus.MustNewConstMetric(c.syncPulledStatesDesc, prometheus.GaugeValue, float64(syncProgress.PulledStates))
		ch <- prometheus.MustNewConstMetric(c.syncKnownStatesDesc, prometheus.GaugeValue, float64(syncProgress.KnownStates))
		ch <- prometheus.MustNewConstMetric(c.syncSyncedAccountsDesc, prometheus.GaugeValue, float64(syncProgress.SyncedAccounts))
		ch <- prometheus.MustNewConstMetric(c.syncSyncedAccountBytesDesc, prometheus.GaugeValue, float64(syncProgress.SyncedAccountBytes))
		ch <- prometheus.MustNewConstMetric(c.syncSyncedBytecodesDesc, prometheus.GaugeValue, float64(syncProgress.SyncedBytecodes))
		ch <- prometheus.MustNewConstMetric(c.syncSyncedBytecodeBytesDesc, prometheus.GaugeValue, float64(syncProgress.SyncedBytecodeBytes))
		ch <- prometheus.MustNewConstMetric(c.syncSyncedStorageDesc, prometheus.GaugeValue, float64(syncProgress.SyncedStorage))
		ch <- prometheus.MustNewConstMetric(c.syncSyncedStorageBytesDesc, prometheus.GaugeValue, float64(syncProgress.SyncedStorageBytes))
		ch <- prometheus.MustNewConstMetric(c.syncHealedTrienodesDesc, prometheus.GaugeValue, float64(syncProgress.HealedTrienodes))
		ch <- prometheus.MustNewConstMetric(c.syncHealedTrienodeBytesDesc, prometheus.GaugeValue, float64(syncProgress.HealedTrienodeBytes))
		ch <- prometheus.MustNewConstMetric(c.syncHealedBytecodesDesc, prometheus.GaugeValue, float64(syncProgress.HealedBytecodes))
		ch <- prometheus.MustNewConstMetric(c.syncHealedBytecodeBytesDesc, prometheus.GaugeValue, float64(syncProgress.HealedBytecodeBytes))
		ch <- prometheus.MustNewConstMetric(c.syncHealingTrienodesDesc, prometheus.GaugeValue, float64(syncProgress.HealingTrienodes))
		ch <- prometheus.MustNewConstMetric(c.syncHealingBytecodeDesc, prometheus.GaugeValue, float64(syncProgress.HealingBytecode))
		ch <- prometheus.MustNewConstMetric(c.syncTxIndexFinishedBlocksDesc, prometheus.GaugeValue, float64(syncProgress.TxIndexFinishedBlocks))
		ch <- prometheus.MustNewConstMetric(c.syncTxIndexRemainingBlocksDesc, prometheus.GaugeValue, float64(syncProgress.TxIndexRemainingBlocks))
		ch <- prometheus.MustNewConstMetric(c.syncStateIndexRemainingDesc, prometheus.GaugeValue, float64(syncProgress.StateIndexRemaining))
	}()
}
