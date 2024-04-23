package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"

	"github.com/jackitaliano/oait-go/cmd/threadsParse"
	"github.com/jackitaliano/oait-go/cmd/imagesParse"
)

func main() {
	const progName = "oait"
	const progDesc = "OpenAI Tools"

	parser := argparse.NewParser(progName, progDesc)
	keyArg := parser.String("k", "key", &argparse.Options{ 
		Required: false, 
		Help: "OpenAI API Key (default to env var 'OPENAI_API_KEY')",
		Default: os.Getenv("OPENAI_API_KEY"),
	})

	threadsService := threadsParse.NewService(parser)
	imagesService := imagesParse.NewService(parser)

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	var commands []*argparse.Command = parser.GetCommands()

	var threadsCommand *argparse.Command = commands[0]
	var imagesCommand *argparse.Command = commands[1]

	if threadsCommand.Happened() {
		err := threadsService.Run(*keyArg)

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
