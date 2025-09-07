package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"slices"
	"strings"
)

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

			if fileName == "vendor" || fileName == "node_modules" {
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

func getDotFilePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	dotFile := usr.HomeDir + "/.gogits"

	return dotFile
}

func addNewSliceElementsToFile(filePath string, newRepos []string) {
	exisitingRepos := parseFileLinesToSlice(filePath)
	repos := joinSlices(newRepos, exisitingRepos)
	dumpStringsSliceToFile(repos, filePath)
}

func parseFileLinesToSlice(filePath string) []string {
	f, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}
		}
		panic(err)
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			panic(err)
		}
	}

	return lines
}

func joinSlices(new []string, existing []string) []string {
	for _, i := range new {
		if !slices.Contains(existing, i) {
			existing = append(existing, i)
		}
	}

	return existing
}

func dumpStringsSliceToFile(repos []string, filePath string) {
	content := strings.Join(repos, "\n")
	err := os.WriteFile(filePath, []byte(content), 0755)
	if err != nil {
		panic(err)
	}
}

func scan(folder string) {
	fmt.Printf("Found folders:\n\n")
	repos := recursiveScanFolder(folder)
	filePath := getDotFilePath()
	addNewSliceElementsToFile(filePath, repos)
	fmt.Printf("\n\nSuccessfully added\n\n")
}
