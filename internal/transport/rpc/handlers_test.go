package rpc

import (
	"context"
	"testing"

	"estoty-test-task/internal/config"
	accountsvc "estoty-test-task/internal/service/account"
	gameconfigsvc "estoty-test-task/internal/service/gameconfig"
	privatesvc "estoty-test-task/internal/service/private"
	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
)

type fakeAccountModule struct {
	account *api.Account
}

type fakeLogger struct{}

func TestAccountHandlerRejectsMissingMetadata(t *testing.T) {
	t.Parallel()

	handler := NewAccountHandler(accountsvc.NewService(&fakeAccountModule{}), fakeLogger{})
	ctx := context.WithValue(context.Background(), runtime.RUNTIME_CTX_USER_ID, "4ec4f126-3f9d-11e7-84ef-b7c182b36521")

	_, err := handler.UpdateAccountMetadataRPC(ctx, fakeLogger{}, nil, nil, `{"metadata":{}}`)
	if err == nil {
		t.Fatal("expected error for empty metadata payload")
	}
}

func TestConfigHandlerReturnsRawJSON(t *testing.T) {
	t.Parallel()

	handler := NewConfigHandler(gameconfigsvc.NewService(&config.GameConfig{
		WelcomeMessage: "hello",
		XPRate:         1.25,
		RarityOptions:  []string{"common", "rare"},
	}), fakeLogger{})

	raw, err := handler.GetGameConfigRPC(context.Background(), fakeLogger{}, nil, nil, "")
	if err != nil {
		t.Fatalf("GetGameConfigRPC returned error: %v", err)
	}

	expected := `{"welcome_message":"hello","xp_rate":1.25,"rarity_options":["common","rare"]}`
	if raw != expected {
		t.Fatalf("unexpected config payload: %s", raw)
	}
}

func TestPrivateHandlerRejectsClientCalls(t *testing.T) {
	t.Parallel()

	handler := NewPrivateHandler(privatesvc.NewService(), fakeLogger{})
	ctx := context.WithValue(context.Background(), runtime.RUNTIME_CTX_USER_ID, "4ec4f126-3f9d-11e7-84ef-b7c182b36521")
	if _, err := handler.PrivateStatusRPC(ctx, fakeLogger{}, nil, nil, ""); err == nil {
		t.Fatal("expected PrivateStatusRPC to reject authenticated client calls")
	}
}

func TestPrivateHandlerAcceptsServerToServerCalls(t *testing.T) {
	t.Parallel()

	handler := NewPrivateHandler(privatesvc.NewService(), fakeLogger{})

	response, err := handler.PrivateStatusRPC(context.Background(), fakeLogger{}, nil, nil, "")
	if err != nil {
		t.Fatalf("PrivateStatusRPC returned error: %v", err)
	}

	if response != "{}" {
		t.Fatalf("unexpected private status response: %s", response)
	}
}

func (f *fakeAccountModule) AccountGetId(ctx context.Context, userID string) (*api.Account, error) {
	if f.account == nil {
		return &api.Account{User: &api.User{}}, nil
	}

	return f.account, nil
}

func (f *fakeAccountModule) AccountUpdateId(
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
	return nil
}

func (fakeLogger) Debug(format string, v ...any)                   {}
func (fakeLogger) Info(format string, v ...any)                    {}
func (fakeLogger) Warn(format string, v ...any)                    {}
func (fakeLogger) Error(format string, v ...any)                   {}
func (fakeLogger) WithField(key string, v any) runtime.Logger      { return nil }
func (fakeLogger) WithFields(fields map[string]any) runtime.Logger { return nil }
func (fakeLogger) Fields() map[string]any                          { return nil }
