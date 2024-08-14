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

	contributors, _ := parsePairsFile()
	primary := contributors[contributorInitials[0]]
	setPrimaryAuthor(primary[0], primary[1])

	writeToCommitTemplate(resolveCoAuthorDetails(contributorInitials[1:], contributors))
}

func parsePairsFile() (map[string][]string, string) {
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

	defer func(pairsFile *os.File) {
		err := pairsFile.Close()
		if err != nil {
			log.Fatalf("Error closing file: %v", err)
		}
	}(pairsFile)

	scanner := bufio.NewScanner(pairsFile)
	var lines []string
	initialsToDetails := make(map[string][]string)
	var domain string
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
		if len(components) == 2 && strings.ToLower(components[0]) == "domain:" {
			domain = components[1]
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Could not scan file: %v", fileOpenErr)
	}

	if len(initialsToDetails) == 0 {
		log.Fatalf("Could not find listings in pairs file. Please ensure following format is used \njd: John Doe; john.doe")
	}

	if domain == "" {
		log.Fatalf("Could not find a domain for pair emails. Please ensure the following format is used \ndomain: google.com")
	}

	return initialsToDetails, domain
}

func setPrimaryAuthor(fullname string, email string) {

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

func writeToCommitTemplate(coAuthorDetails map[string][]string) {

}
