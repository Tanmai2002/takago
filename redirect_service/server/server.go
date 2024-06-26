package server

import (
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Tanmai2002/takago/redirect_service/providers"
	"github.com/gin-gonic/gin"
)

func init() {
	mime.AddExtensionType(".css", "text/css")
	mime.AddExtensionType(".js", "text/javascript")
	mime.AddExtensionType(".html", "text/html")
	mime.AddExtensionType(".png", "image/png")
	mime.AddExtensionType(".jpg", "image/jpg")
	mime.AddExtensionType(".jpeg", "image/jpeg")
	mime.AddExtensionType(".svg", "image/svg+xml")
	mime.AddExtensionType(".json", "application/json")
	mime.AddExtensionType(".woff", "font/woff")
	mime.AddExtensionType(".woff2", "font/woff2")
	mime.AddExtensionType(".ttf", "font/ttf")
	mime.AddExtensionType(".otf", "font/otf")
	mime.AddExtensionType(".eot", "font/eot")
	mime.AddExtensionType(".ico", "image/x-icon")
	mime.AddExtensionType(".mp4", "video/mp4")
	mime.AddExtensionType(".webm", "video/webm")
	mime.AddExtensionType(".webp", "image/webp")
	mime.AddExtensionType(".pdf", "application/pdf")
	mime.AddExtensionType(".zip", "application/zip")
	mime.AddExtensionType(".rar", "application/x-rar-compressed")
	mime.AddExtensionType(".tar", "application/x-tar")
	mime.AddExtensionType(".gz", "application/gzip")
	mime.AddExtensionType(".bz2", "application/x-bzip2")
	mime.AddExtensionType(".json", "application/json")
}

func InitializeServer() *gin.Engine {
	server := gin.Default()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	address := ":" + port
	server.GET("/*t", redirectionHandler)
	server.Run(address) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	return server
}

func redirectionHandler(c *gin.Context) {
	id := strings.Split(c.Request.Host, ".")[0]
	path := c.Request.URL.Path
	log.Println(id)
	log.Println(path)
	if len(strings.Split(filepath.Base(path), ".")) < 2 {
		path = "/index.html"
	}

	reader, ctype, err := providers.GetFileFromS3(filepath.Join(id, path))
	if err != nil {

		panic(err)
	}
	defer reader.Close()

	//mime from extenstion of filepath

	*ctype = mime.TypeByExtension(filepath.Ext(path))
	log.Println(filepath.Base(path))
	log.Println(*ctype)
	c.Header("Content-Type", *ctype)
	if _, err := io.Copy(c.Writer, reader); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

}
