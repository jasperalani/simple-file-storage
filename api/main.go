package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

var config Config

type File struct {
	Name      string `yaml:"name"`
	Extension string `yaml:"extension"`
}

func main() {

	config = getConfig()
	createApplicationFolders()

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.MaxMultipartMemory = 8 << 20

	router.POST("/upload", Upload)
	router.GET("/retrieve/:id", Retrieve)

	err := router.Run(":8080")
	if err != nil {
		fmt.Print(err.Error())
		return
	}

}
