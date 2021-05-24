package main

import (
	"fmt"
	"github.com/fireflyblock/ticket-stats/cmd"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
	"os"
)

type PrintHelpErr struct {
	Err error
	Ctx *cli.Context
}


func main() {

	wallet := cli.App{
		Name:        "ticket-stats",
		Usage:       "Get ticket stats of miner",
		//Version:     "version 0.1",
		//Description: "control wallet",
		Commands: []*cli.Command{
			cmd.StatsCmd,
		},
		EnableBashCompletion: true,
	}
	if err := wallet.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n\n", err)
		var phe *PrintHelpErr
		if xerrors.As(err, &phe) {
			_ = cli.ShowCommandHelp(phe.Ctx, phe.Ctx.Command.Name)
		}
		os.Exit(1)
	}

}
