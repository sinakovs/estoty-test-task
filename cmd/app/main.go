package main

import (
	"context"
	"database/sql"

	"estoty-test-task/internal/app"
	"estoty-test-task/internal/assets"
	"github.com/heroiclabs/nakama-common/runtime"
)

// InitModule is the Nakama Go runtime entrypoint for module registration.
func InitModule(
	_ context.Context,
	logger runtime.Logger,
	_ *sql.DB,
	nk runtime.NakamaModule,
	initializer runtime.Initializer,
) error {
	application, err := app.New(assets.GameConfig, nk, logger)
	if err != nil {
		logger.Error("failed to initialize app dependencies: %v", err)

		return err
	}

	err = initializer.RegisterRpc(
		"update_account_metadata",
		application.AccountRPC.UpdateAccountMetadataRPC,
	)
	if err != nil {
		logger.Error("failed to register update_account_metadata rpc: %v", err)

		return err
	}

	err = initializer.RegisterRpc("get_game_config", application.ConfigRPC.GetGameConfigRPC)
	if err != nil {
		logger.Error("failed to register get_game_config rpc: %v", err)

		return err
	}

	err = initializer.RegisterRpc("private_status", application.PrivateRPC.PrivateStatusRPC)
	if err != nil {
		logger.Error("failed to register private_status rpc: %v", err)

		return err
	}

	logger.Info("registered rpc handlers: update_account_metadata, get_game_config, private_status")

	return nil
}
