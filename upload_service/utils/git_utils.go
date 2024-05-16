package utils

import (
	"log"

	"github.com/go-git/go-git/v5"
)

// To Clone Repositories from URL
func CloneRepo(repoURL string, id string) string {
	outputDIR := GetNewUploadDir(id)
	log.Default().Println("Cloning repo: ", repoURL, " to:", outputDIR)

	git.PlainClone(outputDIR, false, &git.CloneOptions{
		URL: repoURL,
	})
	return outputDIR
}
