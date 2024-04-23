package threadsParse

import (
	"errors"
	"fmt"

	"github.com/akamensky/argparse"
)

type DelCommand struct {
	name    string
	desc    string
	command *argparse.Command

	threadsArg *[]string
	inputArg   *string
}

func NewDelCommand(command *argparse.Command) *DelCommand {
	const name = "del"
	const desc = "Del Thread Tools"

	subCommand := command.NewCommand(name, desc)

	threadsArg := subCommand.List("t", "threads", &argparse.Options{Required: false, Help: "List of Thread IDs"})
	inputArg := subCommand.String("i", "input", &argparse.Options{Required: false, Help: "Thread File Input"})

	return &DelCommand{
		name,
		desc,
		subCommand,
		threadsArg,
		inputArg,
	}
}

func (d *DelCommand) Happened() bool {
	return d.command.Happened()
}

func (d *DelCommand) Run(key string) error {
	args := d.command.GetArgs()
	threadsParsed := args[1].GetParsed()
	inputParsed := args[2].GetParsed()

	// Input flow
	if threadsParsed {
		fmt.Println(d.threadsArg)

	} else if inputParsed {
		fmt.Println(d.inputArg)

	} else {
		errMsg := fmt.Sprintf("No input options passed to `%v`\n", d.name)
		helpMsg := d.command.Help(errMsg)

		err := errors.New(helpMsg)
		return err
	}

	return nil
}
