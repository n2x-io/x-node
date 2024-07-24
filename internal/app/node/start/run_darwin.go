//go:build darwin
// +build darwin

package start

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/kardianos/service"
	"github.com/spf13/viper"
	"n2x.dev/x-lib/pkg/version"
	"n2x.dev/x-lib/pkg/xlog"
)

type serviceAction int

const (
	actionConsoleRun serviceAction = iota
	actionServiceStart
	actionServiceInstall
	actionServiceUninstall
)

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()

	return nil
}

func (p *program) run() {
	start()
}

func (p *program) Stop(s service.Service) error {
	go finish()

	return nil
}

func runAsService(action serviceAction) {
	svcConfig := &service.Config{
		Name:             fmt.Sprintf("io.n2x.%s", version.NODE_NAME),
		DisplayName:      version.NODE_NAME,
		Description:      "n2x-node",
		Arguments:        []string{"service-start"},
		WorkingDirectory: "/opt/n2x/var/tmp",
		Option: service.KeyValue{
			"KeepAlive":         true,
			"RunAtLoad":         true,
			"StandardErrorPath": "/opt/n2x/var/log/io.n2x.n2x-node.err.log",
			"StandardOutPath":   "/opt/n2x/var/log/io.n2x.n2x-node.out.log",
		},
	}

	prg := &program{}

	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	switch action {
	case actionConsoleRun:
		err = s.Run()
	case actionServiceStart:
		err = s.Run()
	case actionServiceInstall:
		err = s.Install()
	case actionServiceUninstall:
		err = s.Uninstall()
	}
	if err != nil {
		logger.Error(err)
	}
}

func Main() {
	xlog.Infof("%s starting on %s :-)", version.NODE_NAME, viper.GetString("host.id"))
	defer xlog.Logger().Close()
	runAsService(actionConsoleRun)
	<-done

	xlog.Infof("%s stopped on %s", version.NODE_NAME, viper.GetString("host.id"))

	os.Exit(0)
}

func ServiceStart() {
	xlog.Infof("Starting %s Service", version.NODE_NAME)
	defer xlog.Logger().Close()
	runAsService(actionServiceStart)

	os.Exit(0)
}

func ServiceInstall() {
	xlog.Infof("Installing %s as Service", version.NODE_NAME)
	runAsService(actionServiceInstall)

	cmd := exec.Command("launchctl", "load", "/Library/LaunchDaemons/io.n2x.n2x-node.plist")
	if err := cmd.Run(); err != nil {
		xlog.Warnf("Unable to load launchctl n2x-node service, please check: %v", err)
	}

	os.Exit(0)
}

func ServiceUninstall() {
	xlog.Infof("Uninstalling %s Service", version.NODE_NAME)
	runAsService(actionServiceUninstall)

	// cmd := exec.Command("launchctl", "unload", "/Library/LaunchDaemons/io.n2x.n2x-node.plist")
	// if err := cmd.Run(); err != nil {
	// 	xlog.Warnf("Unable to unload launchctl n2x-node service, please check: %v", err)
	// }

	os.Exit(0)
}
