package threadsParse

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackitaliano/oait-go/internal/openai"
	"github.com/jackitaliano/oait-go/internal/threads"

	"github.com/akamensky/argparse"
)

type GetCommand struct {
	name    string
	desc    string
	command *argparse.Command

	threadsArg   *[]string
	inputArg     *string
	sessionArg   *string
	orgArg       *string
	outputArg    *string
	rawFlag      *bool
	timeLTEArg   *float64
	timeGTArg    *float64
	lengthLTEArg *float64
	lengthGTArg  *float64
}

func NewGetCommand(command *argparse.Command) *GetCommand {
	const name = "get"
	const desc = "Get Thread Tools"

	subCommand := command.NewCommand(name, desc)

	threadsArg := subCommand.StringList("t", "threads", &argparse.Options{Required: false, Help: "List of Thread IDs"})
	inputArg := subCommand.String("f", "file-input", &argparse.Options{Required: false, Help: "Thread File Input"})
	sessionArg := subCommand.String("s", "session", &argparse.Options{Required: false, Help: "Retrieve Threads from session-id"})
	orgArg := subCommand.String("O", "org", &argparse.Options{Required: false, Help: "Set Organization Id"})
	outputArg := subCommand.String("o", "output", &argparse.Options{Required: false, Help: "Thread File Output"})
	rawFlag := subCommand.Flag("r", "raw", &argparse.Options{Required: false, Help: "Output raw Threads"})
	timeLTEArg := subCommand.Float("d", "days", &argparse.Options{Required: false, Help: "Filter by LTE to days"})
	timeGTArg := subCommand.Float("D", "Days", &argparse.Options{Required: false, Help: "Filter by GT days"})
	lengthLTEArg := subCommand.Float("l", "length", &argparse.Options{Required: false, Help: "Filter by LTE to length"})
	lengthGTArg := subCommand.Float("L", "Length", &argparse.Options{Required: false, Help: "Filter by GT length"})

	return &GetCommand{
		name,
		desc,
		subCommand,
		threadsArg,
		inputArg,
		sessionArg,
		orgArg,
		outputArg,
		rawFlag,
		timeLTEArg,
		timeGTArg,
		lengthLTEArg,
		lengthGTArg,
	}
}

func (g *GetCommand) Happened() bool {

	return g.command.Happened()
}

func (g *GetCommand) Run(key string) error {
	args := g.command.GetArgs()

	fmt.Printf("Retrieving thread ids...\t")
	threadIds, err := g.getThreadIds(&args)

	if err != nil {
		return err
	}

	fmt.Printf("✓\n")

	fmt.Printf("Retrieving threads...\t\t")
	rawThreads := threads.RetrieveThreads(key, threadIds, *g.orgArg)
	fmt.Printf("✓\n")

	fmt.Printf("Filtering threads...\t\t")
	filteredThreads := g.filterThreads(&args, rawThreads)
	fmt.Printf("✓\n")

	fmt.Printf("Formatting thread output...\t")
	threadsOutput, err := g.getThreadsOutput(&args, threadIds, filteredThreads)

	if err != nil {
		return err
	}
	fmt.Printf("✓\n")

	fmt.Printf("Outputting threads... \n\n")
	err = g.outputThreads(&args, threadsOutput)

	if err != nil {
		return err
	}

	return nil
}

func (g *GetCommand) getThreadIds(args *[]argparse.Arg) ([]string, error) {
	threadsParsed := (*args)[1].GetParsed()
	inputParsed := (*args)[2].GetParsed()
	sessionParsed := (*args)[3].GetParsed()

	if threadsParsed { // List passed
		threadIds, err := threads.ListInput(*g.threadsArg)

		if err != nil {
			return nil, err
		}

		return threadIds, nil

	}

	if inputParsed { // File input passed
		threadIds, err := threads.FileInput(*g.inputArg)

		if err != nil {
			return nil, err
		}

		return threadIds, nil
	}

	if sessionParsed {
		threadIds, err := threads.SessionInput(*g.sessionArg, *g.orgArg)

		if err != nil {
			return nil, err
		}

		return threadIds, nil

	}

	errMsg := fmt.Sprintf("No input options passed to `%v`\n", g.name)
	err := errors.New(errMsg)

	return nil, err
}

func (g *GetCommand) filterThreads(args *[]argparse.Arg, rawThreads *[][]openai.Message) *[][]openai.Message {
	timeLTEParsed := (*args)[7].GetParsed()
	timeGTParsed := (*args)[8].GetParsed()
	lengthLTEParsed := (*args)[9].GetParsed()
	lengthGTParsed := (*args)[10].GetParsed()

	if timeLTEParsed {
		return threads.FilterByDaysLTE(rawThreads, *g.timeLTEArg)

	} else if timeGTParsed {
		return threads.FilterByDaysGT(rawThreads, *g.timeGTArg)
	}

	// Filter length flow
	if lengthLTEParsed {
		return threads.FilterByLengthLTE(rawThreads, *g.lengthLTEArg)

	} else if lengthGTParsed {
		return threads.FilterByLengthGT(rawThreads, *g.lengthGTArg)
	}

	return rawThreads
}

func (g *GetCommand) getThreadsOutput(args *[]argparse.Arg, threadIds []string, filteredThreads *[][]openai.Message) (*[]byte, error) {
	rawParsed := (*args)[6].GetParsed()

	if rawParsed && *(g.rawFlag) {
		threadOutput, err := json.MarshalIndent(*filteredThreads, "", "\t")

		if err != nil {
			errMsg := fmt.Sprintf("Error marshalling json: %v\n", err)
			err := errors.New(errMsg)

			return nil, err
		}

		return &threadOutput, nil

	}

	parsedThreads := threads.ParseThreads(threadIds, filteredThreads)
	threadOutput, err := threads.ThreadsToJson(parsedThreads)

	if err != nil {
		return nil, err
	}

	return &threadOutput, nil
}

func (g *GetCommand) outputThreads(args *[]argparse.Arg, output *[]byte) error {
	outputParsed := (*args)[5].GetParsed()

	if outputParsed {
		err := threads.FileOutput(*g.outputArg, output)

		if err != nil {
			return err
		}

	} else {
		fmt.Printf("%v\n", string(*output))
	}

	return nil
}
