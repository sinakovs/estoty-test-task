package account

import (
	"context"
	"encoding/json"
	"maps"

	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
)

type Module interface {
	AccountGetId(ctx context.Context, userID string) (*api.Account, error)
	AccountUpdateId(
		ctx context.Context,
		userID string,
		username string,
		metadata map[string]any,
		displayName string,
		timezone string,
		location string,
		langTag string,
		avatarURL string,
	) error
}

type Service struct {
	nk Module
}

func NewService(nk Module) *Service {
	return &Service{nk: nk}
}

func (s *Service) UpdateMetadata(ctx context.Context, userID string, incoming map[string]any) (map[string]any, error) {
	account, err := s.nk.AccountGetId(ctx, userID)
	if err != nil {
		return nil, runtime.NewError("could not fetch account", 13)
	}

	mergedMetadata := mergeMetadata(account.GetUser().GetMetadata(), incoming)
	err = s.nk.AccountUpdateId(ctx, userID, "", mergedMetadata, "", "", "", "", "")
	if err != nil {
		return nil, runtime.NewError("could not update account metadata", 13)
	}

	return mergedMetadata, nil
}

func mergeMetadata(existingJSON string, incoming map[string]any) map[string]any {
	merged := make(map[string]any)

	if existingJSON != "" {
		var existing map[string]any
		err := json.Unmarshal([]byte(existingJSON), &existing)
		if err == nil && len(existing) > 0 {
			maps.Copy(merged, existing)
		}
	}

	maps.Copy(merged, incoming)

	return merged
}
