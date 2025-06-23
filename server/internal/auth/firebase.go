package auth

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// InitFirebaseAuthClient はFirebase Admin SDKを初期化し、
// Firebase Authenticationクライアントを返します。
// この関数はアプリケーションの起動時に一度だけ呼び出します
func InitFirebaseAuthClient() *auth.Client {
	serviceAccountKeyPath := os.Getenv("FIREBASE_SERVICE_ACCOUNT_KEY_PATH")
	if serviceAccountKeyPath == "" {
		log.Fatalf("FATAL: FIREBASE_SERVICE_ACCOUNT_KEY_PATH environment variable is not set. Please set the path to the Firebase service account key JSON file.")
	}

	opt := option.WithCredentialsFile(serviceAccountKeyPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("FATAL: Failed to initialize Firebase App: %v", err)
	}

	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("FATAL: Failed to obtain Firebase Authentication client: %v", err)
	}

	log.Println("Firebase Admin SDK initialized successfully.")
	return authClient
}
