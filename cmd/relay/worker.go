package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func startWorker() {
	for {
		processQueue()
		time.Sleep(5 * time.Second)
	}
}

func processQueue() {
	files, err := filepath.Glob(filepath.Join(queueDir, "*.eml"))
	if err != nil {
		fmt.Println("queue scan error:", err)
		return
	}

	for _, f := range files {
		processMessage(f)
	}
}

func processMessage(path string) {

	// lock message to prevent double processing
	lockedPath := path + ".sending"

	err := os.Rename(path, lockedPath)
	if err != nil {
		return
	}

	fmt.Println("processing message:", lockedPath)

	err = forwardMessage(lockedPath)
	if err != nil {
		fmt.Println("send error:", err)
		return
	}

	err = os.Remove(lockedPath)
	if err != nil {
		fmt.Println("delete error:", err)
		return
	}

	fmt.Println("processed and removed:", lockedPath)
}