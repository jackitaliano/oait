package files

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

func NewGetCommand(command *argparse.Command) *GetCommand {
	const name = "get"
	const desc = "Get Files Tools"

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

	var fileObjects *[]openai.FileObject

	if allParsed && *g.allFlag {
		fmt.Printf("Retrieving all files...\t\t")
		fileObjects = openai.RetrieveAllFiles(key, *g.orgArg)
		fmt.Printf("✓\n")

	} else {

		fmt.Printf("Retrieving file ids...\t")
		fileIDs, err := g.getFileIDs(&args)

		if err != nil {
			fmt.Printf("X\n")
			return err
		}
		fmt.Printf("✓\n")

		fmt.Printf("Retrieving files...\t\t")
		fileObjects = openai.RetrieveFiles(key, fileIDs, *g.orgArg)
		fmt.Printf("✓\n")
	}

	fmt.Printf("Filtering files...\t\t")
	filteredFileObjects, err := g.filterFiles(&args, fileObjects)

	if err != nil {
		fmt.Printf("X\n")
		return err
	}
	fmt.Printf("✓\n")

	fmt.Printf("Formatting files output...\t")
	filesOutput, err := g.getFilesOutput(&args, filteredFileObjects)

	if err != nil {
		fmt.Printf("X\n")
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

func (g *GetCommand) getFileIDs(args *[]argparse.Arg) ([]string, error) {
	filesParsed := (*args)[1].GetParsed()
	inputParsed := (*args)[2].GetParsed()

	if filesParsed { // List passed
		fileIDs, err := io.ListInput(*g.filesArg)

		if err != nil {
			return nil, err
		}

		return fileIDs, nil

	}

	if inputParsed { // File input passed
		fileIDs, err := io.FileInput(*g.inputArg)

		if err != nil {
			return nil, err
		}

		return fileIDs, nil
	}

	errMsg := fmt.Sprintf("No input options passed to `%v`\n", g.name)
	err := errors.New(errMsg)

	return nil, err
}

func (g *GetCommand) filterFiles(args *[]argparse.Arg, fileObjects *[]openai.FileObject) (*[]openai.FileObject, error) {
	timeLTEParsed := (*args)[6].GetParsed()
	timeGTParsed := (*args)[7].GetParsed()
	nameContainsParsed := (*args)[8].GetParsed()
	nameNotContainsParsed := (*args)[9].GetParsed()

	filtered := fileObjects
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

func (g *GetCommand) getFilesOutput(args *[]argparse.Arg, filteredFileObjects *[]openai.FileObject) (*[]byte, error) {

	filesOutput, err := io.ListToJSON(filteredFileObjects)

	if err != nil {
		return nil, err
	}

	return &filesOutput, nil
}

func (g *GetCommand) outputFiles(args *[]argparse.Arg, output *[]byte) error {
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
