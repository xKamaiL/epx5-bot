package user

import (
	"context"
	"errors"

	"github.com/diamondburned/arikawa/v3/discord"
	"google.golang.org/api/iterator"

	"github.com/xkamail/epx5-bot/fsctx"
)

// FsPath represents firestore path
const FsPath = "users"

type User struct {
	// RefID of document
	RefID string `json:"refId" firestore:"-"`
	// ID of discord.UserID
	ID            string `firestore:"id"`
	Username      string `firestore:"username"`
	Avatar        string `firestore:"avatar"`
	Exp           int64  `firestore:"exp"`
	Verified      bool   `firestore:"verified"`
	Discriminator string `firestore:"discriminator"`
	Coins         int64  `firestore:"coins"`
}

func Create(ctx context.Context, id discord.UserID, username, avatar string, discriminator string) (string, error) {
	result, _, err := fsctx.Collection(ctx, FsPath).Add(ctx, User{
		ID:            id.String(),
		Username:      username,
		Avatar:        avatar,
		Exp:           0,
		Verified:      false,
		Discriminator: discriminator,
		Coins:         0,
	})
	if err != nil {
		return "", err
	}
	return result.ID, nil
}

func Find(ctx context.Context, id discord.UserID) (*User, error) {
	iter := fsctx.Collection(ctx, FsPath).Where("id", "==", id.String()).Documents(ctx)
	doc, err := iter.Next()
	defer iter.Stop()

	if errors.Is(err, iterator.Done) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	var u User
	if err := doc.DataTo(&u); err != nil {
		return nil, err
	}
	u.ID = doc.Ref.ID
	return &u, nil
}
