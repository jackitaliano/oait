package threadsParse

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackitaliano/oait-go/internal/threads"

	"github.com/akamensky/argparse"
)

type GetCommand struct {
	name    string
	desc    string
	command *argparse.Command

	threadsArg          *[]string
	inputArg            *string
	outputArg           *string
	rawFlag             *bool
	filterFlag          *bool
	filterTimeLTEFlag   *bool
	filterTimeGTFlag    *bool
	filterLengthLTEFlag *bool
	filterLengthGTFlag  *bool
	filterValue         *float64
}

func NewGetCommand(command *argparse.Command) *GetCommand {
	const name = "get"
	const desc = "Get Thread Tools"

	subCommand := command.NewCommand(name, desc)

	threadsArg := subCommand.StringList("n", "threads", &argparse.Options{Required: false, Help: "List of Thread IDs"})
	inputArg := subCommand.String("i", "input", &argparse.Options{Required: false, Help: "Thread File Input"})
	outputArg := subCommand.String("o", "output", &argparse.Options{Required: false, Help: "Thread File Output"})
	rawFlag := subCommand.Flag("r", "raw", &argparse.Options{Required: false, Help: "Output raw Threads"})
	filterFlag := subCommand.Flag("f", "filter", &argparse.Options{Required: false, Help: "Filter threads, accompanied by another filter flag\n\tEx: `-ft -v 5` yields threads that are less than or equal to 5 days old."})
	filterTimeLTEFlag := subCommand.Flag("t", "filter-time", &argparse.Options{Required: false, Help: "Filter time (days) LTE"})
	filterTimeGTFlag := subCommand.Flag("T", "filter-Time", &argparse.Options{Required: false, Help: "Filter time (days) GT"})
	filterLengthLTEFlag := subCommand.Flag("l", "filter-length", &argparse.Options{Required: false, Help: "Filter length LTE"})
	filterLengthGTFlag := subCommand.Flag("L", "filter-Length", &argparse.Options{Required: false, Help: "Filter length GT"})
	filterValue := subCommand.Float("v", "filter-value", &argparse.Options{Required: false, Help: "Filter Value (float)"})

	return &GetCommand{
		name,
		desc,
		subCommand,
		threadsArg,
		inputArg,
		outputArg,
		rawFlag,
		filterFlag,
		filterTimeLTEFlag,
		filterTimeGTFlag,
		filterLengthLTEFlag,
		filterLengthGTFlag,
		filterValue,
	}
}

func (g *GetCommand) Happened() bool {

	return g.command.Happened()
}

func (g *GetCommand) Run(key string) error {
	args := g.command.GetArgs()
	threadsParsed := args[1].GetParsed()
	inputParsed := args[2].GetParsed()
	outputParsed := args[3].GetParsed()
	rawParsed := args[4].GetParsed()
	filterParsed := args[5].GetParsed()
	filterTimeLTEParsed := args[6].GetParsed()
	filterTimeGTParsed := args[7].GetParsed()
	filterLengthLTEParsed := args[8].GetParsed()
	filterLengthGTParsed := args[9].GetParsed()
	filterValueParsed := args[10].GetParsed()

	// outputParsed := args[3].GetParsed()

	var threadIds []string
	var err error

	// Input flow
	if threadsParsed { // List passed
		threadIds, err = threads.ListInput(*g.threadsArg)

		if err != nil {
			panic(err)
		}

	} else if inputParsed { // File input passed
		threadIds, err = threads.FileInput(*g.inputArg)

		if err != nil {
			panic(err)
		}

	} else { // No input passed
		errMsg := fmt.Sprintf("No input options passed to `%v`\n", g.name)
		helpMsg := g.command.Help(errMsg)

		err := errors.New(helpMsg)
		return err
	}

	// Retrieval flow
	rawThreads := threads.RetrieveThreads(key, &threadIds)

	// Filter flow
	if filterParsed {
		if !filterValueParsed {
			err = errors.New("Error: cannot pass filter without filter value")
			panic(err)
		}

		if filterTimeLTEParsed {
			rawThreads = threads.FilterByDaysLTE(rawThreads, *g.filterValue)

		} else if filterTimeGTParsed {
			rawThreads = threads.FilterByDaysGT(rawThreads, *g.filterValue)

		}

		// Filter length flow
		if filterLengthLTEParsed {
			rawThreads = threads.FilterByLengthLTE(rawThreads, *g.filterValue)

		} else if filterLengthGTParsed {
			rawThreads = threads.FilterByLengthGT(rawThreads, *g.filterValue)

		}
	}
	var threadOutput []byte

	// Raw flow
	if rawParsed && *(g.rawFlag) {
		threadOutput, err = json.MarshalIndent(*rawThreads, "", "\t")

		if err != nil {
			fmt.Printf("Error marshalling json: %v\n", err)
			return nil
		}

	} else {
		parsedThreads := threads.ParseThreads(rawThreads)
		threadOutput, err = threads.ThreadsToJson(parsedThreads)
	}

	// Parse flow

	// Output flow

	if outputParsed {
		err = threads.FileOutput(*g.outputArg, &threadOutput)
	} else {
		fmt.Printf("%v\n", string(threadOutput))
	}

	return nil
}
