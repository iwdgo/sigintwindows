//go:build windows
// +build windows

package sigintwindows

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
		if err := os.Remove(filepath.Join(processToSignal, exe)); err != nil {
			t.Log(err)
		}
	}()
	o, err := exec.Command("go", "install", filepath.Join(processToSignal, src)).CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to compile: %v\n%v", err, string(o))
	}
	cmd := exec.Command(exe)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
	if err := cmd.Start(); err != nil {
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
		if err := SendCtrlBreak(cmd.Process.Pid); err != nil {
			t.Log(err)
		}
	}()
	if err := cmd.Wait(); err != nil {
		t.Errorf("Wait failed: %v", err)
	}
	if testing.Verbose() {
		f, err := os.Open("ctrlbreak.log")
		if err == nil {
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}
		} else {
			t.Error("cannot access log")
		}
	}
	if err != nil {
		t.Fatalf("Program exited with error: %v\n", err)
	}
}

func TestSendCtrlBreakNoPid(t *testing.T) {
	var nonExistingPid uint32 = 999000
	const READ_CONTROL uint32 = 0x00020000
	var err error
	for err == nil {
		_, err = syscall.OpenProcess(READ_CONTROL, false, nonExistingPid)
		nonExistingPid++
		if nonExistingPid > 999999 {
			t.Skipf("no invalid process id found")
		}
	}
	err = SendCtrlBreak(int(nonExistingPid))
	if err == nil {
		t.Fatalf("Sending Ctrl Break with an invalid process id did not fail")
	}
}
