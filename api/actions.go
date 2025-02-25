package main

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Upload action
// Upload a file into a compression queue or storage.
func Upload(c *gin.Context) {

	// Validate upload

	file, err := c.FormFile("file")
	if nil != err {
		respond(c, 400, gin.H{
			"error":   "bad request",
			"message": "A file must be supplied as value in form data with key 'file'",
		})
		return
	}

	// file.Filename must be a file name and an extension
	matched, err := regexp.MatchString(`[^\\]*\.(\w+)$`, file.Filename)
	if nil != err {
		respond(c, 400, gin.H{
			"error":   "bad request",
			"message": "File name must be a file name and an extension",
		})
		return
	}

	if !matched {
		respond(c, 400, gin.H{
			"error":   "bad request",
			"message": "Content Disposition Filename must be supplied",
		})
		return
	}

	var (
		explodedFilename []string = strings.Split(file.Filename, ".")
		extension                 = explodedFilename[1]
	)

	if !stringInSlice(strings.ToLower(extension), config.permittedFileExtensions) {
		var pretty = strings.Join(config.permittedFileExtensions, ", ")
		respond(c, 400, gin.H{
			"error":   "bad request",
			"message": "File extension must be one of: " + pretty,
		})
		return
	}

	var (
		filename     = generateIdentifier()
		saveFilePath = config.dir.files + "/" + filename + "." + extension
		FileOptions  = File{
			Name:      explodedFilename[0],
			Extension: extension,
		}
	)

	// Move file to queue location
	err = c.SaveUploadedFile(file, saveFilePath)
	if nil != err {
		respond(c, 500, gin.H{
			"error":   "server failure",
			"message": "Failed to save file.",
		})
		return
	}

	// Create options file
	optionsFile, err := os.Create(saveFilePath + ".yaml")
	if err != nil {
		respond(c, 500, gin.H{
			"error":   "server failure",
			"message": err,
		})
		return
	}

	// Marshal the object into YAML format
	yamlData, err := yaml.Marshal(&FileOptions)
	if err != nil {
		respond(c, 500, gin.H{
			"error":   "server failure",
			"message": err,
		})
		return
	}

	// Write the YAML data to the file
	_, err = optionsFile.Write(yamlData)
	if err != nil {
		respond(c, 500, gin.H{
			"error":   "server failure",
			"message": err,
		})
		return
	}

	defer func(optionsFile *os.File) {
		err := optionsFile.Close()
		if err != nil {
			respond(c, 500, gin.H{
				"error":   "server failure",
				"message": err,
			})
		}
	}(optionsFile)

	respond(c, http.StatusOK, gin.H{
		"message": "File has been uploaded successfully.",
		"id":      filename,
	})

	log.Println(time.Now().Format("2006-01-02 15:04:05") + " File uploaded: " + saveFilePath)
}

// Retrieve action
// Retrieve an uploaded file
func Retrieve(c *gin.Context) {
	id := c.Param("id")
	if !validIdentifier(id) {
		respond(c, 400, gin.H{
			"error":   "bad request",
			"message": "Invalid file identifier",
		})
		return
	}

	var files []string

	err := filepath.Walk(config.dir.files, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		respond(c, 500, gin.H{
			"error":   "server failure",
			"message": err,
		})
	}

	var foundFile = "not found"

	for _, file := range files {
		if strings.Contains(file, id) {
			foundFile = file
		}
	}

	if "not found" == foundFile {
		respond(c, http.StatusOK, gin.H{
			"error":   "file not found",
			"message": "No file matching the supplied identifier was found.",
		})
		return
	}

	var (
		explode1 []string
		explode2 []string
	)

	explode1 = strings.Split(foundFile, "/")
	explode2 = strings.Split(explode1[1], ".")

	var (
		fileName        string = explode2[0]
		fileExt         string = explode2[1]
		fileOptionsPath string = fileName + "." + fileExt + "." + explode2[2]
		fileOptions     File
	)

	// Open the YAML file for reading
	file, err := os.ReadFile("./uploads/" + fileOptionsPath)
	if err != nil {
		respond(c, 500, gin.H{
			"error":   "server failure",
			"message": err,
		})
		return
	}

	err = yaml.Unmarshal(file, &fileOptions)
	if err != nil {
		respond(c, 500, gin.H{
			"error":   "server failure",
			"message": err,
		})
		return
	}

	respond(c, http.StatusOK, gin.H{
		"id":   id,
		"file": fileOptions.Name + "." + fileOptions.Extension,
		"link": "localhost/uploads/" + id + "." + fileExt,
	})

	log.Println(time.Now().Format("2006-01-02 15:04:05") + " File retrieved: " + fileName + "." + fileExt)
}
