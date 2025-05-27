package controllers

import (
	"net/http"
	"time"
	"fmt"
	
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"os"
	"encoding/json"
	"bytes"
	"firebase.google.com/go/v4/auth"
	"painpal_api/middleware"
	"painpal_api/models"
)

type LoginRequest struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}

type LoginResponse struct {
	IDToken      string `json:"idToken"`      // ✅ Token que enviarás al frontend/backend
	RefreshToken string `json:"refreshToken"` // Opcional si implementas sesión prolongada
	ExpiresIn    string `json:"expiresIn"`
	LocalID      string `json:"localId"`      // UID del usuario
	Email        string `json:"email"`
}


var jwtSecret = []byte(os.Getenv("JWT_SECRET")) // ⚠️ Reemplaza por una variable de entorno

// POST /auth/register
func Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing failed"})
		return
	}

	authClient, fsClient, ctx, err := middleware.GetFirebaseClients()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Firebase init failed"})
		return
	}
	defer fsClient.Close()

	params := (&auth.UserToCreate{}).
		Email(req.Email).
		Password(string(hashedPassword)) // no se usa realmente en login

	createdUser, err := authClient.CreateUser(ctx, params)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists or invalid email"})
		return
	}

	_, err = fsClient.Collection("users").Doc(createdUser.UID).Set(ctx, map[string]interface{}{
		"email":     req.Email,
		"password":  string(hashedPassword), // ⚠️ este hash se usará en login
		"createdAt": time.Now(),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving user to Firestore"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// POST /auth/login
func Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	authClient, fsClient, ctx, err := middleware.GetFirebaseClients()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Firebase init failed"})
		return
	}
	defer fsClient.Close()
	userRecord, err := authClient.GetUserByEmail(ctx, req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	fmt.Println("User Record:", userRecord.UID)
	doc, err := fsClient.Collection("users").Doc(userRecord.UID).Get(ctx)
	if err != nil || !doc.Exists() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User data not found"})
		return
	}

	data := doc.Data()
	hashedPassword := data["password"].(string)

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": userRecord.UID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT signing error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": signedToken})
}

// LoginWithEmailAndPassword realiza login vía Firebase Auth REST API
func LoginWithEmailAndPassword(email, password string) (*LoginResponse, error) {

	 apiKey := os.Getenv("GOOGLE_API_KEY")
	//apiKey := "AIzaSyBpLWHBH9Gg5qG9ecwicIsYDj4vyO5A_Vg" // Reemplaza con tu API Key real
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s", apiKey)
	
	fmt.Println("API url:", url)
	payload := LoginRequest{
		Email:             email,
		Password:          password,
		ReturnSecureToken: true,
	}
	fmt.Println("API jsonBody:", payload)

	jsonBody, err := json.Marshal(payload)

	if err != nil {
		return nil, fmt.Errorf("error serializing login payload: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error calling Firebase Auth REST API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var msg map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&msg)
		return nil, fmt.Errorf("login failed (status %d): %v", resp.StatusCode, msg)
	}

	var result LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &result, nil
}