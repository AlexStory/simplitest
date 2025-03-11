package main

import (
	"io/fs"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func main() {
	log.Println("listening for changes...")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Has(fsnotify.Write) && strings.HasSuffix(event.Name, ".go") {
					runTests()
				}

				dir := isDir(event.Name)

				if dir && event.Has(fsnotify.Create) {
					addWatcher(watcher, event.Name)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				log.Println("error:", err)
			}
		}
	}()

	fileSystem := os.DirFS(".")

	err = fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if strings.HasPrefix(path, ".git") {
			return nil
		}

		if d.IsDir() {
			addWatcher(watcher, path)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	<-make(chan struct{})
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

func addWatcher(watcher *fsnotify.Watcher, path string) {
	err := watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}
}

func runTests() {
	log.Println("files changed, running tests...")
	cmd := exec.Command("go", "test", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
