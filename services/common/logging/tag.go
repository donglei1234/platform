package logging

import (
	"go.uber.org/zap"
)

const (
	tagsKey = "tags"
)

// This utility function does not guard against duplicate tags
func AddTagsToLogMessage(tags ...string) zap.Field {
	tagsSlice := make([]string, len(tags))

	for i, v := range tags {
		tagsSlice[i] = v
	}

	return zap.Strings(tagsKey, tagsSlice)
}
