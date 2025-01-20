package utils

// func getDockerClient() (*client.Client, error) {
// 	return client.NewClientWithOpts(client.FromEnv)
// }

// func getContainerIP(containerName string) string {
// 	cli, err := getDockerClient()
// 	if err != nil {
// 		fmt.Println("Error creating Docker client:", err)
// 		return ""
// 	}

// 	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
// 	if err != nil {
// 		fmt.Println("Error listing containers:", err)
// 		return ""
// 	}

// 	for _, container := range containers {
// 		for _, name := range container.Names {
// 			if name == "/"+containerName {
// 				inspect, err := cli.ContainerInspect(context.Background(), container.ID)
// 				if err != nil {
// 					fmt.Println("Error inspecting container:", err)
// 					return ""
// 				}
// 				return inspect.NetworkSettings.IPAddress
// 			}
// 		}
// 	}
// 	return ""
// }

// func getDockerClient() (*client.Client, error) {
// 	return client.NewClientWithOpts(client.FromEnv)
// }

// func listRunningContainers() []string {
// 	cli, err := getDockerClient()
// 	if err != nil {
// 		fmt.Println("Error creating Docker client:", err)
// 		return nil
// 	}

// 	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
// 	if err != nil {
// 		fmt.Println("Error listing containers:", err)
// 		return nil
// 	}

// 	containerNames := []string{}
// 	for _, container := range containers {
// 		name := strings.TrimPrefix(container.Names[0], "/")
// 		containerNames = append(containerNames, name)
// 	}
// 	return containerNames
// }

// func getContainerIP(containerName string) string {
// 	cli, err := getDockerClient()
// 	if err != nil {
// 		fmt.Println("Error creating Docker client:", err)
// 		return ""
// 	}

// 	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
// 	if err != nil {
// 		fmt.Println("Error listing containers:", err)
// 		return ""
// 	}

// 	for _, container := range containers {
// 		for _, name := range container.Names {
// 			if strings.TrimPrefix(name, "/") == containerName {
// 				inspect, err := cli.ContainerInspect(context.Background(), container.ID)
// 				if err != nil {
// 					fmt.Println("Error inspecting container:", err)
// 					return ""
// 				}
// 				return inspect.NetworkSettings.IPAddress
// 			}
// 		}
// 	}
// 	return ""
// }