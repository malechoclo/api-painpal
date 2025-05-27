package routes

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "painpal_api/controllers"
    "painpal_api/middleware"
)

func RegisterRoutes(r *gin.Engine, jwtSecret []byte) {
    r.GET("/health", func(c *gin.Context) {
        c.String(http.StatusOK, "hello world")
    })
    r.POST("/register", controllers.Register)
   // r.POST("/login", controllers.Login)
    r.POST("/login", func(c *gin.Context) {
        var req struct {
            Email    string `json:"email" binding:"required"`
            Password string `json:"password" binding:"required"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error ShouldBindJSON": err.Error()})
            return
        }
        resp, err := controllers.LoginWithEmailAndPassword(req.Email, req.Password)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error LoginWithEmailAndPassword": err.Error()})
            return
        }
        c.JSON(http.StatusOK, resp)
    })

    // Grupo protegido
    protected := r.Group("/api")
    protected.Use(middleware.AuthMiddleware()) 
    {
        protected.POST("/survey", controllers.SaveSurvey)
        protected.GET("/survey/:survey_id", controllers.GetSurvey)
        protected.DELETE("/survey/:survey_id", controllers.DeleteSurvey)
        
        protected.GET("/surveys", controllers.GetSurveys)

        protected.GET("/user/surveys", controllers.GetSurveysByUser)

        

        
    }
    
	}