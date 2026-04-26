package rpc

import (
	"context"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

const runtimeCtxUserID = runtime.RUNTIME_CTX_USER_ID

func userIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(runtimeCtxUserID).(string)
	if !ok || userID == "" {
		return "", runtime.NewError("user id is required", 16)
	}

	return userID, nil
}

func decodeMetadataPayload(payload string) (*updateAccountMetadataRequest, error) {
	if payload == "" {
		return nil, runtime.NewError("payload is required", 3)
	}

	var request updateAccountMetadataRequest
	if err := json.Unmarshal([]byte(payload), &request); err != nil {
		return nil, runtime.NewError("payload must be valid json", 3)
	}

	if len(request.Metadata) == 0 {
		return nil, runtime.NewError("metadata is required", 3)
	}

	return &request, nil
}
