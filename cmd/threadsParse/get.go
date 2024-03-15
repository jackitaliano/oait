package threadsParse

import (
	"errors"
	"fmt"

	"github.com/jackitaliano/oait-go/internal/threads"

	"github.com/akamensky/argparse"
)

type GetCommand struct {
	name    string
	desc    string
	command *argparse.Command

	threadsArg *[]string
	inputArg   *string
	outputArg  *string
}

func NewGetCommand(command *argparse.Command) *GetCommand {
	const name = "get"
	const desc = "Get Thread Tools"

	subCommand := command.NewCommand(name, desc)

	threadsArg := subCommand.StringList("t", "threads", &argparse.Options{Required: false, Help: "List of Thread IDs"})
	inputArg := subCommand.String("i", "input", &argparse.Options{Required: false, Help: "Thread File Input"})
	outputArg := subCommand.String("o", "output", &argparse.Options{Required: false, Help: "Thread File Output"})

	return &GetCommand{
		name,
		desc,
		subCommand,
		threadsArg,
		inputArg,
		outputArg,
	}
}

func (g *GetCommand) Happened() bool {
	return g.command.Happened()
}

func (g *GetCommand) Run() error {
	args := g.command.GetArgs()
	threadsParsed := args[1].GetParsed()
	inputParsed := args[2].GetParsed()
	outputParsed := args[3].GetParsed()

	var parsedThreads *[]*threads.Thread
	var err error

	// Input flow
	if threadsParsed { // List passed
		parsedThreads, err = threads.ListInput(g.threadsArg)

		if err != nil {

		}

	} else if inputParsed { // File input passed
		parsedThreads, err = threads.FileInput((*g.inputArg))

		if err != nil {

		}

	} else { // No input passed
		errMsg := fmt.Sprintf("No input options passed to `%v`\n", g.name)
		helpMsg := g.command.Help(errMsg)

		err := errors.New(helpMsg)
		return err
	}

	// Output flow
	if outputParsed {
		b, err := threads.ThreadsToJson(parsedThreads)

		if err != nil {
			panic(err)
		}

		threads.FileOutput(*g.outputArg, b)

	} else {
		fmt.Printf("Printing messages.\n")
		fmt.Printf("Messages:\n")

		for _, msg := range *parsedThreads {
			fmt.Printf("%v\n", *msg)
		}
	}

	return nil
}
