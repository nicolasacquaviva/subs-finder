package utils

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type TermSize struct {
	X int
	Y int
}

func GetTerminalSize() (*TermSize, error) {
	cmd := exec.Command("stty", "size")

	cmd.Stdin = os.Stdin

	out, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	size := strings.Split(string(out), " ")
	y, err := strconv.Atoi(size[0])

	if err != nil {
		return nil, err
	}

	x, err := strconv.Atoi(strings.TrimSuffix(size[1], "\n"))

	if err != nil {
		return nil, err
	}

	termSize := &TermSize{
		X: x,
		Y: y,
	}

	return termSize, nil
}
