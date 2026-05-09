package diff

// Annotation holds metadata attached to a Change for display or processing.
type Annotation struct {
	Key   string
	Value string
}

// Annotate returns a copy of the Change with the given annotation appended.
func Annotate(c Change, key, value string) Change {
	c.Annotations = append(append([]Annotation(nil), c.Annotations...), Annotation{Key: key, Value: value})
	return c
}

// GetAnnotation returns the value of the first annotation with the given key,
// and whether it was found.
func GetAnnotation(c Change, key string) (string, bool) {
	for _, a := range c.Annotations {
		if a.Key == key {
			return a.Value, true
		}
	}
	return "", false
}

// AnnotateResult returns a new Result where every Change that matches the
// predicate is annotated with the given key/value pair.
func AnnotateResult(r *Result, key, value string, match func(Change) bool) *Result {
	if r == nil {
		return nil
	}
	out := &Result{}
	for _, c := range r.Changes {
		if match(c) {
			c = Annotate(c, key, value)
		}
		out.Changes = append(out.Changes, c)
	}
	return out
}
