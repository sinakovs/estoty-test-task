package rpc

import (
	"context"
	"database/sql"

	gameconfigsvc "estoty-test-task/internal/service/gameconfig"
	"github.com/heroiclabs/nakama-common/runtime"
)

type ConfigHandler struct {
	service *gameconfigsvc.Service
	logger  runtime.Logger
}

func NewConfigHandler(service *gameconfigsvc.Service, logger runtime.Logger) *ConfigHandler {
	return &ConfigHandler{service: service, logger: logger}
}

func (h *ConfigHandler) GetGameConfigRPC(
	_ context.Context,
	logger runtime.Logger,
	_ *sql.DB,
	_ runtime.NakamaModule,
	_ string,
) (string, error) {
	activeLogger := h.logger
	if activeLogger == nil {
		activeLogger = logger
	}

	activeLogger.Debug("rpc get_game_config called")

	response, err := h.service.Get()
	if err != nil {
		return "", err
	}

	activeLogger.Debug("rpc get_game_config completed successfully")

	return response, nil
}
