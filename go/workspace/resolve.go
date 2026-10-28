package workspace

import "slices"

// ResolvedPaths holds deduplicated absolute paths, split by kind.
type ResolvedPaths struct {
	Readable  []string
	Writable  []string
	Creatable []string
}

// All returns every path, readable first.
func (p ResolvedPaths) All() []string {
	return slices.Concat(p.Readable, p.Writable, p.Creatable)
}

// Editable returns the writable and creatable paths.
func (p ResolvedPaths) Editable() []string {
	return slices.Concat(p.Writable, p.Creatable)
}

type resolvedCache struct {
	paths ResolvedPaths
	err   error
}
