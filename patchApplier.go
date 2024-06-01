package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"
)

func applyPatches(srcDir string, patchesDir string) {
	applyPatchesRecurse(srcDir, patchesDir)
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
			originalFile := strings.ReplaceAll(FName, srcDir+"/", "src/")
			fmt.Printf("git diff --no-index --default-prefix -u --output \"%s\" \"%s\" \"%s\"\n",
				patchFile,
				originalFile,
				FName, // modified
			)
			exec.Command("git", "diff", "--no-index", "--default-prefix", "-u", "--output",
				patchFile,
				originalFile,
				FName, // modified
			).Run()

			f, _ := os.ReadFile(patchFile)
			if len(f) == 0 {
				os.RemoveAll(patchFile)
			}
		}
	}
}
