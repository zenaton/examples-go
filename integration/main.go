package integration

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func Listen() error {
	args := []string{"listen", "--boot=../boot/boot.go"}
	out, err := exec.Command("zenaton", args...).CombinedOutput()
	fmt.Println("out1: ", string(out))

	//try again
	if err != nil {
		out, err = exec.Command("zenaton", args...).CombinedOutput()
		fmt.Println("out2: ", string(out))
	}

	return err
}
