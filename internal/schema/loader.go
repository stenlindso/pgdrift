package schema

import (
	"context"
	"database/sql"
	"fmt"
)

const queryColumns = `
SELECT
	t.table_name,
	c.column_name,
	c.data_type,
	c.is_nullable,
	c.column_default
FROM information_schema.tables t
JOIN information_schema.columns c
	ON c.table_name = t.table_name AND c.table_schema = t.table_schema
WHERE t.table_schema = 'public'
	AND t.table_type = 'BASE TABLE'
ORDER BY t.table_name, c.ordinal_position;
`

// Load queries a PostgreSQL database and returns its Schema.
func Load(ctx context.Context, db *sql.DB) (*Schema, error) {
	rows, err := db.QueryContext(ctx, queryColumns)
	if err != nil {
		return nil, fmt.Errorf("schema load query: %w", err)
	}
	defer rows.Close()

	schema := NewSchema()

	for rows.Next() {
		var (
			tableName  string
			colName    string
			dataType   string
			isNullable string
			colDefault sql.NullString
		)
		if err := rows.Scan(&tableName, &colName, &dataType, &isNullable, &colDefault); err != nil {
			return nil, fmt.Errorf("schema load scan: %w", err)
		}

		tbl, ok := schema.Tables[tableName]
		if !ok {
			tbl = Table{Name: tableName, Columns: make(map[string]Column)}
		}

		var def *string
		if colDefault.Valid {
			v := colDefault.String
			def = &v
		}

		tbl.Columns[colName] = Column{
			Name:       colName,
			DataType:   dataType,
			IsNullable: isNullable == "YES",
			Default:    def,
		}
		schema.Tables[tableName] = tbl
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("schema load rows: %w", err)
	}

	return schema, nil
}
