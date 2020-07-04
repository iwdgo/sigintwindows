// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func sendCtrlBreak(pid int) {
	d, e := syscall.LoadDLL("kernel32.dll")
	if e != nil {
		log.Fatalf("LoadDLL: %v\n", e)
	}
	p, e := d.FindProc("GenerateConsoleCtrlEvent")
	if e != nil {
		log.Fatalf("FindProc: %v\n", e)
	}
	r, _, e := p.Call(syscall.CTRL_BREAK_EVENT, uintptr(pid))
	if r == 0 {
		log.Fatalf("GenerateConsoleCtrlEvent: %v\n", e)
	}
}

func main() {

	// compile it
	name := "ctrlbreak"
	src := name + ".go"
	exe := name + ".exe"
	defer os.Remove(exe)
	o, err := exec.Command("go", "build", "-o", exe, src).CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to compile: %v\n%v", err, string(o))
	}

	// run it
	cmd := exec.Command(exe)
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
	err = cmd.Start()
	if err != nil {
		log.Fatalf("Start failed: %v", err)
	}
	t := time.Duration(5)
	go func() {
		log.Printf("waiting %d second\n", t)
		time.Sleep(t * time.Second)
		sendCtrlBreak(cmd.Process.Pid)
	}()
	err = cmd.Wait()
	if err != nil {
		log.Fatalf("Program exited with error: %v\n%v", err, string(b.Bytes()))
	}
}
