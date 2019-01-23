package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/docker/docker/pkg/reexec"
)

func init() {
	reexec.Register("hninit", hninit)
	if reexec.Init() {
		os.Exit(0)
	}
}

func hninit() {
	fmt.Println("Setting hostname")
	syscall.Sethostname([]byte("check123"))
	newrootPath := os.Args[1]
	fmt.Println("newRoot:", newrootPath)
	mountProc(newrootPath)
	if err := mountPivot(newrootPath); err != nil {
		log.Fatalln(err)
	}
	run()
}

func run() {
	cmd := exec.Command("/bin/sh")
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	cmd.Env = []string{"PS1=abhi-doc#"}
	fmt.Println("in run function")
	if err := cmd.Run(); err != nil {
		log.Fatalln("Error while running cmd")
	}
}
func main() {
	var rootfsPath string
	flag.StringVar(&rootfsPath, "rootfs", "/home/abhi/work/go/src/github.com/Golang_play/clone_docker/busybox", "Path to the root filesystem to use")
	flag.Parse()
	cmd := reexec.Command("hninit", rootfsPath)

	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	//cmd.Env = []string{"PS1=abhi-doc#"}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWUSER |
			syscall.CLONE_NEWUTS,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}
	fmt.Println("11111: Set the uid, gid now cmd.RUN from mainsss")
	if err := cmd.Run(); err != nil {
		log.Fatalln("Error while running reexec cmd")
	}
}
