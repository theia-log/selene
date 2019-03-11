package cli

// Context in which the text coloration takes place.
type Context map[string]interface{}

// GetString looks up a string stored under a key in the Context.
func (c Context) GetString(key string) string {
	val, ok := c[key]
	if !ok {
		return ""
	}
	strVal, ok := val.(string)
	if !ok {
		return ""
	}
	return strVal
}

// Colors defines methods for coloring text in the command line interface.
type Colors interface {
	// ColoredText colors the given text based on the provided Context.
	// The text may or may not end up being colored - this is controlled by the
	// values given in the context. How the values in the Context are
	// interpreted, depends completely on the underlying implementation.
	// For example, the function may be called with Context: {'type': 'info'}
	// which may be interpreted that the text is an info message and will be
	// colored in the 'info' color from an exiting palette.
	// If the current terminal is not interactive (tty), then no coloring takes
	// place.
	ColoredText(ctx Context, text string) string

	// ColoredTag colors the given tag name. This is a handy wrapper around
	// ColoredText, which tweaks the context to provide logic for coloring an
	// Event Tag.
	// The tags may be colored all in the same color, or may be colored
	// differently, each tag in different color. The bookkeeping of the tags
	// colors is left to the underlying implementation.
	// If the current terminal is not interactive (tty), then no coloring takes
	// place.
	ColoredTag(ctx Context, tag string) string

	// ColoredContent colors the event content. This is a handy wrapper for
	// ColredText, which adds heuristics for determining the color of the
	// event content - for example, the context may provide the list of
	// tags and if there exists a tag 'error' in that list, then the content
	// may be colored with the 'error' color of the palette. Similarly
	// a heuristic to look up the words 'error'  or 'info' may be applied to
	// determine the type of the content, and the content may be colored
	// accordingly.
	// If the current terminal is not interactive (tty), then no coloring takes
	// place.
	ColoredContent(ctx Context, content string) string
}
