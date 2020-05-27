package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	path := "."
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	path, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	repositories := make([]string, 0)
	findRepositories(path, &repositories)
	count := len(repositories)
	switch count {
	case 0:
		fmt.Printf("No git repositories found inside %s\n", path)
	case 1:
		fmt.Println("Fetching one found repository...")
	default:
		fmt.Printf("Fetching %d found repositories...\n", count)
	}

	start := time.Now()
	done := make(chan string)
	for _, repository := range repositories {
		go fetchRepository(repository, done)
	}
	for range repositories {
		<-done
	}
	elapsed := time.Since(start)
	if count > 0 {
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
		fmt.Fprintf(os.Stderr, "fetch %v: %v\n", path, err)
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
