package assts

import (
	"errors"
	"fmt"

	"github.com/jackitaliano/oait/internal/filter"
	"github.com/jackitaliano/oait/internal/io"
	"github.com/jackitaliano/oait/internal/openai"

	"github.com/akamensky/argparse"
)

type GetCommand struct {
	name    string
	desc    string
	command *argparse.Command

	asstsArg           *[]string
	inputArg           *string
	allFlag            *bool
	orgArg             *string
	outputArg          *string
	timeLTEArg         *float64
	timeGTArg          *float64
	nameContainsArg    *[]string
	nameNotContainsArg *[]string
}

func NewGetCommand(command *argparse.Command) *GetCommand {
	const name = "get"
	const desc = "Get Assistants Tools"

	subCommand := command.NewCommand(name, desc)

	asstsArg := subCommand.StringList("i", "ids", &argparse.Options{Required: false, Help: "List of Asst IDs"})
	inputArg := subCommand.String("f", "file-input", &argparse.Options{Required: false, Help: "Asst File Input (of ids)"})
	allFlag := subCommand.Flag("A", "all", &argparse.Options{Required: false, Help: "Get all assts"})
	orgArg := subCommand.String("O", "org", &argparse.Options{Required: false, Help: "Set Organization ID"})
	outputArg := subCommand.String("o", "output", &argparse.Options{Required: false, Help: "Asst File Output"})
	timeLTEArg := subCommand.Float("d", "days", &argparse.Options{Required: false, Help: "Filter by LTE to days"})
	timeGTArg := subCommand.Float("D", "Days", &argparse.Options{Required: false, Help: "Filter by GT days"})
	nameContainsArg := subCommand.StringList("n", "name", &argparse.Options{Required: false, Help: "Filter by Asst containing name"})
	nameNotContainsArg := subCommand.StringList("N", "Name", &argparse.Options{Required: false, Help: "Filter by Asst not containing name"})

	return &GetCommand{
		name,
		desc,
		subCommand,
		asstsArg,
		inputArg,
		allFlag,
		orgArg,
		outputArg,
		timeLTEArg,
		timeGTArg,
		nameContainsArg,
		nameNotContainsArg,
	}
}

func (g *GetCommand) Happened() bool {

	return g.command.Happened()
}

func (g *GetCommand) Run(key string) error {
	args := g.command.GetArgs()
	allParsed := args[3].GetParsed()

	var asstObjects *[]openai.AsstObject
	var err error

	if allParsed && *g.allFlag {
		fmt.Printf("Retrieving all assts...\t\t")
		asstObjects, err = openai.RetrieveAllAssts(key, *g.orgArg)

		if err != nil {
			return err
		}

		fmt.Printf("✓\n")

	} else {

		fmt.Printf("Retrieving asst ids...\t")
		asstIDs, err := g.getAsstIDs(&args)

		if err != nil {
			fmt.Printf("X\n")
			return err
		}
		fmt.Printf("✓\n")

		fmt.Printf("Retrieving assts...\t\t")
		asstObjects = openai.RetrieveAssts(key, asstIDs, *g.orgArg)
		fmt.Printf("✓\n")
	}

	fmt.Printf("Filtering assts...\t\t")
	filteredAsstObjects, err := g.filterAssts(&args, asstObjects)

	if err != nil {
		fmt.Printf("X\n")
		return err
	}
	fmt.Printf("✓\n")

	fmt.Printf("Formatting assts output...\t")
	asstsOutput, err := g.getAsstsOutput(&args, filteredAsstObjects)

	if err != nil {
		fmt.Printf("X\n")
		return err
	}
	fmt.Printf("✓\n")

	fmt.Printf("Outputting assts... \n\n")
	err = g.outputAssts(&args, asstsOutput)

	if err != nil {
		return err
	}

	return nil
}

func (g *GetCommand) getAsstIDs(args *[]argparse.Arg) ([]string, error) {
	asstsParsed := (*args)[1].GetParsed()
	inputParsed := (*args)[2].GetParsed()

	if asstsParsed { // List passed
		asstIDs, err := io.ListInput(*g.asstsArg)

		if err != nil {
			return nil, err
		}

		return asstIDs, nil

	}

	if inputParsed { // Asst input passed
		asstIDs, err := io.FileInput(*g.inputArg)

		if err != nil {
			return nil, err
		}

		return asstIDs, nil
	}

	errMsg := fmt.Sprintf("No input options passed to `%v`\n", g.name)
	err := errors.New(errMsg)

	return nil, err
}

func (g *GetCommand) filterAssts(args *[]argparse.Arg, asstObjects *[]openai.AsstObject) (*[]openai.AsstObject, error) {
	timeLTEParsed := (*args)[6].GetParsed()
	timeGTParsed := (*args)[7].GetParsed()
	nameContainsParsed := (*args)[8].GetParsed()
	nameNotContainsParsed := (*args)[9].GetParsed()

	filtered := asstObjects
	var err error

	if timeLTEParsed {
		filtered, err = filter.DaysLTE(filtered, *g.timeLTEArg)

		if err != nil {
			return nil, err
		}
	}

	if timeGTParsed {
		filtered, err = filter.DaysGT(filtered, *g.timeGTArg)

		if err != nil {
			return nil, err
		}
	}

	if nameContainsParsed {
		filtered = filter.ContainsName(filtered, *g.nameContainsArg)

	}

	if nameNotContainsParsed {
		filtered = filter.NotContainsName(filtered, *g.nameNotContainsArg)
	}

	return filtered, nil
}

func (g *GetCommand) getAsstsOutput(args *[]argparse.Arg, filteredAsstObjects *[]openai.AsstObject) (*[]byte, error) {

	asstsOutput, err := io.ListToJSON(filteredAsstObjects)

	if err != nil {
		return nil, err
	}

	return &asstsOutput, nil
}

func (g *GetCommand) outputAssts(args *[]argparse.Arg, output *[]byte) error {
	outputParsed := (*args)[5].GetParsed()

	if outputParsed {
		err := io.FileOutput(*g.outputArg, output)

		if err != nil {
			return err
		}

	} else {
		fmt.Printf("%v\n", string(*output))
	}

	return nil
}
