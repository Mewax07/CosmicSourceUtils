package main

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
)

func recursiveSort(name string, btld string) {
	recursiveSort1(name, name, btld)
}

func recursiveSort1(tld string, name string, btld string) {
	entries, _ := os.ReadDir(name)
	for _, entry := range entries {
		if entry.IsDir() {
			recursiveSort1(tld, name+"/"+entry.Name(), btld)
		} else {
			if name == tld {
				if !check2(name+"/"+entry.Name(), ".java") {
					os.Rename(name+"/"+entry.Name(), strings.ReplaceAll(name+"/"+entry.Name(), tld+"/", btld+"/main/resources/"))
				}
			}
			if name != tld {
				if check2(name+"/"+entry.Name(), ".java") {
					fmt.Println(strings.ReplaceAll(name, tld+"/", btld+"/main/java/"))
					os.MkdirAll(strings.ReplaceAll(name, tld+"/", btld+"/main/java/"), fs.ModePerm)
					os.Rename(name+"/"+entry.Name(), strings.ReplaceAll(name+"/"+entry.Name(), tld+"/", btld+"/main/java/"))
				} else {
					fmt.Println(strings.ReplaceAll(name, tld+"/", btld+"/main/resources/"))
					os.MkdirAll(strings.ReplaceAll(name, tld+"/", btld+"/main/resources/"), fs.ModePerm)
					os.Rename(name+"/"+entry.Name(), strings.ReplaceAll(name+"/"+entry.Name(), tld+"/", btld+"/main/resources/"))
				}
			}
		}
	}
}
