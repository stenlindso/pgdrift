package diff

// ChangeKind enumerates the types of schema changes pgdrift can detect.
type ChangeKind string

const (
	// Table-level changes
	ChangeKindTableAdded   ChangeKind = "table_added"
	ChangeKindTableRemoved ChangeKind = "table_removed"

	// Column-level changes
	ChangeKindColumnAdded      ChangeKind = "column_added"
	ChangeKindColumnRemoved    ChangeKind = "column_removed"
	ChangeKindColumnTypeChanged ChangeKind = "column_type_changed"
	ChangeKindColumnNullChanged ChangeKind = "column_null_changed"
	ChangeKindColumnDefault    ChangeKind = "column_default_changed"
)

// AllKinds returns all known ChangeKind values in a stable order.
func AllKinds() []ChangeKind {
	return []ChangeKind{
		ChangeKindTableAdded,
		ChangeKindTableRemoved,
		ChangeKindColumnAdded,
		ChangeKindColumnRemoved,
		ChangeKindColumnTypeChanged,
		ChangeKindColumnNullChanged,
		ChangeKindColumnDefault,
	}
}

// IsTableLevel reports whether the kind relates to a table-level change.
func (k ChangeKind) IsTableLevel() bool {
	return k == ChangeKindTableAdded || k == ChangeKindTableRemoved
}

// IsColumnLevel reports whether the kind relates to a column-level change.
func (k ChangeKind) IsColumnLevel() bool {
	return !k.IsTableLevel()
}
