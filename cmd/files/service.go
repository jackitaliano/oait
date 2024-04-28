package files

import (
	"errors"
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

type FilesService struct {
	name    string
	desc    string
	command *argparse.Command

	getCommand *GetCommand
	delCommand *DelCommand
}

func NewService(parser *argparse.Parser) *FilesService {
	const name = "files"
	const desc = "Files Tools"

	service := parser.NewCommand(name, desc)

	get := NewGetCommand(service)
	del := NewDelCommand(service)

	return &FilesService{
		name,
		desc,
		service,
		get,
		del,
	}
}

func (f *FilesService) Run(key string) error {

	if f.getCommand.Happened() {
		err := f.getCommand.Run(key)

		if err != nil {
			fmt.Printf("ERROR: %v", err.Error())
			os.Exit(1)
		}

	} else if f.delCommand.Happened() {
		err := f.delCommand.Run(key)

		if err != nil {
			fmt.Printf("ERROR: %v", err.Error())
			os.Exit(1)
		}

	} else {
		errMsg := fmt.Sprintf("No command given to `%v`\n", f.name)
		helpMsg := f.command.Help(errMsg)
		err := errors.New(helpMsg)
		return err
	}

	return nil
}
