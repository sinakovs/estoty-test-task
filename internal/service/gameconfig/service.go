package gameconfig

import (
	"encoding/json"

	"estoty-test-task/internal/config"
	"github.com/heroiclabs/nakama-common/runtime"
)

type Service struct {
	config *config.GameConfig
}

func NewService(cfg *config.GameConfig) *Service {
	return &Service{config: cfg}
}

func (s *Service) Get() (string, error) {
	payload, err := json.Marshal(s.config)
	if err != nil {
		return "", runtime.NewError("could not encode game config", 13)
	}

	return string(payload), nil
}
