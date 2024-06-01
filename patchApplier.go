package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"
)

func applyPatches(srcDir string, patchesDir string) {
	applyPatchesRecurse(srcDir, srcDir)
}

func applyPatchesRecurse(name string, srcDir string) {
	entries, _ := os.ReadDir(name)

	for _, entry := range entries {
		FName := name + "/" + entry.Name()
		if entry.IsDir() {
			os.MkdirAll(strings.ReplaceAll(FName, srcDir+"/", "patches/"), fs.ModePerm)
			applyPatchesRecurse(FName, srcDir)
		} else {
			patchFile := strings.ReplaceAll(FName, srcDir+"/", "patches/") + ".patch"
			if _, err := os.Stat(patchFile); !errors.Is(err, os.ErrNotExist) {
				fmt.Printf("git apply --ignore-whitespace \"%s\"\n",
					patchFile,
				)
				CMD := exec.Command("git", "apply", "--ignore-whitespace",
					patchFile,
				)
				CMD.Stdout = os.Stdout
				CMD.Stdin = os.Stdin
				CMD.Run()
			}
		}
	}
}
