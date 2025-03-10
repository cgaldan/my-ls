package utils

import (
	"fmt"
	"strconv"
	"strings"
	"syscall"
)

func GetTerminalWidth() (int, error) {
	// Execute "stty size" to get terminal dimensions
	cmd := "/bin/stty"
	args := []string{"stty", "size"}

	// Create a pipe to capture output
	var pipe [2]int
	err := syscall.Pipe(pipe[:])
	if err != nil {
		return 0, err
	}
	read, write := pipe[0], pipe[1]

	// Fork and execute the command
	pid, err := syscall.ForkExec(cmd, args, &syscall.ProcAttr{
		Files: []uintptr{0, uintptr(write), uintptr(write)}, // Redirect stdout and stderr to the pipe
	})
	if err != nil {
		return 0, err
	}

	// Close the write end of the pipe
	syscall.Close(write)

	// Read output from the read end of the pipe
	buf := make([]byte, 32)
	n, err := syscall.Read(read, buf)
	if err != nil {
		return 0, err
	}

	// Close the read end of the pipe
	syscall.Close(read)

	// Wait for the process to finish
	_, err = syscall.Wait4(pid, nil, 0, nil)
	if err != nil {
		return 0, err
	}

	// Parse output
	output := strings.TrimSpace(string(buf[:n]))
	parts := strings.Fields(output)
	if len(parts) < 2 {
		return 0, fmt.Errorf("unexpected output format: %s", output)
	}

	// Convert the second value (columns) to an integer
	width, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}

	return width, nil
}
