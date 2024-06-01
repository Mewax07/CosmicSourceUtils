package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"strings"
)

func mkPatchRecurse(name string, srcDir string) {
	entries, _ := os.ReadDir(name)

	for _, entry := range entries {
		FName := name + "/" + entry.Name()
		if entry.IsDir() {
			os.MkdirAll(strings.ReplaceAll(FName, srcDir+"/", "patches/"), fs.ModePerm)
			mkPatchRecurse(FName, srcDir)
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

func createPatches(version string) {
	modified_source := "modifiedsources/Cosmic_Reach-" + version + "-Source.jar"
	if _, err := os.Stat(modified_source); errors.Is(err, os.ErrNotExist) {
		modified_source = "modifiedsources/Cosmic Reach-" + version + "-Source.jar"
		if _, err := os.Stat(modified_source); errors.Is(err, os.ErrNotExist) {
			log.Fatalf("Could not find any modified sources for cosmic reach %s which should be located in \"./modifiedsources\" as Cosmic Reach-%s-Source.jar",
				version,
				version,
			)
		}
	}
	unzipSource(modified_source, "temp_ext")
	recursiveSort("temp_ext", "temp_src_mod")
	os.RemoveAll("temp_ext")

	mkPatchRecurse("temp_src_mod", "temp_src_mod")
	os.RemoveAll("temp_src_mod")
	deleteEmptyDirs("patches")
}
