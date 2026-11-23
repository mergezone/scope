package scope

// Manager owns a tree of scopes and the items associated with them.
type Manager struct {
	roots []*Scope
}

// NewManager returns an empty Manager.
func NewManager() *Manager {
	return &Manager{
		roots: make([]*Scope, 0),
	}
}

// ScopeCount returns the total number of scopes owned by m.
func (m *Manager) ScopeCount() int {
	n := 0

	for _, root := range m.roots {
		n += root.ScopeCount()
	}

	return n
}

// ItemCount returns the number of unique items owned by m.
func (m *Manager) ItemCount() int {
	items := make([]Pathable, 0)

	for _, root := range m.roots {
		items = root.appendSubtreeItems(items)
	}

	return len(unique(items))
}

// GetScope returns the scope at path, or nil if path is nil or does not exist.
func (m *Manager) GetScope(path *Path) *Scope {
	var current *Scope

	for segment := path; segment != nil; segment = segment.child {
		var candidates []*Scope
		if current == nil {
			candidates = m.roots
		} else {
			candidates = current.children
		}

		var next *Scope
		for _, c := range candidates {
			if c.name == segment.Value {
				next = c
				break
			}
		}

		// Segment didn't match any candidates; fail
		if next == nil {
			return nil
		}
		current = next
	}

	// Found a match for all segments; current contains the deepest matching scope
	return current
}

// GetItems returns items assigned to any scope whose path ends with prefix, or
// to any of that scope's descendants. It returns nil if prefix is nil or no
// matching scope exists.
func (m *Manager) GetItems(prefix *Path) []Pathable {
	if prefix == nil {
		return nil
	}

	found := make([]Pathable, 0)
	var visit func(*Scope)
	visit = func(scope *Scope) {
		if scope.name == prefix.Value {
			match := scope
			for segment := prefix.child; segment != nil && match != nil; segment = segment.child {
				var next *Scope
				for _, child := range match.children {
					if child.name == segment.Value {
						next = child
						break
					}
				}
				match = next
			}
			if match != nil {
				found = match.appendSubtreeItems(found)
			}
		}

		for _, child := range scope.children {
			visit(child)
		}
	}

	for _, root := range m.roots {
		visit(root)
	}

	return unique(found)
}

// AddItem adds an item to all of the scopes it belongs to. It panics if item is
// nil, including when item contains a typed nil value.
func (m *Manager) AddItem(item Pathable) *Manager {
	if isNilValue(item) {
		panic("scope: AddItem called with nil Pathable")
	}

	scopes := m.resolveScopes(item)

	for _, scope := range scopes {
		scope.addItem(item)
	}

	return m
}

// resolveScopes finds or creates the scopes that the item belongs to.
func (m *Manager) resolveScopes(item Pathable) []*Scope {
	found := make([]*Scope, 0)

	for _, root := range item.Paths() {
		var current *Scope

		for segment := root; segment != nil; segment = segment.child {
			var candidates []*Scope
			if current == nil {
				candidates = m.roots
			} else {
				candidates = current.children
			}

			var next *Scope
			for _, candidate := range candidates {
				if candidate.name == segment.Value {
					next = candidate
					break
				}
			}

			// Create a new scope if no candidates matched
			if next == nil {
				scope := newScope(segment.Value)
				scope.parent = current
				next = &scope

				if current == nil {
					m.roots = append(m.roots, next)
				} else {
					current.children = append(current.children, next)
				}
			}

			current = next
		}

		if current != nil {
			found = append(found, current)
		}
	}

	return found
}
