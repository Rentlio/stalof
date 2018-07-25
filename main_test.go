package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestHandleLogLine(t *testing.T) {
	// First test line that will not be able to unmarshal
	output := captureStdOutAndStdErr(func() {
		handleLogLine("This will not unmarshal")
	})

	expectedError := "Log entry unmarshaling error: invalid character 'T' looking for beginning of value\n"
	if output != expectedError {
		t.Errorf("Incorrect output, got: %s, want: %s.", output, expectedError)
	}

	// Test line that will unmarshal. We are trimming newline when formming json
	// Output adds newline to stdOut
	expectedOutput := "This is some logged text\n"
	output = captureStdOutAndStdErr(func() {
		handleLogLine(fmt.Sprintf(`{"textPayload":"%s"}`, strings.TrimRight(expectedOutput, "\n")))
	})

	if output != expectedOutput {
		t.Errorf("Incorrect output, got: %s, want: %s.", output, expectedOutput)
	}

	// Test empty line
	output = captureStdOutAndStdErr(func() {
		handleLogLine("")
	})

	if output != "" {
		t.Error("Expected empty output")
	}
}

func TestHandleLogLines(t *testing.T) {
	lines := []string{
		`{"textPayload":"firstLine"}`,
		`{"textPayload":"secondLine"}`,
		`Unable to marshal`,
	}
	output := captureStdOutAndStdErr(func() {
		handleLogLines(lines)
	})

	if !strings.Contains(output, "firstLine") {
		t.Errorf("Incorrect output, got: %s, want: %s.", output, "firstLine")
	}

	if !strings.Contains(output, "secondLine") {
		t.Errorf("Incorrect output, got: %s, want: %s.", output, "secondLine")
	}

	if !strings.Contains(output, "Log entry unmarshaling error") {
		t.Errorf("Incorrect output, got: %s, want: %s.", output, "Log entry unmarshaling error")
	}
}

func captureStdOutAndStdErr(f func()) string {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	stdout := os.Stdout
	os.Stdout = w
	defer func() {
		os.Stdout = stdout
	}()

	stderr := os.Stderr
	os.Stderr = w
	defer func() {
		os.Stderr = stderr
	}()

	f()
	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
}
