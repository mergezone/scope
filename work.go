package scope

// Work is an item that can be assigned to one or more scopes. Paths should
// return the scope paths for the item; returning nil leaves it unscoped.
type Work interface {
	// Paths returns the scope paths associated with the work item.
	Paths() []*Path
}

type work struct {
	title string
}

func newWork(title string) work {
	return work{title: title}
}

func (w *work) Paths() []*Path {
	return ExtractPaths(w.title)
}
