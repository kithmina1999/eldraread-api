package config

import (
	"context"
	"sync"

	"cloud.google.com/go/storage"
	"firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

var (
	firebaseApp     *firebase.App
	initFirebaseApp sync.Once
)

// Initializes and returns Firebase Auth client
func FirebaseAuth() (*auth.Client, error) {
	if err := initializeFirebaseApp(); err != nil {
		return nil, err
	}
	return firebaseApp.Auth(context.Background())
}

func FirebaseStorage() (*storage.BucketHandle, error) {
	if err := initializeFirebaseApp(); err != nil {
		return nil, err
	}
	storageClient, err := firebaseApp.Storage(context.Background())
	if err != nil {
		return nil, err
	}

	bucketName := "eldraread.firebasestorage.app" // ⚠️ Replace with your actual bucket name
	bucket, err := storageClient.Bucket(bucketName)
	if err != nil {
		return nil, err
	}

	return bucket, nil
}

// Internal shared initializer
func initializeFirebaseApp() error {
	var err error
	initFirebaseApp.Do(func() {
		opt := option.WithCredentialsFile("eldraread-firebase-adminsdk-fbsvc-453019b735.json")
		firebaseApp, err = firebase.NewApp(context.Background(), nil, opt)
	})
	return err
}
