package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Tanmai2002/takago/deployment_service/utils"
)

func uploadBuildFilesToS3(dir_of_project string, build_folder string, id string) bool {
	local_build_folder := filepath.Join(dir_of_project, build_folder)
	files, err := utils.GetAllFilesInDir(local_build_folder)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		utils.UploadS3File(file[len(local_build_folder)+1:], file, "dist/"+id)
	}
	return true
}
func runBuildCommands(path string, commands *[]string) {
	os.Chdir(path)
	for _, command := range *commands {
		log.Println("Running Command", command)
		cmd := exec.Command("pwsh", "-c", command)
		out, err := cmd.Output()
		if err != nil {
			panic(err.Error())
		}
		log.Println(string(out))

	}

}
func buildForID(id string) {
	log.Println("Starting Build")
	dir := utils.GetDownloadDir(id)
	uploaded_dir := utils.GetAWSUploadFolder(id)
	log.Println("Downloading Files")
	local_dir := utils.DownloadFilesFromAWSToLocal(dir, uploaded_dir)
	log.Println(local_dir)
	commands := []string{"npm i", "npm run build"}

	runBuildCommands(local_dir, &commands)
	uploadBuildFilesToS3(local_dir, "build", id)

}
func main() {
	for {
		// id := "QLoFNQm"
		id, err := utils.PullFromRedisBuildQueue()
		if err != nil {
			log.Println("No ID Found")
			time.Sleep(time.Second)
			// return
			continue
		}
		log.Println("Id is " + id)
		buildForID(id)

	}
}
