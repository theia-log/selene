package cli

func EventCommand(args []string) error {
	eventFlags, flagSet := SetupEventGeneratorFlags()
	if err := flagSet.Parse(args); err != nil {
		return err
	}
	return RunEventGenerator(eventFlags)
}

func RunEventGenerator(flags *EventFlags) error {
	return nil
}
