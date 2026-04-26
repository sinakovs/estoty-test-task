//go:build integration

package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestNakamaRPCEndToEnd(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Minute)
	defer cancel()

	baseURL := StartTestNakama(t, ctx)

	t.Run("get_game_config", func(t *testing.T) {
		raw := callRPCWithHTTPKey(t, baseURL, "get_game_config", "{}")

		var response struct {
			WelcomeMessage string   `json:"welcome_message"`
			XPRate         float64  `json:"xp_rate"`
			RarityOptions  []string `json:"rarity_options"`
		}
		mustUnmarshal(t, raw, &response)

		if response.WelcomeMessage == "" {
			t.Fatal("expected welcome_message to be present")
		}
		if response.XPRate != 1.25 {
			t.Fatalf("unexpected xp_rate: %v", response.XPRate)
		}
		if len(response.RarityOptions) != 4 {
			t.Fatalf("unexpected rarity options count: %d", len(response.RarityOptions))
		}
	})

	t.Run("private_status", func(t *testing.T) {
		raw := callRPCWithHTTPKey(t, baseURL, "private_status", "{}")
		if strings.TrimSpace(raw) != "{}" {
			t.Fatalf("unexpected private_status response: %s", raw)
		}
	})

	t.Run("update_account_metadata", func(t *testing.T) {
		sessionToken := authenticateDevice(t, baseURL)

		updatePayload := `{"metadata":{"favorite_color":"blue","xp":10,"rarity":"rare"}}`
		raw := callRPCWithSession(t, baseURL, "update_account_metadata", sessionToken, updatePayload)

		var response struct {
			Metadata map[string]any `json:"metadata"`
		}
		mustUnmarshal(t, raw, &response)

		if response.Metadata["favorite_color"] != "blue" {
			t.Fatalf("unexpected rpc response metadata: %#v", response.Metadata)
		}

		accountMetadata := fetchAccountMetadata(t, baseURL, sessionToken)
		if accountMetadata["favorite_color"] != "blue" {
			t.Fatalf("metadata was not persisted: %#v", accountMetadata)
		}
		if accountMetadata["rarity"] != "rare" {
			t.Fatalf("expected rarity metadata to persist: %#v", accountMetadata)
		}
	})

	t.Run("private_status_rejects_user_session", func(t *testing.T) {
		sessionToken := authenticateDevice(t, baseURL)

		statusCode, body := callRPCExpectingStatus(t, baseURL, "private_status?unwrap", sessionToken, "{}")
		if statusCode == http.StatusOK {
			t.Fatalf("expected private_status to reject authenticated session, body=%s", body)
		}
		if !strings.Contains(body, "server-to-server") {
			t.Fatalf("unexpected private_status rejection body: %s", body)
		}
	})
}

func authenticateDevice(t *testing.T, baseURL string) string {
	t.Helper()

	client := &http.Client{Timeout: 10 * time.Second}
	url := fmt.Sprintf("%s/v2/account/authenticate/device?create=true&username=nakama-test-user", baseURL)
	body := `{"id":"nakama-test-device-001"}`

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("build auth request failed: %v", err)
	}
	req.SetBasicAuth(testServerKey, "")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("authenticate device failed: %v", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read auth response failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("authenticate device returned status %d: %s", resp.StatusCode, string(raw))
	}

	var response struct {
		Token string `json:"token"`
	}
	mustUnmarshal(t, string(raw), &response)
	if response.Token == "" {
		t.Fatal("authenticate device returned empty token")
	}

	return response.Token
}

func fetchAccountMetadata(t *testing.T, baseURL string, sessionToken string) map[string]any {
	t.Helper()

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v2/account", baseURL), nil)
	if err != nil {
		t.Fatalf("build account request failed: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+sessionToken)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("get account failed: %v", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read account response failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("get account returned status %d: %s", resp.StatusCode, string(raw))
	}

	var response struct {
		User struct {
			Metadata any `json:"metadata"`
		} `json:"user"`
	}
	mustUnmarshal(t, string(raw), &response)

	switch metadata := response.User.Metadata.(type) {
	case map[string]any:
		return metadata
	case string:
		var decoded map[string]any
		mustUnmarshal(t, metadata, &decoded)
		return decoded
	default:
		t.Fatalf("unexpected metadata type: %T", response.User.Metadata)
		return nil
	}
}

func callRPCWithHTTPKey(t *testing.T, baseURL string, rpcID string, body string) string {
	t.Helper()

	client := &http.Client{Timeout: 10 * time.Second}
	url := fmt.Sprintf("%s/v2/rpc/%s?http_key=%s&unwrap", baseURL, rpcID, testHTTPKey)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("build rpc request failed: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("rpc %s failed: %v", rpcID, err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read rpc response failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("rpc %s returned status %d: %s", rpcID, resp.StatusCode, string(raw))
	}

	return string(raw)
}

func callRPCWithSession(t *testing.T, baseURL string, rpcID string, sessionToken string, body string) string {
	t.Helper()

	statusCode, raw := callRPCExpectingStatus(t, baseURL, rpcID+"?unwrap", sessionToken, body)
	if statusCode != http.StatusOK {
		t.Fatalf("rpc %s returned status %d: %s", rpcID, statusCode, raw)
	}

	return raw
}

func callRPCExpectingStatus(t *testing.T, baseURL string, rpcPath string, sessionToken string, body string) (int, string) {
	t.Helper()

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v2/rpc/%s", baseURL, rpcPath), bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("build rpc request failed: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+sessionToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("rpc %s failed: %v", rpcPath, err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read rpc response failed: %v", err)
	}

	return resp.StatusCode, string(raw)
}

func mustUnmarshal(t *testing.T, raw string, target any) {
	t.Helper()

	if err := json.Unmarshal([]byte(raw), target); err != nil {
		t.Fatalf("json unmarshal failed: %v\npayload=%s", err, raw)
	}
}
