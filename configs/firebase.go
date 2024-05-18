package configs

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"

	"google.golang.org/api/option"

	"cloud.google.com/go/firestore"
)

var AuthClient *auth.Client

var FirestoreClient *firestore.Client

func InitFirebase() {
	opt := option.WithCredentialsFile("serviceAccountKey.json")

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	AuthClient, err = app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	FirestoreClient, err = app.Firestore(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
}
