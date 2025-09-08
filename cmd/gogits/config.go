package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/user"
	"slices"
	"strings"
)

type GogitsConfig struct {
	pathToDirectories string
	pathToEmails      string
}

const (
	mode    = int(0755)
	baseDir = "/.gogits"
)

func getBasePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return usr.HomeDir + baseDir
}

func getConfig() GogitsConfig {
	basePath := getBasePath()

	config := GogitsConfig{
		pathToDirectories: basePath + "/.directories",
		pathToEmails:      basePath + "/.emails",
	}

	return config
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

func addNewSliceElementsToFile(filePath string, newRepos []string) {
	exisitingRepos := parseFileLinesToSlice(filePath)
	repos := joinSlices(newRepos, exisitingRepos)
	dumpStringsSliceToFile(repos, filePath)
}

func dumpStringsSliceToFile(content []string, filePath string) {
	lines := strings.Join(content, "\n")
	basePath := getBasePath()

	os.Mkdir(basePath, os.FileMode(mode))
	err := os.WriteFile(filePath, []byte(lines), os.FileMode(mode))
	if err != nil {
		panic(err)
	}
}

func addEmailToConfig(email string) {
	filePath := getConfig().pathToEmails
	addNewSliceElementsToFile(filePath, []string{email})
}

func addReposToConfig(paths []string) {
	filePath := getConfig().pathToDirectories
	addNewSliceElementsToFile(filePath, paths)
}

func getReposFromConfig() []string {
	filePath := getConfig().pathToDirectories
	return parseFileLinesToSlice(filePath)
}

func getEmailsFromConfig() []string {
	filePath := getConfig().pathToEmails
	return parseFileLinesToSlice(filePath)
}
