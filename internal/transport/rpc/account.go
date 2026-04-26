package rpc

import (
	"context"
	"database/sql"
	"encoding/json"

	accountsvc "estoty-test-task/internal/service/account"
	"github.com/heroiclabs/nakama-common/runtime"
)

type AccountHandler struct {
	service *accountsvc.Service
	logger  runtime.Logger
}

type updateAccountMetadataRequest struct {
	Metadata map[string]any `json:"metadata"`
}

type updateAccountMetadataResponse struct {
	Metadata map[string]any `json:"metadata"`
}

func NewAccountHandler(service *accountsvc.Service, logger runtime.Logger) *AccountHandler {
	return &AccountHandler{service: service, logger: logger}
}

func (h *AccountHandler) UpdateAccountMetadataRPC(
	ctx context.Context,
	_ runtime.Logger,
	_ *sql.DB,
	_ runtime.NakamaModule,
	payload string,
) (string, error) {
	h.logger.Debug("rpc update_account_metadata called")

	userID, err := userIDFromContext(ctx)
	if err != nil {
		return "", err
	}

	request, err := decodeMetadataPayload(payload)
	if err != nil {
		return "", err
	}

	h.logger.Debug(
		"rpc update_account_metadata validated payload for user %s with %d metadata fields",
		userID,
		len(request.Metadata),
	)

	mergedMetadata, err := h.service.UpdateMetadata(ctx, userID, request.Metadata)
	if err != nil {
		h.logger.Error("failed to update account metadata for user %s: %v", userID, err)
		return "", err
	}

	response, err := json.Marshal(updateAccountMetadataResponse{Metadata: mergedMetadata})
	if err != nil {
		return "", runtime.NewError("could not encode metadata response", 13)
	}

	h.logger.Debug("rpc update_account_metadata completed successfully for user %s", userID)

	return string(response), nil
}
