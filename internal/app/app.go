package app

import (
	"estoty-test-task/internal/config"
	accountsvc "estoty-test-task/internal/service/account"
	gameconfigsvc "estoty-test-task/internal/service/gameconfig"
	privatesvc "estoty-test-task/internal/service/private"
	rpctransport "estoty-test-task/internal/transport/rpc"
	"github.com/heroiclabs/nakama-common/runtime"
)

// App wires transport handlers used by the Nakama runtime entrypoint.
type App struct {
	AccountRPC *rpctransport.AccountHandler
	ConfigRPC  *rpctransport.ConfigHandler
	PrivateRPC *rpctransport.PrivateHandler
}

// New builds the application dependency graph used by the runtime entrypoint.
func New(gameConfigData []byte, nk runtime.NakamaModule, logger runtime.Logger) (*App, error) {
	gameConfig, err := config.LoadGameConfig(gameConfigData)
	if err != nil {
		return nil, err
	}

	accountService := accountsvc.NewService(nk)
	configService := gameconfigsvc.NewService(gameConfig)
	privateService := privatesvc.NewService()

	return &App{
		AccountRPC: rpctransport.NewAccountHandler(accountService, logger),
		ConfigRPC:  rpctransport.NewConfigHandler(configService, logger),
		PrivateRPC: rpctransport.NewPrivateHandler(privateService, logger),
	}, nil
}
