package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const InitialsDelimiter = ':'
const PairsConfigFile = ".pairs"
const CommitTemplateFile = ".commitPairsTemplate"
const DaysPairFile = ".daysPair"

var homeDirectory = resolveHomeDirectory()

func main() {

	messageFlag := flag.String("m", "", "Commit message to include inline")
	pairsFlag := flag.Bool("p", false, "Provide a list of contributor initials")
	helpFlag := flag.Bool("h", false, "Display help message")

	flag.Parse()
	args := flag.Args()

	if *helpFlag {
		fmt.Println("Usage:")
		fmt.Println("  git pc [-m 'inline message'] [list of initials e.g. JD AS TJ]")
		fmt.Println("  git pc [-p] [list of initials]")
		flag.PrintDefaults()
		os.Exit(0)
	}

	// Handle the -p flag
	if *pairsFlag {
		if len(flag.Args()) == 0 {
			log.Fatalf("Error: No contributor initials provided after -p")
		}
		writePairsToFile(flag.Args())
		os.Exit(0)
	}

	var commitMessage string
	var contributorInitials []string

	// Handle the -m flag for commit message
	if *messageFlag != "" {
		commitMessage = *messageFlag
		if len(flag.Args()) > 0 {
			contributorInitials = flag.Args()
		}
	} else {
		// If no message flag, check if the arguments contain "-m"
		for i, arg := range args {
			if arg == "-m" {
				// If "-m" is found, the next argument should be the commit message
				if i+1 < len(args) {
					commitMessage = args[i+1]
				} else {
					log.Fatal("Error: No commit message provided after -m")
				}
				// The initials are the arguments before "-m"
				contributorInitials = args[:i]
				break
			}
		}
	}

	// If "-m" wasn't found in the args, treat all args as contributor initials
	if commitMessage == "" {
		contributorInitials = args
	}

	// If no contributors specified - read from file
	if len(contributorInitials) == 0 {
		contributorInitials = loadContributorsFromFile()
	}

	coAuthors := resolveContributorDetails(contributorInitials)
	commit(commitMessage, coAuthors)
}

func resolveHomeDirectory() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Could not resolve home directory: %v", err)
	}
	return homeDir
}
