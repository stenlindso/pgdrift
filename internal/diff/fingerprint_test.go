package diff

import (
	"strings"
	"testing"

	"github.com/pgdrift/pgdrift/internal/schema"
)

func buildFingerprintSchema() *schema.Schema {
	s := schema.NewSchema()
	t := s.AddTable("public", "users")
	t.AddColumn("id", "bigint", false, "")
	t.AddColumn("email", "text", false, "")
	t2 := s.AddTable("public", "orders")
	t2.AddColumn("id", "bigint", false, "")
	t2.AddColumn("user_id", "bigint", false, "")
	return s
}

func TestSchemaFingerprint_Nil(t *testing.T) {
	f := SchemaFingerprint(nil)
	if f.Hash != "" {
		t.Errorf("expected empty hash for nil schema, got %q", f.Hash)
	}
	if f.Tables != 0 || f.Columns != 0 {
		t.Errorf("expected zero counts for nil schema")
	}
}

func TestSchemaFingerprint_Stable(t *testing.T) {
	s := buildFingerprintSchema()
	f1 := SchemaFingerprint(s)
	f2 := SchemaFingerprint(s)
	if f1.Hash != f2.Hash {
		t.Errorf("fingerprint not stable: %q vs %q", f1.Hash, f2.Hash)
	}
}

func TestSchemaFingerprint_Counts(t *testing.T) {
	s := buildFingerprintSchema()
	f := SchemaFingerprint(s)
	if f.Tables != 2 {
		t.Errorf("expected 2 tables, got %d", f.Tables)
	}
	if f.Columns != 4 {
		t.Errorf("expected 4 columns, got %d", f.Columns)
	}
}

func TestSchemaFingerprint_DifferentSchemas(t *testing.T) {
	s1 := buildFingerprintSchema()
	s2 := buildFingerprintSchema()
	s2.AddTable("public", "products").AddColumn("id", "bigint", false, "")

	f1 := SchemaFingerprint(s1)
	f2 := SchemaFingerprint(s2)
	if f1.Equal(f2) {
		t.Error("expected different fingerprints for different schemas")
	}
}

func TestFingerprint_String_Short(t *testing.T) {
	f := Fingerprint{Hash: "abcdef1234567890"}
	if got := f.String(); got != "abcdef12" {
		t.Errorf("expected 8-char prefix, got %q", got)
	}
}

func TestCompareFingerprints_NoDrift(t *testing.T) {
	s := buildFingerprintSchema()
	fd := CompareFingerprints(s, s)
	if fd.Changed {
		t.Error("expected no change when comparing schema to itself")
	}
}

func TestCompareFingerprints_WithDrift(t *testing.T) {
	s1 := buildFingerprintSchema()
	s2 := buildFingerprintSchema()
	s2.AddTable("public", "events").AddColumn("id", "bigint", false, "")

	fd := CompareFingerprints(s1, s2)
	if !fd.Changed {
		t.Error("expected change detected")
	}
}

func TestFingerprintSummary_NoDrift(t *testing.T) {
	s := buildFingerprintSchema()
	fd := CompareFingerprints(s, s)
	summary := FingerprintSummary(fd)
	if !strings.Contains(summary, "schemas match") {
		t.Errorf("unexpected summary: %q", summary)
	}
}

func TestFingerprintSummary_WithDrift(t *testing.T) {
	s1 := buildFingerprintSchema()
	s2 := buildFingerprintSchema()
	s2.AddTable("public", "events").AddColumn("id", "bigint", false, "")

	fd := CompareFingerprints(s1, s2)
	summary := FingerprintSummary(fd)
	if !strings.Contains(summary, "schema drift detected") {
		t.Errorf("unexpected summary: %q", summary)
	}
	if !strings.Contains(summary, "+1 tables") {
		t.Errorf("expected table delta in summary: %q", summary)
	}
}
