package filter_test

import (
	"testing"

	"github.com/yourorg/pgdrift/internal/filter"
)

func TestAllowSchema_NoRules(t *testing.T) {
	f := filter.New(filter.Options{})
	if !f.AllowSchema("public") {
		t.Error("expected public schema to be allowed with no rules")
	}
}

func TestAllowSchema_Include(t *testing.T) {
	f := filter.New(filter.Options{IncludeSchemas: []string{"public"}})
	if !f.AllowSchema("public") {
		t.Error("expected public to be allowed")
	}
	if f.AllowSchema("private") {
		t.Error("expected private to be excluded")
	}
}

func TestAllowSchema_Exclude(t *testing.T) {
	f := filter.New(filter.Options{ExcludeSchemas: []string{"internal"}})
	if f.AllowSchema("internal") {
		t.Error("expected internal to be excluded")
	}
	if !f.AllowSchema("public") {
		t.Error("expected public to be allowed")
	}
}

func TestAllowTable_NoRules(t *testing.T) {
	f := filter.New(filter.Options{})
	if !f.AllowTable("users") {
		t.Error("expected users table to be allowed with no rules")
	}
}

func TestAllowTable_Include(t *testing.T) {
	f := filter.New(filter.Options{IncludeTables: []string{"users", "orders"}})
	if !f.AllowTable("users") {
		t.Error("expected users to be allowed")
	}
	if f.AllowTable("products") {
		t.Error("expected products to be excluded")
	}
}

func TestAllowTable_Exclude(t *testing.T) {
	f := filter.New(filter.Options{ExcludeTables: []string{"audit_log"}})
	if f.AllowTable("audit_log") {
		t.Error("expected audit_log to be excluded")
	}
	if !f.AllowTable("users") {
		t.Error("expected users to be allowed")
	}
}

func TestAllowTable_CaseInsensitive(t *testing.T) {
	f := filter.New(filter.Options{ExcludeTables: []string{"AuditLog"}})
	if f.AllowTable("auditlog") {
		t.Error("expected case-insensitive match to exclude auditlog")
	}
}
