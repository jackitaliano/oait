package assts

import (
	"errors"
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

type AsstsService struct {
	name    string
	desc    string
	command *argparse.Command

	getCommand    *GetCommand
	delCommand    *DelCommand
	createCommand *CreateCommand
}

func NewService(parser *argparse.Parser) *AsstsService {
	const name = "assts"
	const desc = "Assistants Tools"

	service := parser.NewCommand(name, desc)

	get := NewGetCommand(service)
	del := NewDelCommand(service)
	create := NewCreateCommand(service)

	return &AsstsService{
		name,
		desc,
		service,
		get,
		del,
		create,
	}
}

func (a *AsstsService) Run(key string) error {

	if a.getCommand.Happened() {
		err := a.getCommand.Run(key)

		if err != nil {
			fmt.Printf("ERROR: %v\n", err.Error())
			os.Exit(1)
		}

	} else if a.delCommand.Happened() {
		err := a.delCommand.Run(key)

		if err != nil {
			fmt.Printf("ERROR: %v\n", err.Error())
			os.Exit(1)
		}

	} else if a.createCommand.Happened() {
		err := a.createCommand.Run(key)

		if err != nil {
			fmt.Printf("ERROR: %v\n", err.Error())
			os.Exit(1)
		}

	} else {
		errMsg := fmt.Sprintf("No command given to `%v`\n", a.name)
		helpMsg := a.command.Help(errMsg)
		err := errors.New(helpMsg)
		return err
	}

	return nil
}
