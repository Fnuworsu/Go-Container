# Build My Own Container in Go

- **Namespaces**: Isolate processes and provide them with their own view of the system, such as process IDs, network interfaces, and file systems.  
  - **What you can see**: Namespaces define what a process can see in terms of system resources.  
  - **Created with syscalls**: Namespaces are created using system calls.  
    - Unix Timesharing System  
    - Process IDs  
    - Mounts  
    - Network  
    - User IDs  
    - InterProcess Communications (IPC)

- **Chroot**: Change the root directory for a process, creating a sandboxed environment for file system access.

- **Cgroups**: Limit and isolate resource usage (CPU, memory, disk I/O) for a group of processes.