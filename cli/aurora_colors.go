package cli

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
)

// TagColors generates the names of the colors used for tags.
// The names are not actual color names (like green, blue) but are labels
// used to lookup the colors in the current palette.
func TagColors() []string {
	colors := []string{}
	for i := 0; i < 16; i++ {
		colors = append(colors, fmt.Sprintf("tag-%d", i+1))
	}
	return colors
}

// RandTagColors generates the names of the colors used for tags, but in random
// order.
func RandTagColors() []string {
	shuffled := []string{}
	colors := TagColors()
	idxs := rand.New(rand.NewSource(time.Now().Unix())).Perm(len(colors))

	for _, i := range idxs {
		shuffled = append(shuffled, colors[i])
	}

	return shuffled
}

// TextTypeHeuristic is a function that checks if the given text matches some
// heuristic. Returns true if the texts matches the conditions in the underlying
// heuristic check.
type TextTypeHeuristic func(text string) bool

// TypeHeuristics contains map of heuristics by name.
type TypeHeuristics map[string]TextTypeHeuristic

// Detect tries to detect the category in which the text belongs to.
// It checks against all registered heuristics.
func (h TypeHeuristics) Detect(text string) string {
	for typeName, heuristic := range h {
		if heuristic(text) {
			return typeName
		}
	}
	return ""
}

// KnownTypes defines base text types like: error, warning, success etc.
// A given text will be checked if the text content is an error, warning
// or other registered type of text.
var KnownTypes = TypeHeuristics{
	"error": func(text string) bool {
		text = strings.TrimSpace(strings.ToLower(text))
		match, err := regexp.MatchString("\\[?\\berr?(or)?\\b\\]?", text)
		if err != nil {
			panic(err)
		}
		return match
	},
	"warning": func(text string) bool {
		text = strings.TrimSpace(strings.ToLower(text))
		match, err := regexp.MatchString("\\[?\\bwarn(ing)?\\b\\]?", text)
		if err != nil {
			panic(err)
		}
		return match
	},
	"success": func(text string) bool {
		text = strings.TrimSpace(strings.ToLower(text))
		match, err := regexp.MatchString("\\[?\\bsuccess\\b\\]?", text)
		if err != nil {
			panic(err)
		}
		return match
	},
	"info": func(text string) bool {
		text = strings.TrimSpace(strings.ToLower(text))
		match, err := regexp.MatchString("\\[?\\binfo\\b\\]?", text)
		if err != nil {
			panic(err)
		}
		return match
	},
}

// AuroraColors implements the Color interface using the aurora library.
type AuroraColors struct {
	aur             aurora.Aurora
	knownTags       map[string]string
	availableColors []string
	palette         map[string]aurora.Color
}

// ColoredText would color a text based on the "color" field in the context.
// If set, the color would be resolved from the underlying palette.
func (a *AuroraColors) ColoredText(ctx Context, text string) string {
	colorName := ctx.GetString("color")
	if colorName == "" {
		// no color set
		return text
	}
	color, ok := a.palette[colorName]
	if !ok {
		// no color with the provided name, so no colorizarion
		return text
	}
	return a.aur.Colorize(text, color).String()
}

// ColoredTag returns a colored text for the given tag.
// Same tags have the same colors, and different tags would get different
// colors. However, if the number of total tags is larger than the max number
// of colors reproducible on the terminal, some of the tags may receive same
// colors.
// The colors of the tags are random.
func (a *AuroraColors) ColoredTag(ctx Context, tag string) string {
	knownType := KnownTypes.Detect(tag)
	if knownType != "" {
		ctx["color"] = knownType
	} else {
		tag = strings.ToLower(tag)
		if knownColor, ok := a.knownTags[tag]; ok {
			ctx["color"] = knownColor
		} else {
			color := a.newTagColor()
			a.knownTags[tag] = color
			ctx["color"] = color
		}
	}
	return a.ColoredText(ctx, tag)
}

// ColoredContent colors a content text.
func (a *AuroraColors) ColoredContent(ctx Context, content string) string {
	if knownType := KnownTypes.Detect(content); knownType != "" {
		ctx["color"] = knownType
	}
	return a.ColoredText(ctx, content)
}

func (a *AuroraColors) newTagColor() string {
	if len(a.availableColors) == 0 {
		a.availableColors = RandTagColors()
	}
	color := a.availableColors[0]
	a.availableColors = a.availableColors[1:]
	return color
}

var defaultPalette map[string]aurora.Color = map[string]aurora.Color{
	// tags
	"tag-1":  aurora.BrownBg,
	"tag-2":  aurora.RedBg,
	"tag-3":  aurora.GreenBg,
	"tag-4":  aurora.GrayBg | aurora.CyanFg,
	"tag-5":  aurora.BlueBg,
	"tag-6":  aurora.MagentaBg,
	"tag-7":  aurora.CyanBg,
	"tag-8":  aurora.BrownBg | aurora.RedFg | aurora.BoldFm,
	"tag-9":  aurora.RedBg | aurora.BrownFg | aurora.BoldFm,
	"tag-10": aurora.GreenBg | aurora.GrayFg | aurora.BoldFm,
	"tag-11": aurora.GrayBg | aurora.BlackFg | aurora.BoldFm,
	"tag-12": aurora.BlueBg | aurora.BrownFg | aurora.BoldFm,
	"tag-13": aurora.MagentaBg | aurora.BlueFg | aurora.BoldFm,
	"tag-14": aurora.CyanBg | aurora.RedFg | aurora.BoldFm,
	"tag-15": aurora.GreenFg | aurora.BoldFm,
	"tag-16": aurora.RedFg | aurora.BoldFm,
	// base color labels
	"error":     aurora.RedFg | aurora.BoldFm,
	"warn":      aurora.BrownFg,
	"alert":     aurora.RedFg,
	"info":      aurora.CyanFg,
	"primary":   aurora.BlackBg | aurora.BlueFg,
	"secondary": aurora.BlackBg | aurora.GrayFg,
	"success":   aurora.GreenFg,
	// color names
	"black":   aurora.BlackFg,
	"red":     aurora.RedFg,
	"green":   aurora.GreenFg,
	"brown":   aurora.BrownFg,
	"blue":    aurora.BlueFg,
	"magenta": aurora.MagentaFg,
	"cyan":    aurora.CyanFg,
	"grey":    aurora.GrayFg,
}

// NewAuroraColors builds new Colors.
func NewAuroraColors() Colors {
	return &AuroraColors{
		aur:             aurora.NewAurora(true),
		availableColors: RandTagColors(),
		knownTags:       map[string]string{},
		palette:         defaultPalette,
	}
}
