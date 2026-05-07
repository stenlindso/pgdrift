package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pgdrift/pgdrift/internal/schema"
	"github.com/pgdrift/pgdrift/internal/snapshot"
)

func buildSchema() *schema.Schema {
	s := schema.NewSchema()
	s.AddTable("public", "users", []schema.Column{
		{Name: "id", Type: "integer", Nullable: false},
		{Name: "email", Type: "text", Nullable: false},
	})
	return s
}

func TestNew(t *testing.T) {
	s := buildSchema()
	before := time.Now().UTC()
	snap := snapshot.New(s, "postgres://localhost/mydb")
	after := time.Now().UTC()

	if snap.Schema != s {
		t.Error("expected schema to be set")
	}
	if snap.DSN != "postgres://localhost/mydb" {
		t.Errorf("unexpected DSN: %s", snap.DSN)
	}
	if snap.CapturedAt.Before(before) || snap.CapturedAt.After(after) {
		t.Errorf("CapturedAt out of range: %v", snap.CapturedAt)
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	orig := snapshot.New(buildSchema(), "postgres://localhost/test")
	if err := snapshot.Save(orig, path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if loaded.DSN != orig.DSN {
		t.Errorf("DSN mismatch: got %q want %q", loaded.DSN, orig.DSN)
	}
	tables := loaded.Schema.TableNames()
	if len(tables) != 1 || tables[0] != "public.users" {
		t.Errorf("unexpected tables: %v", tables)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("not json{"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := snapshot.Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
