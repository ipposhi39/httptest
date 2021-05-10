package util

import (
	"context"
	"database/sql"
)

// https://deeeet.com/writing/2017/02/23/go-context-value/
type contextKey string

var (
	orgIDContextKey     contextKey = "orgID"
	userIDContextKey    contextKey = "userID"
	deviceIDContextKey  contextKey = "deviceID"
	ipAddressContextKey contextKey = "ipAddress"
	dbTxContextKey      contextKey = "dbTx"
)

type withoutCancel struct {
	context.Context
}

// GetWithoutCancelContext :
func GetWithoutCancelContext(ctx context.Context) context.Context {
	// dbのtransactionは外す
	ctx = context.WithValue(ctx, dbTxContextKey, nil)
	return withoutCancel{ctx}
}

// SetIPAddress :
func SetIPAddress(ctx context.Context, ipAddress string) context.Context {
	return context.WithValue(ctx, ipAddressContextKey, ipAddress)
}

// GetIPAddress :
func GetIPAddress(ctx context.Context) string {
	v := ctx.Value(ipAddressContextKey)
	ipAddress, ok := v.(string)
	if !ok {
		return ""
	}
	return ipAddress
}

// SetDBTx :
func SetDBTx(ctx context.Context, dbTx *sql.Tx) context.Context {
	return context.WithValue(ctx, dbTxContextKey, dbTx)
}

// GetDBTx :
func GetDBTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(dbTxContextKey).(*sql.Tx); ok {
		return tx
	}
	return nil
}
