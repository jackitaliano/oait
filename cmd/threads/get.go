package threads

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

	threadsArg   *[]string
	inputArg     *string
	sessionArg   *string
	orgArg       *string
	outputArg    *string
	prettyFlag   *bool
	timeLTEArg   *float64
	timeGTArg    *float64
	lengthLTEArg *float64
	lengthGTArg  *float64
}

func NewGetCommand(command *argparse.Command) *GetCommand {
	const name = "get"
	const desc = "Get Thread Tools"

	subCommand := command.NewCommand(name, desc)

	threadsArg := subCommand.StringList("i", "ids", &argparse.Options{Required: false, Help: "List of Thread IDs"})
	inputArg := subCommand.String("f", "file-input", &argparse.Options{Required: false, Help: "Thread File Input"})
	sessionArg := subCommand.String("s", "session", &argparse.Options{Required: false, Help: "Retrieve Threads from session-id"})
	orgArg := subCommand.String("O", "org", &argparse.Options{Required: false, Help: "Set Organization ID"})
	outputArg := subCommand.String("o", "output", &argparse.Options{Required: false, Help: "Thread File Output"})
	prettyFlag := subCommand.Flag("p", "pretty", &argparse.Options{Required: false, Help: "Pretty print threads"})
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
		prettyFlag,
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
	threadIDs, err := g.getThreadIDs(&args)

	if err != nil {
		fmt.Printf("X\n")
		return err
	}

	fmt.Printf("✓\n")

	fmt.Printf("Retrieving threads...\t\t")
	rawThreads := openai.RetrieveThreads(key, threadIDs, *g.orgArg)
	fmt.Printf("✓\n")

	fmt.Printf("Filtering threads...\t\t")
	filteredThreads, err := g.filterThreads(&args, rawThreads)

	if err != nil {
		fmt.Printf("X\n")
		return err
	}
	fmt.Printf("✓\n")

	fmt.Printf("Formatting thread output...\t")
	threadsOutput, err := g.getThreadsOutput(&args, threadIDs, filteredThreads)

	if err != nil {
		fmt.Printf("X\n")
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

func (g *GetCommand) getThreadIDs(args *[]argparse.Arg) ([]string, error) {
	threadsParsed := (*args)[1].GetParsed()
	inputParsed := (*args)[2].GetParsed()
	sessionParsed := (*args)[3].GetParsed()

	if threadsParsed { // List passed
		threadIDs, err := io.ListInput(*g.threadsArg)

		if err != nil {
			return nil, err
		}

		return threadIDs, nil

	}

	if inputParsed { // File input passed
		threadIDs, err := io.FileInput(*g.inputArg)

		if err != nil {
			return nil, err
		}

		return threadIDs, nil
	}

	if sessionParsed {
		threadIDs, err := io.SessionInput(*g.sessionArg, *g.orgArg)

		if err != nil {
			return nil, err
		}

		return threadIDs, nil

	}

	errMsg := fmt.Sprintf("No input options passed to `%v`\n", g.name)
	err := errors.New(errMsg)

	return nil, err
}

func (g *GetCommand) filterThreads(args *[]argparse.Arg, rawThreads *[]openai.Messages) (*[]openai.Messages, error) {
	timeLTEParsed := (*args)[7].GetParsed()
	timeGTParsed := (*args)[8].GetParsed()
	lengthLTEParsed := (*args)[9].GetParsed()
	lengthGTParsed := (*args)[10].GetParsed()

	if timeLTEParsed {
		filtered, err := filter.DaysLTE(rawThreads, *g.timeLTEArg)

		if err != nil {
			return nil, err
		}

		return filtered, nil

	} else if timeGTParsed {
		filtered, err := filter.DaysGT(rawThreads, *g.timeGTArg)

		if err != nil {
			return nil, err
		}

		return filtered, nil
	}

	// Filter length flow
	if lengthLTEParsed {
		filtered, err := filter.LengthLTE(rawThreads, *g.lengthLTEArg)

		if err != nil {
			return nil, err
		}

		return filtered, nil

	} else if lengthGTParsed {
		filtered, err := filter.LengthGT(rawThreads, *g.lengthGTArg)

		if err != nil {
			return nil, err
		}

		return filtered, nil
	}

	return rawThreads, nil
}

func (g *GetCommand) getThreadsOutput(args *[]argparse.Arg, threadIDs []string, filteredThreads *[]openai.Messages) (*[]byte, error) {
	prettyParsed := (*args)[6].GetParsed()

	if prettyParsed && *(g.prettyFlag) {
		parsedThreads := io.ParseThreads(threadIDs, filteredThreads)
		threadOutput, err := io.ListToJSON(parsedThreads)

		if err != nil {
			return nil, err
		}

		return &threadOutput, nil
	}

	threadOutput, err := io.ListToJSON(filteredThreads)

	if err != nil {
		return nil, err
	}

	return &threadOutput, nil
}

func (g *GetCommand) outputThreads(args *[]argparse.Arg, output *[]byte) error {
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
