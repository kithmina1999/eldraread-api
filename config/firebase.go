package config

import (
	"context"

	"firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

var firebaseApp *firebase.App

func FirebaseAuth() (*auth.Client, error) {
	if firebaseApp == nil {
		opt := option.WithCredentialsFile("eldraread-firebase-adminsdk-fbsvc-453019b735.json")
		app, err := firebase.NewApp(context.Background(), nil, opt)
		if err != nil {
			return nil, err
		}
		firebaseApp = app
	}
	return firebaseApp.Auth(context.Background())
}
