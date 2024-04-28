package threadsParse

import (
	"errors"
	"fmt"

	"github.com/jackitaliano/oait/internal/io"
	"github.com/jackitaliano/oait/internal/openai"
	"github.com/jackitaliano/oait/internal/threads"
	"github.com/jackitaliano/oait/internal/tui"

	"github.com/akamensky/argparse"
)

type AddCommand struct {
	name    string
	desc    string
	command *argparse.Command

	threadArg  *string
	inputArg   *string
	orgArg     *string
	messageArg *string
	roleArg    *string
}

func NewAddCommand(command *argparse.Command) *AddCommand {
	const name = "add"
	const desc = "Add message(s) to Thread Tools"

	subCommand := command.NewCommand(name, desc)

	threadArg := subCommand.String("i", "ids", &argparse.Options{Required: true, Help: "Thread ID to add message to"})
	inputArg := subCommand.String("f", "file-input", &argparse.Options{Required: false, Help: "Thread File Input"})
	orgArg := subCommand.String("O", "org", &argparse.Options{Required: false, Help: "Set Organization ID"})
	messageArg := subCommand.String("m", "msg", &argparse.Options{Required: false, Help: "Message text to add"})
	roleArg := subCommand.String("r", "role", &argparse.Options{Required: false, Help: "Message role to add", Default: "user"})

	return &AddCommand{
		name,
		desc,
		subCommand,
		threadArg,
		inputArg,
		orgArg,
		messageArg,
		roleArg,
	}
}

func (a *AddCommand) Happened() bool {
	return a.command.Happened()
}

func (a *AddCommand) Run(key string) error {
	args := a.command.GetArgs()

	threadID, err := a.getThreadID(&args)

	if err != nil {
		return err
	}

	message, err := a.createMessage(&args, *threadID)

	if err != nil {
		return err
	}

	verify := verifyBeforeAdd()

	if verify {
		messageJson, err := threads.CreatedMessageToJson(message)

		if err != nil {
			return err
		}

		fmt.Printf("%v\n", string(*&messageJson))
	}

	confirmed := confirmAdd()

	if confirmed {
		fmt.Printf("Adding message...\t\t")
		threads.AddMessage(key, *threadID, message, *a.orgArg)
		fmt.Printf("✓\n")
	} else {
		fmt.Printf("Canceled.\n")
	}

	return nil
}

func verifyBeforeAdd() bool {
	return tui.YesNoLoop("Verify before adding message?")
}

func confirmAdd() bool {
	return tui.YesNoLoop("Confirm addition")
}

func (a *AddCommand) getThreadID(args *[]argparse.Arg) (*string, error) {
	threadParsed := (*args)[1].GetParsed()

	if threadParsed { // List passed
		threadID, err := io.SingleInput(*a.threadArg)

		if err != nil {
			return nil, err
		}

		return &threadID, nil
	}

	errMsg := fmt.Sprintf("No input options passed to `%v`\n", a.name)
	err := errors.New(errMsg)

	return nil, err
}

func (a *AddCommand) createMessage(args *[]argparse.Arg, threadID string) (*openai.CreatedMessage, error) {
	messageParsed := (*args)[4].GetParsed()

	if messageParsed {
		message := threads.CreateMessage(*a.messageArg, *a.roleArg)

		return message, nil
	}

	errMsg := fmt.Sprintf("No input options passed to `%v`\n", a.name)
	err := errors.New(errMsg)

	return nil, err
}
