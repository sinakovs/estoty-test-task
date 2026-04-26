package private

import (
	"context"

	"github.com/heroiclabs/nakama-common/runtime"
)

const (
	runtimeCtxUserID = runtime.RUNTIME_CTX_USER_ID
	emptyJSONObject  = "{}"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (*Service) Status(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(runtimeCtxUserID).(string)
	if ok && userID != "" {
		return "", runtime.NewError("rpc is only callable via server-to-server http key", 7)
	}

	return emptyJSONObject, nil
}
