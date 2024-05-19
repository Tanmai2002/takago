package server

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/Tanmai2002/takago/utils"
	"github.com/gin-gonic/gin"
)

func AddDeploymentRoutes(s *gin.Engine) {
	s.POST("/deploy", createDeployment)
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
	id := utils.GenerateRandomString(7)
	utils.InsertOne(utils.RepoCollection, utils.TakaGoProject{ID: id, RepoURL: jsonBody.RepoURL})

	go uploadFilesToS3(id, jsonBody.RepoURL)

	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"repoURL": jsonBody.RepoURL,
	})
}

func uploadFilesToS3(id string, repoUrl string) {
	outputDir := utils.CloneRepo(repoUrl, id)

	//Get All the File List Format

	files, err := utils.GetAllFilesInDir(outputDir)
	if err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		utils.UpdateOne(utils.RepoCollection, utils.TakaGoProject{ID: id}, utils.TakaGoProject{ID: id, RepoURL: repoUrl, Branch: "main", Status: "cancelledbyError"})

		panic(err)
	}
	for _, file := range files {
		log.Default().Println(file, filepath.Join(utils.GetCWD(), file))
		utils.UploadS3File(file, filepath.Join(utils.GetCWD(), file), "test")
	}
	log.Default().Println(files)
	//Push to Redis
	utils.PushToRedisBuildQueue(id)
	utils.UpdateOne(utils.RepoCollection, utils.TakaGoProject{ID: id}, utils.TakaGoProject{ID: id, RepoURL: repoUrl, Branch: "main", Status: "queued"})

}
