package cmd

import (
	"errors"
	"fmt"
	"github.com/fireflyblock/ticket-stats/block"
	"github.com/urfave/cli/v2"
)

var StatsCmd = &cli.Command{
	Name: "stats",
	Usage: "Get ticket stats of miner",
	Subcommands: []*cli.Command{
		TimeCmd,
		EpochCmd,
	},
}

var TimeCmd = &cli.Command{
	Name: "time",
	Usage: "Get ticket stats of miner by time; Time format:2006-01-02T15:04:05",
	//ArgsUsage: "[miner] [begin time] [end time]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "miner",
			Usage: "The miner you want to check",
		},
		&cli.StringFlag{
			Name: "start",
			Usage: "The begin time you want to check",
		},
		&cli.StringFlag{
			Name: "end",
			Usage: "The end time you want to check,must be later than start time",
		},
	},
	Action: func(ctx *cli.Context) error {
		var miner,start,end string
		if ctx.IsSet("miner"){
			miner=ctx.String("miner")
		}else {
			return errors.New("Fail to get miner id")
		}
		if ctx.IsSet("start"){
			start=ctx.String("start")
		}else {
			return errors.New("Fail to get start time")
		}
		if ctx.IsSet("end"){
			end=ctx.String("end")
		}else {
			return errors.New("Fail to get end time")
		}
		query:=&block.QueryParams{
			Start: start,
			End:   end,
			Miner: miner,
		}
		fmt.Println("Checking...")
		if err:=query.GetTicketStatsByTime();err!=nil{
			return err
		}
		fmt.Println("Completed.")
		return nil
	},
}

var EpochCmd = &cli.Command{
	Name: "epoch",
	Usage: "Get ticket stats of miner by epoch",
	//ArgsUsage: "[miner] [begin epoch] [end epoch]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "miner",
			Usage: "The miner you want to check",
		},
		&cli.StringFlag{
			Name: "start",
			Usage: "The start epoch you want to check",
		},
		&cli.StringFlag{
			Name: "end",
			Usage: "The end epoch you want to check,must be bigger than start epoch",
		},
	},
	Action: func(ctx *cli.Context) error {
		var miner,start,end string
		if ctx.IsSet("miner"){
			miner=ctx.String("miner")
		}else {
			return errors.New("Fail to get miner id")
		}
		if ctx.IsSet("start"){
			start=ctx.String("start")
		}else {
			return errors.New("Fail to get start epoch")
		}
		if ctx.IsSet("end"){
			end=ctx.String("end")
		}else {
			return errors.New("Fail to get end epoch")
		}
		query:=&block.QueryParams{
			Start: start,
			End:   end,
			Miner: miner,
		}
		fmt.Println("Checking...")
		if err:=query.GetTicketStatsByEpoch();err!=nil{
			return err
		}
		fmt.Println("Completed.")
		return nil
	},
}


