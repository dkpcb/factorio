package main

import (
	"bytes"
	"strings"
	"testing"
)

const example = "{\"schema_version\":1,\"world_id\":\"unverified\",\"run_id\":\"unverified\",\"sequence\":1,\"tick\":123,\"event_type\":\"research_finished\",\"payload\":{\"technology\":\"electronics\"}}"

func TestConsumeCompleteAndPartial(t *testing.T) {
	data := example + "\n" + example[:len(example)-3]
	reader := strings.NewReader(data)
	var output bytes.Buffer
	var offset int64
	n, err := consume(reader, &offset, &output)
	if err != nil || n != 1 {
		t.Fatalf("first read: got n=%d err=%v", n, err)
	}
	if output.String() != example+"\n" {
		t.Fatalf("unexpected output: %s", output.String())
	}
	firstOffset := offset
	n, err = consume(reader, &offset, &output)
	if err != nil || n != 0 || offset != firstOffset {
		t.Fatalf("partial line should be retried: n=%d err=%v offset=%d", n, err, offset)
	}
}

func TestConsumeBadRecordDoesNotAdvance(t *testing.T) {
	var offset int64
	var out bytes.Buffer
	_, err := consume(strings.NewReader("not-json\n"), &offset, &out)
	if err == nil || offset != 0 || out.Len() != 0 {
		t.Fatalf("expected rejection without checkpoint advance: %v %d", err, offset)
	}
}

func TestEnvelope(t *testing.T) {
	if err := validate([]byte(example)); err != nil {
		t.Fatal(err)
	}
	if err := validate([]byte("{}")); err == nil {
		t.Fatal("expected invalid envelope")
	}
}
