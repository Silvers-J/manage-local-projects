package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	// "os/exec"

	"gopkg.in/yaml.v2"
)

type ProcessConfig struct {
	Name    string   `yaml:"name"`
	Command string   `yaml:"command"`
	Args    []string `yaml:"args"`
}

type Config struct {
	Processes []ProcessConfig `yaml:"processes"`
}

var config Config

func main() {
	loadConfig()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Choose an option:")
		fmt.Println("1. List processes")
		fmt.Println("2. Start process")
		fmt.Println("3. Stop process")
		fmt.Println("4. Exit")

		choice, _ := reader.ReadString('\n')

		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			listProcesses()
		case "2":
			startProcess()
		case "3":
			stopProcess()
		case "4":
			os.Exit(0)
		default:
			fmt.Println("Invalid choice. Please enter a number between 1 and 5.")
		}
	}
	// cmd := exec.Command("cat", "main.go")
	//
	// stdOut, _ := cmd.StdoutPipe()
	// stdErr, _ := cmd.StderrPipe()
	//
	// if err := cmd.Start(); err != nil {
	// 	fmt.Printf("Error starting command: %v\n", err)
	// 	return
	// }
	//
	// printLogs(stdOut, "STDOUT")
	// printLogs(stdErr, "STDERR")
}

func stopProcess() {
	panic("unimplemented")
}

func startProcess() {
	panic("unimplemented")
}

func listProcesses() {
	fmt.Println("Available processes:")
	for _, p := range config.Processes {
		fmt.Printf("- %s: %s %v\n", p.Name, p.Command, p.Args)
	}
}

func loadConfig() {
	file, err := os.ReadFile("configs/config.yaml")
	if err != nil {
		fmt.Printf("Failed to read config file: %v\n", err)
		os.Exit(1)
	}

	if err := yaml.Unmarshal(file, &config); err != nil {
		fmt.Printf("Failed to parse config file: %v\n", err)
		os.Exit(1)
	}
}

func printLogs(pipe io.ReadCloser, logLvl string) {
	buffer := make([]byte, 1024)
	for {
		n, err := pipe.Read(buffer)
		if n > 0 {
			fmt.Printf("[%s] %s", logLvl, buffer[:n])
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("[%s] Error reading logs: %v\n", logLvl, err)
			break
		}
	}
}
