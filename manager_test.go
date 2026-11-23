package scope

import "testing"

func TestAddItemPanicsOnNil(t *testing.T) {
	var typedNil *item
	for _, test := range []struct {
		name string
		item Pathable
	}{
		{name: "nil interface", item: nil},
		{name: "typed nil", item: typedNil},
	} {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("AddItem did not panic for nil Item")
				}
			}()
			NewManager().AddItem(test.item)
		})
	}
}

func TestUnscopedItem(t *testing.T) {
	item := newItem("Unscoped item")

	m := NewManager().AddItem(&item)

	if m.ScopeCount() != 0 {
		t.Fatalf("Scope added for unscoped item")
	}

	if m.ItemCount() != 0 {
		t.Fatalf("Item added for unscoped item")
	}
}

func TestDuplicateItem(t *testing.T) {
	item1 := newItem("test: Item.")

	m := NewManager().AddItem(&item1).AddItem(&item1)

	if m.ScopeCount() != 1 {
		t.Fatalf("Incorrect number of scopes created: got %d, want 1", m.ScopeCount())
	}

	if m.ItemCount() != 1 {
		t.Fatalf("Incorrect number of items created: got %d, want 1", m.ItemCount())
	}

	test := m.GetItems(NewPath("test"))

	if len(test) != 1 {
		t.Fatalf("Incorrect number of items returned for test: got %d, want 1", len(test))
	}

	// GetScope returns the scope containing the item.
	s1 := m.GetScope(NewPath("test"))
	if s1 == nil {
		t.Fatal("GetScope failed to find scope")
	}
	if s1.Name() != "test" {
		t.Fatalf("Scope.Name() = %q, want %q", s1.Name(), "test")
	}
	items := s1.Items()
	if len(items) != 1 || items[0] != test[0] {
		t.Fatal("Scope does not contain the item returned by GetItems")
	}
	items[0] = nil
	if len(s1.Items()) != 1 || s1.Items()[0] != test[0] {
		t.Fatal("Scope.Items() exposed its internal slice")
	}

	// Test if scope nodes deduplicate by pointer
	if len(s1.Items()) != 1 {
		t.Fatalf("Incorrect number of items in test scope: got %d, want 1", len(s1.Items()))
	}

	item2 := newItem("test: Item.")
	m.AddItem(&item2)

	if m.ScopeCount() != 1 {
		t.Fatalf("Incorrect number of scopes created: got %d, want 1", m.ScopeCount())
	}

	if m.ItemCount() != 2 {
		t.Fatalf("Incorrect number of items created: got %d, want 2", m.ItemCount())
	}

	test = m.GetItems(NewPath("test"))

	if len(test) != 2 {
		t.Fatalf("Incorrect number of items returned for test: got %d, want 2", len(test))
	}
}

func TestItemFromMatchingRootScope(t *testing.T) {
	item1 := newItem("cmd: First command.")
	item2 := newItem("cmd: Second command.")
	item3 := newItem("api: API.")

	m := NewManager().AddItem(&item1).AddItem(&item2).AddItem(&item3)

	if m.ScopeCount() != 2 {
		t.Fatalf("Incorrect number of scopes created: got %d, want 2", m.ScopeCount())
	}

	if m.ItemCount() != 3 {
		t.Fatalf("Incorrect number of items created: got %d, want 3", m.ItemCount())
	}

	cmd := m.GetItems(NewPath("cmd"))

	if len(cmd) != 2 {
		t.Fatalf("Incorrect number of items returned for cmd: got %d, want 2", len(cmd))
	}

	api := m.GetItems(NewPath("api"))

	if len(api) != 1 {
		t.Fatalf("Incorrect number of items returned for api: got %d, want 1", len(api))
	}
}

func TestNestedScope(t *testing.T) {
	item1 := newItem("parent: Parent item.")
	item2 := newItem("parent/child: Nested child item.")
	item3 := newItem("parent/other: Nested other item.")

	m := NewManager().AddItem(&item1).AddItem(&item2).AddItem(&item3)

	if m.ScopeCount() != 3 {
		t.Fatalf("Incorrect number of scopes created: got %d, want 3", m.ScopeCount())
	}

	if m.ItemCount() != 3 {
		t.Fatalf("Incorrect number of items created: got %d, want 3", m.ItemCount())
	}

	parent := m.GetItems(NewPath("parent"))

	if len(parent) != 3 {
		t.Fatalf("Incorrect number of items returned for parent: got %d, want 3", len(parent))
	}

	child := m.GetItems(NewPath("child"))

	if len(child) != 1 {
		t.Fatalf("Incorrect number of items returned for child: got %d, want 1", len(child))
	}

	other := m.GetItems(NewPath("other"))

	if len(other) != 1 {
		t.Fatalf("Incorrect number of items returned for other: got %d, want 1", len(child))
	}

	par := m.GetItems(NewPath("par"))

	if len(par) != 0 {
		t.Fatalf("Incorrect number of items returned for par: got %d, want 0", len(child))
	}
}

func TestItemFromMatchingMultiRootScope(t *testing.T) {
	item1 := newItem("[cmd, api] Combined item.")
	item2 := newItem("cmd/api: Nested item.")

	m := NewManager().AddItem(&item1).AddItem(&item2)

	if m.ScopeCount() != 3 {
		t.Fatalf("Incorrect number of scopes created: got %d, want 3", m.ScopeCount())
	}

	if m.ItemCount() != 2 {
		t.Fatalf("Incorrect number of items created: got %d, want 2", m.ItemCount())
	}

	cmd := m.GetItems(NewPath("cmd"))

	if len(cmd) != 2 {
		t.Fatalf("Incorrect number of items returned for cmd: got %d, want 2", len(cmd))
	}

	api := m.GetItems(NewPath("api"))

	if len(api) != 2 {
		t.Fatalf("Incorrect number of items returned for api: got %d, want 2", len(api))
	}

	cmdApi := m.GetItems(NewPath("cmd/api"))

	if len(cmdApi) != 1 {
		t.Fatalf("Incorrect number of items returned for cmd/api: got %d, want 1", len(api))
	}
}

func TestCaseSensitivity(t *testing.T) {
	item1 := newItem("api: Item.")
	item2 := newItem("API: ITEM.")

	m := NewManager().AddItem(&item1).AddItem(&item2)

	if m.ScopeCount() != 1 {
		t.Fatalf("Incorrect number of scopes created: got %d, want 1", m.ScopeCount())
	}

	if m.ItemCount() != 2 {
		t.Fatalf("Incorrect number of items created: got %d, want 2", m.ItemCount())
	}

	api := m.GetItems(NewPath("api"))

	if len(api) != 2 {
		t.Fatalf("Incorrect number of items returned for api: got %d, want 2", len(api))
	}

	API := m.GetItems(NewPath("API"))

	if len(API) != 2 {
		t.Fatalf("Incorrect number of items returned for API: got %d, want 2", len(API))
	}
}

// TODO(hayden): Add tests for internal scope resolution
// Shared prefixes and scope reuse: Add a/b, a/c, and another a/b; verify the hierarchy is reused and lookups return the expected items.
// Repeated ResolveScopes: Call it more than once for the same item and verify it reuses existing scopes without adding duplicates to the manager’s tree.
