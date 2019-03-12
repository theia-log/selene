package cli

import (
	"regexp"
	"strings"

	"github.com/logrusorgru/aurora"
)

var TagsColors []aurora.Color = []aurora.Color{
	aurora.RedBg,
	aurora.GreenBg,
	aurora.BrownBg,
	aurora.BlueBg,
	aurora.MagentaBg,
	aurora.CyanBg,
	aurora.GrayBg,
}

type TextTypeHeuristic func(text string) bool

type TypeHeuristics map[string]TextTypeHeuristic

func (h TypeHeuristics) Detect(text string) string {
	for typeName, heuristic := range h {
		if heuristic(text) {
			return typeName
		}
	}
	return ""
}

var KnownTypes TypeHeuristics = TypeHeuristics{
	"error": func(text string) bool {
		text = strings.TrimSpace(strings.ToLower(text))
		match, err := regexp.MatchString("\\[?err?(or)?\\]?", text)
		if err != nil {
			panic(err)
		}
		return match
	},
}

type AuroraColors struct {
	aur             aurora.Aurora
	knownTags       map[string]string
	availableColors []aurora.Color
	palette         map[string]aurora.Color
}

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

func (a *AuroraColors) ColoredTag(ctx Context, tag string) string {
	knownType := KnownTypes.Detect(tag)
	if knownType != "" {
		ctx["color"] = knownType
	}
	return ""
}
