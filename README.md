## Container Implementation in Go

This project demonstrates a lightweight container written in Go using Linux namespaces and cgroups:
- **Namespaces** (`UTS`, `PID`, `mount`) isolate process information, hostnames, and mounted filesystems.  
- **Chroot** sets a new root directory (`/mycontainerroot`) for the isolated environment.  
- **Cgroups** (in `/sys/fs/cgroup/pids`) limit the maximum number of processes a container can spawn.

### How It Works
1. **run():**  
   - Creates a new process that calls itself (`child`) with new UTS, PID, and mount namespaces.  
2. **child():**  
   - Configures a control group (cgroup), setting a maximum of 20 processes.  
   - Sets a custom hostname.  
   - Mounts a proc filesystem in `/mycontainerroot/proc`.  
   - Changes the root to `/mycontainerroot`.  
   - Executes the requested command.  
   - Unmounts `/proc` on exit.  
3. **cg()**
    The `cg()` function sets up a cgroup to limit the container’s total number of processes:
   - Creates a directory for this cgroup (`felix`) in `/sys/fs/cgroup/pids`.  
   - Writes a max process limit (“20”) to `pids.max`.  
   - Configures the cgroup to remove itself on process exit (`notify_on_release`).  
   - Writes the current process ID into `cgroup.procs`, attaching the child process to the new group.  

   Example usage inside `child()`:
    ```go
    cg() // Enforce cgroup limits
    syscall.Sethostname([]byte("container"))
    // ...mount /mycontainerroot/proc, chroot, etc....

Run the container:
```bash
go run main.go run /bin/bash