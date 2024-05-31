package main

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math/big"
	"net/http"
	url2 "net/url"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func downloadFile(url string) string {

	fileURL, err := url2.Parse(url)
	if err != nil {
		log.Fatal(err)
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]

	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)

	defer file.Close()

	return fileName

}

func cleanFile(name string) {
	if !strings.HasSuffix(name, ".java") {
		return
	}

	bytez, _ := os.ReadFile(name)
	p0, _ := regexp.Compile("[^`0-9A-Za-z\\s{}\"<>.=?:();&,\\\\!'|+-@#$%^*\\]\\[_]")
	content := string(bytez)
	namePrint := true
	for _, s := range strings.Split(content, "") {
		if p0.Match([]byte(s)) {
			if namePrint {
				fmt.Println(name)
				namePrint = false
			}
			char := string(p0.Find([]byte(s)))
			charInt := strings.ReplaceAll(strings.ReplaceAll(strconv.QuoteToASCII(char), "\\u", ""), "\"", "")
			n := new(big.Int)
			n.SetString(charInt, 16)

			fmt.Println(char, "->", n.Int64())
			content = strings.ReplaceAll(content, "'"+char+"'", strconv.Itoa(int(n.Int64())))
			os.WriteFile(name, []byte(content), fs.ModePerm)
		}
	}
}

func recursiveCleanDir(name string) {
	entries, _ := os.ReadDir(name)

	for _, entry := range entries {
		if entry.IsDir() {
			recursiveCleanDir(name + "/" + entry.Name())
		} else {
			cleanFile(name + "/" + entry.Name())
		}
	}
}

func check(arr []string, name string, suff string) []string {
	if strings.HasSuffix(name, suff) {
		arr = append(arr, name)
	}
	return arr
}

func recursiveSearch(name string) []string {
	return recursiveSearch2(name, name)
}

func recursiveRemoveEmptyDirs(name string) {
	entries, _ := os.ReadDir(name)

	for _, entry := range entries {
		if entry.IsDir() {
			if entry.Name() == "com" || entry.Name() == "META-INF" {
				os.RemoveAll(name + "/" + entry.Name())
			}
			recursiveRemoveEmptyDirs(name + "/" + entry.Name())
		}
	}
	entries, _ = os.ReadDir(name)

	if len(entries) == 0 {
		os.RemoveAll(name)
	}
}

func recursiveSearch2(tld string, name string) []string {
	removableItems := make([]string, 0)

	entries, _ := os.ReadDir(name)

	for _, entry := range entries {
		if entry.IsDir() {
			removableItems = append(removableItems, recursiveSearch2(tld, name+"/"+entry.Name())...)
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

func check2(name string, suff string) bool {
	return strings.HasSuffix(name, suff)
}

func recursiveSort(name string) {
	os.MkdirAll("src/main/java", fs.ModePerm)
	os.MkdirAll("src/main/resources", fs.ModePerm)

	recursiveSort1(name, name)
}

func recursiveSort1(tld string, name string) {
	entries, _ := os.ReadDir(name)
	for _, entry := range entries {
		if entry.IsDir() {
			recursiveSort1(tld, name+"/"+entry.Name())
		} else {
			if name == tld {
				if !check2(name+"/"+entry.Name(), ".java") {
					os.Rename(name+"/"+entry.Name(), strings.ReplaceAll(name+"/"+entry.Name(), tld+"/", "src/main/resources/"))
				}
			}
			if name != tld {
				if check2(name+"/"+entry.Name(), ".java") {
					fmt.Println(strings.ReplaceAll(name, tld+"/", "src/main/java/"))
					os.MkdirAll(strings.ReplaceAll(name, tld+"/", "src/main/java/"), fs.ModePerm)
					os.Rename(name+"/"+entry.Name(), strings.ReplaceAll(name+"/"+entry.Name(), tld+"/", "src/main/java/"))
				} else {
					fmt.Println(strings.ReplaceAll(name, tld+"/", "src/main/resources/"))
					os.MkdirAll(strings.ReplaceAll(name, tld+"/", "src/main/resources/"), fs.ModePerm)
					os.Rename(name+"/"+entry.Name(), strings.ReplaceAll(name+"/"+entry.Name(), tld+"/", "src/main/resources/"))
				}
			}
		}
	}
}

//go:embed lib
var lib embed.FS

func main() {
	version := "0.1.33"
	if len(os.Args) == 2 {
		version = os.Args[1]
	}

	fmt.Println("Downloading \"Cosmic Reach-" + version + ".jar\"")
	name := downloadFile("https://cosmic-archive.netlify.app/Cosmic%20Reach-" + version + ".jar")
	fmt.Println("Renaming \"" + name + "\" to \"cosmic-reach.jar\"")
	os.Rename(name, "cosmic-reach.jar")

	if _, err := os.Stat("quiltflower.jar"); errors.Is(err, os.ErrNotExist) {
		content, _ := lib.ReadFile("lib/quiltflower.jar")
		os.WriteFile("quiltflower.jar", content, fs.ModePerm)
	}

	cmd := exec.Command("java", "-jar", "quiltflower.jar", "cosmic-reach.jar", "test")
	cmd.Stdout = os.Stdout
	cmd.Run()

	for _, s := range recursiveSearch("test") {
		fmt.Println(s)
		os.RemoveAll(s)
	}
	recursiveRemoveEmptyDirs("test")
	recursiveCleanDir("test")
	recursiveSort("test")
	os.RemoveAll("test")
	os.RemoveAll("cosmic-reach.jar")
	os.RemoveAll("quiltflower.jar")

}
