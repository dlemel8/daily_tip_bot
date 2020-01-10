package main

var topicSources = map[string][]tipSource{
	"bash": {&inMemorySource{
		tips: []string{
			"bash placeholder 1",
			"bash placeholder 2",
			"bash placeholder 3",
			"bash placeholder 4",
			"bash placeholder 5",
		},
	}},
	"vim": {&inMemorySource{
		tips: []string{
			"vim placeholder 1",
			"vim placeholder 2",
			"vim placeholder 3",
			"vim placeholder 4",
			"vim placeholder 5",
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
