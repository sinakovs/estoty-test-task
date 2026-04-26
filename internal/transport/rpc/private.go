package rpc

import (
	"context"
	"database/sql"

	privatesvc "estoty-test-task/internal/service/private"
	"github.com/heroiclabs/nakama-common/runtime"
)

type PrivateHandler struct {
	service *privatesvc.Service
	logger  runtime.Logger
}

func NewPrivateHandler(service *privatesvc.Service, logger runtime.Logger) *PrivateHandler {
	return &PrivateHandler{service: service, logger: logger}
}

func (h *PrivateHandler) PrivateStatusRPC(
	ctx context.Context,
	logger runtime.Logger,
	_ *sql.DB,
	_ runtime.NakamaModule,
	_ string,
) (string, error) {
	activeLogger := h.logger
	if activeLogger == nil {
		activeLogger = logger
	}

	activeLogger.Debug("rpc private_status called")

	response, err := h.service.Status(ctx)
	if err != nil {
		return "", err
	}

	activeLogger.Debug("rpc private_status completed successfully")

	return response, nil
}
