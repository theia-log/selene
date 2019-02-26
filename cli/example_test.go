package cli_test

import (
	"flag"
	"fmt"

	"github.com/theia-log/selene/cli"
)

// ExampleParrot implements a hypothetical parrot that repeats a phrase given
// as an argument.
// Invoking with help would give the following output:
// 		Usage: parrot say [ARGS]
//
// 		Hypothetical parrot.
//
// 		Commands:
// 			say - Say a phrase.
// And then the help for the subcommand can be invoked like this:
// 	Usage of say:
// 	-phrase string
// 		  Phrase to say.
// 	-times int
// 		  How many times to repeat the phrase. (default 1)
func ExampleCommandsCLI() {
	parrot := cli.NewCommandCLI("parrot", "Hypothetical parrot.")

	parrot.AddCommand("say", func(args []string) error {
		fs := flag.NewFlagSet("say", flag.ExitOnError)
		phrase := fs.String("phrase", "", "Phrase to say.")
		times := fs.Int("times", 1, "How many times to repeat the phrase.")

		fs.Parse(args)

		for i := 0; i < *times; i++ {
			fmt.Println(*phrase)
		}

		return nil
	}, "Say a phrase.")

	parrot.Execute()
	// Our hypothetical parrot can then be invoked like this:
	//		parrot say -phrase Hello -times 3
	// and would print:
	//		Hello
	// 		Hello
	// 		Hello
	//
}
