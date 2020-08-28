// Copyright (c) 2017 The ifishnet developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package netsync

import (
	"github.com/ifishnet/hdfd/blockchain"
	"github.com/ifishnet/hdfd/chaincfg"
	"github.com/ifishnet/hdfd/chaincfg/chainhash"
	"github.com/ifishnet/hdfd/mempool"
	"github.com/ifishnet/hdfd/peer"
	"github.com/ifishnet/hdfd/wire"
	"github.com/ifishnet/hdfutil"
)

// PeerNotifier exposes methods to notify peers of status changes to
// transactions, blocks, etc. Currently server (in the main package) implements
// this interface.
type PeerNotifier interface {
	AnnounceNewTransactions(newTxs []*mempool.TxDesc)

	UpdatePeerHeights(latestBlkHash *chainhash.Hash, latestHeight int32, updateSource *peer.Peer)

	RelayInventory(invVect *wire.InvVect, data interface{})

	TransactionConfirmed(tx *hdfutil.Tx)
}

// Config is a configuration struct used to initialize a new SyncManager.
type Config struct {
	PeerNotifier PeerNotifier
	Chain        *blockchain.BlockChain
	TxMemPool    *mempool.TxPool
	ChainParams  *chaincfg.Params

	DisableCheckpoints bool
	MaxPeers           int

	FeeEstimator *mempool.FeeEstimator
}
