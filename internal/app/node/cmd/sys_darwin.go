//go:build darwin
// +build darwin

package cmd

func ConsoleInit() error {
	return nil
}

func defaultConfigFile() string {
	return "/opt/n2x/etc/n2x-node.yml"
}

/*
func logFile() string {
	return "/opt/n2x/var/log/n2x-node.log"
}
*/
