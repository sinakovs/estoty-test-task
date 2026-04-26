package account

import (
	"context"
	"errors"
	"testing"

	"github.com/heroiclabs/nakama-common/api"
)

type fakeModule struct {
	account         *api.Account
	getErr          error
	updateErr       error
	updatedUserID   string
	updatedMetadata map[string]any
}

func TestUpdateMetadataMergesExistingAndIncomingMetadata(t *testing.T) {
	t.Parallel()

	nk := &fakeModule{
		account: &api.Account{
			User: &api.User{
				Metadata: `{"level": 5, "title": "rookie"}`,
			},
		},
	}

	service := NewService(nk)

	updated, err := service.UpdateMetadata(context.Background(), "user-id", map[string]any{
		"xp":     125.5,
		"title":  "veteran",
		"rarity": "epic",
	})
	if err != nil {
		t.Fatalf("UpdateMetadata returned error: %v", err)
	}

	if nk.updatedUserID != "user-id" {
		t.Fatalf("unexpected updated user id: %s", nk.updatedUserID)
	}

	if updated["title"] != "veteran" {
		t.Fatalf("expected title to be overwritten, got %#v", updated["title"])
	}

	if updated["level"] != float64(5) {
		t.Fatalf("expected existing metadata to be preserved, got %#v", updated["level"])
	}
}

func (f *fakeModule) AccountGetId(ctx context.Context, userID string) (*api.Account, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}

	if f.account == nil {
		return nil, errors.New("account not configured")
	}

	return f.account, nil
}

func (f *fakeModule) AccountUpdateId(
	ctx context.Context,
	userID string,
	username string,
	metadata map[string]any,
	displayName string,
	timezone string,
	location string,
	langTag string,
	avatarURL string,
) error {
	f.updatedUserID = userID
	f.updatedMetadata = metadata

	return f.updateErr
}
