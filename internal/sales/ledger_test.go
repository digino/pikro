package sales

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func at(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestLoadMergesChunksAndIgnoresLegacy(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	// The pre-v2 .json ledger (no voucher names) must be ignored: reports start from v2.
	legacy := `[{"at":"2026-01-05T10:00:00Z","profile":"1h","price":"100","currency":"XOF","count":1},
	{"at":"2026-01-06T10:00:00Z","profile":"1d","price":"500","currency":"XOF","count":3}]`
	if err := os.WriteFile(filepath.Join(home, ".pikro-sales-r1-2026.json"), []byte(legacy), 0600); err != nil {
		t.Fatal(err)
	}
	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	must(AppendBatch("r1", "b1", at("2026-02-01T10:05:00Z"), "24hours", "500", "XOF", []string{"a", "b"}))
	must(AppendBatch("r1", "b1", at("2026-02-01T10:00:00Z"), "24hours", "500", "XOF", []string{"c"}))
	must(AppendBatch("r1", "b2", at("2026-02-02T10:00:00Z"), "1h", "100", "XOF", []string{"d"}))

	got, err := Load("r1", 2026)
	if err != nil {
		t.Fatal(err)
	}
	type row struct {
		ID    string
		At    string
		Count int
	}
	var rows []row
	for _, b := range got.Batches {
		rows = append(rows, row{b.ID, b.At.Format(time.RFC3339), b.Count})
	}
	want := []row{
		{"b1", "2026-02-01T10:00:00Z", 3},
		{"b2", "2026-02-02T10:00:00Z", 1},
	}
	if !reflect.DeepEqual(rows, want) {
		t.Errorf("batches = %+v, want %+v", rows, want)
	}
	if got.Activations == nil || len(got.Activations) != 0 {
		t.Errorf("activations = %#v, want empty non-nil", got.Activations)
	}
	if _, err := os.Stat(filepath.Join(home, ".pikro-sales-r1-2026.jsonl")); err != nil {
		t.Errorf("jsonl not written: %v", err)
	}
}

func TestLoadMissingFilesIsEmpty(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	got, err := Load("none", 2026)
	if err != nil || len(got.Batches) != 0 || len(got.Activations) != 0 {
		t.Errorf("got %+v, %v; want empty", got, err)
	}
}
