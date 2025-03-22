package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// docker         run image <cmd> <params>
// go run main.go run       <cmd> <params>

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("bad command")
	}
}

func run() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Adding namespaces
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}

	cmd.Run()
}

func child() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	// Setting up hostname
	syscall.Sethostname([]byte("container"))

	// Mount pseudo filesystems : kernel for the user and user space to share info
	// Mount dir asproc pseudo filesystem so that kernel knows i'm going to populate
	// that with all the information about these running processes
	syscall.Mount("proc", "/mycontainerroot/proc", "proc", 0, "")

	// Setting root directory(my own /proc)
	syscall.Chroot("/mycontainerroot")
	syscall.Chdir("/")

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

	syscall.Unmount("/proc", 0)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
