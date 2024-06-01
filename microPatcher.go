package main

import (
	"fmt"
	"io/fs"
	"math/big"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func recursiveMicroPatchDir(name string) {
	entries, _ := os.ReadDir(name)

	for _, entry := range entries {
		if entry.IsDir() {
			recursiveMicroPatchDir(name + "/" + entry.Name())
		} else {
			microFilePatcher(name + "/" + entry.Name())
		}
	}
}

func microFilePatcher(name string) {
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
