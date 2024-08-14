package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const InitialsDelimiter = ':'

func main() {
	contributorInitials := os.Args[1:]

	if len(contributorInitials) < 1 {
		log.Fatalf("Missing command line arguments. Please use the format 'git pcommit [primary intials] [co author initials]'")
	}

	contributors := resolveContributors()
	writeToCommitTemplate(resolveCoAuthorDetails(contributorInitials[1:], contributors))
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

func resolveCoAuthorDetails(contributorInitials []string, contributors map[string][]string) map[string][]string {
	coAuthorDetails := make(map[string][]string)

	for _, value := range contributorInitials {
		if _, ok := contributors[value]; ok {
			coAuthorDetails[value] = contributors[value]
		} else {
			log.Printf("Warning: could not find user: %v ... skipping", value)
		}
	}

	return coAuthorDetails
}

func writeToCommitTemplate(contributorDetails map[string][]string) {

}
