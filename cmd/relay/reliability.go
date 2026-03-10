package main

import (
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	strandedThreshold = 10 * time.Minute
	reaperInterval    = 5 * time.Minute
	sessionTimeout    = 5 * time.Minute
)

// reapStranded finds abandoned *.sending files and returns them to *.eml
// so the worker can retry them.
func reapStranded(queueDir string) {
	entries, err := os.ReadDir(queueDir)
	if err != nil {
		log.Printf("reaper: cannot read queue dir %s: %v", queueDir, err)
		return
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sending") {
			continue
		}

		info, err := e.Info()
		if err != nil {
			continue
		}

		age := time.Since(info.ModTime())
		if age < strandedThreshold {
			continue
		}

		oldPath := filepath.Join(queueDir, e.Name())
		newPath := strings.TrimSuffix(oldPath, ".sending")

		if err := os.Rename(oldPath, newPath); err != nil {
			log.Printf("reaper: failed to recover %s: %v", e.Name(), err)
			continue
		}

		log.Printf("reaper: recovered stranded file %s (age %s)", e.Name(), age.Round(time.Second))
	}
}

// startReaper runs the stranded-file recovery immediately and then periodically.
func startReaper(queueDir string) {
	reapStranded(queueDir)

	ticker := time.NewTicker(reaperInterval)
	defer ticker.Stop()

	for range ticker.C {
		reapStranded(queueDir)
	}
}

// logQueueState logs the current queue depth.
func logQueueState(queueDir string) {
	entries, err := os.ReadDir(queueDir)
	if err != nil {
		log.Printf("queue state: cannot read queue dir %s: %v", queueDir, err)
		return
	}

	var pending, sending int
	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		switch {
		case strings.HasSuffix(e.Name(), ".eml"):
			pending++
		case strings.HasSuffix(e.Name(), ".sending"):
			sending++
		}
	}

	log.Printf("queue: pending=%d sending=%d", pending, sending)
}

// maybeStartReaper should be called once before the worker loop starts.
func maybeStartReaper(queueDir string) {
	go startReaper(queueDir)
}

// beforeWorkerCycle should be called once per worker scan cycle.
func beforeWorkerCycle(queueDir string) {
	logQueueState(queueDir)
}

// setSessionDeadline prevents hung SMTP connections.
func setSessionDeadline(conn net.Conn) {
	if err := conn.SetDeadline(time.Now().Add(sessionTimeout)); err != nil {
		log.Printf("smtp: failed to set deadline: %v", err)
	}
}

// refreshSessionDeadline can be called after each client command if needed.
func refreshSessionDeadline(conn net.Conn) {
	setSessionDeadline(conn)
}