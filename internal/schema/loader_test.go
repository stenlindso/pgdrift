package schema

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestLoad(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{
		"table_name", "column_name", "data_type", "is_nullable", "column_default",
	}).
		AddRow("users", "id", "integer", "NO", nil).
		AddRow("users", "email", "text", "YES", nil).
		AddRow("orders", "id", "integer", "NO", nil).
		AddRow("orders", "total", "numeric", "YES", "0.00")

	mock.ExpectQuery(`SELECT`).WillReturnRows(rows)

	schema, err := Load(context.Background(), db)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if len(schema.Tables) != 2 {
		t.Errorf("expected 2 tables, got %d", len(schema.Tables))
	}

	users, ok := schema.Tables["users"]
	if !ok {
		t.Fatal("expected table 'users'")
	}
	if len(users.Columns) != 2 {
		t.Errorf("expected 2 columns in users, got %d", len(users.Columns))
	}

	orders, ok := schema.Tables["orders"]
	if !ok {
		t.Fatal("expected table 'orders'")
	}
	total := orders.Columns["total"]
	if total.Default == nil || *total.Default != "0.00" {
		t.Error("expected default '0.00' for orders.total")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet mock expectations: %v", err)
	}
}
