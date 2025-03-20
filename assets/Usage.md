#### Listing Running Containers
```sh
docker-ssh-client list
```

#### Running Commands Inside a Container
```sh
docker-ssh-client exec -c <container_id> -- "ls -al"
```

#### Copying Files to a Container
```sh
docker-ssh-client cp local_file.txt <container_id>:/path/in/container
```

#### Copying Files from a Container
```sh
docker-ssh-client cp <container_id>:/path/in/container local_file.txt
```

#### Running a Script on Multiple Containers
```sh
docker-ssh-client run-script -f script.sh -c <container_id1> <container_id2>
```

### **1️⃣ Listing Running Containers**
```sh
go run client.go list
```
🟢 Returns:
```json
{
  "container_ids": ["123abc", "456def"]
}
```

### **2️⃣ Pushing a File to a Container**
```sh
go run client.go push --container=123abc --file=test.txt
```
🟢 Returns **real-time file transfer updates**
```json
{
  "status": "uploading",
  "transferred_bytes": 1024
}
```

### **3️⃣ Pulling a File from a Container**
```sh
go run client.go pull --container=123abc --file=/tmp/test.txt
```
🟢 Returns **real-time download progress**
```json
{
  "status": "downloading",
  "transferred_bytes": 2048
}
```

### **4️⃣ Executing a Command in a Container**
```sh
go run client.go exec --container=123abc --command="ls -la"
```
🟢 Returns **streaming command output**
```json
{
  "status": "running",
  "output": "total 4\ndrwxr-xr-x 2 root root 4096 Feb 10 10:00 ."
}
```

- **List running containers**
- **Push a file to a container**
- **Pull a file from a container**
- **Execute a command inside a container** with real-time streaming output