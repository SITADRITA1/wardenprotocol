package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"gitlab.qredo.com/qrdochain/fusionchain/x/treasury/types"
)

func (k msgServer) FulfillSignatureRequest(goCtx context.Context, msg *types.MsgFulfillSignatureRequest) (*types.MsgFulfillSignatureRequestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// if !isAllowedToCreateSignatures(msg.Creator) {
	// 	return nil, fmt.Errorf("only MPC can sign data")
	// }

	req, found := k.GetSignRequest(ctx, msg.RequestId)
	if !found {
		return nil, fmt.Errorf("request not found")
	}

	if req.Status != types.SignRequestStatus_SIGN_REQUEST_STATUS_PENDING {
		return nil, fmt.Errorf("request is not pending, can't be updated")
	}

	switch msg.Status {
	case types.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED:
		signed := types.SignedPayload{
			WorkspaceId: req.WorkspaceId,
			Type:        req.SignType,
			SignedData:  (msg.Result.(*types.MsgFulfillSignatureRequest_Payload)).Payload.SignedData,
		}
		sigID := k.AppendSignedPayload(ctx, signed)

		// update WalletRequest with newly created wallet id
		req.Status = types.SignRequestStatus_SIGN_REQUEST_STATUS_FULFILLED
		req.Result = &types.SignRequest_SignedPayloadId{
			SignedPayloadId: sigID,
		}
		k.SetSignRequest(ctx, req)

		return &types.MsgFulfillSignatureRequestResponse{}, nil

	case types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED:
		req.Status = types.SignRequestStatus_SIGN_REQUEST_STATUS_REJECTED
		req.Result = &types.SignRequest_RejectReason{
			RejectReason: msg.Result.(*types.MsgFulfillSignatureRequest_RejectReason).RejectReason,
		}
		k.SetSignRequest(ctx, req)

	default:
		return nil, fmt.Errorf("invalid status field, should be one of approved/rejected")
	}

	return &types.MsgFulfillSignatureRequestResponse{}, nil
}
