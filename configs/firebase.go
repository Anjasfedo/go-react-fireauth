package configs

import (
	"context"
	"log"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"

	"cloud.google.com/go/firestore"
)

var (
	AuthClient        *auth.Client
	FirestoreClient   *firestore.Client
	StorageClient     *storage.Client
	StorageBucket     *storage.BucketHandle
	StorageBucketName = "friendlyeats-codelab-2663a.appspot.com"
)

func InitFirebase(ctx context.Context) {
	opt := option.WithCredentialsFile("serviceAccountKey.json")

	config := &firebase.Config{
		StorageBucket: StorageBucketName,
	}

	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	AuthClient, err = app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Fireauth client: %v\n", err)
	}

	FirestoreClient, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Error getting Firestore client: %v\n", err)
	}

	StorageClient, err = storage.NewClient(ctx, opt)
	if err != nil {
		log.Fatalf("Error getting Storage client: %v\n", err)
	}

	StorageBucket = StorageClient.Bucket(StorageBucketName)
}

func CloseFirebase() {
	if FirestoreClient != nil {
		if err := FirestoreClient.Close(); err != nil {
			log.Fatalf("Failed to close Firestore client: %v", err)
		}
	}

    if StorageClient != nil {
        if err := StorageClient.Close(); err != nil {
            log.Fatalf("Failed to close Storage client: %v", err)
        }
    }
}
