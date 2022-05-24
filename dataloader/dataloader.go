package dataloader

import (
	"context"
	"github.com/angelorc/sinfonia-indexer/generated/dataloaden"
	"net/http"
)

// ContextKey ...
const ContextKey = "DATALOADER"

// Loaders ...
type Loaders struct {
	TransactionLoader *dataloaden.TransactionLoader
}

// DataLoaderMiddleware ...
func DataLoaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loaders := Loaders{}
		loaders.TransactionLoader = getTransactionLoader()

		ctx := context.WithValue(r.Context(), ContextKey, loaders)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetContextLoaders ...
func GetContextLoaders(ctx context.Context) Loaders {
	return ctx.Value(ContextKey).(Loaders)
}
