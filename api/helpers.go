package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"log"
	"os"
	"regexp"
)

func createApplicationFolders() {
	var applicationFolders = []string{
		config.dir.storage,
		config.dir.files,
		config.dir.errors,
		config.dir.logs,
	}

	for _, dir := range applicationFolders {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.Mkdir(dir, config.dir.permissions)
			if err != nil {
				log.Fatal("Failed to create dir '" + dir + "'.")
			}
		}
	}
}

func generateIdentifier() string {
	guid := xid.New()
	return guid.String()
}

func validIdentifier(identifier string) bool {
	if len(identifier) != 20 {
		return false
	}

	matched, err := regexp.MatchString(`([A-Za-z0-9\-]+)`, identifier)
	if nil != err {
		log.Fatal(err)
	}

	if !matched {
		return false
	}
	return true
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func respond(c *gin.Context, code int, res gin.H) {
	c.JSON(code, res)
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
