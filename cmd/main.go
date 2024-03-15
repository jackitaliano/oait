package main

import (
	"fmt"
	"os"

	"github.com/jackitaliano/oait-go/cmd/threadsParse"
	"github.com/jackitaliano/oait-go/cmd/imagesParse"

	"github.com/akamensky/argparse"
)

func main() {
	const progName = "oait"
	const progDesc = "OpenAI Tools"

	parser := argparse.NewParser(progName, progDesc)

	threadsService := threadsParse.NewService(parser)
	imagesService := imagesParse.NewService(parser)

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	var args []*argparse.Command = parser.GetCommands()

	var threadsCommand *argparse.Command = args[0]
	var imagesCommand *argparse.Command = args[1]

	if threadsCommand.Happened() {
		err := threadsService.Run()

		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

	} else if imagesCommand.Happened() {
		err := imagesService.Run()

		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}
	}
}
