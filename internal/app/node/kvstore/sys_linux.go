//go:build linux
// +build linux

package kvstore

func dbDir() string {
	return "/var/lib/n2x/db"
}
