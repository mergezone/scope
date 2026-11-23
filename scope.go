package scope

// Scope represents one named node in the manager's scope tree.
type Scope struct {
	parent   *Scope
	children []*Scope

	name  string
	items []Pathable
}

func newScope(name string) Scope {
	return Scope{
		parent:   nil,
		children: make([]*Scope, 0),
		name:     name,
		items:    make([]Pathable, 0),
	}
}

const pathDelim = "/"

// Name returns the name of the scope.
func (s *Scope) Name() string {
	return s.name
}

// Items returns a copy of the items directly assigned to the scope.
func (s *Scope) Items() []Pathable {
	return append([]Pathable(nil), s.items...)
}

// String returns the full slash-separated path of the scope.
func (s *Scope) String() string {
	if s.parent == nil {
		return s.name
	}

	return s.parent.String() + pathDelim + s.name
}

// ScopeCount returns the number of scopes in this scope's subtree, including itself.
func (s *Scope) ScopeCount() int {
	// Base case: Include self in count
	n := 1

	for _, child := range s.children {
		n += child.ScopeCount()
	}

	return n
}

func (s *Scope) appendSubtreeItems(items []Pathable) []Pathable {
	for _, item := range s.items {
		items = append(items, item)
	}

	for _, child := range s.children {
		items = child.appendSubtreeItems(items)
	}

	return items
}

func (s *Scope) addItem(item Pathable) {
	s.items = append(s.items, item)
	s.items = unique(s.items)
}
