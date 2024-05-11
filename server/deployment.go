package server

import (
	"log"
	"net/http"

	"github.com/Tanmai2002/takago/utils"
	"github.com/gin-gonic/gin"
)

func AddDeploymentRoutes(s *gin.Engine) {
	s.POST("/deploy", createDeployment)
	// s.POST("/deployments", createDeployment)
	// s.GET("/deployments/:id", getDeployment)
	// s.PUT("/deployments/:id", updateDeployment)
	// s.DELETE("/deployments/:id", deleteDeployment)
}

type DeploymentStruct struct {
	RepoURL string `json:"repoURL"`
}

func createDeployment(c *gin.Context) {
	//get repoURL from body of request
	var jsonBody DeploymentStruct
	if err := c.BindJSON(&jsonBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Default().Println(jsonBody)
	outputDir := utils.CloneRepo(jsonBody.RepoURL)

	//Get All the File List Format

	files, err := utils.GetAllFilesInDir(outputDir)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		panic(err)
		return
	}
	log.Default().Println(files)

	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"repoURL": jsonBody.RepoURL,
	})
}
