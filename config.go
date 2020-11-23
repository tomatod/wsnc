package main

import (
	"github.com/urfave/cli/v2"
	"strconv"
)

var appConfig Config

type Config struct {
	Port        int
	PortStr     string // Any libraly require string of port.
	Url         string
	Message     string
	Path        string
	Headers     []string
	IsServer    bool
	IsOneShot   bool
	IsTime      bool
	IsLevel     bool
	IsDebug     bool
	IsBroadcast bool
}

func setAppConfig(c *cli.Context) error {
	appConfig = Config{
		Port:        c.Int("port"),
		PortStr:     strconv.Itoa(c.Int("port")),
		Message:     c.String("message"),
		Path:        c.String("path"),
		IsServer:    c.Bool("server"),
		IsOneShot:   c.Bool("oneshot"),
		IsTime:      c.Bool("logtime"),
		IsLevel:     c.Bool("loglevel"),
		IsDebug:     c.Bool("debug"),
		IsBroadcast: c.Bool("broadcast"),
		Headers:     c.StringSlice("header"),
	}

	if appConfig.IsServer {
		return serverConfigCheck()
	}

	// If client mode, default first parameter is destination URL.
	appConfig.Url = c.Args().Get(0)
	return nil
}

// Check command option parameter used in server mode.
func serverConfigCheck() error {
	if appConfig.Path[:1] != "/" {
		return errLogf("first character of path option must be \"/\".")
	}
	return nil
}
