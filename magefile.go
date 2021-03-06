// +build mage

package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/pkg/errors"
)

var (
	wslGoRoot string = "/home/tobias/go1.11"
)

func UI(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "ng", "build", "--prod")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = "ui"
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "statik")
	}
	return nil
}

func Statik(ctx context.Context) error {
	mg.Deps(UI)
	cmd := exec.CommandContext(ctx, "statik", "-src", filepath.Join("ui", "dist", "ui"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "statik")
	}
	return nil
}

func Rockrobo() error {
	mg.Deps(Statik)
	fmt.Println("building and installing on vacuum...")
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	path := filepath.Dir(filename)
	path = strings.Replace(path, "\\", "/", -1)
	path = strings.Replace(path, ":", "", -1)
	path = strings.ToLower("/mnt/" + path)
	fmt.Println(path)
	cmd := exec.Command("wsl", "$GOROOT/bin/go", "build", "-ldflags=-s -w")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(),
		"WSLENV=GOOS:GOARCH:GOARM:CGO_CFLAGS:CGO_LDFLAGS:CGO_ENABLED:CXX:CXX_FOR_TARGET:CC:CC_FOR_TARGET:GOROOT",
		fmt.Sprintf("GOROOT=%s", wslGoRoot),
		"GOOS=linux",
		"GOARCH=arm",
		"GOARM=7",
		fmt.Sprintf("CGO_CFLAGS=-I%s/libs/arm/include", path),
		fmt.Sprintf("CGO_LDFLAGS=-L%s/libs/arm/lib -lasound -lm -ldl -lpthread -lrt", path),
		"CGO_ENABLED=1",
		"CXX=arm-linux-gnueabihf-g++",
		"CXX_FOR_TARGET=arm-linux-gnueabihf-g++",
		"CC=arm-linux-gnueabihf-gcc",
		"CC_FOR_TARGET=arm-linux-gnueabihf-cc",
	)
	if err := cmd.Run(); err != nil {
		return err
	}
	cmd = exec.Command("upx", "dustrobo")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("scp", "dustrobo", "rockrobo:/mnt/data/tmp/dustrobo")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func Build() error {
	mg.Deps(Statik)
	fmt.Println("building...")
	cmd := exec.Command("go", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
