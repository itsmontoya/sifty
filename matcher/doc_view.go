package matcher

// DocView provides field-path access for matcher evaluation.
// Implementations should return (nil, false, nil) when the path is missing.
type DocView interface {
	// Get returns the value for a dot-delimited path.
	// It returns ok=false with err=nil when the path cannot be resolved.
	Get(path string) (any, bool, error)
}
