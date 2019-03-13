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

func TestAuroraColorsManager(t *testing.T) {
	colors := NewAuroraColors()

	fmt.Println(colors.ColoredContent(Context{"color": "green"}, "test green"))

	for i := 0; i < 20; i++ {
		fmt.Println(colors.ColoredTag(Context{}, fmt.Sprintf("test-%d", i)))
	}

	for i := 0; i < 16; i++ {
		fmt.Println(colors.ColoredTag(Context{}, fmt.Sprintf("test-%d", i)))
	}

	fmt.Println(colors.ColoredContent(Context{}, "Some ordinary content"))
	fmt.Println(colors.ColoredContent(Context{}, "Some error content"))
	fmt.Println(colors.ColoredContent(Context{}, "Some success content"))
	fmt.Println(colors.ColoredContent(Context{}, "Some info content"))
}
