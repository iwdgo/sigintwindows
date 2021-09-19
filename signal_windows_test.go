//go:build windows
// +build windows

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

const processToSignal = "ctrlbreak"

func TestSendCtrlBreak(t *testing.T) {
	_ = os.Remove(processToSignal + ".log")
	src := processToSignal + ".go"
	exe := processToSignal + ".exe"
	defer func() {
		if err := os.Remove(exe); err != nil {
			t.Fatal(err)
		}
	}()
	o, err := exec.Command("go", "build", "-o", exe, filepath.Join(processToSignal, src)).CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to compile: %v\n%v", err, string(o))
	}
	cmd := exec.Command(exe)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
	err = cmd.Start()
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	// If interrupted here while cmd waits, "exit status STATUS_CONTROL_C_EXIT" displays
	// See go/src/os/exec_posix.go:109
	d := time.Duration(5)
	t.Logf("waiting %d seconds before goroutine. No log to find.\n", d)
	time.Sleep(d * time.Second)
	go func() {
		t.Logf("waiting %d seconds in goroutine. Log displays unless interrupted.\n", d)
		time.Sleep(d * time.Second)
		err = SendCtrlBreak(cmd.Process.Pid)
		if err != nil {
			t.Log(err)
		}
	}()
	err = cmd.Wait()
	if testing.Verbose() {
		f, err := os.Open("ctrlbreak.log")
		if err == nil {
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				fmt.Println(scanner.Text()) // token in unicode-char
			}
		} else {
			t.Error("cannot access log")
		}
	}
	if err != nil {
		t.Fatalf("Program exited with error: %v\n", err)
	}
}
