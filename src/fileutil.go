package main

import (
	"log"
	"os"
	"path/filepath"
)

func writeFile(filename string, textToWrite string) {
	fullPath := filepath.Join(homeDirectory, filename)
	err := os.WriteFile(fullPath, []byte(textToWrite), 0666)
	if err != nil {
		log.Fatalf("Failed to write file %s %v", filename, err)
	}
}

func deleteFile(filename string) {
	if err := os.Remove(filepath.Join(homeDirectory, filename)); err != nil {
		log.Printf("Error deleting file: %v", err)
	}
}
