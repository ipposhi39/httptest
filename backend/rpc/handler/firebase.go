package handler

import (
	"context"
	"net/http"

	"github.com/httptest/backend/pkg/config"
	"github.com/httptest/backend/pkg/errof"
	"github.com/httptest/backend/pkg/jsonrpc"
	"github.com/httptest/backend/pkg/util"
	"github.com/httptest/backend/rpc/usecase"
	"github.com/inconshreveable/log15"
	"github.com/pkg/errors"
)

type firebaseHandler struct {
	allowOrigin string
	funcMap     map[string]Func
}

// NewFirebaseHandler :
func NewFirebaseHandler(
	c config.HTTP,
	successUsecase usecase.Success,
) http.Handler {
	return firebaseHandler{
		c.Cors,
		GetFirebaseFuncMap(
			successUsecase,
		),
	}
}

func (h firebaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	writeSecurityHeaders(w)
	writeCORSHeaders(w, h.allowOrigin)
	if preflightCheck(w, r, true) {
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// ctx := r.Context()
	ctx := context.Background()
	returnsCh := make(chan []*jsonrpc.Return, 1)
	go func() {
		sourceIps := r.Header.Values(XForwardedFor.String())
		if 0 < len(sourceIps) {
			ctx = util.SetIPAddress(ctx, sourceIps[len(sourceIps)-1])
		}

		requests, errs := jsonrpc.Parse(r)
		if errs != nil {
			returnsCh <- handleReturn(ctx, nil, nil, err)
			return
		}

		// rpcを呼ぶ時、空の配列だった場合の処理
		if len(requests) == 0 {
			returnsCh <- handleReturn(ctx, nil, nil, errors.Wrap(errof.ErrParse, "empty request"))
			return
		}

		authorization := r.Header.Get(Authorization.String())
		if authorization == "" {
			returnsCh <- handleReturn(ctx, nil, nil, errof.ErrInvalidRequest)
			return
		}

		var returns []*jsonrpc.Return
		for _, request := range requests {
			result, err := h.Exec(ctx, request.Method, request.Params)
			returns = append(returns, handleReturn(ctx, request.ID, result, err)...)
		}
		returnsCh <- returns
	}()

	select {
	case r := <-returnsCh:
		if err = jsonrpc.WriteResponses(w, r...); err != nil {
			log15.Crit("Failed to write success response ", "err", err, "request", r)
			return
		}
		return
		//	case <-ctx.Done():
		//		_ = jsonrpc.WriteResponses(w, handleReturn(ctx, nil, nil, jsonrpc.ErrInternal())...)
		//		return
	}
}

func (h *firebaseHandler) Exec(ctx context.Context, methodName string, params []byte) (result interface{}, err error) {
	defer func() {
		if p := recover(); p != nil {
			err = errors.Wrap(errof.ErrInternal, errof.PanicToErr(p).Error())
		}
		if err != nil {
			err = errors.Wrapf(err, "Method Front Failed methodName: %s, params: %s", methodName, string(params))
		}
		go func() {
			ctx = util.GetWithoutCancelContext(ctx)
		}()
	}()
	if f, ok := h.funcMap[methodName]; ok {
		return f.Call(ctx, params)
	}
	return nil, errors.WithStack(errof.ErrMethodNotFound)
}

// GetFirebaseFuncMap :
func GetFirebaseFuncMap(
	successUsecase usecase.Success,
) map[string]Func {
	return map[string]Func{
		"getSuccess": {Name: "成功", Method: successUsecase.GetSuccess},
	}
}
