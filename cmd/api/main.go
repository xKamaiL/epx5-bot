package main

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"github.com/diamondburned/arikawa/v3/discord"

	"github.com/xkamail/epx5-bot/fsctx"
	"github.com/xkamail/epx5-bot/user"
)

func main() {
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: projectID}

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()
	ctx = fsctx.NewContext(ctx, client)

	snowflake, err := discord.ParseSnowflake("395561779368951811")
	if err != nil {
		return
	}
	id := discord.UserID(snowflake)

	//if _, err := user.Create(ctx, id, "xkamail", "as", "4882"); err != nil {
	//	log.Fatalln(err)
	//
	//}
	acc, err := user.Find(ctx, id)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(acc)
	select {}
}
