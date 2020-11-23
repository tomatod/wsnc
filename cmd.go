package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

var app *cli.App = &cli.App{
	Name:      "wsnc",
	Usage:     "Simple websocket client and server tool.",
	Action:    cliRun,
	UsageText: "wsnc [global options] [arguments...]",
}

// All flags are set up in init function.
func init() {
	appendTemplate()
	/* if change help flag.
	cli.HelpFlag = &cli.BoolFlag {
		Name: "help",
		Aliases: []string{"h"},
		Value: false,
		Usage: "If you want to see help, use this option.",
	}
	*/
	serverFlag := &cli.BoolFlag{
		Name:    "server",
		Aliases: []string{"s"},
		Value:   false,
		Usage:   "Run websocket server mode. Default is client mode (= false).",
	}

	pathFlag := &cli.StringFlag{
		Name:    "path",
		Aliases: []string{"u"},
		Value:   "/",
		Usage:   "Specify standby path of server. This option is ignored in client mode.",
	}

	portFlag := &cli.IntFlag{
		Name:    "port",
		Aliases: []string{"p"},
		Value:   8080,
		Usage:   "Specify listen port. In client mode, this option is ignored and argument (ws://... or wss://...) is refered to.",
	}

	messageFlag := &cli.StringFlag{
		Name:    "message",
		Aliases: []string{"m"},
		Usage:   "Specify message sent by server in text or binary message or by one-shot mode client.",
	}

	oneShotFlag := &cli.BoolFlag{
		Name:    "oneshot",
		Aliases: []string{"o"},
		Value:   false,
		Usage:   "Run client with one shot  mode. This option is ignored in server mode.",
	}

	logTimeFlag := &cli.BoolFlag{
		Name:    "logtime",
		Aliases: []string{"t"},
		Value:   false,
		Usage:   "Attach timestamp to server logs.",
	}

	logLevelFlag := &cli.BoolFlag{
		Name:    "loglevel",
		Aliases: []string{"l"},
		Value:   false,
		Usage:   "Attach log level to server logs.",
	}

	broadcastFlag := &cli.BoolFlag{
		Name:    "broadcast",
		Aliases: []string{"b"},
		Value:   false,
		Usage:   "If you broadcast from server, specify this option.",
	}

	debugFlag := &cli.BoolFlag{
		Name:    "debug",
		Aliases: []string{"d"},
		Value:   false,
		Usage:   "Active additional debug logs. Note that this option is unimplemented now.",
	}

	app.Flags = []cli.Flag{
		serverFlag,
		pathFlag,
		portFlag,
		oneShotFlag,
		broadcastFlag,
		messageFlag,
		logTimeFlag,
		logLevelFlag,
		debugFlag,
	}
}

func main() {
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// 1. Set up configration depending on command options
// 2. Judged server or client mode and jump to either start function of that.
func cliRun(c *cli.Context) error {
	err := setAppConfig(c)
	if err != nil {
		return err
	}
	if appConfig.IsServer {
		err = startServer()
	} else {
		err = startClient()
	}
	return err
}

func appendTemplate() {
	cli.AppHelpTemplate = fmt.Sprintf(`%s
EXAMPLE:
   // connect to test server of websocket.org
   wsnc wss://echo.websocket.org
   // send only one message and close connection (one-shot mode)
   wsnc -o -m hoge wss://echo.websocket.org
   // open websocket server (port: 8080 and path: "/")
   wsnc -s 
   // open websocket server (port: 80 and path: "/bar/")
   wsnc -s -p 80 -u /bar/ 
   // output time (-t) and loglevel (-l) in server usage.
   wsnc -s -t -l

HTML MANUAL: https://github.com/tomatod/wsnc/blob/master/README.md

`, cli.AppHelpTemplate)
}
