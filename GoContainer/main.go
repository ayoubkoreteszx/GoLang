package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

//docker run container commands arg
// go run main.go run commands arg

func main() {
	switch os.Args[1] {
	case "run":
		{
			run()
		}
	case "child":
		child()
	default:
		panic("what?")
	}

}

func run() {
	cmd := exec.Command("proc/self/exec", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}
	must(cmd.Run())

}

func child() {
	fmt.Printf("runing %v as PID %d \n", os.Args[2:], os.Getpid())
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	must(syscall.Chroot("/home/rootfs"))
	must(syscall.Mount("proc", "proc", "proc", 0, ""))
	must(os.Chdir("/"))
	must(cmd.Run())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
