package diff

import (
	"fmt"

	"github.com/pgdrift/pgdrift/internal/schema"
)

// ChangeType represents the kind of schema change detected.
type ChangeType string

const (
	ChangeAdded   ChangeType = "added"
	ChangeRemoved ChangeType = "removed"
	ChangeAltered ChangeType = "altered"
)

// Change represents a single schema drift item.
type Change struct {
	Object     string
	ChangeType ChangeType
	Detail     string
}

func (c Change) String() string {
	return fmt.Sprintf("[%s] %s: %s", c.ChangeType, c.Object, c.Detail)
}

// Result holds all detected changes between two schemas.
type Result struct {
	Changes []Change
}

func (r *Result) HasDrift() bool {
	return len(r.Changes) > 0
}

func (r *Result) Add(c Change) {
	r.Changes = append(r.Changes, c)
}

// Compare detects schema drift between a source and target schema.
func Compare(source, target *schema.Schema) *Result {
	result := &Result{}

	sourceNames := toSet(source.TableNames())
	targetNames := toSet(target.TableNames())

	for name := range sourceNames {
		if _, ok := targetNames[name]; !ok {
			result.Add(Change{
				Object:     "table:" + name,
				ChangeType: ChangeRemoved,
				Detail:     fmt.Sprintf("table %q exists in source but not in target", name),
			})
			continue
		}
		compareTables(name, source, target, result)
	}

	for name := range targetNames {
		if _, ok := sourceNames[name]; !ok {
			result.Add(Change{
				Object:     "table:" + name,
				ChangeType: ChangeAdded,
				Detail:     fmt.Sprintf("table %q exists in target but not in source", name),
			})
		}
	}

	return result
}

func compareTables(tableName string, source, target *schema.Schema, result *Result) {
	srcTable, _ := source.Table(tableName)
	tgtTable, _ := target.Table(tableName)

	srcCols := toSet(srcTable.ColumnNames())
	tgtCols := toSet(tgtTable.ColumnNames())

	for col := range srcCols {
		if _, ok := tgtCols[col]; !ok {
			result.Add(Change{
				Object:     fmt.Sprintf("column:%s.%s", tableName, col),
				ChangeType: ChangeRemoved,
				Detail:     fmt.Sprintf("column %q removed from table %q", col, tableName),
			})
			continue
		}
		compareColumns(tableName, col, srcTable, tgtTable, result)
	}

	for col := range tgtCols {
		if _, ok := srcCols[col]; !ok {
			result.Add(Change{
				Object:     fmt.Sprintf("column:%s.%s", tableName, col),
				ChangeType: ChangeAdded,
				Detail:     fmt.Sprintf("column %q added to table %q", col, tableName),
			})
		}
	}
}

func compareColumns(tableName, colName string, src, tgt *schema.Table, result *Result) {
	srcCol, _ := src.Column(colName)
	tgtCol, _ := tgt.Column(colName)

	if srcCol.DataType != tgtCol.DataType {
		result.Add(Change{
			Object:     fmt.Sprintf("column:%s.%s", tableName, colName),
			ChangeType: ChangeAltered,
			Detail:     fmt.Sprintf("data type changed from %q to %q", srcCol.DataType, tgtCol.DataType),
		})
	}

	if srcCol.Nullable != tgtCol.Nullable {
		result.Add(Change{
			Object:     fmt.Sprintf("column:%s.%s", tableName, colName),
			ChangeType: ChangeAltered,
			Detail:     fmt.Sprintf("nullable changed from %v to %v", srcCol.Nullable, tgtCol.Nullable),
		})
	}

	if srcCol.Default != tgtCol.Default {
		result.Add(Change{
			Object:     fmt.Sprintf("column:%s.%s", tableName, colName),
			ChangeType: ChangeAltered,
			Detail:     fmt.Sprintf("default changed from %q to %q", srcCol.Default, tgtCol.Default),
		})
	}
}

func toSet(items []string) map[string]struct{} {
	s := make(map[string]struct{}, len(items))
	for _, item := range items {
		s[item] = struct{}{}
	}
	return s
}
