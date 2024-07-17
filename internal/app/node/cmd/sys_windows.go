//go:build windows
// +build windows

package cmd

import (
	"fmt"
	"os"
	"runtime"

	"n2x.dev/x-lib/pkg/utils/colors"
)

func ConsoleInit() error {
	if runtime.GOOS == "windows" {
		if err := colors.EnableWindowsVirtualTerminalProcessing(); err != nil {
			return err
		}
	}

	return nil
}

func defaultConfigFile() string {
	programFiles := os.Getenv("ProgramFiles")

	if len(programFiles) == 0 {
		programFiles = `C:\Program Files`
	}

	return fmt.Sprintf(`%s\n2x\n2x-node.yml`, programFiles)
}

func logFile() string {
	programFiles := os.Getenv("ProgramFiles")

	if len(programFiles) == 0 {
		programFiles = `C:\Program Files`
	}

	return fmt.Sprintf(`%s\n2x\n2x-node.log`, programFiles)
}
