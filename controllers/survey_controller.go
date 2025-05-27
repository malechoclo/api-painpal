package controllers

import (
	"net/http"
	"time"
	"fmt"
	"github.com/gin-gonic/gin"
	"painpal_api/firebase" // Replace with the actual import path to your firebase package
	"github.com/google/uuid"
	"strings"
	"google.golang.org/api/iterator")


func SaveSurvey(c *gin.Context) {
	survey_id := uuid.New().String()
	fmt.Println("SaveSurvey called")
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Asociar el UID en el documento
	data["uid"] = c.MustGet("uid").(string)
	data["survey_id"] = survey_id
	data["createdAt"] = time.Now()
	fmt.Println("SaveSurvey called",survey_id)
	_, _, err := firebase.FirestoreClient.Collection("surveys").Add(c, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Formulario guardado exitosamente"})

}

func GetSurvey(c *gin.Context) {
	survey_id := strings.TrimSpace(c.Param("survey_id"))
	fmt.Println("GetSurvey called", survey_id)
	iter := firebase.FirestoreClient.Collection("surveys").Where("survey_id", "==", survey_id).Documents(c)
	doc, err := iter.Next()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Survey not found"})
		return
	}
	c.JSON(http.StatusOK, doc.Data())
}

func DeleteSurvey(c *gin.Context) {
    // 1. Obtener y limpiar el survey_id de la ruta
	surveyID := strings.TrimSpace(c.Param("survey_id"))

    // 2. Contexto de la petición
    ctx := c.Request.Context()

    // 3. Cliente de Firestore (debe estar inicializado previamente)
    client := firebase.FirestoreClient

    // 4. Consulta los documentos que tengan survey_id == surveyID
    iter := client.
        Collection("surveys").
        Where("survey_id", "==", surveyID).
        Documents(ctx)
    defer iter.Stop()

    deleted := 0
    for {
        docSnap, err := iter.Next()
        if err == iterator.Done {
            break
        }
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error":   "Error querying documents",
                "details": err.Error(),
            })
            return
        }

        // 5. Elimina cada documento
        if _, err := docSnap.Ref.Delete(ctx); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error":   "Failed to delete a survey document",
                "details": err.Error(),
            })
            return
        }
        deleted++
    }

    // 6. Responder con éxito
    c.JSON(http.StatusOK, gin.H{
        "message":     "Surveys deleted successfully",
        "survey_id":   surveyID,
        "deletedCount": deleted,
    })
}


func GetSurveys(c *gin.Context) {
	iter := firebase.FirestoreClient.Collection("surveys").Documents(c)
	var surveys []map[string]interface{}
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		survey := doc.Data()
		survey["id"] = doc.Ref.ID
		surveys = append(surveys, survey)
	}
	c.JSON(http.StatusOK, surveys)
}

func GetSurveysByUser(c *gin.Context) {
	user_id := strings.TrimSpace(c.MustGet("uid").(string))
	
	fmt.Println("GetSurveysByUser",user_id)	
	iter := firebase.FirestoreClient.Collection("surveys").Where("uid", "==", user_id).Documents(c)
	var surveys []map[string]interface{}
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		survey := doc.Data()
		survey["id"] = doc.Ref.ID
		surveys = append(surveys, survey)
	}
	c.JSON(http.StatusOK, surveys)
}
