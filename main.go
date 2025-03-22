package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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

	// Adding namespaces to kernel
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		// not sharing new namespace with host
		Unshareflags: syscall.CLONE_NEWNS,
	}

	cmd.Run()
}

func child() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	cg()
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

// limiting the number of (caping memory)
func cg() {
	cgroups := "/sys/fs/cgroup"
	pids := filepath.Join(cgroups, "pids")
	err := os.Mkdir(filepath.Join(pids, "felix"), 0755)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
	// Limit control group to 20 processes
	must(os.WriteFile(filepath.Join(pids, "felix/pids.max"), []byte("20"), 0700))
	// Removes the new cgroup in place after the container exits
	must(os.WriteFile(filepath.Join(pids, "felix/notify_on_release"), []byte("1"), 0700))
	// Write the pid of the current process to the cgroup's cgroup.procs file
	must(os.WriteFile(filepath.Join(pids, "felix/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
