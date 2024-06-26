package threads

import (
	"errors"
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

type ThreadsService struct {
	name    string
	desc    string
	command *argparse.Command

	getCommand *GetCommand
	delCommand *DelCommand
	addCommand *AddCommand
}

func NewService(parser *argparse.Parser) *ThreadsService {
	const name = "threads"
	const desc = "Thread Tools"

	service := parser.NewCommand(name, desc)

	get := NewGetCommand(service)
	del := NewDelCommand(service)
	add := NewAddCommand(service)

	return &ThreadsService{
		name,
		desc,
		service,
		get,
		del,
		add,
	}
}

func (t *ThreadsService) Run(key string) error {

	if t.getCommand.Happened() {
		err := t.getCommand.Run(key)

		if err != nil {
			fmt.Printf("ERROR: %v", err.Error())
			os.Exit(1)
		}

	} else if t.delCommand.Happened() {
		err := t.delCommand.Run(key)

		if err != nil {
			fmt.Printf("ERROR: %v", err.Error())
			os.Exit(1)
		}

	} else if t.addCommand.Happened() {
		err := t.addCommand.Run(key)

		if err != nil {
			fmt.Printf("ERROR: %v", err.Error())
			os.Exit(1)
		}

	} else {
		errMsg := fmt.Sprintf("No command given to `%v`\n", t.name)
		helpMsg := t.command.Help(errMsg)
		err := errors.New(helpMsg)
		return err
	}

	return nil
}
