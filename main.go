package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const InitialsDelimiter = ':'

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

	scanner := bufio.NewScanner(pairsFile)
	var lines []string
	initialsToDetails := make(map[string][]string)
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
		components := strings.Split(strings.TrimSpace(scanner.Text()), " ")
		if len(components) == 4 {
			initials := components[0]
			initials = initials[:strings.IndexByte(initials, InitialsDelimiter)]

			name := components[1] + " " + components[2]
			name = name[:len(name)-1]

			emailName := components[3]

			initialsToDetails[initials] = []string{name, emailName}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Could not scan file: %v", fileOpenErr)
	}

	return initialsToDetails
}
