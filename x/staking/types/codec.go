package types

import (
	"github.com/line/lbm-sdk/codec"
	"github.com/line/lbm-sdk/codec/legacy"
	"github.com/line/lbm-sdk/codec/types"
	cryptocodec "github.com/line/lbm-sdk/crypto/codec"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/msgservice"
	"github.com/line/lbm-sdk/x/authz"
	authzcodec "github.com/line/lbm-sdk/x/authz/codec"
	fdncodec "github.com/line/lbm-sdk/x/foundation/codec"
	govcodec "github.com/line/lbm-sdk/x/gov/codec"
)

// RegisterLegacyAminoCodec registers the necessary x/staking interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgCreateValidator{}, "cosmos-sdk/MsgCreateValidator")
	legacy.RegisterAminoMsg(cdc, &MsgEditValidator{}, "cosmos-sdk/MsgEditValidator")
	legacy.RegisterAminoMsg(cdc, &MsgDelegate{}, "cosmos-sdk/MsgDelegate")
	legacy.RegisterAminoMsg(cdc, &MsgUndelegate{}, "cosmos-sdk/MsgUndelegate")
	legacy.RegisterAminoMsg(cdc, &MsgBeginRedelegate{}, "cosmos-sdk/MsgBeginRedelegate")

	cdc.RegisterInterface((*isStakeAuthorization_Validators)(nil), nil)
	cdc.RegisterConcrete(&StakeAuthorization_AllowList{}, "cosmos-sdk/StakeAuthorization/AllowList", nil)
	cdc.RegisterConcrete(&StakeAuthorization_DenyList{}, "cosmos-sdk/StakeAuthorization/DenyList", nil)
	cdc.RegisterConcrete(&StakeAuthorization{}, "cosmos-sdk/StakeAuthorization", nil)
}

// RegisterInterfaces registers the x/staking interfaces types with the interface registry
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateValidator{},
		&MsgEditValidator{},
		&MsgDelegate{},
		&MsgUndelegate{},
		&MsgBeginRedelegate{},
	)
	registry.RegisterImplementations(
		(*authz.Authorization)(nil),
		&StakeAuthorization{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	sdk.RegisterLegacyAminoCodec(amino)

	// Register all Amino interfaces and concrete types on the authz  and gov Amino codec so that this can later be
	// used to properly serialize MsgGrant, MsgExec and MsgSubmitProposal instances
	RegisterLegacyAminoCodec(authzcodec.Amino)
	RegisterLegacyAminoCodec(govcodec.Amino)
	RegisterLegacyAminoCodec(fdncodec.Amino)
}
