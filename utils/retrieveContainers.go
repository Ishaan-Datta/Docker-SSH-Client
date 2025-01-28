package utils

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// Define ContainerDetails and ContainerList
type ContainerDetails struct {
	PortMappings string
	CreatedTime  string
	Names        string
	ID           string
	Image        string
	Status       string
	Command      string
}

type EnvironmentVariables struct {
	DockerHost string
	DockerCertPath string
	DockerTLSVerify string
}


// attribute container list for each config so you know how to access it later, need type to store according config
type ContainerList struct {
	Configuration Config
	Containers []ContainerDetails
}

func RetrieveAllContainers(configurations []Config) ([]ContainerList, error) {
	originalVars, err := RetrieveEnvironmentVariables()

	if err != nil {
		return []ContainerList{}, err
	}

	// 
	var list []ContainerList

	for _, selectedConfig := range configurations{
		// Set environment variables for Docker client
		var clientVars EnvironmentVariables 
		clientVars.DockerHost = selectedConfig.Host
		clientVars.DockerCertPath = selectedConfig.CertPath
		clientVars.DockerTLSVerify = selectedConfig.TLSVerify
		SetEnvironmentVariables(clientVars)

		// Create Docker client using environment variables
		apiClient, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			return nil, fmt.Errorf("error creating Docker client: %w", err)
		}
		defer apiClient.Close()

		// Retrieve container list
		containers, err := apiClient.ContainerList(context.Background(), container.ListOptions{All: true})
		if err != nil {
			return nil, fmt.Errorf("error listing containers: %w", err)
		}

		var configContainerList ContainerList
		var configContainers []ContainerDetails

		for _, ctr := range containers {
			details := ContainerDetails{
				PortMappings: formatPorts(ctr.Ports),
				CreatedTime:  time.Unix(ctr.Created, 0).Format(time.RFC3339),
				Names:        formatContainerNames(ctr.Names),
				ID:           ctr.ID[:8],
				Image:        ctr.Image,
				Status:       ctr.Status,
				Command:      ctr.Command,
			}
			configContainers = append(configContainers, details)
		}

		configContainerList.Configuration = selectedConfig
		configContainerList.Containers = configContainers
		list = append(list, configContainerList)
	}

	err = SetEnvironmentVariables(originalVars)

	if err != nil {
		return []ContainerList{}, err
	}

	return list, nil
}

func RetrieveEnvironmentVariables() (EnvironmentVariables, error){
	// check the details of the Getenv function, for difference between unset and empty, should clarify difference
	// Store the current environment variables
	var originalVars EnvironmentVariables
	originalVars.DockerHost = os.Getenv("DOCKER_HOST")
	originalVars.DockerCertPath = os.Getenv("DOCKER_CERT_PATH")
	originalVars.DockerTLSVerify = os.Getenv("DOCKER_TLS_VERIFY")

	return originalVars, nil
}

// rename
func SetEnvironmentVariables(originalVars EnvironmentVariables) (error){
	// Restore the original environment variables
	if err := os.Setenv("DOCKER_HOST", originalVars.DockerHost); err != nil {
		return fmt.Errorf("error restoring DOCKER_HOST: %w", err)
	}
	if err := os.Setenv("DOCKER_CERT_PATH", originalVars.DockerCertPath); err != nil {
		return fmt.Errorf("error restoring DOCKER_CERT_PATH: %w", err)
	}
	if err := os.Setenv("DOCKER_TLS_VERIFY", originalVars.DockerTLSVerify); err != nil {
		return fmt.Errorf("error restoring DOCKER_TLS_VERIFY: %w", err)
	}
	return nil
}

// formatPorts converts Docker port mappings to a readable string
func formatPorts(ports []types.Port) string {
	if len(ports) == 0 {
		return "-"
	}

	portStr := ""
	for _, port := range ports {
		if port.PublicPort > 0 {
			portStr += fmt.Sprintf("%d->%d/%s, ", port.PublicPort, port.PrivatePort, port.Type)
		} else {
			portStr += fmt.Sprintf("%d/%s, ", port.PrivatePort, port.Type)
		}
	}
	return portStr[:len(portStr)-2] // Remove trailing comma and space
}

// formatContainerNames removes leading '/' from container names
func formatContainerNames(names []string) string {
	formattedNames := make([]string, len(names))
	for i, name := range names {
		formattedNames[i] = name[1:] // Remove leading '/'
	}
	return formatSlice(formattedNames)
}

// formatSlice converts a slice to a comma-separated string
func formatSlice(slice []string) string {
	if len(slice) == 0 {
		return "-"
	}
	result := ""
	for _, s := range slice {
		result += s + ", "
	}
	return result[:len(result)-2] // Remove trailing comma and space
}