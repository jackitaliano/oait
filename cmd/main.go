package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"

	"github.com/jackitaliano/oait-go/cmd/imagesParse"
	"github.com/jackitaliano/oait-go/cmd/threadsParse"
	"github.com/jackitaliano/oait-go/cmd/filesParse"
)

func main() {
	const progName = "oait"
	const progDesc = "OpenAI Tools"

	parser := argparse.NewParser(progName, progDesc)
	keyArg := parser.String("k", "key", &argparse.Options{
		Required: false,
		Help:     "OpenAI API Key (default to env var 'OPENAI_API_KEY')",
	})

	threadsService := threadsParse.NewService(parser)
	filesService := filesParse.NewService(parser)
	imagesService := imagesParse.NewService(parser)

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	if *keyArg == "" {
		*keyArg = os.Getenv("OPENAI_API_KEY")
	}

	var commands []*argparse.Command = parser.GetCommands()

	var threadsCommand *argparse.Command = commands[0]
	var filesCommand *argparse.Command = commands[1]
	var imagesCommand *argparse.Command = commands[2]

	if threadsCommand.Happened() {
		err := threadsService.Run(*keyArg)

		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

	} else if filesCommand.Happened() {
		err := filesService.Run(*keyArg)

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
