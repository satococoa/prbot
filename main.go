package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("Hello, world")
	repos, exists := os.LookupEnv("GITHUB_REPOSITORY")
	if !exists {
		log.Fatalln("Please set GITHUB_REPOSITORY.")
	}
	fmt.Println("REPOSITORY: " + repos)
}
