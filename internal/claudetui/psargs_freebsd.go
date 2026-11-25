package claudetui

// psOutputArgs lists every process with no headers, one column each.
// FreeBSD's ps takes one keyword per -o; a comma list after the first
// "=" becomes that column's header instead of more columns.
var psOutputArgs = []string{"-A", "-o", "pid=", "-o", "ppid=", "-o", "rss=", "-o", "etime=", "-o", "command="}
