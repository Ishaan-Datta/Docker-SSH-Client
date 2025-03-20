package utils

// use file to select it, should check file system on container, note if there or not w/ prompt to override before cont
// progress animated for download
// check docs on the action to see if sends progress or anything

import (
	"context"
	"fmt"
	"io"
	"os"

	// "path/filepath"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var p *tea.Program

type progressTracker struct {
	total      int64
	transferred int64
	onProgress func(float64)
}

type model struct {
	progress progress.Model
}

type progressMsg float64

func (m model) Init() tea.Cmd {
	return nil
}

func (pt *progressTracker) Write(p []byte) (int, error) {
	n := len(p)
	pt.transferred += int64(n)
	if pt.total > 0 && pt.onProgress != nil {
		pt.onProgress(float64(pt.transferred) / float64(pt.total))
	}
	return n, nil
}

func uploadToContainer(cli *client.Client, containerID, srcPath, destPath string) error {
	file, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	tracker := &progressTracker{
		total: fileInfo.Size(),
		onProgress: func(ratio float64) {
			p.Send(progressMsg(ratio))
		},
	}

	reader := io.TeeReader(file, tracker)
	ctx := context.Background()
	err = cli.CopyToContainer(ctx, containerID, destPath, reader, types.CopyToContainerOptions{})
	fmt.Println()
	return err
}

// func pog() {
// 	containerID := flag.String("container", "", "ID of the Docker container")
// 	srcPath := flag.String("src", "", "Path to the local file to upload")
// 	destPath := flag.String("dest", "", "Destination path inside the container")
// 	flag.Parse()

// 	if *containerID == "" || *srcPath == "" || *destPath == "" {
// 		flag.Usage()
// 		os.Exit(1)
// 	}

// 	cli, err := client.NewClientWithOpts(client.FromEnv)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer cli.Close()

// 	m := model{
// 		progress: progress.New(progress.WithDefaultGradient()),
// 	}
// 	p = tea.NewProgram(m)

// 	go func() {
// 		if err := uploadToContainer(cli, *containerID, *srcPath, *destPath); err != nil {
// 			log.Fatal(err)
// 		}
// 	}()

// 	if _, err := p.Run(); err != nil {
// 		log.Fatal(err)
// 	}
// }
