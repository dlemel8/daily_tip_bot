package main

// TODO - load from config file
var topicSources = map[string][]tipSource{
	"bash": {&inMemorySource{
		tips: []string{
			"Use [Ctrl][A] to move to the beginning of the current line.",
			"Use [Ctrl[E] to move to the end of the current line. ",
			"Use [Alt][F] to move the cursor forward one word on the current line.",
			"Use [Alt][B] to move the cursor backwards one word on the current line.",
			"Use [Ctrl][U] to clear the characters on the line before the current cursor position",
			"Use [Ctrl][K] to clear the characters on the line after the current cursor position.",
			"Use [Ctrl][W] to delete the word in front of the cursor.",
			"Use [Alt][D] to delete the word after the cursor.",
			"Use [Alt][U] to make the current word after the cursor uppercase.",
			"Use [Alt][L] to make the current word after the cursor lowercase.",
			"Use [Alt][C] to capitalize the current word after the cursor.",
		},
	}},
	"vim": {&inMemorySource{
		tips: []string{
			"h j k l: Basic movement keys. Most useful when prefixed with a number.",
			"b w B W: Move back by token/forward by token/back by word/forward by word.",
			"0 ^ $: Jump to first column/first non-whitespace character/end of line.",
			"ctrl+u ctrl+d: Moves by half a screenful and doesn’t lose your cursor position.",
			"<line-number>G: Jump directly to a specific line number.",
			"H M L: Move to the top/middle/bottom of the screen (i.e. High/Middle/Low).",
			"# *: Find the previous/next occurrence of the token under the cursor.",
			"n N: Repeat the last find command forward/backward.",
			"“ (two back-ticks): Jump back to where you just were.",
			"ctrl+o ctrl+i: Move backward/forward through the jump history.",
			"i a I A: Enter insert mode at cursor/append after cursor/insert at beginning of line/append to end of line.",
			"o O: Open new line (below the current line/above the current line).",
			"cw cW: Correct (delete and insert) the token(s)/word(s) following the cursor.",
			"cc: Correct line(s) by clearing and then entering insert mode.",
			"dd: Delete line(s).",
			"yy: Copy line(s).  The “y” is for “yank.”",
			"yw yW: Copy token(s)/word(s).",
			"p P: Paste the last thing that was deleted or copied before/after cursor.",
			"u ctrl+r: Undo and redo,",
		},
	}},
}

func listTopics() []string {
	keys := make([]string, 0, len(topicSources))
	for k := range topicSources {
		keys = append(keys, k)
	}
	return keys
}

func getTip(topic string) string {
	if sources, ok := topicSources[topic]; !ok {
		return ""
	} else {
		for _, source := range sources {
			tip := source.randomTip()
			if tip != "" {
				return tip
			}
		}
	}
	return ""
}
