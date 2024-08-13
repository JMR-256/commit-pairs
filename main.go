package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	contributors := resolveContributors()
	fmt.Println(contributors)
}

func resolveContributors() map[string][]string {
	pairsFilePath, err := func() (string, error) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		return filepath.Join(homeDir, ".pairs"), nil
	}()

	if err != nil {
		log.Fatalf("Failed to get the home directory or join the path: %v", err)
	}
	pairsFile, fileOpenErr := os.Open(pairsFilePath)

	if fileOpenErr != nil {
		log.Fatalf("Open file failure %v", fileOpenErr)
	}

	defer pairsFile.Close()

	return nil
}
