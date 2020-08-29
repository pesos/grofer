Grofer
======

![build](https://api.travis-ci.org/pesos/grofer.svg?branch=master&status=started)

A clean system monitor and profiler written purely in golang using [termui](https://github.com/gizak/termui) and [gopsutil](https://github.com/shirou/gopsutil)!

Installation
------------

Using go get:

```
go get -u github.com/pesos/grofer
```

As an executable:

```
curl -sSL https://github.com/pesos/grofer/releases/download/<version tag>/grofer --output grofer
chmod +x grofer
```

For system wide usage, install `grofer` to a location on `$PATH`, e.g. `/usr/local/bin`

```
mv grofer /usr/local/bin
```

Building from source:

```
git clone https://github.com/pesos/grofer
cd grofer
go build grofer.go
```

Usage
-----

```
grofer is a system profiler written in golang

Usage:
  grofer [flags]
  grofer [command]

Available Commands:
  about       about is a command that gives information about the project in a cute way
  help        Help about any command
  proc        proc command is used to get per-process information

Flags:
      --config string   config file (default is $HOME/.grofer.yaml)
  -c, --cpuinfo         Info about the CPU Load over all CPUs
  -h, --help            help for grofer
  -r, --refresh int32   Overall stats UI refreshes rate in milliseconds greater than 1000 (default 1000)
  -t, --toggle          Help message for toggle

Use "grofer [command] --help" for more information about a command.

```

Examples
--------

`grofer [-r refreshRate][--cpuinfo]`
-------------------------

This gives overall utilization stats refreshed every `refreshRate` milliseconds. Default and minimum value of the refresh rate is `1000 ms`. 

![grofer](images/README/grofer.png)

Information provided:  
- CPU utilization per core  
- Memory (RAM) usage  
- Network usage  
- Disk storage

The `--cpuinfo` flag displayes finer details about the CPU load such as percentage of the time spent servicing software interrupts, hardware interrupts, etc.

![grofer-cpu](images/README/cpuload.png)

Information provided:
- Usr   : % of time spent executing user level applications.
- Sys   : % of time spent executing kernel level processes.
- Irq   : % of time spent servicing hardware interrupts.
- Idle  : % of time CPU was idle.
- Nice  : % of time spent by CPU executing user level processes with a nice priority.
- Iowait: % of time spent by CPU waiting for an outstanding disk I/O.
- Soft  : % of time spent by the CPU servicing software interrupts.
- Steal : % of time spent in involuntary waiting by logical CPUs.
---

`grofer proc [-p PID] [-r refreshRate]`
---------------------------------------

If the `-r` flag is specified then the UI will refresh and display new information every `refreshRate` milliseconds. The minimum and default value for `refreshRate` is `1000 ms`.

### `grofer proc`

This lists all running processes and relevant information.

![grofer-proc](images/README/grofer-proc.png)

---

### `grofer proc -p PID`

This gives information specific to a process, specified by a valid PID.

![grofer-proc-pid](images/README/grofer-proc-pid.png)

Information provided:  
 - CPU utilization %  
 - Memory utilization %  
 - Child processes  
 - Number of voluntary and involuntary context switches  
 - Memory usage (RSS, Data, Stack, Swap)
