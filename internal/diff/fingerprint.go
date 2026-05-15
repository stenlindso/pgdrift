package diff

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/pgdrift/pgdrift/internal/schema"
)

// Fingerprint is a stable hash representing a schema's structure.
type Fingerprint struct {
	Hash   string `json:"hash"`
	Tables int    `json:"tables"`
	Columns int   `json:"columns"`
}

// String returns the short (8-char) hex prefix of the fingerprint hash.
func (f Fingerprint) String() string {
	if len(f.Hash) >= 8 {
		return f.Hash[:8]
	}
	return f.Hash
}

// Equal reports whether two fingerprints represent identical schemas.
func (f Fingerprint) Equal(other Fingerprint) bool {
	return f.Hash == other.Hash
}

// SchemaFingerprint computes a deterministic Fingerprint for the given schema.
// The hash is derived from the sorted set of table and column definitions so
// that two schemas with identical structure always produce the same value.
func SchemaFingerprint(s *schema.Schema) Fingerprint {
	if s == nil {
		return Fingerprint{}
	}

	h := sha256.New()
	names := s.TableNames()
	sort.Strings(names)

	totalCols := 0
	for _, tname := range names {
		t, ok := s.Table(tname)
		if !ok {
			continue
		}
		cols := t.ColumnNames()
		sort.Strings(cols)
		totalCols += len(cols)
		for _, cname := range cols {
			c, _ := t.Column(cname)
			line := fmt.Sprintf("%s.%s:%s:nullable=%v:default=%s\n",
				tname, cname, c.DataType, c.Nullable, c.Default)
			h.Write([]byte(line))
		}
	}

	return Fingerprint{
		Hash:    hex.EncodeToString(h.Sum(nil)),
		Tables:  len(names),
		Columns: totalCols,
	}
}

// FingerprintDiff describes the change between two schema fingerprints.
type FingerprintDiff struct {
	Source Fingerprint `json:"source"`
	Target Fingerprint `json:"target"`
	Changed bool       `json:"changed"`
}

// CompareFingerprints returns a FingerprintDiff for two schemas.
func CompareFingerprints(src, tgt *schema.Schema) FingerprintDiff {
	s := SchemaFingerprint(src)
	t := SchemaFingerprint(tgt)
	return FingerprintDiff{
		Source:  s,
		Target:  t,
		Changed: !s.Equal(t),
	}
}

// FingerprintSummary returns a human-readable one-line summary.
func FingerprintSummary(fd FingerprintDiff) string {
	if !fd.Changed {
		return fmt.Sprintf("schemas match (fingerprint %s, %d tables, %d columns)",
			fd.Source.String(), fd.Source.Tables, fd.Source.Columns)
	}
	return fmt.Sprintf("schema drift detected: %s → %s (%s)",
		fd.Source.String(), fd.Target.String(),
		deltaDesc(fd.Source, fd.Target))
}

func deltaDesc(src, tgt Fingerprint) string {
	parts := []string{}
	dt := tgt.Tables - src.Tables
	dc := tgt.Columns - src.Columns
	if dt != 0 {
		parts = append(parts, fmt.Sprintf("%+d tables", dt))
	}
	if dc != 0 {
		parts = append(parts, fmt.Sprintf("%+d columns", dc))
	}
	if len(parts) == 0 {
		return "structure changed"
	}
	return strings.Join(parts, ", ")
}
