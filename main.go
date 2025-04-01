package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

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

type Process struct {
	Cmd        *exec.Cmd
	StdOutPipe io.ReadCloser
	StdErrPipe io.ReadCloser
	LogBuffer  []string
}

var config Config
var ProcessDict = make(map[string]ProcessConfig)
var processes = make(map[string]*Process)

func main() {
	loadConfig()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Choose an option:")
		fmt.Println("1. List processes")
		fmt.Println("2. Start process")
		fmt.Println("3. Stop process")
		fmt.Println("4. Get logs")
		fmt.Println("5. Exit")

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
			printLogs()
		case "5":
			os.Exit(0)
		default:
			fmt.Println("Invalid choice. Please enter a number between 1 and 5.")
		}
	}
}

func printLogs() {
	fmt.Println("Enter process name to get logs:")
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	p := processes[choice]

	if len(p.LogBuffer) == 0 {
		fmt.Println("No logs available")
	}
	for _, log := range p.LogBuffer {
		fmt.Print(log)
	}
}

func stopProcess() {
	fmt.Println("Enter process name to kill:")
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	p, exists := processes[choice]
	if !exists {
		fmt.Printf("%s process is not running\n", choice)
	}
	p.Cmd.Process.Kill()
	fmt.Printf("%s process stopped", choice)
}

func startProcess() {
	fmt.Println("Enter process name to start:")
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	process := ProcessDict[choice]

	if _, exists := processes[process.Name]; exists {
		fmt.Printf("%s process is already running", process.Name)
		return
	}

	cmd := exec.Command(process.Command, process.Args...)

	stdOut, _ := cmd.StdoutPipe()
	stdErr, _ := cmd.StderrPipe()

	p := &Process{
		Cmd:        cmd,
		StdOutPipe: stdOut,
		StdErrPipe: stdErr,
	}

	processes[process.Name] = p

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting command: %v\n", err)
		return
	}

	go captureLogs(p, stdOut, "STDOUT")
	go captureLogs(p, stdErr, "STDERR")
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
	for _, process := range config.Processes {
		ProcessDict[process.Name] = process
	}
}

func captureLogs(process *Process, pipe io.ReadCloser, logLvl string) {
	reader := bufio.NewReader(pipe)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		logLine := fmt.Sprintf("[%s] %s", logLvl, line)
		process.LogBuffer = append(process.LogBuffer, logLine)
	}
}
