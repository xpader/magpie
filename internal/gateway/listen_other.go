//go:build !darwin && !linux && !freebsd

package gateway

import "syscall"

// Elsewhere a second process can't share the port: the handover is a
// plain restart.
func reusePort(_, _ string, _ syscall.RawConn) error { return nil }
