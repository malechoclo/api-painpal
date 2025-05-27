package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	
)

var firebaseAuth *auth.Client

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func InitFirestore() error {
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: os.Getenv("FIREBASE_PROJECT_ID"), // Reemplaza por tu ID real
	}, option.WithCredentialsFile(os.Getenv("GOOGLE_CREDENTIALS")))
	if err != nil {
		return fmt.Errorf("error initializing Firebase App: %v", err)
	}

	firebaseAuth, err = app.Auth(ctx)
	if err != nil {
		return fmt.Errorf("error initializing Firebase Auth: %v", err)
	}

	return nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			c.Abort()
			return
		}

		idToken := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := firebaseAuth.VerifyIDToken(c, idToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("uid", token.UID)
		c.Next()
	}
}

func GetFirebaseClients() (*auth.Client, *firestore.Client, context.Context, error) {
	ctx := context.Background()
	credentialPath := os.Getenv("GOOGLE_CREDENTIALS")

	fmt.Println("Credential Path:", os.Getenv("FIREBASE_PROJECT_ID"))

	opt := option.WithCredentialsFile(credentialPath)
	app, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: os.Getenv("FIREBASE_PROJECT_ID"),
	}, opt)
	if err != nil {
		return nil, nil, nil, err
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	fsClient, err := app.Firestore(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	return authClient, fsClient, ctx, nil
}