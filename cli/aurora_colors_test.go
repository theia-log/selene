package cli

import (
	"fmt"
	"testing"

	"github.com/logrusorgru/aurora"
)

func TestAuroraColors(t *testing.T) {
	aur := aurora.NewAurora(true)

	fmt.Println(aur.Green("test"))
	fmt.Println(aur.Colorize("test", aurora.GreenBg))
}
