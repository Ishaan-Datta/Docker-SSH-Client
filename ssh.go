package main

import (
	"fmt"
	"net"

	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	client *ssh.Client
}

func NewSSHClient(containerName string) (*SSHClient, error) {
	// Retrieve Docker container's IP address
	// containerIP := getContainerIP(containerName)
	containerIP := "bruh"
	if containerIP == "" {
		return nil, fmt.Errorf("could not find container: %s", containerName)
	}

	// SSH Configurations with multiplexing
	config := &ssh.ClientConfig{
		User:            "root", // assuming root access
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		// Add additional auth methods (password, public key, etc.)
		Auth: []ssh.AuthMethod{
			ssh.Password("yourpassword"),
			// Add more advanced auth methods here
		},
	}

	// Establish the SSH connection
	conn, err := ssh.Dial("tcp", net.JoinHostPort(containerIP, "22"), config)
	if err != nil {
		return nil, err
	}

	return &SSHClient{client: conn}, nil
}

func (c *SSHClient) ExecuteCommand(command string) (string, error) {
	session, err := c.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	// Implement pipelining here
	output, err := session.CombinedOutput(command)
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func (c *SSHClient) Close() {
	if c.client != nil {
		c.client.Close()
	}
}

// Mock function to retrieve the Docker container's IP address
// func getContainerIP(containerName string) string {
// 	// Implement Docker API calls to get the IP address
// 	return "172.17.0.2"
// }
