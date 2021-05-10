package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/friendsofgo/errors"
	"github.com/httptest/backend/pkg/errof"
	"github.com/httptest/backend/pkg/jsonrpc"
	"github.com/httptest/backend/pkg/util"
	"github.com/inconshreveable/log15"
)

type Header string

func (h Header) String() string {
	return string(h)
}

// Func :
type Func struct {
	Name        string
	Method      interface{}
	Permissions []string
}

const (
	OrgCode       Header = "OrgCode"
	CardSecret    Header = "CardSecret"
	DeviceCode    Header = "DeviceCode"
	Authorization Header = "Authorization"
	XForwardedFor Header = "X-Forwarded-For"
)

// Call :
func (f Func) Call(ctx context.Context, paramJSON []byte) (result interface{}, err error) {
	// func(context.Context, input.AddLotCount) error
	args := []reflect.Value{
		reflect.ValueOf(ctx),
	}
	validate := util.NewValidator()
	funcType := reflect.TypeOf(f.Method)
	if 1 < funcType.NumIn() {
		// 引数2つ目が input params
		inputType := funcType.In(1)
		params := reflect.New(inputType).Interface()
		if err := json.Unmarshal(paramJSON, params); err != nil {
			return nil, errors.Wrapf(errof.ErrParse, err.Error())
		}
		if inputType.Kind() != reflect.Slice {
			if err = validate.Struct(params); err != nil {
				return nil, errors.Wrapf(errof.ErrInvalidParams, "input :%+v, err: %s", inputType.String(), err)
			}
		}
		// indirectを使って、値を参照する
		args = append(args, reflect.Indirect(reflect.ValueOf(params)))
	}

	fv := reflect.ValueOf(f.Method)
	results := fv.Call(args)
	// 返り値は1つ or 2つ
	switch len(results) {
	case 1:
		errResult := results[0]
		if errResult.Interface() == nil {
			return nil, nil
		}
		return nil, errResult.Interface().(error)
	case 2:
		nomalResult, errResult := results[0], results[1]
		if errResult.Interface() == nil {
			return nomalResult.Interface(), nil
		}
		return nomalResult.Interface(), errResult.Interface().(error)
	default:
		return nil, fmt.Errorf("Invalid result, results: %+v", results)
	}
}

func writeCORSHeaders(w http.ResponseWriter, allowOrigin string) {
	header := w.Header()
	header.Set("Access-Control-Allow-Origin", allowOrigin)
	header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	header.Set("Access-Control-Allow-Headers", fmt.Sprintf("Content-Type, %s, %s", Authorization, OrgCode))
	header.Set("Access-Control-Allow-Credentials", "true")
	header.Set("Access-Control-Max-Age", "86400")
	header.Set("Content-Type", "application/json; charset=utf-8")
}

func writeSecurityHeaders(w http.ResponseWriter) {
	header := w.Header()
	header.Set("X-Frame-Options", "DENY")
	header.Set("X-Content-Type-Options", "nosniff")
	header.Set("X-XSS-Protection", "1")
	header.Set("Cache-Control", "no-store")
	header.Set("Pragma:", "no-cache")
}

func preflightCheck(w http.ResponseWriter, r *http.Request, needAuthorization bool) bool {
	if r.Method == http.MethodOptions {
		s := r.Header.Get("Access-Control-Request-Headers")
		if needAuthorization {
			if strings.Contains(s, "authorization") || strings.Contains(s, "Authorization") {
				w.WriteHeader(http.StatusNoContent)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
		return true
	}
	return false
}

func handleReturn(ctx context.Context, id, result interface{}, err error) (returns []*jsonrpc.Return) {
	if ctx.Err() == context.Canceled {
		return nil
	}

	r := &jsonrpc.Return{ID: id, Result: result}
	if err == nil {
		return []*jsonrpc.Return{r}
	}

	cause := errors.Cause(err)
	var skipLog bool
	for _, originErr := range []error{errof.ErrAuthentication} {
		if cause == originErr {
			log15.Warn(cause.Error(), "err", strings.Replace(fmt.Sprintf("%+v", err), "'", "*", -1))
			skipLog = true
		}
	}
	if !skipLog {
		log15.Error(cause.Error(), "err", strings.Replace(fmt.Sprintf("%+v", err), "'", "*", -1))
	}

	switch cause {
	case errof.ErrDatabase:
		if strings.Contains(err.Error(), "value too long") {
			err = errors.Wrap(errof.ErrTooLongParameter, err.Error())
			log15.Error(cause.Error(), "err", fmt.Sprintf("%+v", err))
			r.Error = jsonrpc.ErrTooLongParameter()
		}
	case errof.ErrParse:
		r.Error = jsonrpc.ErrParse()
	case errof.ErrInvalidRequest:
		r.Error = jsonrpc.ErrInvalidRequest()
	case errof.ErrMethodNotFound:
		r.Error = jsonrpc.ErrMethodNotFound()
	case errof.ErrInvalidParams:
		r.Error = jsonrpc.ErrInvalidParams()
	case errof.ErrInternal:
		r.Error = jsonrpc.ErrInternal()
	default:
		r.Error = jsonrpc.ErrServer(cause)
	}
	return []*jsonrpc.Return{r}
}
