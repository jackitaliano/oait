package assts

import (
	"errors"
	"fmt"

	"github.com/jackitaliano/oait/internal/filter"
	"github.com/jackitaliano/oait/internal/io"
	"github.com/jackitaliano/oait/internal/openai"
	"github.com/jackitaliano/oait/internal/tui"

	"github.com/akamensky/argparse"
)

type DelCommand struct {
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

func NewDelCommand(command *argparse.Command) *DelCommand {
	const name = "del"
	const desc = "Del Assistants Tools"

	subCommand := command.NewCommand(name, desc)

	asstsArg := subCommand.StringList("i", "ids", &argparse.Options{Required: false, Help: "List of Asst IDs"})
	inputArg := subCommand.String("f", "file-input", &argparse.Options{Required: false, Help: "Asst Asst Input"})
	allFlag := subCommand.Flag("A", "all", &argparse.Options{Required: false, Help: "Get all assts"})
	orgArg := subCommand.String("O", "org", &argparse.Options{Required: false, Help: "Set Organization ID"})
	outputArg := subCommand.String("o", "output", &argparse.Options{Required: false, Help: "Asst File Output"})
	timeLTEArg := subCommand.Float("d", "days", &argparse.Options{Required: false, Help: "Filter by LTE to days"})
	timeGTArg := subCommand.Float("D", "Days", &argparse.Options{Required: false, Help: "Filter by GT days"})
	nameContainsArg := subCommand.StringList("n", "name", &argparse.Options{Required: false, Help: "Filter by Asst containing name"})
	nameNotContainsArg := subCommand.StringList("N", "Name", &argparse.Options{Required: false, Help: "Filter by Asst not containing name"})

	return &DelCommand{
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

func (d *DelCommand) Happened() bool {

	return d.command.Happened()
}

func (d *DelCommand) Run(key string) error {
	args := d.command.GetArgs()
	allParsed := args[3].GetParsed()

	var asstObjects *[]openai.AsstObject
	var err error

	if allParsed && *d.allFlag {
		fmt.Printf("Retrieving all assts...\t\t")
		asstObjects, err = openai.RetrieveAllAssts(key, *d.orgArg)

		if err != nil {
			return err
		}

		fmt.Printf("✓\n")

	} else {

		fmt.Printf("Retrieving asst ids...\t")
		asstIDs, err := d.getAsstIDs(&args)

		if err != nil {
			fmt.Printf("X\n")
			return err
		}
		fmt.Printf("✓\n")

		fmt.Printf("Retrieving assts...\t\t")
		asstObjects = openai.RetrieveAssts(key, asstIDs, *d.orgArg)
		fmt.Printf("✓\n")
	}

	fmt.Printf("Filtering assts...\t\t")
	filteredAsstObjects, err := d.filterAssts(&args, asstObjects)

	if err != nil {
		fmt.Printf("X\n")
		return err
	}
	fmt.Printf("✓\n")

	deleteAsstIDs := getAsstIDsFromObjects(filteredAsstObjects)

	verify := verifyBeforeDelete()

	if verify {
		fmt.Printf("Formatting assts output...\t")
		asstsOutput, err := d.getAsstsOutput(&args, filteredAsstObjects)

		if err != nil {
			fmt.Printf("X\n")
			return err
		}
		fmt.Printf("✓\n")

		fmt.Printf("Outputting assts... \n\n")
		err = d.outputAssts(&args, asstsOutput)

		if err != nil {
			return err
		}
	}

	confirmed := confirmDelete()

	if confirmed {
		fmt.Printf("Deleting assts...\t\t")
		numDeleted := openai.DeleteAssts(key, deleteAsstIDs, *d.orgArg)
		fmt.Printf("✓\n")
		fmt.Printf("Deleted %v assts.\n", numDeleted)
	} else {
		fmt.Printf("Cancelled.\n")
	}

	return nil
}

func verifyBeforeDelete() bool {
	return tui.YesNoLoop("Verify assistants before deletion?")
}

func confirmDelete() bool {
	return tui.YesNoLoop("Confirm deletion")
}

func (d *DelCommand) getAsstIDs(args *[]argparse.Arg) ([]string, error) {
	asstsParsed := (*args)[1].GetParsed()
	inputParsed := (*args)[2].GetParsed()

	if asstsParsed { // List passed
		asstIDs, err := io.ListInput(*d.asstsArg)

		if err != nil {
			return nil, err
		}

		return asstIDs, nil

	}

	if inputParsed { // Asst input passed
		asstIDs, err := io.FileInput(*d.inputArg)

		if err != nil {
			return nil, err
		}

		return asstIDs, nil
	}

	errMsg := fmt.Sprintf("No input options passed to `%v`\n", d.name)
	err := errors.New(errMsg)

	return nil, err
}

func (d *DelCommand) filterAssts(args *[]argparse.Arg, asstObjects *[]openai.AsstObject) (*[]openai.AsstObject, error) {
	timeLTEParsed := (*args)[6].GetParsed()
	timeGTParsed := (*args)[7].GetParsed()
	nameContainsParsed := (*args)[8].GetParsed()
	nameNotContainsParsed := (*args)[9].GetParsed()

	filtered := asstObjects
	var err error

	if timeLTEParsed {
		filtered, err = filter.DaysLTE(filtered, *d.timeLTEArg)

		if err != nil {
			return nil, err
		}
	}

	if timeGTParsed {
		filtered, err = filter.DaysGT(filtered, *d.timeGTArg)

		if err != nil {
			return nil, err
		}
	}

	if nameContainsParsed {
		filtered = filter.ContainsName(filtered, *d.nameContainsArg)

	}

	if nameNotContainsParsed {
		filtered = filter.NotContainsName(filtered, *d.nameNotContainsArg)
	}

	return filtered, nil
}

func (d *DelCommand) getAsstsOutput(args *[]argparse.Arg, filteredAsstObjects *[]openai.AsstObject) (*[]byte, error) {

	asstsOutput, err := io.ListToJSON(filteredAsstObjects)

	if err != nil {
		return nil, err
	}

	return &asstsOutput, nil
}

func (d *DelCommand) outputAssts(args *[]argparse.Arg, output *[]byte) error {
	outputParsed := (*args)[5].GetParsed()

	if outputParsed {
		err := io.FileOutput(*d.outputArg, output)

		if err != nil {
			return err
		}

	} else {
		fmt.Printf("%v\n", string(*output))
	}

	return nil
}

func getAsstIDsFromObjects(asstObjects *[]openai.AsstObject) []string {
	asstIDs := make([]string, len(*asstObjects))

	for i, asstObject := range *asstObjects {
		asstIDs[i] = asstObject.ID
	}

	return asstIDs
}
