package proc

import (
	"context"
	"errors"
	"strconv"
	"strings"
)

// listeningOn asks sockstat, which lists the listening sockets this user
// may see, with the process holding each. The columns are read from the
// right: a command's name may hold spaces.
func listeningOn(ctx context.Context, port int) ([]int, error) {
	out, err := CommandContext(ctx, "sockstat", "-l", "-4", "-6", "-p", strconv.Itoa(port)).Output()
	if err != nil {
		return nil, err
	}
	var pids []int
	seen := map[int]bool{}
	for _, line := range strings.Split(string(out), "\n") {
		f := strings.Fields(line)
		// USER COMMAND… PID FD PROTO LOCAL FOREIGN
		if len(f) < 7 || !strings.HasPrefix(f[len(f)-3], "tcp") {
			continue
		}
		addr := f[len(f)-2]
		i := strings.LastIndex(addr, ":")
		if i < 0 || addr[i+1:] != strconv.Itoa(port) {
			continue
		}
		pid, err := strconv.Atoi(f[len(f)-5])
		if err != nil || pid <= 0 || seen[pid] {
			continue
		}
		seen[pid] = true
		pids = append(pids, pid)
	}
	return pids, nil
}

// executable asks procstat for the path the process was started from.
func executable(pid int) (string, error) {
	out, err := CommandContext(context.Background(), "procstat", "-b", strconv.Itoa(pid)).Output()
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(string(out), "\n") {
		f := strings.Fields(line)
		// PID COMM OSREL PATH
		if len(f) >= 4 && f[0] == strconv.Itoa(pid) {
			return strings.Join(f[3:], " "), nil
		}
	}
	return "", errors.New("no program path")
}

// processName is the program's name by ps, which names any user's process
// where the path (procstat) is denied.
func processName(pid int) string {
	out, err := CommandContext(context.Background(), "ps", "-o", "comm=", "-p", strconv.Itoa(pid)).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
