package rootmulti_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/Finschia/ostracon/libs/log"

	"github.com/Finschia/finschia-rdk/l2app"
)

func setup(withGenesis bool, invCheckPeriod uint, db dbm.DB) (*l2app.SimApp, l2app.GenesisState) {
	encCdc := l2app.MakeTestEncodingConfig()
	app := l2app.NewSimApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, l2app.DefaultNodeHome, invCheckPeriod, encCdc, l2app.EmptyAppOptions{})
	if withGenesis {
		return app, l2app.NewDefaultGenesisState(encCdc.Marshaler)
	}
	return app, l2app.GenesisState{}
}

// Setup initializes a new SimApp. A Nop logger is set in SimApp.
func SetupWithDB(isCheckTx bool, db dbm.DB) *l2app.SimApp {
	app, genesisState := setup(!isCheckTx, 5, db)
	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: l2app.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return app
}

func TestRollback(t *testing.T) {
	t.Skip()
	db := dbm.NewMemDB()
	app := SetupWithDB(false, db)
	app.Commit()
	ver0 := app.LastBlockHeight()
	// commit 10 blocks
	for i := int64(1); i <= 10; i++ {
		header := tmproto.Header{
			Height:  ver0 + i,
			AppHash: app.LastCommitID().Hash,
		}
		app.BeginBlock(abci.RequestBeginBlock{Header: header})
		ctx := app.NewContext(false, header)
		store := ctx.KVStore(app.GetKey("bank"))
		store.Set([]byte("key"), []byte(fmt.Sprintf("value%d", i)))
		app.Commit()
	}

	require.Equal(t, ver0+10, app.LastBlockHeight())
	store := app.NewContext(true, tmproto.Header{}).KVStore(app.GetKey("bank"))
	require.Equal(t, []byte("value10"), store.Get([]byte("key")))

	// rollback 5 blocks
	target := ver0 + 5
	require.NoError(t, app.CommitMultiStore().RollbackToVersion(target))
	require.Equal(t, target, app.LastBlockHeight())

	// recreate app to have clean check state
	encCdc := l2app.MakeTestEncodingConfig()
	app = l2app.NewSimApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, l2app.DefaultNodeHome, 5, encCdc, l2app.EmptyAppOptions{})
	store = app.NewContext(true, tmproto.Header{}).KVStore(app.GetKey("bank"))
	require.Equal(t, []byte("value5"), store.Get([]byte("key")))

	// commit another 5 blocks with different values
	for i := int64(6); i <= 10; i++ {
		header := tmproto.Header{
			Height:  ver0 + i,
			AppHash: app.LastCommitID().Hash,
		}
		app.BeginBlock(abci.RequestBeginBlock{Header: header})
		ctx := app.NewContext(false, header)
		store := ctx.KVStore(app.GetKey("bank"))
		store.Set([]byte("key"), []byte(fmt.Sprintf("VALUE%d", i)))
		app.Commit()
	}

	require.Equal(t, ver0+10, app.LastBlockHeight())
	store = app.NewContext(true, tmproto.Header{}).KVStore(app.GetKey("bank"))
	require.Equal(t, []byte("VALUE10"), store.Get([]byte("key")))
}
