package util

import (
	"os/exec"
	"path"
)

func GetProjectRoot() (string, error) {
	cmd := exec.Command("go", "env", "GOMOD")

	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return path.Dir(string(stdout)), nil
}
