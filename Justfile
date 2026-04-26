# On Windows, `just` defaults to `sh`, so force PowerShell explicitly.
set shell := ["powershell.exe", "-NoLogo", "-Command"]

compose_file := "docker-compose-postgres.yml"
nakama_container := "nakama"

nakama-up:
  docker compose -f {{compose_file}} -p {{nakama_container}} up --build -d

nakama-down:
  docker compose -f {{compose_file}} -p {{nakama_container}} down

nakama-logs:
  docker compose -f {{compose_file}} -p {{nakama_container}} logs -f nakama

nakama-ps:
  docker compose -f {{compose_file}} -p {{nakama_container}} ps

lint:
  just lint-gofmt
  golangci-lint run ./...

lint-gofmt:
  $files = @(go list -f '{{ "{{.Dir}}" }}' ./... | ForEach-Object { Get-ChildItem -Path $_ -Filter *.go -File -Recurse | ForEach-Object { $_.FullName } }); $unformatted = @($files | ForEach-Object { gofmt -l $_ } | Where-Object { $_ }); if ($unformatted.Count -gt 0) { $unformatted; throw "gofmt found unformatted files" }

lint-gofmt-fix:
  $files = @(go list -f '{{ "{{.Dir}}" }}' ./... | ForEach-Object { Get-ChildItem -Path $_ -Filter *.go -File -Recurse | ForEach-Object { $_.FullName } }); if ($files.Count -gt 0) { gofmt -w $files }

test:
  go test -count=1 ./...
  go test -count=1 -tags=integration ./test/integration/...

test-e2e:
  go test -tags=integration -v ./test/integration/...
