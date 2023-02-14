package utility

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func StringToTimestamp(timestamp string) *timestamppb.Timestamp {
	parsed, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return nil
	}
	return timestamppb.New(parsed)
}
