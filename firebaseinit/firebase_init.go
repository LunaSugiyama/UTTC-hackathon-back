package firebaseinit

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

var AuthClient *auth.Client

func InitFirebase() error {
	opt := option.WithCredentialsFile("/home/denjo/ダウンロード/term4-luna-sugiyama-firebase-adminsdk-1joai-b0f371c4d8.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return err
	}

	AuthClient, err = app.Auth(context.Background())
	if err != nil {
		return err
	}

	return nil
}
