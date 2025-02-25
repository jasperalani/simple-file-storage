package main

import "os"

// Config - Configuration
type Config struct {
	permittedFileExtensions []string
	dir                     ApplicationFolders
}

// ApplicationFolders - Folders used by the application.
type ApplicationFolders struct {
	permissions os.FileMode
	storage     string
	files       string
	errors      string
	logs        string
}

func getConfig() Config {

	/*
		Edit config here
	*/
	var config = Config{
		permittedFileExtensions: []string{
			"mp3",
			"wav",
			"flac",
		},
		dir: ApplicationFolders{ // These folders will be created if they do not exist.
			permissions: 0777, // Permissions the folders will be created with.
			storage:     "uploads",
			files:       "uploads",
			errors:      "errors",
			logs:        "logs",
			// temporary:   "storage/temporary",
		},
	}

	return config
}
