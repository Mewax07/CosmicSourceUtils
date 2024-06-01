package main

import (
	"os"
	"strings"
)

func deleteEmptyDirs(name string) {
	entries, _ := os.ReadDir(name)

	for _, entry := range entries {
		if entry.IsDir() {
			if entry.Name() == "com" || entry.Name() == "META-INF" {
				os.RemoveAll(name + "/" + entry.Name())
			}
			deleteEmptyDirs(name + "/" + entry.Name())
		}
	}
	entries, _ = os.ReadDir(name)

	if len(entries) == 0 {
		os.RemoveAll(name)
	}
}

func recursiveCleanUnNeededFiles(name string) {
	for _, path := range recursiveCleanUnNeededFiles2(name, name) {
		os.RemoveAll(path)
	}
}

func recursiveCleanUnNeededFiles2(tld string, name string) []string {
	removableItems := make([]string, 0)

	entries, _ := os.ReadDir(name)

	for _, entry := range entries {
		if entry.IsDir() {
			removableItems = append(removableItems, recursiveCleanUnNeededFiles2(tld, name+"/"+entry.Name())...)
		} else {
			if !strings.HasPrefix(strings.ReplaceAll(name, tld+"/", ""), "finalforeach") &&
				!strings.HasPrefix(strings.ReplaceAll(name, tld+"/", ""), "libs") &&
				!strings.HasPrefix(strings.ReplaceAll(name, tld+"/", ""), "opensimplex2") {
				removableItems = check(removableItems, name+"/"+entry.Name(), ".java")
				removableItems = check(removableItems, name+"/"+entry.Name(), ".dll")
				removableItems = check(removableItems, name+"/"+entry.Name(), ".so")
				removableItems = check(removableItems, name+"/"+entry.Name(), ".dylib")
				removableItems = check(removableItems, name+"/"+entry.Name(), "gamecontrollerdb.txt")
				removableItems = check(removableItems, name+"/"+entry.Name(), "sfd.ser")
				removableItems = check(removableItems, name+"/"+entry.Name(), ".xml")
			}
			if name == tld {
				removableItems = check(removableItems, name+"/"+entry.Name(), ".png")
			}

		}
	}

	return removableItems
}
