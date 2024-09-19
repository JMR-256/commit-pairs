package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func loadContributorsFromFile() []string {
	contributorInitials, err := parseDaysPairFile()
	if err != nil {
		log.Fatalf("%v\nNo pair set:\nPlease set a pair to write to file with 'git pc -p [primary initials] [co-author initials]'\nOR\nProvide a one-time pair with 'git pc [primary initials] [co-author initials]'", err)
	}
	return contributorInitials
}

func writePairsToFile(initials []string) {
	writeFile(DaysPairFile, strings.Join(initials, " "))
}

func resolveContributorDetails(contributorInitials []string) (coAuthors string) {
	contributors, domain := parsePairsFile()

	primary := contributors[contributorInitials[0]]
	if primary == nil {
		log.Fatalf("Could not find mapping for initials '%v' in %v/%v", contributorInitials[0], homeDirectory, PairsConfigFile)
	}

	//doing this here is ugly and is annoying me
	setPrimaryUsername(primary[0])
	setPrimaryEmail(primary[1], domain)

	coAuthors = resolveCoAuthorDetails(contributorInitials[1:], contributors, domain)
	return coAuthors
}

func setPrimaryUsername(fullName string) {
	cmd := exec.Command("git", "config", "--global", "user.name", fullName)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Something went wrong setting the git username: %v", err)
	}
	fmt.Printf("Successfully updated git config user.name: %v\n", fullName)
}

func setPrimaryEmail(emailName string, domain string) {
	email := emailName + "@" + domain
	cmd := exec.Command("git", "config", "--global", "user.email", email)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Something went wrong setting the git email: %v", err)
	}
	fmt.Printf("Successfully updated git config user.email: %v\n", email)
}

func resolveCoAuthorDetails(contributorInitials []string, contributors map[string][]string, domain string) string {
	coAuthorDetails := make(map[string][]string)

	for _, value := range contributorInitials {
		if _, ok := contributors[value]; ok {
			coAuthorDetails[value] = contributors[value]
		} else {
			log.Printf("Warning: could not find user: %v ... skipping", value)
		}
	}

	var formattedCoAuthors strings.Builder
	formattedCoAuthors.WriteString("\n")
	for _, value := range coAuthorDetails {
		line := "Co-authored-by: " + value[0] + " <" + value[1] + "@" + domain + ">\n"
		formattedCoAuthors.WriteString(line)
	}

	return formattedCoAuthors.String()
}

func parseDaysPairFile() ([]string, error) {
	daysPairFilePath := filepath.Join(homeDirectory, DaysPairFile)
	daysPairFile, fileOpenErr := os.Open(daysPairFilePath)

	if fileOpenErr != nil {
		return nil, fileOpenErr
	}

	defer func(pairsFile *os.File) {
		err := pairsFile.Close()
		if err != nil {
			log.Fatalf("Error closing file: %v", err)
		}
	}(daysPairFile)

	scanner := bufio.NewScanner(daysPairFile)

	var initials []string
	if scanner.Scan() {
		initials = strings.Split(strings.TrimSpace(scanner.Text()), " ")
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Could not scan file: %v", fileOpenErr)
	}

	return initials, nil
}

func parsePairsFile() (map[string][]string, string) {
	pairsFilePath := filepath.Join(homeDirectory, PairsConfigFile)
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
