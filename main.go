package main

import (
    "log"
    "os"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "painpal_api/routes"
    "painpal_api/middleware"
)

var jwtSecret []byte

func main() {
    // Cargar variables desde .env
    if err := godotenv.Load(); err != nil {
        log.Fatalf("Error cargando archivo .env: %v", err)
    }

    jwtSecret = []byte(os.Getenv("JWT_SECRET"))
    if len(jwtSecret) == 0 {
        log.Fatal("JWT_SECRET no definido en el archivo .env")
    }

    credentialPath := os.Getenv("GOOGLE_CREDENTIALS")
    if credentialPath == "" {
        log.Fatal("GOOGLE_CREDENTIALS no definido en el archivo .env")
    }

    // Inicializar Firestore
    middleware.InitFirestore()
    // defer middleware.Client.Close()

    // Inicializar router
    r := gin.Default()
    routes.RegisterRoutes(r, []byte(jwtSecret)) // 👈 pasamos jwtSecret aquí
    r.Run("0.0.0.0:5001")

}
