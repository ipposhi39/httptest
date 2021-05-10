// +build wireinject

package injector

import (
	"net/http"

	"github.com/google/wire"
	"github.com/httptest/backend/pkg/config"
	"github.com/httptest/backend/pkg/db"
	"github.com/httptest/backend/rpc/handler"
)

// FirebaseFuncMap :
var FirebaseFuncMap = wire.NewSet(
	db.NewPSQL,
	db.NewDB,
)

// InitializeFirebaseMap :
func InitializeFirebaseMap(config.Postgres, config.Firebase, string) (_ map[string]handler.Func) {
	wire.Build(
		handler.GetFirebaseFuncMap,
		FirebaseFuncMap,
	)
	return
}

// InitializeFirebaseHandler :
func InitializeFirebaseHandler(config.HTTP, config.Postgres, config.Firebase, string) (_ http.Handler) {
	wire.Build(
		handler.NewFirebaseHandler,
		FirebaseFuncMap,
	)
	return
}
