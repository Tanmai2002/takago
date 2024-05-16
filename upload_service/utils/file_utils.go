package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// Get the Root Directory for the Project
func GetCWD() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return wd

}

// Whenever we need to clone the files to local,
// we need new folder to clone the code
func GetNewUploadDir(id string) string {
	outputDIR := fmt.Sprintf("%v%v/", GetOutputDir(), id)
	return outputDIR

}

// Any Prefixes to be added for uploading
func GetOutputDir() string {

	OUTPUT_DIR := "./output/"
	return OUTPUT_DIR

}

// To get all files in a directory as List
func GetAllFilesInDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
