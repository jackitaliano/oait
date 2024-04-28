package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"

	"github.com/jackitaliano/oait/cmd/assts"
	"github.com/jackitaliano/oait/cmd/files"
	"github.com/jackitaliano/oait/cmd/images"
	"github.com/jackitaliano/oait/cmd/threads"
)

func main() {
	const progName = "oait"
	const progDesc = "OpenAI Tools"

	parser := argparse.NewParser(progName, progDesc)
	keyArg := parser.String("k", "key", &argparse.Options{
		Required: false,
		Help:     "OpenAI API Key (default to env var 'OPENAI_API_KEY')",
	})

	threadsService := threads.NewService(parser)
	filesService := files.NewService(parser)
	asstsService := assts.NewService(parser)
	imagesService := images.NewService(parser)

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	if *keyArg == "" {
		*keyArg = os.Getenv("OPENAI_API_KEY")
	}

	commands := parser.GetCommands()

	threadsCommand := commands[0]
	filesCommand := commands[1]
	asstsCommand := commands[2]
	imagesCommand := commands[3]

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

	} else if asstsCommand.Happened() {
		err := asstsService.Run(*keyArg)

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
