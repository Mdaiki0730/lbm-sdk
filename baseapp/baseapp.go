package baseapp

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/gogo/protobuf/proto"

	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	ocabci "github.com/line/ostracon/abci/types"
	"github.com/line/ostracon/crypto/tmhash"
	"github.com/line/ostracon/libs/log"
	ocproto "github.com/line/ostracon/proto/ostracon/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/line/lbm-sdk/codec/types"
	"github.com/line/lbm-sdk/server/config"
	"github.com/line/lbm-sdk/snapshots"
	"github.com/line/lbm-sdk/store"
	"github.com/line/lbm-sdk/store/rootmulti"
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/auth/legacy/legacytx"
)

var _ ocabci.Application = (*BaseApp)(nil)

type (
	// StoreLoader defines a customizable function to control how we load the CommitMultiStore
	// from disk. This is useful for state migration, when loading a datastore written with
	// an older version of the software. In particular, if a module changed the substore key name
	// (or removed a substore) between two versions of the software.
	StoreLoader func(ms sdk.CommitMultiStore) error
)

// BaseApp reflects the ABCI application implementation.
type BaseApp struct { //nolint: maligned
	// initialized on creation
	logger            log.Logger
	name              string // application name from abci.Info
	interfaceRegistry types.InterfaceRegistry
	txDecoder         sdk.TxDecoder // unmarshal []byte into sdk.Tx

	anteHandler sdk.AnteHandler // ante handler for fee and auth

	appStore
	baseappVersions
	peerFilters
	snapshotData
	abciData
	moduleRouter

	// volatile states:
	//
	// checkState is set on InitChain and reset on Commit
	// deliverState is set on InitChain and BeginBlock and set to nil on Commit
	checkState   *state // for CheckTx
	deliverState *state // for DeliverTx

	checkStateMtx sync.RWMutex

	checkAccountWGs *AccountWGs
	chCheckTx       chan *RequestCheckTxAsync
	chCheckTxSize   uint // chCheckTxSize is the initial size for chCheckTx

	// paramStore is used to query for ABCI consensus parameters from an
	// application parameter store.
	paramStore ParamStore

	// The minimum gas prices a validator is willing to accept for processing a
	// transaction. This is mainly used for DoS and spam prevention.
	minGasPrices sdk.DecCoins

	// initialHeight is the initial height at which we start the baseapp
	initialHeight int64

	// flag for sealing options and parameters to a BaseApp
	sealed bool

	// block height at which to halt the chain and gracefully shutdown
	haltHeight uint64

	// minimum block time (in Unix seconds) at which to halt the chain and gracefully shutdown
	haltTime uint64

	// minRetainBlocks defines the minimum block height offset from the current
	// block being committed, such that all blocks past this offset are pruned
	// from Tendermint. It is used as part of the process of determining the
	// ResponseCommit.RetainHeight value during ABCI Commit. A value of 0 indicates
	// that no blocks should be pruned.
	//
	// Note: Tendermint block pruning is dependant on this parameter in conunction
	// with the unbonding (safety threshold) period, state pruning and state sync
	// snapshot parameters to determine the correct minimum value of
	// ResponseCommit.RetainHeight.
	minRetainBlocks uint64

	// recovery handler for app.runTx method
	runTxRecoveryMiddleware recoveryMiddleware

	// trace set will return full stack traces for errors in ABCI Log field
	trace bool

	// indexEvents defines the set of events in the form {eventType}.{attributeKey},
	// which informs Tendermint what to index. If empty, all events will be indexed.
	indexEvents map[string]struct{}

	// abciListeners for hooking into the ABCI message processing of the BaseApp
	// and exposing the requests and responses to external consumers
	abciListeners []ABCIListener
}

type appStore struct {
	db          dbm.DB               // common DB backend
	cms         sdk.CommitMultiStore // Main (uncached) state
	storeLoader StoreLoader          // function to handle store loading, may be overridden with SetStoreLoader()

	// an inter-block write-through cache provided to the context during deliverState
	interBlockCache sdk.MultiStorePersistentCache

	fauxMerkleMode bool // if true, IAVL MountStores uses MountStoresDB for simulation speed.
}

type moduleRouter struct {
	router           sdk.Router        // handle any kind of message
	queryRouter      sdk.QueryRouter   // router for redirecting query calls
	grpcQueryRouter  *GRPCQueryRouter  // router for redirecting gRPC query calls
	msgServiceRouter *MsgServiceRouter // router for redirecting Msg service messages
}

type abciData struct {
	initChainer  sdk.InitChainer  // initialize state with validators and state blob
	beginBlocker sdk.BeginBlocker // logic to run before any txs
	endBlocker   sdk.EndBlocker   // logic to run after all txs, and to determine valset changes

	// absent validators from begin block
	voteInfos []abci.VoteInfo
}

type baseappVersions struct {
	// application's version string
	version string

	// application's protocol version that increments on every upgrade
	// if BaseApp is passed to the upgrade keeper's NewKeeper method.
	appVersion uint64
}

// should really get handled in some db struct
// which then has a sub-item, persistence fields
type snapshotData struct {
	// manages snapshots, i.e. dumps of app state at certain intervals
	snapshotManager    *snapshots.Manager
	snapshotInterval   uint64 // block interval between state sync snapshots
	snapshotKeepRecent uint32 // recent state sync snapshots to keep
}

// NewBaseApp returns a reference to an initialized BaseApp. It accepts a
// variadic number of option functions, which act on the BaseApp to set
// configuration choices.
//
// NOTE: The db is used to store the version number for now.
func NewBaseApp(
	name string, logger log.Logger, db dbm.DB, txDecoder sdk.TxDecoder, options ...func(*BaseApp),
) *BaseApp {
	app := &BaseApp{
		logger: logger,
		name:   name,
		appStore: appStore{
			db:             db,
			cms:            store.NewCommitMultiStore(db),
			storeLoader:    DefaultStoreLoader,
			fauxMerkleMode: false,
		},
		moduleRouter: moduleRouter{
			router:           NewRouter(),
			queryRouter:      NewQueryRouter(),
			grpcQueryRouter:  NewGRPCQueryRouter(),
			msgServiceRouter: NewMsgServiceRouter(),
		},
		txDecoder:       txDecoder,
		checkAccountWGs: NewAccountWGs(),
	}

	for _, option := range options {
		option(app)
	}

	chCheckTxSize := app.chCheckTxSize
	if chCheckTxSize == 0 {
		chCheckTxSize = config.DefaultChanCheckTxSize
	}
	app.chCheckTx = make(chan *RequestCheckTxAsync, chCheckTxSize)

	if app.interBlockCache != nil {
		app.cms.SetInterBlockCache(app.interBlockCache)
	}

	app.runTxRecoveryMiddleware = newDefaultRecoveryMiddleware()

	app.startReactors()

	return app
}

// Name returns the name of the BaseApp.
func (app *BaseApp) Name() string {
	return app.name
}

// AppVersion returns the application's protocol version.
func (app *BaseApp) AppVersion() uint64 {
	return app.appVersion
}

// Version returns the application's version string.
func (app *BaseApp) Version() string {
	return app.version
}

// Logger returns the logger of the BaseApp.
func (app *BaseApp) Logger() log.Logger {
	return app.logger
}

// Trace returns the boolean value for logging error stack traces.
func (app *BaseApp) Trace() bool {
	return app.trace
}

// MsgServiceRouter returns the MsgServiceRouter of a BaseApp.
func (app *BaseApp) MsgServiceRouter() *MsgServiceRouter { return app.msgServiceRouter }

// MountStores mounts all IAVL or DB stores to the provided keys in the BaseApp
// multistore.
func (app *BaseApp) MountStores(keys ...sdk.StoreKey) {
	for _, key := range keys {
		switch key.(type) {
		case *sdk.KVStoreKey:
			if !app.fauxMerkleMode {
				app.MountStore(key, sdk.StoreTypeIAVL)
			} else {
				// StoreTypeDB doesn't do anything upon commit, and it doesn't
				// retain history, but it's useful for faster simulation.
				app.MountStore(key, sdk.StoreTypeDB)
			}

		default:
			panic("Unrecognized store key type " + reflect.TypeOf(key).Name())
		}
	}
}

// MountKVStores mounts all IAVL or DB stores to the provided keys in the
// BaseApp multistore.
func (app *BaseApp) MountKVStores(keys map[string]*sdk.KVStoreKey) {
	for _, key := range keys {
		if !app.fauxMerkleMode {
			app.MountStore(key, sdk.StoreTypeIAVL)
		} else {
			// StoreTypeDB doesn't do anything upon commit, and it doesn't
			// retain history, but it's useful for faster simulation.
			app.MountStore(key, sdk.StoreTypeDB)
		}
	}
}

// MountMemoryStores mounts all in-memory KVStores with the BaseApp's internal
// commit multi-store.
func (app *BaseApp) MountMemoryStores(keys map[string]*sdk.MemoryStoreKey) {
	for _, memKey := range keys {
		app.MountStore(memKey, sdk.StoreTypeMemory)
	}
}

// MountStore mounts a store to the provided key in the BaseApp multistore,
// using the default DB.
func (app *BaseApp) MountStore(key sdk.StoreKey, typ sdk.StoreType) {
	app.cms.MountStoreWithDB(key, typ, nil)
}

// LoadLatestVersion loads the latest application version. It will panic if
// called more than once on a running BaseApp.
func (app *BaseApp) LoadLatestVersion() error {
	err := app.storeLoader(app.cms)
	if err != nil {
		return fmt.Errorf("failed to load latest version: %w", err)
	}

	return app.init()
}

// DefaultStoreLoader will be used by default and loads the latest version
func DefaultStoreLoader(ms sdk.CommitMultiStore) error {
	return ms.LoadLatestVersion()
}

// CommitMultiStore returns the root multi-store.
// App constructor can use this to access the `cms`.
// UNSAFE: must not be used during the abci life cycle.
func (app *BaseApp) CommitMultiStore() sdk.CommitMultiStore {
	return app.cms
}

// SnapshotManager returns the snapshot manager.
// application use this to register extra extension snapshotters.
func (app *BaseApp) SnapshotManager() *snapshots.Manager {
	return app.snapshotManager
}

// LoadVersion loads the BaseApp application version. It will panic if called
// more than once on a running baseapp.
func (app *BaseApp) LoadVersion(version int64) error {
	err := app.cms.LoadVersion(version)
	if err != nil {
		return fmt.Errorf("failed to load version %d: %w", version, err)
	}

	return app.init()
}

// LastCommitID returns the last CommitID of the multistore.
func (app *BaseApp) LastCommitID() sdk.CommitID {
	return app.cms.LastCommitID()
}

// LastBlockHeight returns the last committed block height.
func (app *BaseApp) LastBlockHeight() int64 {
	return app.cms.LastCommitID().Version
}

func (app *BaseApp) init() error {
	if app.sealed {
		panic("cannot call initFromMainStore: baseapp already sealed")
	}

	// needed for the export command which inits from store but never calls initchain
	app.setCheckState(ocproto.Header{})
	app.Seal()

	// make sure the snapshot interval is a multiple of the pruning KeepEvery interval
	if app.snapshotManager != nil && app.snapshotInterval > 0 {
		rms, ok := app.cms.(*rootmulti.Store)
		if !ok {
			return errors.New("state sync snapshots require a rootmulti store")
		}
		pruningOpts := rms.GetPruning()
		if pruningOpts.KeepEvery > 0 && app.snapshotInterval%pruningOpts.KeepEvery != 0 {
			return fmt.Errorf(
				"state sync snapshot interval %v must be a multiple of pruning keep every interval %v",
				app.snapshotInterval, pruningOpts.KeepEvery)
		}
	}

	return nil
}

func (app *BaseApp) setMinGasPrices(gasPrices sdk.DecCoins) {
	app.minGasPrices = gasPrices
}

func (app *BaseApp) setHaltHeight(haltHeight uint64) {
	app.haltHeight = haltHeight
}

func (app *BaseApp) setHaltTime(haltTime uint64) {
	app.haltTime = haltTime
}

func (app *BaseApp) setMinRetainBlocks(minRetainBlocks uint64) {
	app.minRetainBlocks = minRetainBlocks
}

func (app *BaseApp) setInterBlockCache(cache sdk.MultiStorePersistentCache) {
	app.interBlockCache = cache
}

func (app *BaseApp) setTrace(trace bool) {
	app.trace = trace
}

func (app *BaseApp) setIndexEvents(ie []string) {
	app.indexEvents = make(map[string]struct{})

	for _, e := range ie {
		app.indexEvents[e] = struct{}{}
	}
}

// Router returns the router of the BaseApp.
func (app *BaseApp) Router() sdk.Router {
	if app.sealed {
		// We cannot return a Router when the app is sealed because we can't have
		// any routes modified which would cause unexpected routing behavior.
		panic("Router() on sealed BaseApp")
	}

	return app.router
}

// QueryRouter returns the QueryRouter of a BaseApp.
func (app *BaseApp) QueryRouter() sdk.QueryRouter { return app.queryRouter }

// Seal seals a BaseApp. It prohibits any further modifications to a BaseApp.
func (app *BaseApp) Seal() { app.sealed = true }

// IsSealed returns true if the BaseApp is sealed and false otherwise.
func (app *BaseApp) IsSealed() bool { return app.sealed }

// setCheckState sets the BaseApp's checkState with a branched multi-store
// (i.e. a CacheMultiStore) and a new Context with the same multi-store branch,
// provided header, and minimum gas prices set. It is set on InitChain and reset
// on Commit.
func (app *BaseApp) setCheckState(header ocproto.Header) {
	ms := app.cms.CacheMultiStore()
	app.checkStateMtx.Lock()
	defer app.checkStateMtx.Unlock()

	ctx := sdk.NewContext(ms, header, true, app.logger).
		WithMinGasPrices(app.minGasPrices).
		WithVoteInfos(app.voteInfos)

	app.checkState = &state{
		ms:  ms,
		ctx: ctx,
	}
}

// setDeliverState sets the BaseApp's deliverState with a branched multi-store
// (i.e. a CacheMultiStore) and a new Context with the same multi-store branch,
// and provided header. It is set on InitChain and BeginBlock and set to nil on
// Commit.
func (app *BaseApp) setDeliverState(header ocproto.Header) {
	ms := app.cms.CacheMultiStore()
	app.deliverState = &state{
		ms:  ms,
		ctx: sdk.NewContext(ms, header, false, app.logger),
	}
}

// GetConsensusParams returns the current consensus parameters from the BaseApp's
// ParamStore. If the BaseApp has no ParamStore defined, nil is returned.
func (app *BaseApp) GetConsensusParams(ctx sdk.Context) *abci.ConsensusParams {
	if app.paramStore == nil {
		return nil
	}

	cp := new(abci.ConsensusParams)

	if app.paramStore.Has(ctx, ParamStoreKeyBlockParams) {
		var bp abci.BlockParams

		app.paramStore.Get(ctx, ParamStoreKeyBlockParams, &bp)
		cp.Block = &bp
	}

	if app.paramStore.Has(ctx, ParamStoreKeyEvidenceParams) {
		var ep tmproto.EvidenceParams

		app.paramStore.Get(ctx, ParamStoreKeyEvidenceParams, &ep)
		cp.Evidence = &ep
	}

	if app.paramStore.Has(ctx, ParamStoreKeyValidatorParams) {
		var vp tmproto.ValidatorParams

		app.paramStore.Get(ctx, ParamStoreKeyValidatorParams, &vp)
		cp.Validator = &vp
	}

	return cp
}

// AddRunTxRecoveryHandler adds custom app.runTx method panic handlers.
func (app *BaseApp) AddRunTxRecoveryHandler(handlers ...RecoveryHandler) {
	for _, h := range handlers {
		app.runTxRecoveryMiddleware = newRecoveryMiddleware(h, app.runTxRecoveryMiddleware)
	}
}

// StoreConsensusParams sets the consensus parameters to the baseapp's param store.
func (app *BaseApp) StoreConsensusParams(ctx sdk.Context, cp *abci.ConsensusParams) {
	if app.paramStore == nil {
		panic("cannot store consensus params with no params store set")
	}

	if cp == nil {
		return
	}

	app.paramStore.Set(ctx, ParamStoreKeyBlockParams, cp.Block)
	app.paramStore.Set(ctx, ParamStoreKeyEvidenceParams, cp.Evidence)
	app.paramStore.Set(ctx, ParamStoreKeyValidatorParams, cp.Validator)
}

// getMaximumBlockGas gets the maximum gas from the consensus params. It panics
// if maximum block gas is less than negative one and returns zero if negative
// one.
func (app *BaseApp) getMaximumBlockGas(ctx sdk.Context) uint64 {
	cp := app.GetConsensusParams(ctx)
	if cp == nil || cp.Block == nil {
		return 0
	}

	maxGas := cp.Block.MaxGas

	switch {
	case maxGas < -1:
		panic(fmt.Sprintf("invalid maximum block gas: %d", maxGas))

	case maxGas == -1:
		return 0

	default:
		return uint64(maxGas)
	}
}

func (app *BaseApp) validateHeight(req ocabci.RequestBeginBlock) error {
	if req.Header.Height < 1 {
		return fmt.Errorf("invalid height: %d", req.Header.Height)
	}

	// expectedHeight holds the expected height to validate.
	var expectedHeight int64
	if app.LastBlockHeight() == 0 && app.initialHeight > 1 {
		// In this case, we're validating the first block of the chain (no
		// previous commit). The height we're expecting is the initial height.
		expectedHeight = app.initialHeight
	} else {
		// This case can means two things:
		// - either there was already a previous commit in the store, in which
		// case we increment the version from there,
		// - or there was no previous commit, and initial version was not set,
		// in which case we start at version 1.
		expectedHeight = app.LastBlockHeight() + 1
	}

	if req.Header.Height != expectedHeight {
		return fmt.Errorf("invalid height: %d; expected: %d", req.Header.Height, expectedHeight)
	}

	return nil
}

// validateBasicTxMsgs executes basic validator calls for messages.
func validateBasicTxMsgs(msgs []sdk.Msg) error {
	if len(msgs) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "must contain at least one message")
	}

	for _, msg := range msgs {
		err := msg.ValidateBasic()
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *BaseApp) getCheckContextForTx(txBytes []byte, recheck bool) sdk.Context {
	app.checkStateMtx.RLock()
	defer app.checkStateMtx.RUnlock()
	return app.getContextForTx(app.checkState, txBytes).WithIsReCheckTx(recheck)
}

// retrieve the context for the tx w/ txBytes and other memoized values.
func (app *BaseApp) getRunContextForTx(txBytes []byte, simulate bool) sdk.Context {
	if !simulate {
		return app.getContextForTx(app.deliverState, txBytes)
	}

	app.checkStateMtx.RLock()
	defer app.checkStateMtx.RUnlock()
	ctx := app.getContextForTx(app.checkState, txBytes)
	ctx, _ = ctx.CacheContext()
	return ctx
}

func (app *BaseApp) getContextForTx(s *state, txBytes []byte) sdk.Context {
	ctx := s.ctx.WithTxBytes(txBytes)
	ctx = ctx.WithConsensusParams(app.GetConsensusParams(ctx))
	return ctx
}

// cacheTxContext returns a new context based off of the provided context with
// a branched multi-store.
func (app *BaseApp) cacheTxContext(ctx sdk.Context, txBytes []byte) (sdk.Context, sdk.CacheMultiStore) {
	ms := ctx.MultiStore()
	// TODO: https://github.com/cosmos/cosmos-sdk/issues/2824
	msCache := ms.CacheMultiStore()
	if msCache.TracingEnabled() {
		msCache = msCache.SetTracingContext(
			sdk.TraceContext(
				map[string]interface{}{
					"txHash": fmt.Sprintf("%X", tmhash.Sum(txBytes)),
				},
			),
		).(sdk.CacheMultiStore)
	}

	return ctx.WithMultiStore(msCache), msCache
}

// stateless checkTx
func (app *BaseApp) preCheckTx(txBytes []byte) (tx sdk.Tx, err error) {
	defer func() {
		if r := recover(); r != nil {
			recoveryMW := newDefaultRecoveryMiddleware()
			err = processRecovery(r, recoveryMW)
		}
	}()

	tx, err = app.txDecoder(txBytes)
	if err != nil {
		return tx, err
	}

	msgs := tx.GetMsgs()
	err = validateBasicTxMsgs(msgs)

	return tx, err
}

func (app *BaseApp) PreCheckTx(txBytes []byte) (sdk.Tx, error) {
	return app.preCheckTx(txBytes)
}

func (app *BaseApp) checkTx(txBytes []byte, tx sdk.Tx, recheck bool) (gInfo sdk.GasInfo, err error) {
	ctx := app.getCheckContextForTx(txBytes, recheck)
	gasCtx := &ctx

	defer func() {
		if r := recover(); r != nil {
			recoveryMW := newDefaultRecoveryMiddleware()
			err = processRecovery(r, recoveryMW)
		}
		gInfo = sdk.GasInfo{GasWanted: gasCtx.GasMeter().Limit(), GasUsed: gasCtx.GasMeter().GasConsumed()}
	}()

	var anteCtx sdk.Context
	anteCtx, err = app.anteTx(ctx, txBytes, tx, false)
	if !anteCtx.IsZero() {
		gasCtx = &anteCtx
	}

	return gInfo, err
}

func (app *BaseApp) anteTx(ctx sdk.Context, txBytes []byte, tx sdk.Tx, simulate bool) (sdk.Context, error) {
	if app.anteHandler == nil {
		return ctx, nil
	}

	// Branch context before AnteHandler call in case it aborts.
	// This is required for both CheckTx and DeliverTx.
	// Ref: https://github.com/cosmos/cosmos-sdk/issues/2772
	//
	// NOTE: Alternatively, we could require that AnteHandler ensures that
	// writes do not happen if aborted/failed.  This may have some
	// performance benefits, but it'll be more difficult to get right.
	anteCtx, msCache := app.cacheTxContext(ctx, txBytes)
	anteCtx = anteCtx.WithEventManager(sdk.NewEventManager())
	newCtx, err := app.anteHandler(anteCtx, tx, simulate)

	if err != nil {
		return newCtx, err
	}

	msCache.Write()
	return newCtx, err
}

// runTx processes a transaction within a given execution mode, encoded transaction
// bytes, and the decoded transaction itself. All state transitions occur through
// a cached Context depending on the mode provided. State only gets persisted
// if all messages get executed successfully and the execution mode is DeliverTx.
// Note, gas execution info is always returned. A reference to a Result is
// returned if the tx does not run out of gas and if all the messages are valid
// and execute successfully. An error is returned otherwise.
func (app *BaseApp) runTx(txBytes []byte, tx sdk.Tx, simulate bool) (gInfo sdk.GasInfo, result *sdk.Result, anteEvents []abci.Event, err error) {
	ctx := app.getRunContextForTx(txBytes, simulate)
	ms := ctx.MultiStore()

	// only run the tx if there is block gas remaining
	if !simulate && ctx.BlockGasMeter().IsOutOfGas() {
		return gInfo, nil, nil, sdkerrors.Wrap(sdkerrors.ErrOutOfGas, "no block gas left to run tx")
	}

	defer func() {
		if r := recover(); r != nil {
			recoveryMW := newOutOfGasRecoveryMiddleware(ctx.GasMeter().Limit(), ctx, app.runTxRecoveryMiddleware)
			err, result = processRecovery(r, recoveryMW), nil
		}

		gInfo = sdk.GasInfo{GasWanted: ctx.GasMeter().Limit(), GasUsed: ctx.GasMeter().GasConsumed()}
	}()

	blockGasConsumed := false
	// consumeBlockGas makes sure block gas is consumed at most once. It must happen after
	// tx processing, and must be execute even if tx processing fails. Hence we use trick with `defer`
	consumeBlockGas := func() {
		if !blockGasConsumed {
			blockGasConsumed = true
			ctx.BlockGasMeter().ConsumeGas(
				ctx.GasMeter().GasConsumedToLimit(), "block gas meter",
			)
		}
	}

	// If BlockGasMeter() panics it will be caught by the above recover and will
	// return an error - in any case BlockGasMeter will consume gas past the limit.
	//
	// NOTE: This must exist in a separate defer function for the above recovery
	// to recover from this one.
	if !simulate {
		defer consumeBlockGas()
	}

	msgs := tx.GetMsgs()
	if err = validateBasicTxMsgs(msgs); err != nil {
		return sdk.GasInfo{}, nil, nil, err
	}

	var newCtx sdk.Context
	newCtx, err = app.anteTx(ctx, txBytes, tx, simulate)
	if !newCtx.IsZero() {
		// At this point, newCtx.MultiStore() is a store branch, or something else
		// replaced by the AnteHandler. We want the original multistore.
		//
		// Also, in the case of the tx aborting, we need to track gas consumed via
		// the instantiated gas meter in the AnteHandler, so we update the context
		// prior to returning.
		ctx = newCtx.WithMultiStore(ms)
	}

	if err != nil {
		return gInfo, nil, nil, err
	}

	events := ctx.EventManager().Events()
	anteEvents = events.ToABCIEvents()

	// Create a new Context based off of the existing Context with a MultiStore branch
	// in case message processing fails. At this point, the MultiStore
	// is a branch of a branch.
	runMsgCtx, msCache := app.cacheTxContext(ctx, txBytes)

	// Attempt to execute all messages and only update state if all messages pass
	// and we're in DeliverTx. Note, runMsgs will never return a reference to a
	// Result if any single message fails or does not have a registered Handler.
	result, err = app.runMsgs(runMsgCtx, msgs)
	if err == nil && !simulate {
		// When block gas exceeds, it'll panic and won't commit the cached store.
		consumeBlockGas()

		msCache.Write()

		if len(anteEvents) > 0 {
			// append the events in the order of occurrence
			result.Events = append(anteEvents, result.Events...)
		}
	}

	return gInfo, result, anteEvents, err
}

// runMsgs iterates through a list of messages and executes them with the provided
// Context and execution mode. Messages will only be executed during simulation
// and DeliverTx. An error is returned if any single message fails or if a
// Handler does not exist for a given message route. Otherwise, a reference to a
// Result is returned. The caller must not commit state if an error is returned.
func (app *BaseApp) runMsgs(ctx sdk.Context, msgs []sdk.Msg) (*sdk.Result, error) {
	msgLogs := make(sdk.ABCIMessageLogs, 0, len(msgs))
	events := sdk.EmptyEvents()
	txMsgData := &sdk.TxMsgData{
		Data: make([]*sdk.MsgData, 0, len(msgs)),
	}

	// NOTE: GasWanted is determined by the AnteHandler and GasUsed by the GasMeter.
	for i, msg := range msgs {
		var (
			msgResult    *sdk.Result
			eventMsgName string // name to use as value in event `message.action`
			err          error
		)

		if handler := app.msgServiceRouter.Handler(msg); handler != nil {
			// ADR 031 request type routing
			msgResult, err = handler(ctx, msg)
			eventMsgName = sdk.MsgTypeURL(msg)
		} else if legacyMsg, ok := msg.(legacytx.LegacyMsg); ok {
			// legacy sdk.Msg routing
			// Assuming that the app developer has migrated all their Msgs to
			// proto messages and has registered all `Msg services`, then this
			// path should never be called, because all those Msgs should be
			// registered within the `msgServiceRouter` already.
			msgRoute := legacyMsg.Route()
			eventMsgName = legacyMsg.Type()
			handler := app.router.Route(ctx, msgRoute)
			if handler == nil {
				return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized message route: %s; message index: %d", msgRoute, i)
			}

			msgResult, err = handler(ctx, msg)
		} else {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "can't route message %+v", msg)
		}

		if err != nil {
			return nil, sdkerrors.Wrapf(err, "failed to execute message; message index: %d", i)
		}

		msgEvents := sdk.Events{
			sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyAction, eventMsgName)),
		}
		msgEvents = msgEvents.AppendEvents(msgResult.GetEvents())

		// append message events, data and logs
		//
		// Note: Each message result's data must be length-prefixed in order to
		// separate each result.
		events = events.AppendEvents(msgEvents)

		txMsgData.Data = append(txMsgData.Data, &sdk.MsgData{MsgType: sdk.MsgTypeURL(msg), Data: msgResult.Data})
		msgLogs = append(msgLogs, sdk.NewABCIMessageLog(uint32(i), msgResult.Log, msgEvents))
	}

	data, err := proto.Marshal(txMsgData)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to marshal tx data")
	}

	return &sdk.Result{
		Data:   data,
		Log:    strings.TrimSpace(msgLogs.String()),
		Events: events.ToABCIEvents(),
	}, nil
}
