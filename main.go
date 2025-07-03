package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/briandowns/spinner"
	"github.com/common-nighthawk/go-figure"
	"github.com/manifoldco/promptui"
)

const (
	cliVersion = "0.0.1"
	doaVersion = "DOA S1"
	repoURL    = "https://github.com/DavidNzube101/DOA-Local"
	cloneDir   = "DOA-Local"
)

func main() {
	myFigure := figure.NewFigure("DOA Launcher", "ogre", true)
	myFigure.Print()
	fmt.Println("CLI Launcher for your favourite web3 pvp game Daughter Of Aether S1")
	fmt.Println()

	for {
		prompt := promptui.Select{
			Label: "Select an option",
			Items: []string{
				"Download & Install DOA",
				"Launch DOA",
				"Update DOA",
				"Check CLI version",
				"Check DOA version",
				"Exit",
			},
		}

		_, result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		switch result {
		case "Download & Install DOA":
			downloadAndInstall()
		case "Launch DOA":
			launchDOA()
		case "Update DOA":
			updateDOA()
		case "Check CLI version":
			fmt.Printf("DOA-Launcher version: %s\n", cliVersion)
		case "Check DOA version":
			fmt.Printf("DOA version: %s\n", doaVersion)
		case "Exit":
			fmt.Println("Exiting DOA-Launcher.")
			return
		}
	}
}

func downloadAndInstall() {
	fmt.Println("Downloading and installing DOA...")

	if err := runCommandWithSpinner("Cloning DOA-Local repository...", "git", "clone", repoURL); err != nil {
		fmt.Println("\nError cloning repository:", err)
		return
	}

	if err := os.Chdir(cloneDir); err != nil {
		fmt.Println("\nError changing directory:", err)
		return
	}
	defer os.Chdir("..")

	_, err := exec.LookPath("node")
	if err != nil {
		fmt.Println("\nNode.js is not installed. Please install it and run DOA-Launcher again.")
		return
	}

	if err := runCommandWithSpinner("Installing pnpm...", "npm", "install", "-g", "pnpm"); err != nil {
		fmt.Println("\nError installing pnpm:", err)
		return
	}

	if err := runCommandWithSpinner("Installing dependencies with pnpm...", "pnpm", "install"); err != nil {
		fmt.Println("\nError installing dependencies:", err)
		return
	}

	fmt.Println("\nDOA has been downloaded and installed successfully.")
}

func launchDOA() {
	if !isDOAInstalled() {
		fmt.Println("DOA is not installed. Please select 'Download & Install DOA' first.")
		return
	}

	fmt.Println("Launching DOA...")

	if err := os.Chdir(cloneDir); err != nil {
		fmt.Println("\nError changing directory:", err)
		return
	}
	defer os.Chdir("..")

	if err := runCommandWithSpinner("Building the project...", "pnpm", "build"); err != nil {
		fmt.Println("\nError building project:", err)
		return
	}

	fmt.Println("Starting the preview server...")
	go func() {
		if err := runCommandWithSpinner("Starting preview server...", "pnpm", "preview"); err != nil {
			fmt.Println("\nError starting preview server:", err)
		}
	}()

	openBrowser("http://localhost:3000")

	fmt.Println("\nDOA is running. You can stop the server with 'Ctrl+C' in the terminal where you launched DOA-Launcher.")
}

func updateDOA() {
	if !isDOAInstalled() {
		fmt.Println("DOA is not installed. Please select 'Download & Install DOA' first.")
		return
	}

	fmt.Println("Updating DOA...")

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Prefix = "Removing old DOA directory... "
	s.Start()
	if err := os.RemoveAll(cloneDir); err != nil {
		s.Stop()
		fmt.Printf("\nError removing old directory: %v\n", err)
		return
	}
	s.Stop()

	downloadAndInstall()
	launchDOA()

	fmt.Println("\nDOA has been updated and launched successfully.")
}

func isDOAInstalled() bool {
	_, err := os.Stat(cloneDir)
	return !os.IsNotExist(err)
}

func runCommandWithSpinner(message string, name string, args ...string) error {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner
	s.Prefix = message
	s.Start()
	defer s.Stop()

	cmd := exec.Command(name, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil

	return cmd.Run()
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		fmt.Println("\nError opening browser:", err)
	}
}