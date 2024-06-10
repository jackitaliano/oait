package files

import (
	"errors"
	"fmt"

	"github.com/akamensky/argparse"

	"github.com/jackitaliano/oait/internal/filter"
	"github.com/jackitaliano/oait/internal/io"
	"github.com/jackitaliano/oait/internal/openai"
	"github.com/jackitaliano/oait/internal/tui"
)

type DelCommand struct {
	name    string
	desc    string
	command *argparse.Command

	filesArg   *[]string
	inputArg   *string
	allFlag    *bool
	orgArg     *string
	outputArg  *string
	timeLTEArg *float64
	timeGTArg  *float64
	nameContainsArg    *[]string
	nameNotContainsArg *[]string
}

func NewDelCommand(command *argparse.Command) *DelCommand {
	const name = "del"
	const desc = "Del Files Tools"

	subCommand := command.NewCommand(name, desc)

	filesArg := subCommand.StringList("i", "ids", &argparse.Options{Required: false, Help: "List of File IDs"})
	inputArg := subCommand.String("f", "file-input", &argparse.Options{Required: false, Help: "File File Input"})
	allFlag := subCommand.Flag("A", "all", &argparse.Options{Required: false, Help: "Get all files"})
	orgArg := subCommand.String("O", "org", &argparse.Options{Required: false, Help: "Set Organization ID"})
	outputArg := subCommand.String("o", "output", &argparse.Options{Required: false, Help: "File File Output"})
	timeLTEArg := subCommand.Float("d", "days", &argparse.Options{Required: false, Help: "Filter by LTE to days"})
	timeGTArg := subCommand.Float("D", "Days", &argparse.Options{Required: false, Help: "Filter by GT days"})
	nameContainsArg := subCommand.StringList("n", "name", &argparse.Options{Required: false, Help: "Filter by File containing name"})
	nameNotContainsArg := subCommand.StringList("N", "Name", &argparse.Options{Required: false, Help: "Filter by File not containing name"})

	return &DelCommand{
		name,
		desc,
		subCommand,
		filesArg,
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

	var fileObjects *[]openai.FileObject

	if allParsed && *d.allFlag {
		fmt.Printf("Retrieving all files...\t\t")
		fileObjects = openai.RetrieveAllFiles(key, *d.orgArg)
		fmt.Printf("✓\n")

	} else {

		fmt.Printf("Retrieving file ids...\t")
		fileIDs, err := d.getFileIDs(&args)

		if err != nil {
			fmt.Printf("X\n")
			return err
		}
		fmt.Printf("✓\n")

		fmt.Printf("Retrieving files...\t\t")
		fileObjects = openai.RetrieveFiles(key, fileIDs, *d.orgArg)
		fmt.Printf("✓\n")
	}

	fmt.Printf("Filtering files...\t\t")
	filteredFileObjects, err := d.filterFiles(&args, fileObjects)

	if err != nil {
		fmt.Printf("X\n")
		return err
	}
	fmt.Printf("✓\n")

	deleteFileIDs := getFileIDsFromObjects(filteredFileObjects)

	verify := verifyBeforeDelete()

	if verify {
		fmt.Printf("Formatting files output...\t")
		filesOutput, err := d.getFilesOutput(&args, filteredFileObjects)

		if err != nil {
			fmt.Printf("X\n")
			return err
		}
		fmt.Printf("✓\n")

		fmt.Printf("Outputting files... \n\n")
		err = d.outputFiles(&args, filesOutput)

		if err != nil {
			return err
		}
	}

	confirmed := confirmDelete()

	if confirmed {
		fmt.Printf("Deleting files...\t\t")
		numDeleted := openai.DeleteFiles(key, deleteFileIDs, *d.orgArg)
		fmt.Printf("✓\n")
		fmt.Printf("Deleted %v files.\n", numDeleted)
	} else {
		fmt.Printf("Cancelled.\n")
	}

	return nil
}

func verifyBeforeDelete() bool {
	return tui.YesNoLoop("Verify files before deletion?")
}

func confirmDelete() bool {
	return tui.YesNoLoop("Confirm deletion")
}

func (d *DelCommand) getFileIDs(args *[]argparse.Arg) ([]string, error) {
	filesParsed := (*args)[1].GetParsed()
	inputParsed := (*args)[2].GetParsed()

	if filesParsed { // List passed
		fileIDs, err := io.ListInput(*d.filesArg)

		if err != nil {
			return nil, err
		}

		return fileIDs, nil

	}

	if inputParsed { // File input passed
		fileIDs, err := io.FileInput(*d.inputArg)

		if err != nil {
			return nil, err
		}

		return fileIDs, nil
	}

	errMsg := fmt.Sprintf("No input options passed to `%v`\n", d.name)
	err := errors.New(errMsg)

	return nil, err
}

func (d *DelCommand) filterFiles(args *[]argparse.Arg, fileObjects *[]openai.FileObject) (*[]openai.FileObject, error) {
	timeLTEParsed := (*args)[6].GetParsed()
	timeGTParsed := (*args)[7].GetParsed()
	nameContainsParsed := (*args)[8].GetParsed()
	nameNotContainsParsed := (*args)[9].GetParsed()

	filtered := fileObjects
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

func (d *DelCommand) getFilesOutput(args *[]argparse.Arg, filteredFileObjects *[]openai.FileObject) (*[]byte, error) {

	filesOutput, err := io.ListToJSON(filteredFileObjects)

	if err != nil {
		return nil, err
	}

	return &filesOutput, nil
}

func (d *DelCommand) outputFiles(args *[]argparse.Arg, output *[]byte) error {
	outputParsed := (*args)[4].GetParsed()

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

func getFileIDsFromObjects(fileObjects *[]openai.FileObject) []string {
	fileIDs := make([]string, len(*fileObjects))

	for i, fileObject := range *fileObjects {
		fileIDs[i] = fileObject.ID
	}

	return fileIDs
}
