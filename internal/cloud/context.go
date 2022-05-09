package cloud

import (
	"context"

	"firebase.google.com/go/storage"
)

type ctxKey struct {
}

func NewContext(parent context.Context, storageClient *storage.Client) context.Context {
	return context.WithValue(parent, ctxKey{}, storageClient)
}
func client(ctx context.Context) *storage.Client {
	return ctx.Value(ctxKey{}).(*storage.Client)
}
