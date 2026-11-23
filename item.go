package scope

// Pathable is an item that can be assigned to one or more scopes. Paths should
// return the scope paths for the item; returning nil leaves it unscoped.
type Pathable interface {
	// Paths returns the scope paths associated with the item.
	Paths() []*Path
}

type item struct {
	title string
}

func newItem(title string) item {
	return item{title: title}
}

func (t *item) Paths() []*Path {
	return ExtractPaths(t.title)
}
