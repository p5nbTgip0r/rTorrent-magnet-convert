package main

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

const directoryEnvironment = "DIRECTORY"
const permissionEnvironment = "PERMISSION"

var watchedDirectory string
var permission os.FileMode

func main() {
	if dirErr := options(); dirErr != nil {
		log.Fatal(dirErr)
	}

	if fi, err := os.Stat(watchedDirectory); err == nil {
		if !fi.IsDir() {
			log.Fatalf("Specified path \"%s\" is not a directory", watchedDirectory)
		}
	} else if os.IsNotExist(err) {
		log.Fatalf("Specified path \"%s\" does not exist", watchedDirectory)
	} else {
		log.Fatal(err)
	}

	log.Printf("Creating watcher for directory '%s'", watchedDirectory)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(watchedDirectory)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Watching directory '%s'", watchedDirectory)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				log.Fatal("Bad event")
			}
			log.Println("FS event:", event)
			if event.Op == fsnotify.Create || event.Op == fsnotify.Write {
				log.Println("File was modified")
				read(event.Name)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				log.Fatal("Bad error")
			}
			log.Println("Error:", err)
		}
	}
}

func options() error {
	var argsLen = len(os.Args)
	var dir string
	var givenPermission string

	if argsLen < 2 {
		dir = os.Getenv(directoryEnvironment)
		givenPermission = os.Getenv(permissionEnvironment)
	} else {
		dir = os.Args[1]
		if argsLen > 2 {
			givenPermission = os.Args[2]
		}
	}

	if dir == "" {
		return errors.New("no directory specified")
	}

	parsedUint, err := strconv.ParseUint(givenPermission, 10, 32)
	if err != nil {
		permission = 0664
	} else {
		permission = os.FileMode(parsedUint)
	}

	dir, dirErr := filepath.Abs(dir)
	watchedDirectory = dir
	return dirErr
}
