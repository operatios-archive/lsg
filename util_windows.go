package main

import (
	"syscall"

	"golang.org/x/sys/windows"
)

var (
	kernel32       = syscall.NewLazyDLL("Kernel32.dll")
	setConsoleMode = kernel32.NewProc("SetConsoleMode")
)

func enableColors() error {
	var mode uint32
	err := syscall.GetConsoleMode(syscall.Stdout, &mode)
	if err != nil {
		return err
	}

	if mode&windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING != 0 {
		return nil
	}

	mode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
	ret, _, err := setConsoleMode.Call(uintptr(syscall.Stdout), uintptr(mode))
	if ret == 0 {
		return err
	}
	return nil
}
