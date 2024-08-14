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

const InitialsDelimiter = ':'
const PairsConfigFile = ".pairs"
const CommitTemplateFile = ".commitPairsTemplate"

func main() {
	contributorInitials := os.Args[1:]

	if len(contributorInitials) < 1 {
		log.Fatalf("Missing command line arguments. Please use the format 'git pcommit [primary intials] [co author initials]'")
	}

	homeDirectory := resolveHomeDirectory()
	contributors, domain := parsePairsFile(homeDirectory)
	primary := contributors[contributorInitials[0]]
	setPrimaryUsername(primary[0])
	setPrimaryEmail(primary[1], domain)

	//TODO write co authors to template.
	//TODO if message provided with -m then we should append our authors to the end of the message
	//TODO if message not provided then we should run the git commit command which should read from template

	writeToCommitTemplate(resolveCoAuthorDetails(contributorInitials[1:], contributors), domain, homeDirectory)
	executeCommitWithTemplate(homeDirectory)
}

func resolveHomeDirectory() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Could not resolve home directory: %v", err)
	}
	return homeDir
}

func parsePairsFile(homeDirectory string) (map[string][]string, string) {
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

// Write co authors to commit template file to be read with -t flag on commit
// We do this so that authors appear automatically in native text editor when not providing message
func writeToCommitTemplate(coAuthorDetails map[string][]string, domain string, path string) {
	var sb strings.Builder
	sb.WriteString("\n")

	for _, value := range coAuthorDetails {
		line := "Co-authored-by: " + value[0] + " <" + value[1] + "@" + domain + ">\n"
		sb.WriteString(line)
	}

	commitTemplatePath := filepath.Join(path, CommitTemplateFile)
	err := os.WriteFile(commitTemplatePath, []byte(sb.String()), 0666)

	if err != nil {
		log.Fatalf("Failed to write commit template file: %v", err)
	}
}

func executeCommitWithTemplate(pathToTemplate string) {
	commitTemplatePath := filepath.Join(pathToTemplate, CommitTemplateFile)
	cmd := exec.Command("git", "commit", "-t", commitTemplatePath)

	// Set the command's standard input/output/error to the current process's
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		log.Fatalf("Something went wrong when running git commit: %v", err)
	}

	fmt.Printf("Opening native text editor to write commit message ...\n")
}
