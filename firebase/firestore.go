package firebase

import (
    "context"
    "log"
    "os"
    "fmt"

    "cloud.google.com/go/firestore"
    "google.golang.org/api/option"
)


var (
  // Cliente global, asegúrate de no re-inicializarlo en cada request
  FirestoreClient *firestore.Client
)
// func InitFirestore() {
 func init() {
    googleCredentials := os.Getenv("GOOGLE_CREDENTIALS")
    firebaseProjectID := os.Getenv("FIREBASE_PROJECT_ID")

    fmt.Println("GOOGLE_CREDENTIALS:", googleCredentials)
    fmt.Println("FIREBASE_PROJECT_ID:", firebaseProjectID)
    ctx := context.Background()
    sa := option.WithCredentialsFile(os.Getenv("GOOGLE_CREDENTIALS"))
    app, err := firestore.NewClient(ctx, os.Getenv("FIREBASE_PROJECT_ID"), sa)
    if err != nil {
        log.Fatalf("Failed to create Firestore client: %v", err)
    }
    FirestoreClient = app
}
