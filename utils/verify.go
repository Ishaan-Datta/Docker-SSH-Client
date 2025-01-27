package utils

// json file present
// docker desktop running
// remote context active

// return err message after completing checks

import (
	"os"
	"os/exec"
	"strings"
)

// IsConfigFilePresent checks if config.json is present in the directory
func IsConfigFilePresent(directory string) bool {
    _, err := os.Stat(directory + "/config.json")
    return !os.IsNotExist(err)
}

// IsDockerEngineRunning checks if Docker engine is running on the system
func IsDockerEngineRunning() bool {
    cmd := exec.Command("docker", "info")
    err := cmd.Run()
    return err == nil
}

// IsDockerRemoteContextActive checks if Docker remote context is configured
func IsDockerRemoteContextActive() bool {
    cmd := exec.Command("docker", "context", "ls", "--format", "{{.Current}}")
    output, err := cmd.Output()
    if err != nil {
        return false
    }
    return strings.TrimSpace(string(output)) != ""
}