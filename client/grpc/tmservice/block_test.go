package tmservice_test

import (
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	
	ocproto "github.com/line/ostracon/proto/ostracon/types"

	"github.com/line/lbm-sdk/client/grpc/tmservice"
	"github.com/line/lbm-sdk/simapp"
)

func (s IntegrationTestSuite) TestGetProtoBlock() {
	val := s.network.Validators[0]
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{})

	height := int64(-1)
	blockID, block, err := tmservice.GetProtoBlock(ctx.Context(), val.ClientCtx, &height)
	s.Require().Equal(tmproto.BlockID{}, blockID)
	s.Require().Nil(block)
	s.Require().Error(err)

	height = int64(1)
	_, _, err = tmservice.GetProtoBlock(ctx.Context(), val.ClientCtx, &height)
	s.Require().NoError(err)
}
