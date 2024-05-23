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
	s.GET("/check_id/:id", checkID)
	s.GET("/deployments/:id", getDeployment)
	// s.GET("/deployments/:id", getDeployment)
	// s.PUT("/deployments/:id", updateDeployment)
	// s.DELETE("/deployments/:id", deleteDeployment)
}

type DeploymentStruct struct {
	RepoURL      string `json:"repoURL"`
	DeploymentId string `json:"id" bson:"id" default:""`
}

func checkID(c *gin.Context) {
	id := c.Param("id")
	exists := utils.CheckIfExist(utils.RepoCollection, utils.TakaGoProjectID{ID: id})
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot Create as ID Already Exists"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Available to Create",
		"id":      id,
	})
}

func getDeployment(c *gin.Context) {
	id := c.Param("id")
	result := utils.FindOne(utils.RepoCollection, utils.TakaGoProjectID{ID: id})
	if result.Err() != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Not Found"})
		return
	}
	var project utils.TakaGoProject
	result.Decode(&project)
	c.JSON(http.StatusOK, project)

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
	if jsonBody.DeploymentId != "" {
		id = jsonBody.DeploymentId
	}
	exists := utils.CheckIfExist(utils.RepoCollection, utils.TakaGoProjectID{ID: id})
	log.Printf("Exists: %v", exists)
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot Create as ID Already Exists"})
		return
	}
	utils.InsertOne(utils.RepoCollection, utils.TakaGoProject{ID: id, RepoURL: jsonBody.RepoURL})

	go uploadFilesToS3(id, jsonBody.RepoURL)

	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"repoURL": jsonBody.RepoURL,
		"id":      id,
	})
}

func uploadFilesToS3(id string, repoUrl string) {
	outputDir := utils.CloneRepo(repoUrl, id)

	//Get All the File List Format

	files, err := utils.GetAllFilesInDir(outputDir)
	if err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		utils.UpdateOne(utils.RepoCollection, utils.TakaGoProjectID{ID: id}, utils.TakaGoProject{ID: id, RepoURL: repoUrl, Branch: "main", Status: "cancelledbyError"})

		panic(err)
	}
	for _, file := range files {
		log.Default().Println(file, filepath.Join(utils.GetCWD(), file))
		utils.UploadS3File(file, filepath.Join(utils.GetCWD(), file), "test")
	}
	log.Default().Println(files)
	//Push to Redis
	utils.PushToRedisBuildQueue(id)
	x, err := utils.UpdateOne(utils.RepoCollection, utils.TakaGoProjectID{ID: id}, utils.TakaGoProject{ID: id, RepoURL: repoUrl, Branch: "main", Status: "queued"})
	log.Println(err)
	log.Println(x.UpsertedID)
	go utils.CleanupFiles(outputDir)
}
