FROM heroiclabs/nakama-pluginbuilder:3.22.0 AS builder

ENV GO111MODULE=on
ENV CGO_ENABLED=1

WORKDIR /backend
COPY . .

RUN go mod download && \
    go build --trimpath --buildmode=plugin -o ./backend.so ./cmd/app

FROM registry.heroiclabs.com/heroiclabs/nakama:3.22.0

COPY --from=builder /backend/backend.so /nakama/data/modules/backend.so
COPY local.yml /nakama/data/local.yml
