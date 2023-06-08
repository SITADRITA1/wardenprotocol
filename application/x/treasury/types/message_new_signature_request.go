package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgNewSignatureRequest = "new_signature_request"

var _ sdk.Msg = &MsgNewSignatureRequest{}

func NewMsgNewSignatureRequest(creator string, wsID uint64, signType SignType, dataForSigning []byte) *MsgNewSignatureRequest {
	return &MsgNewSignatureRequest{
		Creator:        creator,
		WorkspaceId:    wsID,
		SignType:       signType,
		DataForSigning: dataForSigning,
	}
}

func (msg *MsgNewSignatureRequest) Route() string {
	return RouterKey
}

func (msg *MsgNewSignatureRequest) Type() string {
	return TypeMsgNewSignatureRequest
}

func (msg *MsgNewSignatureRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgNewSignatureRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgNewSignatureRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
