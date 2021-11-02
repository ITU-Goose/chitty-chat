package message

import (
	"github.com/goose-alt/chitty-chat/internal/logging"
	"github.com/goose-alt/chitty-chat/internal/time"
)

func PrintMessage(prefix string, logger logging.Log, content string, timestamp *time.VectorTimestamp, id string) {
	logger.IPrintf(
		"%s: \"%s\", at: %s, from %s\n",
		prefix,
		content,
		timestamp.GetDisplayableContent(),
		id,
	)
}
