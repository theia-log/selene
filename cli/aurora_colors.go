package cli

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
)

func TagColors() []string {
	colors := []string{}
	for i := 0; i < 16; i++ {
		colors = append(colors, fmt.Sprintf("tag-%d", i+1))
	}
	return colors
}

func RandTagColors() []string {
	shuffled := []string{}
	colors := TagColors()
	idxs := rand.New(rand.NewSource(time.Now().Unix())).Perm(len(colors))

	for _, i := range idxs {
		shuffled = append(shuffled, colors[i])
	}

	return shuffled
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
	availableColors []string
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

func (a *AuroraColors) newTagColor() string {
	if len(a.availableColors) == 0 {
		a.availableColors = RandTagColors()
	}
	color := a.availableColors[0]
	a.availableColors = a.availableColors[1:]
	return color
}
