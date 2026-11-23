package scope

// Scope represents one named node in the manager's scope tree.
type Scope struct {
	parent   *Scope
	children []*Scope

	name string
	work []Work
}

func newScope(name string) Scope {
	return Scope{
		parent:   nil,
		children: make([]*Scope, 0),
		name:     name,
		work:     make([]Work, 0),
	}
}

const pathDelim = "/"

// Name returns the name of the scope.
func (s *Scope) Name() string {
	return s.name
}

// Works returns a copy of the work directly assigned to the scope.
func (s *Scope) Works() []Work {
	return append([]Work(nil), s.work...)
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

func (s *Scope) appendSubtreeWork(work []Work) []Work {
	for _, w := range s.work {
		work = append(work, w)
	}

	for _, child := range s.children {
		work = child.appendSubtreeWork(work)
	}

	return work
}

func (s *Scope) addWork(w Work) {
	s.work = append(s.work, w)
	s.work = unique(s.work)
}
