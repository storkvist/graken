package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

const version = "0.1.0"

func main() {
	rootPath := "."
	if len(os.Args) > 1 {
		rootPath = os.Args[1]

		if rootPath == "version" {
			fmt.Printf("graken v%v %v/%v\n", version, runtime.GOOS, runtime.GOARCH)
			os.Exit(0)
		}
	}
	rootPath, err := filepath.Abs(rootPath)
	if err != nil {
		panic(err)
	}

	repositories := make([]string, 0)
	findRepositories(rootPath, &repositories)
	switch count := len(repositories); count {
	case 0:
		fmt.Printf("No git repositories found inside %s\n", rootPath)
	case 1:
		fmt.Println("Fetching one found repository...")
	default:
		fmt.Printf("Fetching %d found repositories...\n", count)
	}

	fetchRepositories(repositories)
}

func fetchRepositories(repositories []string) {
	start := time.Now()
	done := make(chan string)
	for _, repository := range repositories {
		go fetchRepository(repository, done)
	}
	for range repositories {
		<-done
	}
	elapsed := time.Since(start)
	if count := len(repositories); count > 0 {
		switch count {
		case 1:
			fmt.Printf("Fetched one repository in %s\n", elapsed)
		default:
			fmt.Printf("Fetched %d repositories in %s\n", count, elapsed)
		}
	}
}

func fetchRepository(path string, done chan<- string) {
	fmt.Printf("Fetching %s\n", path)
	fetchCmd := exec.Command("git", "fetch", "--all")
	fetchCmd.Dir = path
	if err := fetchCmd.Run(); err != nil {
		_, errStderr := fmt.Fprintf(os.Stderr, "fetch %v: %v\n", path, err)
		if errStderr != nil {
			fmt.Printf("fetch %v: %v\n", path, err)
			fmt.Printf("output to stderr: %v\n", errStderr)
		}
	}
	done <- path
}

func findRepositories(path string, repositories *[]string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.IsDir() {
			if file.Name() == ".git" {
				*repositories = append(*repositories, path)
				return
			}

			findRepositories(filepath.Join(path, file.Name()), repositories)
		}
	}
}
