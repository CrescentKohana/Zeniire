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

// Clamp clamps the given value to the given range.
func Clamp(value, min, max int64) int64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
