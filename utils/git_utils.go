package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
)

func GetNewUploadDir() string {
	outputDIR := fmt.Sprintf("%v/%v/", GetOutputDir(), generateRamdomString(7))
	return outputDIR

}

func GetOutputDir() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	OUTPUT_DIR := wd + "/output/"
	return OUTPUT_DIR

}

func CloneRepo(repoURL string) string {
	outputDIR := GetNewUploadDir()
	log.Default().Println("Cloning repo: ", repoURL, " to:", outputDIR)

	git.PlainClone(outputDIR, false, &git.CloneOptions{
		URL: repoURL,
	})
	return outputDIR
}
