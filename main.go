package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: sshclient <container_name> <auth_method>")
		fmt.Println("auth_method: kerberos | oauth2 | saml")
		return
	}
	containerName := os.Args[1]
	authMethod := os.Args[2]

	var sshAuthMethod ssh.AuthMethod
	var err error

	switch authMethod {
	case "kerberos":
		sshAuthMethod, err = kerberosAuth()
	case "oauth2":
		sshAuthMethod, err = oauth2Auth()
	case "saml":
		sshAuthMethod, err = samlAuth()
	default:
		fmt.Println("Invalid authentication method")
		return
	}

	if err != nil {
		fmt.Println("Error in authentication:", err)
		return
	}

	// Initialize SSH connection with chosen auth method
	client, err := NewSSHClient(containerName, sshAuthMethod)
	if err != nil {
		fmt.Println("Error initializing SSH client:", err)
		return
	}
	defer client.Close()

	// Example command execution with pipelining
	output, err := client.ExecuteCommand("ls -l /app")
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}
	fmt.Println(output)
}
