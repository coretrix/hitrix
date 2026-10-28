package redissearch

import (
	"testing"

	"github.com/coretrix/trixorm"
)

type recordedDocument struct {
	key    string
	fields map[string]string
}

type recordingPusher struct {
	trixorm.RedisSearchIndexPusher
	key       string
	fields    map[string]string
	queued    *recordedDocument
	persisted *[]recordedDocument
}

func (pusher *recordingPusher) NewDocument(key string) {
	pusher.key = key
	pusher.fields = make(map[string]string)
}

func (pusher *recordingPusher) SetString(key, value string) {
	pusher.fields[key] = value
}

func (pusher *recordingPusher) PushDocument() {
	pusher.queued = &recordedDocument{key: pusher.key, fields: pusher.fields}
}

func (pusher *recordingPusher) Flush() {
	if pusher.queued != nil {
		*pusher.persisted = append(*pusher.persisted, *pusher.queued)
		pusher.queued = nil
	}
}

func TestEntityReindexPusherWritesOnlySelectedRecord(t *testing.T) {
	var persisted []recordedDocument
	pusher := &entityReindexPusher{
		targetKey: "vehicle:42",
		newPusher: func() trixorm.RedisSearchIndexPusher {
			return &recordingPusher{persisted: &persisted}
		},
	}

	// Exercise skipped documents both before and after the selected row.
	for _, key := range []string{"vehicle:41", "vehicle:42", "vehicle:43"} {
		pusher.NewDocument(key)
		pusher.SetString("Name", key)
		if key != "vehicle:42" {
			pusher.SetString("UnrelatedField", "must not leak")
		}
		pusher.PushDocument()
	}
	if len(persisted) != 0 {
		t.Fatal("documents were written before the operation was flushed")
	}
	pusher.Flush()
	if len(persisted) != 1 || persisted[0].key != "vehicle:42" {
		t.Fatalf("expected only vehicle:42, got %#v", persisted)
	}
	if len(persisted[0].fields) != 1 || persisted[0].fields["Name"] != "vehicle:42" {
		t.Fatalf("fields leaked from another document: %#v", persisted[0].fields)
	}
}

func TestEntityReindexPusherMissingRecordWritesNothing(t *testing.T) {
	var persisted []recordedDocument
	pusher := &entityReindexPusher{
		targetKey: "vehicle:42",
		newPusher: func() trixorm.RedisSearchIndexPusher {
			return &recordingPusher{persisted: &persisted}
		},
	}
	pusher.NewDocument("vehicle:43")
	pusher.SetString("Name", "following row")
	pusher.PushDocument()
	pusher.Flush()
	if pusher.selectedPusher != nil || len(persisted) != 0 {
		t.Fatal("a missing ID must never reindex the following row")
	}
}
