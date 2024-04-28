package images

import (
	"errors"
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

type ImagesService struct {
	name    string
	desc    string
	command *argparse.Command

	genCommand *GenCommand
}

func NewService(parser *argparse.Parser) *ImagesService {
	const name = "images"
	const desc = "Images Tools"

	service := parser.NewCommand(name, desc)

	gen := NewGenCommand(service)

	return &ImagesService{
		name,
		desc,
		service,
		gen,
	}
}

func (i *ImagesService) Run() error {

	if i.genCommand.Happened() {
		err := i.genCommand.Run()

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

	} else {
		errMsg := fmt.Sprintf("No command given to `%v`", i.name)
		helpMsg := i.command.Help(errMsg)
		err := errors.New(helpMsg)
		return err
	}

	return nil
}
