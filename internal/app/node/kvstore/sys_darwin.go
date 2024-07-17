//go:build darwin
// +build darwin

package kvstore

func dbDir() string {
	return "/opt/n2x/var/lib/db"
}
