package fsctx

import (
	"context"

	"cloud.google.com/go/firestore"
)

func client(ctx context.Context) *firestore.Client {
	if c, ok := ctx.Value(ctxKey{}).(*firestore.Client); ok {
		return c
	}
	panic("no firestore client in context")
}

func Collection(ctx context.Context, path string) *firestore.CollectionRef {
	return client(ctx).Collection(path)
}
func Collections(ctx context.Context) *firestore.CollectionIterator {
	return client(ctx).Collections(ctx)
}
