package collection

import (
	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

const collectionCodespace = ModuleName

var (
	ErrTokenNotExist                 = sdkerrors.Register(collectionCodespace, 2, "token symbol, token-id does not exist")
	ErrTokenNotMintable              = sdkerrors.Register(collectionCodespace, 3, "token symbol, token-id is not mintable")
	ErrInvalidTokenName              = sdkerrors.Register(collectionCodespace, 4, "token name should not be empty")
	ErrInvalidTokenID                = sdkerrors.Register(collectionCodespace, 5, "invalid token id")
	ErrInvalidTokenDecimals          = sdkerrors.Register(collectionCodespace, 6, "token decimals should be within the range in 0 ~ 18")
	ErrInvalidIssueFT                = sdkerrors.Register(collectionCodespace, 7, "Issuing token with amount[1], decimals[0], mintable[false] prohibited. Issue nft token instead.")
	ErrInvalidAmount                 = sdkerrors.Register(collectionCodespace, 8, "invalid token amount")
	ErrInvalidBaseImgURILength       = sdkerrors.Register(collectionCodespace, 9, "invalid base_img_uri length")
	ErrInvalidNameLength             = sdkerrors.Register(collectionCodespace, 10, "invalid name length")
	ErrInvalidTokenType              = sdkerrors.Register(collectionCodespace, 11, "invalid token type pattern found")
	ErrInvalidTokenIndex             = sdkerrors.Register(collectionCodespace, 12, "invalid token index pattern found")
	ErrCollectionExist               = sdkerrors.Register(collectionCodespace, 13, "collection already exists")
	ErrCollectionNotExist            = sdkerrors.Register(collectionCodespace, 14, "collection does not exists")
	ErrTokenTypeExist                = sdkerrors.Register(collectionCodespace, 15, "token type for contract_id, token-type already exists")
	ErrTokenTypeNotExist             = sdkerrors.Register(collectionCodespace, 16, "token type for contract_id, token-type does not exist")
	ErrTokenTypeFull                 = sdkerrors.Register(collectionCodespace, 17, "all token type for contract_id are occupied")
	ErrTokenIndexFull                = sdkerrors.Register(collectionCodespace, 18, "all non-fungible token index for contract_id, token-type are occupied")
	ErrTokenIDFull                   = sdkerrors.Register(collectionCodespace, 19, "all fungible token-id for contract_id are occupied")
	ErrTokenNoPermission             = sdkerrors.Register(collectionCodespace, 20, "account does not have the permission")
	ErrTokenAlreadyAChild            = sdkerrors.Register(collectionCodespace, 21, "token is already a child of some other")
	ErrTokenNotAChild                = sdkerrors.Register(collectionCodespace, 22, "token is not a child of some other")
	ErrTokenNotOwnedBy               = sdkerrors.Register(collectionCodespace, 23, "token is being not owned by")
	ErrTokenCannotTransferChildToken = sdkerrors.Register(collectionCodespace, 24, "cannot transfer a child token")
	ErrTokenNotNFT                   = sdkerrors.Register(collectionCodespace, 25, "token is not a NFT")
	ErrCannotAttachToItself          = sdkerrors.Register(collectionCodespace, 26, "cannot attach token to itself")
	ErrCannotAttachToADescendant     = sdkerrors.Register(collectionCodespace, 27, "cannot attach token to a descendant")
	ErrApproverProxySame             = sdkerrors.Register(collectionCodespace, 28, "approver is same with proxy")
	ErrCollectionNotApproved         = sdkerrors.Register(collectionCodespace, 29, "proxy is not approved on the collection")
	ErrCollectionAlreadyApproved     = sdkerrors.Register(collectionCodespace, 30, "proxy is already approved on the collection")
	ErrAccountExist                  = sdkerrors.Register(collectionCodespace, 31, "account already exists")
	ErrAccountNotExist               = sdkerrors.Register(collectionCodespace, 32, "account does not exists")
	ErrInsufficientSupply            = sdkerrors.Register(collectionCodespace, 33, "insufficient supply")
	ErrInvalidCoin                   = sdkerrors.Register(collectionCodespace, 34, "invalid coin")
	ErrInvalidChangesFieldCount      = sdkerrors.Register(collectionCodespace, 35, "invalid count of field changes")
	ErrEmptyChanges                  = sdkerrors.Register(collectionCodespace, 36, "changes is empty")
	ErrInvalidChangesField           = sdkerrors.Register(collectionCodespace, 37, "invalid field of changes")
	ErrTokenIndexWithoutType         = sdkerrors.Register(collectionCodespace, 38, "There is a token index but no token type")
	ErrTokenTypeFTWithoutIndex       = sdkerrors.Register(collectionCodespace, 39, "There is a token type of ft but no token index")
	ErrInsufficientToken             = sdkerrors.Register(collectionCodespace, 40, "insufficient token")
	ErrDuplicateChangesField         = sdkerrors.Register(collectionCodespace, 41, "duplicate field of changes")
	ErrInvalidMetaLength             = sdkerrors.Register(collectionCodespace, 42, "invalid meta length")
	ErrSupplyOverflow                = sdkerrors.Register(collectionCodespace, 43, "supply for collection reached maximum")
	ErrEmptyField                    = sdkerrors.Register(collectionCodespace, 44, "required field cannot be empty")
	ErrCompositionTooDeep            = sdkerrors.Register(collectionCodespace, 45, "cannot attach token (composition too deep)")
	ErrCompositionTooWide            = sdkerrors.Register(collectionCodespace, 46, "cannot attach token (composition too wide)")
	ErrBurnNonRootNFT                = sdkerrors.Register(collectionCodespace, 47, "cannot burn non-root NFTs")
)
