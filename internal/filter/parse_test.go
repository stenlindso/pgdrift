package filter_test

import (
	"testing"

	"github.com/yourorg/pgdrift/internal/filter"
)

func TestParseOptions_Empty(t *testing.T) {
	opts := filter.ParseOptions("", "", "", "")
	if len(opts.IncludeSchemas) != 0 {
		t.Errorf("expected empty IncludeSchemas, got %v", opts.IncludeSchemas)
	}
	if len(opts.ExcludeTables) != 0 {
		t.Errorf("expected empty ExcludeTables, got %v", opts.ExcludeTables)
	}
}

func TestParseOptions_CSV(t *testing.T) {
	opts := filter.ParseOptions("public,app", "", "users, orders", "")
	if len(opts.IncludeSchemas) != 2 {
		t.Fatalf("expected 2 IncludeSchemas, got %d", len(opts.IncludeSchemas))
	}
	if opts.IncludeSchemas[0] != "public" || opts.IncludeSchemas[1] != "app" {
		t.Errorf("unexpected IncludeSchemas: %v", opts.IncludeSchemas)
	}
	if len(opts.IncludeTables) != 2 {
		t.Fatalf("expected 2 IncludeTables, got %d", len(opts.IncludeTables))
	}
	if opts.IncludeTables[1] != "orders" {
		t.Errorf("expected orders, got %s", opts.IncludeTables[1])
	}
}

func TestParseOptions_Whitespace(t *testing.T) {
	opts := filter.ParseOptions(" public , app ", "", "", "")
	if len(opts.IncludeSchemas) != 2 {
		t.Fatalf("expected 2 schemas, got %d", len(opts.IncludeSchemas))
	}
	if opts.IncludeSchemas[0] != "public" {
		t.Errorf("expected trimmed value 'public', got '%s'", opts.IncludeSchemas[0])
	}
}

func TestParseOptions_TrailingComma(t *testing.T) {
	opts := filter.ParseOptions("public,", "", "", "")
	if len(opts.IncludeSchemas) != 1 {
		t.Errorf("expected 1 schema, got %d: %v", len(opts.IncludeSchemas), opts.IncludeSchemas)
	}
}
