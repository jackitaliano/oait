package threads

import (
	"errors"
	"fmt"
	"strings"

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

	threadsArg            *[]string
	inputArg              *string
	sessionArg            *string
	orgArg                *string
	outputArg             *string
	prettyFlag            *bool
	timeLTEArg            *float64
	timeGTArg             *float64
	lengthLTEArg          *float64
	lengthGTArg           *float64
	contentContainsArg    *[]string
	contentNotContainsArg *[]string
	metadataArg           *[]string
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
	prettyFlag := subCommand.Flag("p", "pretty", &argparse.Options{Required: false, Help: "Pretty print threads"})
	timeLTEArg := subCommand.Float("d", "days", &argparse.Options{Required: false, Help: "Filter by LTE to days"})
	timeGTArg := subCommand.Float("D", "Days", &argparse.Options{Required: false, Help: "Filter by GT days"})
	lengthLTEArg := subCommand.Float("l", "length", &argparse.Options{Required: false, Help: "Filter by LTE to length"})
	lengthGTArg := subCommand.Float("L", "Length", &argparse.Options{Required: false, Help: "Filter by GT length"})
	contentContainsArg := subCommand.StringList("c", "content", &argparse.Options{Required: false, Help: "Filter by thread content contains"})
	contentNotContainsArg := subCommand.StringList("C", "Content", &argparse.Options{Required: false, Help: "Filter by thread content not contains"})
	metadataArg := subCommand.StringList("m", "meta", &argparse.Options{Required: false, Help: "Filter by thread metadata"})

	return &DelCommand{
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
		contentContainsArg,
		contentNotContainsArg,
		metadataArg,
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

	fmt.Printf("Filtering thread ids...\t\t")
	filteredThreadIDs, err := d.filterThreadsIds(&args, key, threadIDs)
	if err != nil {
		fmt.Printf("X\n")
		return err
	}
	fmt.Printf("✓\n")

	fmt.Printf("Retrieving threads...\t\t")
	rawThreads := openai.RetrieveThreadsMessages(key, filteredThreadIDs, *d.orgArg)
	fmt.Printf("✓\n")

	fmt.Printf("Filtering threads...\t\t")
	filteredThreads, err := d.filterThreads(&args, rawThreads)

	if err != nil {
		fmt.Printf("X\n")
		return err
	}
	fmt.Printf("✓\n")

	verify := verifyBeforeDelete()

	deleteThreadIDs := getThreadIDsFromObjects(filteredThreads)

	if verify {
		fmt.Printf("Formatting thread output...\t")
		threadsOutput, err := d.getThreadsOutput(&args, filteredThreadIDs, filteredThreads)

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
		openai.DeleteThreads(key, deleteThreadIDs, *d.orgArg)
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
		threadIDs, err := io.ListInput(*d.threadsArg)

		if err != nil {
			return nil, err
		}

		return threadIDs, nil

	}

	if inputParsed { // File input passed
		threadIDs, err := io.FileInput(*d.inputArg)

		if err != nil {
			return nil, err
		}

		return threadIDs, nil
	}

	if sessionParsed {
		threadIDs, err := io.SessionInput(*d.sessionArg, *d.orgArg)

		if err != nil {
			return nil, err
		}

		return threadIDs, nil

	}

	errMsg := fmt.Sprintf("No input options passed to `%v`\n", d.name)
	err := errors.New(errMsg)

	return nil, err
}

func (d *DelCommand) filterThreadsIds(args *[]argparse.Arg, key string, threadIds []string) ([]string, error) {
	metadataParsed := (*args)[13].GetParsed()

	if !metadataParsed {
		return threadIds, nil
	}

	filtered := threadIds
	var err error

	threads := openai.RetrieveThreads(key, threadIds, *d.orgArg)

	if metadataParsed {
		metadata := make(map[string]string, len(*d.metadataArg))

		for _, metadataStr := range *d.metadataArg {
			metadataSplit := strings.Split(metadataStr, "=")

			if len(metadataSplit) < 2 {
				errMsg := fmt.Sprintf("invalid metadata: '%s'. (should be '<key>=<value>')", metadataStr)
				err = errors.New(errMsg)
				return nil, err
			}

			metadataKey := metadataSplit[0]
			metadataVal := metadataSplit[1]

			metadata[metadataKey] = metadataVal
		}

		filteredThreads := filter.MetadataEquals(threads, metadata)

		filtered = make([]string, len(*filteredThreads))

		for i, thread := range *filteredThreads {
			id := thread.ID

			filtered[i] = id
		}
	}

	return filtered, nil
}

func (d *DelCommand) filterThreads(args *[]argparse.Arg, rawThreads *[]openai.Messages) (*[]openai.Messages, error) {
	timeLTEParsed := (*args)[7].GetParsed()
	timeGTParsed := (*args)[8].GetParsed()
	lengthLTEParsed := (*args)[9].GetParsed()
	lengthGTParsed := (*args)[10].GetParsed()
	contentContainsParsed := (*args)[11].GetParsed()
	contentNotContainsParsed := (*args)[12].GetParsed()

	filtered := rawThreads
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

	// Filter length flow
	if lengthLTEParsed {
		filtered, err = filter.LengthLTE(filtered, *d.lengthLTEArg)

		if err != nil {
			return nil, err
		}
	}

	if lengthGTParsed {
		filtered, err = filter.LengthGT(filtered, *d.lengthGTArg)

		if err != nil {
			return nil, err
		}
	}

	if contentContainsParsed {
		filtered = filter.ContainsContent(filtered, *d.contentContainsArg)
	}

	if contentNotContainsParsed {
		filtered = filter.NotContainsContent(filtered, *d.contentNotContainsArg)
	}

	return filtered, nil
}

func (d *DelCommand) getThreadsOutput(args *[]argparse.Arg, threadIDs []string, filteredThreads *[]openai.Messages) (*[]byte, error) {
	prettyParsed := (*args)[6].GetParsed()

	if prettyParsed && *(d.prettyFlag) {
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

func (d *DelCommand) outputThreads(args *[]argparse.Arg, output *[]byte) error {
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

func getThreadIDsFromObjects(threads *[]openai.Messages) []string {
	threadIDs := []string{}

	for _, thread := range *threads {
		if thread.GetLen() > 0 {
			threadIDs = append(threadIDs, thread.Messages[0].ThreadID)
		}
	}

	return threadIDs
}
