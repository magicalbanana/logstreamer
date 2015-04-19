package logstreamer

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"
)

func TestLogstreamerOk(t *testing.T) {
	// Create a logger (your app probably already has one)
	logger := log.New(os.Stdout, "--> ", log.Ldate|log.Ltime)

	// Setup a streamer that we'll pipe cmd.Stdout to
	logStreamerOut := NewLogstreamer(logger, "stdout", false)
	// Setup a streamer that we'll pipe cmd.Stderr to.
	// We want to record/buffer anything that's written to this (3rd argument true)
	logStreamerErr := NewLogstreamer(logger, "stderr", true)

	// Execute something that succeeds
	cmd := exec.Command(
		"ls",
		"-al",
	)
	cmd.Stderr = logStreamerErr
	cmd.Stdout = logStreamerOut

	// Reset any error we recorded
	logStreamerErr.FlushRecord()

	// Execute command
	err := cmd.Start()

	// Failed to spawn?
	if err != nil {
		t.Fatal("ERROR could not spawn command.", err.Error())
	}

	// Failed to execute?
	err = cmd.Wait()
	if err != nil {
		t.Fatal("ERROR command finished with error. ", err.Error(), logStreamerErr.FlushRecord())
	}
}

func TestLogstreamerErr(t *testing.T) {
	// Create a logger (your app probably already has one)
	logger := log.New(os.Stdout, "--> ", log.Ldate|log.Ltime)

	// Setup a streamer that we'll pipe cmd.Stdout to
	logStreamerOut := NewLogstreamer(logger, "stdout", false)
	// Setup a streamer that we'll pipe cmd.Stderr to.
	// We want to record/buffer anything that's written to this (3rd argument true)
	logStreamerErr := NewLogstreamer(logger, "stderr", true)

	// Execute something that succeeds
	cmd := exec.Command(
		"ls",
		"nonexisting",
	)
	cmd.Stderr = logStreamerErr
	cmd.Stdout = logStreamerOut

	// Reset any error we recorded
	logStreamerErr.FlushRecord()

	// Execute command
	err := cmd.Start()

	// Failed to spawn?
	if err != nil {
		logger.Print("ERROR could not spawn command. ")
	}

	// Failed to execute?
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("Good. command finished with %s. %s. \n", err.Error(), logStreamerErr.FlushRecord())
	} else {
		t.Fatal("This command should have failed")
	}
}
