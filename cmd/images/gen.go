package images

import (
	"errors"
	"fmt"

	"github.com/akamensky/argparse"
)

type GenCommand struct {
	name    string
	desc    string
	command *argparse.Command

	promptArg *string
	outputArg *string
}

func NewGenCommand(command *argparse.Command) *GenCommand {
	const name = "get"
	const desc = "Get Thread Tools"

	subCommand := command.NewCommand(name, desc)

	promptArg := subCommand.String("p", "prompt", &argparse.Options{Required: false, Help: "Image Generation Prompt"})
	outputArg := subCommand.String("o", "output", &argparse.Options{Required: false, Help: "Image File Output"})

	return &GenCommand{
		name,
		desc,
		subCommand,
		promptArg,
		outputArg,
	}
}

func (g *GenCommand) Happened() bool {
	return g.command.Happened()
}

func (g *GenCommand) Run() error {
	args := g.command.GetArgs()
	promptParsed := args[1].GetParsed()
	outputParsed := args[2].GetParsed()

	// Input flow
	if promptParsed {
		fmt.Println(*g.promptArg)

	} else {
		errMsg := fmt.Sprintf("No input options passed to `%v`\n", g.name)
		helpMsg := g.command.Help(errMsg)

		err := errors.New(helpMsg)
		return err
	}

	// Output flow
	if outputParsed {
		fmt.Println(g.outputArg)
	}

	return nil
}
