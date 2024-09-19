package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func commit(commitMessage string, coAuthors string) {
	if commitMessage == "" {
		// Write co authors to commit template file to be read with -t flag on commit
		// We do this so that authors appear automatically in native text editor when not providing message
		// Using --edit -m instead will not abort an unsaved commit message so this is the alternative
		writeFile(CommitTemplateFile, coAuthors)

		// Ensure that the file is deleted after the function completes
		defer deleteFile(CommitTemplateFile)

		executeCommitWithTemplate()
	} else {
		executeCommitWithoutTemplate(commitMessage, coAuthors)
	}
}

func executeCommitWithTemplate() {
	fmt.Printf("Opening native text editor to write commit message ...\n")
	commitTemplatePath := filepath.Join(homeDirectory, CommitTemplateFile)
	cmd := exec.Command("git", "commit", "-t", commitTemplatePath)

	// Set the command's standard input/output/error to the current process's
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	//TODO improve error handling
	if err != nil {
		os.Exit(1)
	}
}

func executeCommitWithoutTemplate(commitMessage string, coAuthors string) {
	commitMessage = commitMessage + "\n\n" + coAuthors
	cmd := exec.Command("git", "commit", "-m", commitMessage)

	// Set the command's standard input/output/error to the current process's
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	//TODO improve error handling
	if err != nil {
		os.Exit(1)
	}
}
