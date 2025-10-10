package lnd

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/connmgr"
	"github.com/btcsuite/btcd/wire"
	"github.com/lightningnetwork/lnd/aliasmgr"
	"github.com/lightningnetwork/lnd/brontide"
	"github.com/lightningnetwork/lnd/chainntnfs"
	"github.com/lightningnetwork/lnd/channeldb"
	"github.com/lightningnetwork/lnd/channelnotifier"
	"github.com/lightningnetwork/lnd/discovery"
	"github.com/lightningnetwork/lnd/feature"
	"github.com/lightningnetwork/lnd/fn/v2"
	graphdb "github.com/lightningnetwork/lnd/graph/db"
	"github.com/lightningnetwork/lnd/graph/db/models"
	"github.com/lightningnetwork/lnd/htlcswitch"
	"github.com/lightningnetwork/lnd/htlcswitch/hodl"
	"github.com/lightningnetwork/lnd/keychain"
	"github.com/lightningnetwork/lnd/kvdb"
	"github.com/lightningnetwork/lnd/lncfg"
	"github.com/lightningnetwork/lnd/lnpeer"
	"github.com/lightningnetwork/lnd/lntest/mock"
	"github.com/lightningnetwork/lnd/lnwallet"
	"github.com/lightningnetwork/lnd/lnwallet/chainfee"
	"github.com/lightningnetwork/lnd/lnwallet/chancloser"
	"github.com/lightningnetwork/lnd/lnwire"
	"github.com/lightningnetwork/lnd/msgmux"
	"github.com/lightningnetwork/lnd/netann"
	"github.com/lightningnetwork/lnd/peer"
	"github.com/lightningnetwork/lnd/pool"
	"github.com/lightningnetwork/lnd/queue"
	"github.com/lightningnetwork/lnd/routing/route"
	"github.com/lightningnetwork/lnd/ticker"
)

const (
	defaultChannelCommitBatchSize = 10
	defaultCoopCloseTargetConfs   = 6
)

var (
	chanEnableTimeout            = 19 * time.Minute
	defaultChannelCommitInterval = 50 * time.Millisecond
	defaultPendingCommitInterval = 1 * time.Minute
)

func ConnectPeer(conn *brontide.Conn, connReq *connmgr.ConnReq,
	netParams *chaincfg.Params,
	identityECDH keychain.SingleKeyECDH) (*peer.Brontide, *channeldb.DB,
	error) {

	featureMgr, err := feature.NewManager(feature.Config{})
	if err != nil {
		return nil, nil, err
	}

	initFeatures := featureMgr.Get(feature.SetInit)
	legacyFeatures := featureMgr.Get(feature.SetLegacyGlobal)

	addr := conn.RemoteAddr()
	pubKey := conn.RemotePub()
	peerAddr := &lnwire.NetAddress{
		IdentityKey: pubKey,
		Address:     addr,
		ChainNet:    netParams.Net,
	}
	errBuffer, err := queue.NewCircularBuffer(500)
	if err != nil {
		return nil, nil, err
	}

	pongBuf := make([]byte, lnwire.MaxPongBytes)

	writeBufferPool := pool.NewWriteBuffer(
		pool.DefaultWriteBufferGCInterval,
		pool.DefaultWriteBufferExpiryInterval,
	)
	writePool := pool.NewWrite(
		writeBufferPool, lncfg.DefaultWriteWorkers,
		pool.DefaultWorkerTimeout,
	)

	readBufferPool := pool.NewReadBuffer(
		pool.DefaultReadBufferGCInterval,
		pool.DefaultReadBufferExpiryInterval,
	)
	readPool := pool.NewRead(
		readBufferPool, lncfg.DefaultWriteWorkers,
		pool.DefaultWorkerTimeout,
	)
	commitFee := chainfee.SatPerKVByte(
		lnwallet.DefaultAnchorsCommitMaxFeeRateSatPerVByte * 1000,
	)

	if err := writePool.Start(); err != nil {
		return nil, nil, fmt.Errorf("unable to start write pool: %w",
			err)
	}
	if err := readPool.Start(); err != nil {
		return nil, nil, fmt.Errorf("unable to start read pool: %w",
			err)
	}

	randNum := rand.Int31()
	backend, err := kvdb.GetBoltBackend(&kvdb.BoltBackendConfig{
		DBPath:            os.TempDir(),
		DBFileName:        fmt.Sprintf("channel-%d.db", randNum),
		NoFreelistSync:    true,
		AutoCompact:       false,
		AutoCompactMinAge: kvdb.DefaultBoltAutoCompactMinAge,
		DBTimeout:         kvdb.DefaultDBTimeout,
	})
	if err != nil {
		return nil, nil, err
	}

	channelDB, err := channeldb.CreateWithBackend(backend)
	if err != nil {
		_ = backend.Close()
		return nil, nil, err
	}

	graphDB, err := graphdb.NewChannelGraph(&graphdb.Config{
		KVDB: backend,
	})
	if err != nil {
		_ = backend.Close()
		_ = channelDB.Close()

		return nil, nil, fmt.Errorf("unable to open graph db: %w",
			err)
	}

	gossiper := discovery.New(discovery.Config{
		ChainHash: *netParams.GenesisHash,
		Broadcast: func(_ map[route.Vertex]struct{},
			_ ...lnwire.Message) error {

			return nil
		},
		NotifyWhenOnline: func([33]byte, chan<- lnpeer.Peer) {
		},
		NotifyWhenOffline: func(_ [33]byte) <-chan struct{} {
			return make(chan struct{})
		},
		FetchSelfAnnouncement: func() lnwire.NodeAnnouncement {
			return lnwire.NodeAnnouncement{}
		},
		ProofMatureDelta:    0,
		TrickleDelay:        time.Millisecond * 50,
		RetransmitTicker:    ticker.New(time.Minute * 30),
		RebroadcastInterval: time.Hour * 24,
		RotateTicker: ticker.New(
			discovery.DefaultSyncerRotationInterval,
		),
		HistoricalSyncTicker: ticker.New(
			discovery.DefaultHistoricalSyncInterval,
		),
		NumActiveSyncers:        0,
		MinimumBatchSize:        10,
		SubBatchDelay:           discovery.DefaultSubBatchDelay,
		IgnoreHistoricalFilters: true,
		PinnedSyncers:           make(map[route.Vertex]struct{}),
		MaxChannelUpdateBurst:   discovery.DefaultMaxChannelUpdateBurst,
		ChannelUpdateInterval:   discovery.DefaultChannelUpdateInterval,
		IsAlias:                 aliasmgr.IsAlias,
		SignAliasUpdate: func(
			*lnwire.ChannelUpdate1) (*ecdsa.Signature, error) {

			return nil, errors.New("unimplemented")
		},
		FindBaseByAlias: func(
			lnwire.ShortChannelID) (lnwire.ShortChannelID, error) {

			return lnwire.ShortChannelID{},
				errors.New("unimplemented")
		},
		GetAlias: func(_ lnwire.ChannelID) (lnwire.ShortChannelID,
			error) {

			return lnwire.ShortChannelID{},
				errors.New("unimplemented")
		},
		FindChannel: func(*btcec.PublicKey,
			lnwire.ChannelID) (*channeldb.OpenChannel, error) {

			return nil, errors.New("unimplemented")
		},
	}, &keychain.KeyDescriptor{
		KeyLocator: keychain.KeyLocator{},
		PubKey:     identityECDH.PubKey(),
	})

	chanStatusMgr, err := netann.NewChanStatusManager(
		&netann.ChanStatusConfig{
			ChanStatusSampleInterval: 30 * time.Second,
			// Enable + Sample Interval must be <= DisableTimeout.
			ChanEnableTimeout:  30 * time.Second,
			ChanDisableTimeout: 2 * time.Minute,
			DB:                 channelDB.ChannelStateDB(),
			Graph:              graphDB,
			OurPubKey:          identityECDH.PubKey(),
			IsChannelActive: func(lnwire.ChannelID) bool {
				return true
			},
			ApplyChannelUpdate: func(*lnwire.ChannelUpdate1,
				*wire.OutPoint, bool) error {

				return nil
			},
		})
	if err != nil {
		_ = channelDB.Close()
		return nil, nil, fmt.Errorf("unable to create channel status "+
			"manager: %w", err)
	}

	channelNotifier := channelnotifier.New(channelDB.ChannelStateDB())
	interceptableSwitchNotifier := &mock.ChainNotifier{
		EpochChan: make(chan *chainntnfs.BlockEpoch, 1),
	}
	interceptableSwitchNotifier.EpochChan <- &chainntnfs.BlockEpoch{
		Height: 1,
	}
	interceptableSwitch, err := htlcswitch.NewInterceptableSwitch(
		&htlcswitch.InterceptableSwitchConfig{
			CltvRejectDelta:    13,
			CltvInterceptDelta: 16,
			Notifier:           interceptableSwitchNotifier,
		},
	)
	if err != nil {
		_ = channelDB.Close()
		return nil, nil, fmt.Errorf("unable to create interceptable "+
			"switch: %w", err)
	}

	pCfg := peer.Config{
		Conn:    conn,
		ConnReq: connReq,
		PubKeyBytes: [33]byte(
			identityECDH.PubKey().SerializeCompressed(),
		),
		Addr:                    peerAddr,
		Features:                initFeatures,
		LegacyFeatures:          legacyFeatures,
		OutgoingCltvRejectDelta: lncfg.DefaultOutgoingCltvRejectDelta,
		ChanActiveTimeout:       chanEnableTimeout,
		ErrorBuffer:             errBuffer,
		WritePool:               writePool,
		ReadPool:                readPool,
		Switch:                  &mockMessageSwitch{},
		InterceptSwitch:         interceptableSwitch,
		ChannelDB:               channelDB.ChannelStateDB(),
		ChainArb:                nil,
		AuthGossiper:            gossiper,
		ChanStatusMgr:           chanStatusMgr,
		ChainIO:                 &mock.ChainIO{},
		FeeEstimator:            nil,
		Signer:                  nil,
		SigPool:                 nil,
		Wallet: &lnwallet.LightningWallet{
			WalletController: &mock.WalletController{},
		},
		ChainNotifier: &mock.ChainNotifier{},
		BestBlockView: chainntnfs.NewBestBlockTracker(
			&mock.ChainNotifier{},
		),
		RoutingPolicy:   models.ForwardingPolicy{},
		Sphinx:          nil,
		WitnessBeacon:   nil,
		Invoices:        nil,
		ChannelNotifier: channelNotifier,
		HtlcNotifier:    nil,
		TowerClient:     nil,
		DisconnectPeer: func(key *btcec.PublicKey) error {
			fmt.Printf("Peer %x disconnected\n",
				key.SerializeCompressed())
			return nil
		},
		GenNodeAnnouncement: func(_ ...netann.NodeAnnModifier) (
			lnwire.NodeAnnouncement, error) {

			return lnwire.NodeAnnouncement{},
				errors.New("unimplemented")
		},
		PrunePersistentPeerConnection: func(_ [33]byte) {},
		FetchLastChanUpdate: func(_ lnwire.ShortChannelID) (
			*lnwire.ChannelUpdate1, error) {

			return nil, errors.New("unimplemented")
		},
		FundingManager:          nil,
		Hodl:                    &hodl.Config{},
		UnsafeReplay:            false,
		MaxOutgoingCltvExpiry:   htlcswitch.DefaultMaxOutgoingCltvExpiry,
		MaxChannelFeeAllocation: htlcswitch.DefaultMaxLinkFeeAllocation,
		MaxAnchorsCommitFeeRate: commitFee.FeePerKWeight(),
		CoopCloseTargetConfs:    defaultCoopCloseTargetConfs,
		ServerPubKey:            [33]byte{},
		ChannelCommitInterval:   defaultChannelCommitInterval,
		PendingCommitInterval:   defaultPendingCommitInterval,
		ChannelCommitBatchSize:  defaultChannelCommitBatchSize,
		HandleCustomMessage: func(peer [33]byte,
			msg *lnwire.Custom) error {

			fmt.Printf("Received custom message from %x: %v\n",
				peer[:], msg)
			return nil
		},
		GetAliases: func(
			_ lnwire.ShortChannelID) []lnwire.ShortChannelID {

			return nil
		},
		RequestAlias: func() (lnwire.ShortChannelID, error) {
			return lnwire.ShortChannelID{}, nil
		},
		AddLocalAlias: func(_, _ lnwire.ShortChannelID, _,
			_ bool) error {

			return nil
		},
		AuxLeafStore:              fn.None[lnwallet.AuxLeafStore](),
		AuxSigner:                 fn.None[lnwallet.AuxSigner](),
		AuxResolver:               fn.None[lnwallet.AuxContractResolver](),
		AuxTrafficShaper:          fn.None[htlcswitch.AuxTrafficShaper](),
		PongBuf:                   pongBuf,
		DisallowRouteBlinding:     false,
		DisallowQuiescence:        false,
		MaxFeeExposure:            0,
		MsgRouter:                 fn.None[msgmux.Router](),
		AuxChanCloser:             fn.None[chancloser.AuxChanCloser](),
		ShouldFwdExpEndorsement:   nil,
		NoDisconnectOnPongFailure: false,
		Quit:                      make(chan struct{}),
	}

	copy(pCfg.PubKeyBytes[:], peerAddr.IdentityKey.SerializeCompressed())
	copy(pCfg.ServerPubKey[:], identityECDH.PubKey().SerializeCompressed())

	p := peer.NewBrontide(pCfg)
	if err := p.Start(); err != nil {
		_ = channelDB.Close()
		return nil, nil, err
	}

	return p, channelDB, nil
}
