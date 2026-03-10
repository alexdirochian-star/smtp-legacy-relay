package main

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	listenAddr = "127.0.0.1:2525"
	queueDir   = "queue"
)

func main() {

	// ensure queue directory exists
	if err := os.MkdirAll(queueDir, 0755); err != nil {
		fmt.Printf("failed to create queue dir: %v\n", err)
		return
	}

	// start stranded-file reaper (reliability layer)
	maybeStartReaper(queueDir)

	// start background worker
	go startWorker()

	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Printf("listen error: %v\n", err)
		return
	}
	defer ln.Close()

	fmt.Printf("SMTP relay listening on %s\n", listenAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("accept error: %v\n", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

addr := conn.RemoteAddr().String()

if !strings.HasPrefix(addr, "127.0.0.1") {
	return

}

	// protect server from hung clients
	setSessionDeadline(conn)

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	writeLine(writer, "220 relay ready")

	var mailFrom string
	var rcptTo []string

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		line = strings.TrimRight(line, "\r\n")
		upper := strings.ToUpper(line)

		switch {

		case strings.HasPrefix(upper, "HELO "):
			writeLine(writer, "250 Hello")

		case strings.HasPrefix(upper, "EHLO "):
			writeLine(writer, "250-Hello")
			writeLine(writer, "250 SIZE 10485760")

		case strings.HasPrefix(upper, "MAIL FROM:"):
			mailFrom = strings.TrimSpace(line[len("MAIL FROM:"):])
			rcptTo = nil
			writeLine(writer, "250 OK")

		case strings.HasPrefix(upper, "RCPT TO:"):
			rcpt := strings.TrimSpace(line[len("RCPT TO:"):])
			rcptTo = append(rcptTo, rcpt)
			writeLine(writer, "250 OK")

		case upper == "DATA":

			if mailFrom == "" || len(rcptTo) == 0 {
				writeLine(writer, "503 Bad sequence of commands")
				continue
			}

			writeLine(writer, "354 End data with <CR><LF>.<CR><LF>")

			data, err := readData(reader)
			if err != nil {
				writeLine(writer, "451 Requested action aborted")
				return
			}

			msg := buildMessage(mailFrom, rcptTo, data)

			filename, err := saveToSpool(msg)
			if err != nil {
				writeLine(writer, "451 Requested action aborted")
				return
			}

			fmt.Printf("queued message: %s\n", filename)

			writeLine(writer, "250 Message accepted")

			mailFrom = ""
			rcptTo = nil

		case upper == "RSET":
			mailFrom = ""
			rcptTo = nil
			writeLine(writer, "250 OK")

		case upper == "NOOP":
			writeLine(writer, "250 OK")

		case upper == "QUIT":
			writeLine(writer, "221 Bye")
			return

		default:
			writeLine(writer, "500 Command not recognized")
		}
	}
}

func writeLine(w *bufio.Writer, s string) {
	w.WriteString(s + "\r\n")
	w.Flush()
}

func readData(reader *bufio.Reader) (string, error) {
	var sb strings.Builder

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		if line == ".\r\n" || line == ".\n" {
			break
		}

		if strings.HasPrefix(line, "..") {
			line = line[1:]
		}

		sb.WriteString(line)
	}

	return sb.String(), nil
}

func buildMessage(mailFrom string, rcptTo []string, body string) string {

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("X-Relay-Received-At: %s\r\n", time.Now().UTC().Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("X-Relay-MailFrom: %s\r\n", mailFrom))
	sb.WriteString(fmt.Sprintf("X-Relay-RcptTo: %s\r\n", strings.Join(rcptTo, ", ")))
	sb.WriteString("\r\n")

	sb.WriteString(body)

	return sb.String()
}

func saveToSpool(message string) (string, error) {

	id, err := newMessageID()
	if err != nil {
		return "", err
	}

	tmpName := filepath.Join(queueDir, id+".tmp")
	finalName := filepath.Join(queueDir, id+".eml")

	f, err := os.OpenFile(tmpName, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
	if err != nil {
		return "", err
	}

	_, writeErr := f.WriteString(message)
	closeErr := f.Close()

	if writeErr != nil {
		os.Remove(tmpName)
		return "", writeErr
	}

	if closeErr != nil {
		os.Remove(tmpName)
		return "", closeErr
	}

	if err := os.Rename(tmpName, finalName); err != nil {
		os.Remove(tmpName)
		return "", err
	}

	return finalName, nil
}

func newMessageID() (string, error) {

	now := time.Now().UTC().Format("20060102T150405.000000000")

	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	return now + "-" + hex.EncodeToString(buf), nil
}



