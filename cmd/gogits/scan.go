package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

var ignore = []string{"vendor", "node_modules"}

func scanDotGitFolders(folders []string, folder string) []string {
	folder = strings.TrimSuffix(folder, "/")

	f, err := os.Open(folder)
	if err != nil {
		log.Fatal(err)
	}

	files, err := f.ReadDir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}

	var path string

	for _, file := range files {
		if file.IsDir() {
			fileName := file.Name()
			path = folder + "/" + fileName

			if fileName == ".git" {
				path = strings.TrimSuffix(path, "/.git")
				fmt.Println(path)
				folders = append(folders, path)
				continue
			}

			if slices.Contains(ignore, fileName) {
				continue
			}

			folders = scanDotGitFolders(folders, path)
		}
	}

	return folders
}

func recursiveScanFolder(folder string) []string {
	return scanDotGitFolders(make([]string, 0), folder)
}

func scan(folder string) {
	fmt.Printf("Found folders:\n\n")
	repos := recursiveScanFolder(folder)
	addReposToConfig(repos)
	fmt.Printf("\n\nSuccessfully added\n\n")
}
