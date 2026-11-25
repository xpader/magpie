package claudetui

import (
	"os"
	"strconv"
	"syscall"

	"golang.org/x/sys/unix"
)

// openPTY opens a pseudo-terminal: its master, and the name of its slave.
// FreeBSD's ptmx needs no grant or unlock; TIOCGPTN gives the pts unit.
// The master is opened blocking and wrapped with os.NewFile, so it stays
// out of Go's kqueue poller, as on macOS (pty_darwin.go).
func openPTY() (*os.File, string, error) {
	fd, err := syscall.Open("/dev/ptmx", syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_CLOEXEC, 0)
	if err != nil {
		return nil, "", &os.PathError{Op: "open", Path: "/dev/ptmx", Err: err}
	}
	m := os.NewFile(uintptr(fd), "/dev/ptmx")
	n, err := unix.IoctlGetInt(fd, unix.TIOCGPTN)
	if err != nil {
		m.Close()
		return nil, "", err
	}
	return m, "/dev/pts/" + strconv.Itoa(n), nil
}

// the requests that read and set a terminal's modes (the tests' fake
// Claude Code sets its own raw, as the real one does)
const ioctlGetTermios, ioctlSetTermios = unix.TIOCGETA, unix.TIOCSETA
