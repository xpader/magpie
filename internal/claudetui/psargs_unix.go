//go:build darwin || linux

package claudetui

// psOutputArgs lists every process with no headers, one column each.
var psOutputArgs = []string{"-A", "-o", "pid=,ppid=,rss=,etime=,command="}
