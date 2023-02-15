package utility

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

func TestStringToTimestampWithValidFormat(t *testing.T) {
	datetime := "2023-01-03T15:45:01+07:00"
	wantTimestamp := timestamppb.Timestamp{Seconds: 0, Nanos: 0}
	gotTimestamp := StringToTimestamp(datetime)

	if wantTimestamp.Seconds == gotTimestamp.Seconds && wantTimestamp.Nanos == gotTimestamp.Nanos {
		t.Error("parsed valid timestamp didn't match the expected timestamp")
		return
	}
}

func TestStringToTimestampWithInvalidFormat(t *testing.T) {
	datetime := "broken payload"
	gotTimestamp := StringToTimestamp(datetime)

	if gotTimestamp != nil {
		t.Error("parsed invalid timestamp wasn't nil as it should have been")
		return
	}
}

func TestStringToTimestampWithInvalidTime(t *testing.T) {
	datetime := "2023-13-44T15:45:01+07:00"
	gotTimestamp := StringToTimestamp(datetime)

	if gotTimestamp != nil {
		t.Error("parsed timestamp with invalid time wasn't nil as it should have been")
		return
	}
}
