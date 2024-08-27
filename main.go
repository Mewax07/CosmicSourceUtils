package main

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"
)

func check(arr []string, name string, suff string) []string {
	if strings.HasSuffix(name, suff) {
		arr = append(arr, name)
	}
	return arr
}

func check2(name string, suff string) bool {
	return strings.HasSuffix(name, suff)
}

//go:embed lib
var lib embed.FS

func findQuiltflower() string {
	if _, err := os.Stat("quiltflower.jar"); errors.Is(err, os.ErrNotExist) {
		embedQuiltflower := "embededQuiltflower.jar"
		if _, err := os.Stat(embedQuiltflower); errors.Is(err, os.ErrNotExist) {
			content, _ := lib.ReadFile("lib/quiltflower.jar")
			os.WriteFile(embedQuiltflower, content, fs.ModePerm)
		}
		return embedQuiltflower
	}
	return "quiltflower.jar"
}

func createSources(version string, srcDir string) {
	fmt.Println("Downloading \"Cosmic Reach-" + version + ".jar\"")
	name := downloadFile("https://cosmic-archive.netlify.app/versions/pre-alpha/Cosmic%20Reach-" + version + ".jar")
	fmt.Println("Renaming \"" + name + "\" to \"cosmic-reach.jar\"")
	os.Rename(name, "cosmic-reach.jar")

	qf := findQuiltflower()

	// Run QuiltFlower
	cmd := exec.Command("java", "-jar", qf, "cosmic-reach.jar", "cr_temp")
	cmd.Stdout = os.Stdout
	cmd.Run()

	recursiveCleanUnNeededFiles("cr_temp")
	deleteEmptyDirs("cr_temp")
	recursiveMicroPatchDir("cr_temp")
	recursiveSort("cr_temp", srcDir)
	os.RemoveAll("cr_temp")
	os.RemoveAll("cosmic-reach.jar")
	os.RemoveAll("embededQuiltflower.jar")
}

func printHelp() {
	fmt.Println("execution format must be:")
	fmt.Println("\tCosmicSRCUtil.exe {flag} {version}")
	fmt.Println("\nflags:")
	fmt.Println("\t\"-d\" Download and Prepare cr {version}, ex: CosmicSRCUtil.exe -d 0.1.33\n")
	fmt.Println("\t\"-c\" Create patches for cr {version} using a file at \"./modifedsources/Cosmic Reach-{version}-Sources.jar\"\n\tand creates the patches in \"./patches\", ex: CosmicSRCUtil.exe -c 0.1.33\n")
	fmt.Println("\t\"-p\" Applies the patches from \"./patches\" onto the \"./src\" dir\n\t ex: CosmicSRCUtil.exe -p 0.1.33\n")
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "-d":
		if len(os.Args) < 3 {
			printHelp()
			os.Exit(0)
		}
		createSources(os.Args[2], "src")
	case "-c":
		if len(os.Args) < 3 {
			printHelp()
			os.Exit(0)
		}
		createSources(os.Args[2], "temp_src")
		createPatches(os.Args[2])
		os.RemoveAll("temp_src")
	case "-p":
		if _, err := os.Stat("patches"); errors.Is(err, os.ErrNotExist) {
			fmt.Println("You must have a dir named \"patches\" containing all the patches.")
			os.Exit(1)
		}
		if _, err := os.Stat("src"); errors.Is(err, os.ErrNotExist) {
			fmt.Println("You must have a dir named \"src\" containing all the generated for create source code from cosmic reach.")
			os.Exit(1)
		}
		applyPatches("src", "patches")
	default:
		printHelp()
	}

}
