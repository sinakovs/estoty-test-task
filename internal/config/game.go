package config

import (
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

type GameConfig struct {
	WelcomeMessage string   `json:"welcome_message"`
	XPRate         float64  `json:"xp_rate"`
	RarityOptions  []string `json:"rarity_options"`
}

func LoadGameConfig(data []byte) (*GameConfig, error) {
	if len(data) == 0 {
		return nil, runtime.NewError("game config is empty", 13)
	}

	var cfg GameConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, runtime.NewError("game config is invalid", 13)
	}

	return &cfg, nil
}
