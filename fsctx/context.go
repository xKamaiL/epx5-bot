package fsctx

import (
	"context"

	"cloud.google.com/go/firestore"
)

func NewContext(parent context.Context, fs *firestore.Client) context.Context {
	return context.WithValue(parent, ctxKey{}, fs)
}

type ctxKey struct {
}
