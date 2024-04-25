package filesParse

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackitaliano/oait-go/internal/files"
	"github.com/jackitaliano/oait-go/internal/openai"

	"github.com/akamensky/argparse"
)

type GetCommand struct {
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

func NewGetCommand(command *argparse.Command) *GetCommand {
	const name = "get"
	const desc = "Get Files Tools"

	subCommand := command.NewCommand(name, desc)

	filesArg := subCommand.StringList("f", "files", &argparse.Options{Required: false, Help: "List of File IDs"})
	inputArg := subCommand.String("f", "file-input", &argparse.Options{Required: false, Help: "File File Input"})
	allFlag := subCommand.Flag("A", "all", &argparse.Options{Required: false, Help: "Get all files"})
	orgArg := subCommand.String("O", "org", &argparse.Options{Required: false, Help: "Set Organization Id"})
	outputArg := subCommand.String("o", "output", &argparse.Options{Required: false, Help: "File File Output"})
	timeLTEArg := subCommand.Float("d", "days", &argparse.Options{Required: false, Help: "Filter by LTE to days"})
	timeGTArg := subCommand.Float("D", "Days", &argparse.Options{Required: false, Help: "Filter by GT days"})

	return &GetCommand{
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

func (g *GetCommand) Happened() bool {

	return g.command.Happened()
}

func (g *GetCommand) Run(key string) error {
	args := g.command.GetArgs()
	allParsed := args[3].GetParsed()

	var fileObjects *[]openai.FileObject

	if allParsed && *g.allFlag {
		fmt.Printf("Retrieving all files...\t\t")
		fileObjects = files.RetrieveAllFiles(key, *g.orgArg)
		fmt.Printf("✓\n")

	} else {

		fmt.Printf("Retrieving file ids...\t")
		fileIds, err := g.getFileIds(&args)

		if err != nil {
			return err
		}
		fmt.Printf("✓\n")

		fmt.Printf("Retrieving files...\t\t")
		fileObjects = files.RetrieveFiles(key, fileIds, *g.orgArg)
		fmt.Printf("✓\n")
	}

	fmt.Printf("Filtering files...\t\t")
	filteredFileObjects, err := g.filterFiles(&args, fileObjects)

	if err != nil {
		return err
	}
	fmt.Printf("✓\n")

	fmt.Printf("Formatting thread output...\t")
	filesOutput, err := g.getFilesOutput(&args, filteredFileObjects)

	if err != nil {
		return err
	}
	fmt.Printf("✓\n")

	fmt.Printf("Outputting files... \n\n")
	err = g.outputFiles(&args, filesOutput)

	if err != nil {
		return err
	}

	return nil
}

func (g *GetCommand) getFileIds(args *[]argparse.Arg) ([]string, error) {
	filesParsed := (*args)[1].GetParsed()
	inputParsed := (*args)[2].GetParsed()

	if filesParsed { // List passed
		fileIds, err := files.ListInput(*g.filesArg)

		if err != nil {
			return nil, err
		}

		return fileIds, nil

	}

	if inputParsed { // File input passed
		fileIds, err := files.FileInput(*g.inputArg)

		if err != nil {
			return nil, err
		}

		return fileIds, nil
	}

	errMsg := fmt.Sprintf("No input options passed to `%v`\n", g.name)
	err := errors.New(errMsg)

	return nil, err
}

func (g *GetCommand) filterFiles(args *[]argparse.Arg, fileObjects *[]openai.FileObject) (*[]openai.FileObject, error) {
	timeLTEParsed := (*args)[6].GetParsed()
	timeGTParsed := (*args)[7].GetParsed()

	if timeLTEParsed {
		filtered, err := files.FilterByDaysLTE(fileObjects, *g.timeLTEArg)

		if err != nil {
			return nil, err
		}

		return filtered, nil

	} else if timeGTParsed {
		filtered, err := files.FilterByDaysGT(fileObjects, *g.timeGTArg)

		if err != nil {
			return nil, err
		}

		return filtered, nil
	}

	return fileObjects, nil
}

func (g *GetCommand) getFilesOutput(args *[]argparse.Arg, filteredFileObjects *[]openai.FileObject) (*[]byte, error) {

	filesOutput, err := json.MarshalIndent(*filteredFileObjects, "", "\t")

	if err != nil {
		errMsg := fmt.Sprintf("Error marshalling json: %v\n", err)
		err := errors.New(errMsg)

		return nil, err
	}

	return &filesOutput, nil
}

func (g *GetCommand) outputFiles(args *[]argparse.Arg, output *[]byte) error {
	outputParsed := (*args)[4].GetParsed()

	if outputParsed {
		err := files.FileOutput(*g.outputArg, output)

		if err != nil {
			return err
		}

	} else {
		fmt.Printf("%v\n", string(*output))
	}

	return nil
}
