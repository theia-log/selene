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
	"warning": func(text string) bool {
		text = strings.TrimSpace(strings.ToLower(text))
		match, err := regexp.MatchString("\\[?warn(ing)?\\]?", text)
		if err != nil {
			panic(err)
		}
		return match
	},
	"success": func(text string) bool {
		text = strings.TrimSpace(strings.ToLower(text))
		match, err := regexp.MatchString("\\[?success\\]?", text)
		if err != nil {
			panic(err)
		}
		return match
	},
	"info": func(text string) bool {
		text = strings.TrimSpace(strings.ToLower(text))
		match, err := regexp.MatchString("\\[?info\\]?", text)
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

func NewAuroraColors() Colors {
	return &AuroraColors{
		aur:             aurora.NewAurora(true),
		availableColors: RandTagColors(),
		knownTags:       map[string]string{},
		palette:         defaultPalette,
	}
}
