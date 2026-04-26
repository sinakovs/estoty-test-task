//go:build integration

package rpc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const (
	testProjectName      = "nakama-test"
	testHTTPPort         = 17350
	testConsolePort      = 17351
	testGRPCPort         = 17349
	testPostgresHostPort = 15432
	testServerKey        = "defaultkey"
	testHTTPKey          = "defaulthttpkey"
	testConsolePassword  = "password"
)

func StartTestNakama(t *testing.T, ctx context.Context) string {
	t.Helper()

	repoRoot := mustRepoRoot(t)
	envFile := writeTestEnvFile(t)

	runCompose(t, ctx, repoRoot, envFile, "up", "--build", "-d")
	t.Cleanup(func() {
		if t.Failed() {
			dumpComposeLogs(t, repoRoot, envFile)
		}
		runCompose(t, context.Background(), repoRoot, envFile, "down", "-v")
	})

	baseURL := fmt.Sprintf("http://127.0.0.1:%d", testHTTPPort)
	waitForRPCReady(t, ctx, baseURL)
	return baseURL
}

func mustRepoRoot(t *testing.T) string {
	t.Helper()

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd failed: %v", err)
	}

	return filepath.Clean(filepath.Join(cwd, "..", "..", ".."))
}

func writeTestEnvFile(t *testing.T) string {
	t.Helper()

	content := strings.Join([]string{
		"POSTGRES_IMAGE=postgres:12.2-alpine",
		"POSTGRES_DB=nakama",
		"POSTGRES_PASSWORD=localdb",
		"POSTGRES_PORT=5432",
		fmt.Sprintf("POSTGRES_HOST_PORT=%d", testPostgresHostPort),
		"NAKAMA_IMAGE=registry.heroiclabs.com/heroiclabs/nakama:3.22.0",
		"NAKAMA_SOCKET_SERVER_KEY=defaultkey",
		"NAKAMA_RUNTIME_HTTP_KEY=defaulthttpkey",
		"NAKAMA_SESSION_ENCRYPTION_KEY=defaultencryptionkey",
		"NAKAMA_SESSION_REFRESH_ENCRYPTION_KEY=defaultrefreshencryptionkey",
		fmt.Sprintf("NAKAMA_CONSOLE_PASSWORD=%s", testConsolePassword),
		"NAKAMA_CONSOLE_SIGNING_KEY=supersecretconsolekey",
		"NAKAMA_GRPC_PORT=7349",
		fmt.Sprintf("NAKAMA_GRPC_HOST_PORT=%d", testGRPCPort),
		"NAKAMA_HTTP_PORT=7350",
		fmt.Sprintf("NAKAMA_HTTP_HOST_PORT=%d", testHTTPPort),
		"NAKAMA_CONSOLE_PORT=7351",
		fmt.Sprintf("NAKAMA_CONSOLE_HOST_PORT=%d", testConsolePort),
		"",
	}, "\n")

	path := filepath.Join(t.TempDir(), ".env.e2e")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write env file failed: %v", err)
	}

	return path
}

func runCompose(t *testing.T, ctx context.Context, repoRoot string, envFile string, args ...string) {
	t.Helper()

	cmdArgs := []string{"compose", "--env-file", envFile, "-f", "docker-compose-postgres.yml", "-p", testProjectName}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.CommandContext(ctx, "docker", cmdArgs...)
	cmd.Dir = repoRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("docker %s failed: %v\n%s", strings.Join(cmdArgs, " "), err, string(output))
	}
}

func dumpComposeLogs(t *testing.T, repoRoot string, envFile string) {
	t.Helper()

	cmd := exec.Command("docker", "compose", "--env-file", envFile, "-f", "docker-compose-postgres.yml", "-p", testProjectName, "logs", "--tail=100")
	cmd.Dir = repoRoot
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Logf("docker compose logs:\n%s", string(output))
	}
}

func waitForRPCReady(t *testing.T, ctx context.Context, baseURL string) {
	t.Helper()

	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("%s/v2/rpc/private_status?http_key=%s&unwrap", baseURL, testHTTPKey)

	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBufferString("{}"))
		if err != nil {
			t.Fatalf("build readiness request failed: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err == nil {
			body, _ := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK && strings.TrimSpace(string(body)) == "{}" {
				return
			}
		}

		select {
		case <-ctx.Done():
			t.Fatalf("nakama test instance did not become ready: %v", ctx.Err())
		case <-time.After(2 * time.Second):
		}
	}
}
