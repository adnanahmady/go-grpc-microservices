package app

import (
	"log"
	"os"
	"regexp"
)

func GetRootDir() string {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current directory: %v", err)
	}
	return regexp.MustCompile("internal|bin|pkg|vendor").Split(currentDir, -1)[0]
}
