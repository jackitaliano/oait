package filesParse

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackitaliano/oait-go/internal/files"
	"github.com/jackitaliano/oait-go/internal/openai"
	"github.com/jackitaliano/oait-go/internal/tui"

	"github.com/akamensky/argparse"
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
}

func NewDelCommand(command *argparse.Command) *DelCommand {
	const name = "del"
	const desc = "Del Files Tools"

	subCommand := command.NewCommand(name, desc)

	filesArg := subCommand.StringList("i", "ids", &argparse.Options{Required: false, Help: "List of File IDs"})
	inputArg := subCommand.String("f", "file-input", &argparse.Options{Required: false, Help: "File File Input"})
	allFlag := subCommand.Flag("A", "all", &argparse.Options{Required: false, Help: "Get all files"})
	orgArg := subCommand.String("O", "org", &argparse.Options{Required: false, Help: "Set Organization Id"})
	outputArg := subCommand.String("o", "output", &argparse.Options{Required: false, Help: "File File Output"})
	timeLTEArg := subCommand.Float("d", "days", &argparse.Options{Required: false, Help: "Filter by LTE to days"})
	timeGTArg := subCommand.Float("D", "Days", &argparse.Options{Required: false, Help: "Filter by GT days"})

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
	}
}

func (d *DelCommand) Happened() bool {
	return d.command.Happened()
}

func (d *DelCommand) Run(key string) error {
	args := d.command.GetArgs()
	allParsed := args[3].GetParsed()

	var fileObjects *[]openai.FileObject
	var fileIds []string
	var err error

	if allParsed && *d.allFlag {
		fmt.Printf("Retrieving all files...\t\t")
		fileObjects = files.RetrieveAllFiles(key, *d.orgArg)
		fileIds = getFileIdsFromObjects(fileObjects)
		fmt.Printf("✓\n")

	} else {

		fmt.Printf("Retrieving file ids...\t")
		fileIds, err = d.getFileIds(&args)

		if err != nil {
			return err
		}
		fmt.Printf("✓\n")

		fmt.Printf("Retrieving files...\t\t")
		fileObjects = files.RetrieveFiles(key, fileIds, *d.orgArg)
		fmt.Printf("✓\n")
	}

	fmt.Printf("Filtering files...\t\t")
	filteredFileObjects, err := d.filterFiles(&args, fileObjects)

	if err != nil {
		return err
	}
	fmt.Printf("✓\n")

	verify := verifyBeforeDelete()

	if verify {
		fmt.Printf("Formatting thread output...\t")
		filesOutput, err := d.getFilesOutput(&args, filteredFileObjects)

		if err != nil {
			return err
		}
		fmt.Printf("✓\n")

		fmt.Printf("Outputting threads... \n\n")
		err = d.outputFiles(&args, filesOutput)

		if err != nil {
			return err
		}
	}

	confirmed := confirmDelete()

	if confirmed {
		fmt.Printf("Deleting threads...\t\t")
		numDeleted := files.DeleteFiles(key, fileIds, *d.orgArg)
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

func (d *DelCommand) getFileIds(args *[]argparse.Arg) ([]string, error) {
	filesParsed := (*args)[1].GetParsed()
	inputParsed := (*args)[2].GetParsed()

	if filesParsed { // List passed
		fileIds, err := files.ListInput(*d.filesArg)

		if err != nil {
			return nil, err
		}

		return fileIds, nil

	}

	if inputParsed { // File input passed
		fileIds, err := files.FileInput(*d.inputArg)

		if err != nil {
			return nil, err
		}

		return fileIds, nil
	}

	errMsg := fmt.Sprintf("No input options passed to `%v`\n", d.name)
	err := errors.New(errMsg)

	return nil, err
}

func (d *DelCommand) filterFiles(args *[]argparse.Arg, fileObjects *[]openai.FileObject) (*[]openai.FileObject, error) {
	timeLTEParsed := (*args)[6].GetParsed()
	timeGTParsed := (*args)[7].GetParsed()

	if timeLTEParsed {
		filtered, err := files.FilterByDaysLTE(fileObjects, *d.timeLTEArg)

		if err != nil {
			return nil, err
		}

		return filtered, nil

	} else if timeGTParsed {
		filtered, err := files.FilterByDaysGT(fileObjects, *d.timeGTArg)

		if err != nil {
			return nil, err
		}

		return filtered, nil
	}

	return fileObjects, nil
}

func (d *DelCommand) getFilesOutput(args *[]argparse.Arg, filteredFileObjects *[]openai.FileObject) (*[]byte, error) {

	filesOutput, err := json.MarshalIndent(*filteredFileObjects, "", "\t")

	if err != nil {
		errMsg := fmt.Sprintf("Error marshalling json: %v\n", err)
		err := errors.New(errMsg)

		return nil, err
	}

	return &filesOutput, nil
}

func (d *DelCommand) outputFiles(args *[]argparse.Arg, output *[]byte) error {
	outputParsed := (*args)[4].GetParsed()

	if outputParsed {
		err := files.FileOutput(*d.outputArg, output)

		if err != nil {
			return err
		}

	} else {
		fmt.Printf("%v\n", string(*output))
	}

	return nil
}

func getFileIdsFromObjects(fileObjects *[]openai.FileObject) ([]string) {
	fileIds := make([]string, len(*fileObjects))

	for i, fileObject := range *fileObjects {
		fileIds[i] = fileObject.Id
	}

	return fileIds
}
