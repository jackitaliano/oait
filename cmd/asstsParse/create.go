package asstsParse

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackitaliano/oait/internal/assts"
	"github.com/jackitaliano/oait/internal/openai"
	"github.com/jackitaliano/oait/internal/tui"

	"github.com/akamensky/argparse"
)

type CreateCommand struct {
	name    string
	desc    string
	command *argparse.Command

	asstName *string
	asstDesc *string
	instruct *string
	model    *string
	temp     *float64
	topP     *float64
	resForm  *string
	inputArg *string
	orgArg   *string
}

func NewCreateCommand(command *argparse.Command) *CreateCommand {
	const name = "create"
	const desc = "Create Assistants Tools"

	subCommand := command.NewCommand(name, desc)

	asstName := subCommand.String("n", "name", &argparse.Options{Required: false, Help: "Name of assistant", Default: ""})
	asstDes := subCommand.String("d", "desc", &argparse.Options{Required: false, Help: "Description of assistant", Default: ""})
	instruct := subCommand.String("i", "instruct", &argparse.Options{Required: false, Help: "Instructions for assistant", Default: ""})
	model := subCommand.String("m", "model", &argparse.Options{Required: false, Help: "OpenAI Model for assistant", Default: "gpt-3.5-turbo"})
	temp := subCommand.Float("t", "temp", &argparse.Options{Required: false, Help: "Temperature of assistant <0.0 - 2.0>", Default: 1.0})
	topP := subCommand.Float("T", "topp", &argparse.Options{Required: false, Help: "Top P of assistant <0.0 - 1.0>", Default: 1.0})
	resForm := subCommand.String("r", "resformat", &argparse.Options{Required: false, Help: "Response format of assistant <'json_object' | 'auto'>", Default: "auto"})
	inputArg := subCommand.String("f", "file-input", &argparse.Options{Required: false, Help: "Asst File Input"})
	orgArg := subCommand.String("O", "org", &argparse.Options{Required: false, Help: "Set Organization Id"})

	return &CreateCommand{
		name,
		desc,
		subCommand,
		asstName,
		asstDes,
		instruct,
		model,
		temp,
		topP,
		resForm,
		inputArg,
		orgArg,
	}
}

func (c *CreateCommand) Happened() bool {

	return c.command.Happened()
}

func (c *CreateCommand) Run(key string) error {
	args := c.command.GetArgs()

	createdAsst, err := c.getCreatedAssistant(&args)

	if err != nil {
		return err
	}

	verify := verifyBeforeCreate()

	if verify {
		fmt.Printf("Formatting assts output...\t")
		asstsOutput, err := c.getCreatedAsstOutput(&args, createdAsst)

		if err != nil {
			fmt.Printf("X\n")
			return err
		}
		fmt.Printf("✓\n")

		fmt.Printf("Outputting assts... \n\n")
		err = c.outputAsst(&args, asstsOutput)

		if err != nil {
			return err
		}
	}

	confirmed := confirmCreate()

	if confirmed {
		fmt.Printf("Creating assistant...\t\t")
		asstObject, err := assts.CreateAssistant(key, createdAsst, *c.orgArg)

		if err != nil {
			fmt.Printf("X\n")
			return err
		}

		fmt.Printf("✓\n")
		fmt.Printf("Formatting assts output...\t")
		asstsOutput, err := c.getAsstOutput(&args, asstObject)

		if err != nil {
			fmt.Printf("X\n")
			return err
		}
		fmt.Printf("✓\n")

		fmt.Printf("Outputting assts... \n\n")
		err = c.outputAsst(&args, asstsOutput)

		if err != nil {
			return err
		}

	} else {
		fmt.Printf("Cancelled.\n")
	}

	return nil
}

func verifyBeforeCreate() bool {
	return tui.YesNoLoop("Verify assistants before creation?")
}

func confirmCreate() bool {
	return tui.YesNoLoop("Confirm creation")
}
func (c *CreateCommand) getCreatedAssistant(args *[]argparse.Arg) (*openai.CreatedAssistant, error) {
	inputParsed := (*args)[8].GetParsed()

	if inputParsed {
		createdAssistant, err := assts.JsonInput[openai.CreatedAssistant](*c.inputArg)

		if err != nil {
			return nil, err
		}

		if createdAssistant.Model == "" {
			err := errors.New("Must provide a model name.")
			return nil, err
		}

		return createdAssistant, nil
	}

	if *c.model == "" {
		err := errors.New("Must provide a model name.")
		return nil, err
	}

	createdAssistant := openai.CreatedAssistant{
		Name:         *c.asstName,
		Description:  *c.asstDesc,
		Instructions: *c.instruct,
		Model:        *c.model,
		Temp:         *c.temp,
		TopP:         *c.topP,
		ResFormat:    *c.resForm,
	}

	return &createdAssistant, nil
}

func (c *CreateCommand) getCreatedAsstOutput(args *[]argparse.Arg, asstObject *openai.CreatedAssistant) (*[]byte, error) {

	asstsOutput, err := json.MarshalIndent(*asstObject, "", "\t")

	if err != nil {
		errMsg := fmt.Sprintf("Error marshalling json: %v\n", err)
		err := errors.New(errMsg)

		return nil, err
	}

	return &asstsOutput, nil
}

func (c *CreateCommand) getAsstOutput(args *[]argparse.Arg, asstObject *openai.AsstObject) (*[]byte, error) {

	asstsOutput, err := json.MarshalIndent(*asstObject, "", "\t")

	if err != nil {
		errMsg := fmt.Sprintf("Error marshalling json: %v\n", err)
		err := errors.New(errMsg)

		return nil, err
	}

	return &asstsOutput, nil
}

func (c *CreateCommand) outputAsst(args *[]argparse.Arg, output *[]byte) error {

	fmt.Printf("%v\n", string(*output))

	return nil
}
