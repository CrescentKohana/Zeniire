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
		t.Fatal("parsed valid timestamp didn't match the expected timestamp")
	}
}

func TestStringToTimestampWithInvalidFormat(t *testing.T) {
	datetime := "broken payload"
	gotTimestamp := StringToTimestamp(datetime)

	if gotTimestamp != nil {
		t.Fatal("parsed invalid timestamp wasn't nil as it should have been")
	}
}

func TestStringToTimestampWithInvalidTime(t *testing.T) {
	datetime := "2023-13-44T15:45:01+07:00"
	gotTimestamp := StringToTimestamp(datetime)

	if gotTimestamp != nil {
		t.Fatal("parsed timestamp with invalid time wasn't nil as it should have been")
	}
}

func TestClamp(t *testing.T) {
	wantNumber := int64(500)
	gotNumber := Clamp(1000, 0, 500)

	if wantNumber != gotNumber {
		t.Fatal("the output of clamp function did not math the expected")
	}
}
