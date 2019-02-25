package cli_test

import (
	"flag"
	"fmt"

	"github.com/theia-log/selene/cli"
)

// ExampleParrot implements a hypothetical parrot that repeats a phrase given
// as an argument.
func ExampleCommandsCLI() {
	parrot := cli.NewCommandCLI("parrot", "Hypothetical parrot.")

	parrot.AddCommand("say", func(args []string) error {
		phrase := flag.String("phrase", "", "Phrase to say.")
		times := flag.Int("times", 1, "How many times to repeat the phrase.")

		flag.Parse()

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
