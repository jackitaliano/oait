package threadsParse

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackitaliano/oait/internal/openai"
	"github.com/jackitaliano/oait/internal/threads"
	"github.com/jackitaliano/oait/internal/tui"

	"github.com/akamensky/argparse"
)

type DelCommand struct {
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

func NewDelCommand(command *argparse.Command) *DelCommand {
	const name = "del"
	const desc = "Del Thread Tools"

	subCommand := command.NewCommand(name, desc)

	threadsArg := subCommand.StringList("i", "ids", &argparse.Options{Required: false, Help: "List of Thread IDs"})
	inputArg := subCommand.String("f", "file-input", &argparse.Options{Required: false, Help: "Thread File Input"})
	sessionArg := subCommand.String("s", "session", &argparse.Options{Required: false, Help: "Retrieve Threads from session-id"})
	orgArg := subCommand.String("O", "org", &argparse.Options{Required: false, Help: "Set Organization ID"})
	outputArg := subCommand.String("o", "output", &argparse.Options{Required: false, Help: "Thread File Output"})
	rawFlag := subCommand.Flag("r", "raw", &argparse.Options{Required: false, Help: "Output raw Threads"})
	timeLTEArg := subCommand.Float("d", "days", &argparse.Options{Required: false, Help: "Filter by LTE to days"})
	timeGTArg := subCommand.Float("D", "Days", &argparse.Options{Required: false, Help: "Filter by GT days"})
	lengthLTEArg := subCommand.Float("l", "length", &argparse.Options{Required: false, Help: "Filter by LTE to length"})
	lengthGTArg := subCommand.Float("L", "Length", &argparse.Options{Required: false, Help: "Filter by GT length"})

	return &DelCommand{
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

func (d *DelCommand) Happened() bool {
	return d.command.Happened()
}

func (d *DelCommand) Run(key string) error {
	args := d.command.GetArgs()

	fmt.Printf("Retrieving thread ids...\t")
	threadIDs, err := d.getThreadIDs(&args)

	if err != nil {
		fmt.Printf("X\n")
		return err
	}

	fmt.Printf("✓\n")

	fmt.Printf("Retrieving threads...\t\t")
	rawThreads := threads.RetrieveThreads(key, threadIDs, *d.orgArg)
	fmt.Printf("✓\n")

	fmt.Printf("Filtering threads...\t\t")
	filteredThreads := d.filterThreads(&args, rawThreads)
	fmt.Printf("✓\n")

	verify := verifyBeforeDelete()

	deleteThreadIDs := getThreadIDsFromObjects(filteredThreads)

	if verify {
		fmt.Printf("Formatting thread output...\t")
		threadsOutput, err := d.getThreadsOutput(&args, threadIDs, filteredThreads)

		if err != nil {
			fmt.Printf("X\n")
			return err
		}
		fmt.Printf("✓\n")

		fmt.Printf("Outputting threads... \n\n")
		err = d.outputThreads(&args, threadsOutput)

		if err != nil {
			fmt.Printf("X\n")
			return err
		}
	}

	confirmed := confirmDelete()

	if confirmed {
		fmt.Printf("Deleting threads...\t\t")
		threads.DeleteThreads(key, deleteThreadIDs, *d.orgArg)
		fmt.Printf("✓\n")
	} else {
		fmt.Printf("Cancelled.\n")
	}

	return nil
}

func verifyBeforeDelete() bool {
	return tui.YesNoLoop("Verify threads before deletion?")
}

func confirmDelete() bool {
	return tui.YesNoLoop("Confirm deletion")
}

func (d *DelCommand) getThreadIDs(args *[]argparse.Arg) ([]string, error) {
	threadsParsed := (*args)[1].GetParsed()
	inputParsed := (*args)[2].GetParsed()
	sessionParsed := (*args)[3].GetParsed()

	if threadsParsed { // List passed
		threadIDs, err := threads.ListInput(*d.threadsArg)

		if err != nil {
			return nil, err
		}

		return threadIDs, nil

	}

	if inputParsed { // File input passed
		threadIDs, err := threads.FileInput(*d.inputArg)

		if err != nil {
			return nil, err
		}

		return threadIDs, nil
	}

	if sessionParsed {
		threadIDs, err := threads.SessionInput(*d.sessionArg, *d.orgArg)

		if err != nil {
			return nil, err
		}

		return threadIDs, nil

	}

	errMsg := fmt.Sprintf("No input options passed to `%v`\n", d.name)
	err := errors.New(errMsg)

	return nil, err
}

func (d *DelCommand) filterThreads(args *[]argparse.Arg, rawThreads *[][]openai.Message) *[][]openai.Message {
	timeLTEParsed := (*args)[7].GetParsed()
	timeGTParsed := (*args)[8].GetParsed()
	lengthLTEParsed := (*args)[9].GetParsed()
	lengthGTParsed := (*args)[10].GetParsed()

	if timeLTEParsed {
		return threads.FilterByDaysLTE(rawThreads, *d.timeLTEArg)

	} else if timeGTParsed {
		return threads.FilterByDaysGT(rawThreads, *d.timeGTArg)
	}

	// Filter length flow
	if lengthLTEParsed {
		return threads.FilterByLengthLTE(rawThreads, *d.lengthLTEArg)

	} else if lengthGTParsed {
		return threads.FilterByLengthGT(rawThreads, *d.lengthGTArg)
	}

	return rawThreads
}

func (d *DelCommand) getThreadsOutput(args *[]argparse.Arg, threadIDs []string, filteredThreads *[][]openai.Message) (*[]byte, error) {
	rawParsed := (*args)[6].GetParsed()

	if rawParsed && *(d.rawFlag) {
		threadOutput, err := json.MarshalIndent(*filteredThreads, "", "\t")

		if err != nil {
			errMsg := fmt.Sprintf("Error marshalling json: %v\n", err)
			err := errors.New(errMsg)

			return nil, err
		}

		return &threadOutput, nil

	}

	parsedThreads := threads.ParseThreads(threadIDs, filteredThreads)
	threadOutput, err := threads.ThreadsToJson(parsedThreads)

	if err != nil {
		return nil, err
	}

	return &threadOutput, nil
}

func (d *DelCommand) outputThreads(args *[]argparse.Arg, output *[]byte) error {
	outputParsed := (*args)[5].GetParsed()

	if outputParsed {
		err := threads.FileOutput(*d.outputArg, output)

		if err != nil {
			return err
		}

	} else {
		fmt.Printf("%v\n", string(*output))
	}

	return nil
}

func getThreadIDsFromObjects(threads *[][]openai.Message) []string {
	threadIDs := []string{}

	for _, fileObject := range *threads {
		if len(fileObject) > 0 {
			threadIDs = append(threadIDs, fileObject[0].ID)
		}
	}

	return threadIDs
}
