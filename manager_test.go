package scope

import "testing"

func TestAddWorkPanicsOnNil(t *testing.T) {
	var typedNil *work
	for _, test := range []struct {
		name string
		work Work
	}{
		{name: "nil interface", work: nil},
		{name: "typed nil", work: typedNil},
	} {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("AddWork did not panic for nil Work")
				}
			}()
			NewManager().AddWork(test.work)
		})
	}
}

func TestUnscopedWork(t *testing.T) {
	w := newWork("Unscoped work")

	m := NewManager().AddWork(&w)

	if m.ScopeCount() != 0 {
		t.Fatalf("Scope added for unscoped work")
	}

	if m.WorkCount() != 0 {
		t.Fatalf("Work added for unscoped work")
	}
}

func TestDuplicateWork(t *testing.T) {
	w1 := newWork("test: Work.")

	m := NewManager().AddWork(&w1).AddWork(&w1)

	if m.ScopeCount() != 1 {
		t.Fatalf("Incorrect number of scopes created: got %d, want 1", m.ScopeCount())
	}

	if m.WorkCount() != 1 {
		t.Fatalf("Incorrect number of work created: got %d, want 1", m.WorkCount())
	}

	test := m.GetWork(NewPath("test"))

	if len(test) != 1 {
		t.Fatalf("Incorrect amount of work returned for test: got %d, want 1", len(test))
	}

	// GetScope returns the scope containing the work.
	s1 := m.GetScope(NewPath("test"))
	if s1 == nil {
		t.Fatal("GetScope failed to find scope")
	}
	if s1.Name() != "test" {
		t.Fatalf("Scope.Name() = %q, want %q", s1.Name(), "test")
	}
	works := s1.Works()
	if len(works) != 1 || works[0] != test[0] {
		t.Fatal("Scope does not contain the work returned by GetWork")
	}
	works[0] = nil
	if len(s1.Works()) != 1 || s1.Works()[0] != test[0] {
		t.Fatal("Scope.Works() exposed its internal slice")
	}

	// Test if scope nodes deduplicate by pointer
	if len(s1.Works()) != 1 {
		t.Fatalf("Incorrect amount of work in test scope: got %d, want 1", len(s1.Works()))
	}

	w2 := newWork("test: Work.")
	m.AddWork(&w2)

	if m.ScopeCount() != 1 {
		t.Fatalf("Incorrect number of scopes created: got %d, want 1", m.ScopeCount())
	}

	if m.WorkCount() != 2 {
		t.Fatalf("Incorrect number of work created: got %d, want 2", m.WorkCount())
	}

	test = m.GetWork(NewPath("test"))

	if len(test) != 2 {
		t.Fatalf("Incorrect amount of work returned for test: got %d, want 2", len(test))
	}
}

func TestWorkFromMatchingRootScope(t *testing.T) {
	w1 := newWork("cmd: First command.")
	w2 := newWork("cmd: Second command.")
	w3 := newWork("api: API.")

	m := NewManager().AddWork(&w1).AddWork(&w2).AddWork(&w3)

	if m.ScopeCount() != 2 {
		t.Fatalf("Incorrect number of scopes created: got %d, want 2", m.ScopeCount())
	}

	if m.WorkCount() != 3 {
		t.Fatalf("Incorrect number of work created: got %d, want 3", m.WorkCount())
	}

	cmd := m.GetWork(NewPath("cmd"))

	if len(cmd) != 2 {
		t.Fatalf("Incorrect amount of work returned for cmd: got %d, want 2", len(cmd))
	}

	api := m.GetWork(NewPath("api"))

	if len(api) != 1 {
		t.Fatalf("Incorrect amount of work returned for api: got %d, want 1", len(api))
	}
}

func TestNestedScope(t *testing.T) {
	w1 := newWork("parent: Parent work.")
	w2 := newWork("parent/child: Nested child work.")
	w3 := newWork("parent/other: Nested other work.")

	m := NewManager().AddWork(&w1).AddWork(&w2).AddWork(&w3)

	if m.ScopeCount() != 3 {
		t.Fatalf("Incorrect number of scopes created: got %d, want 3", m.ScopeCount())
	}

	if m.WorkCount() != 3 {
		t.Fatalf("Incorrect number of work created: got %d, want 3", m.WorkCount())
	}

	parent := m.GetWork(NewPath("parent"))

	if len(parent) != 3 {
		t.Fatalf("Incorrect amount of work returned for parent: got %d, want 3", len(parent))
	}

	child := m.GetWork(NewPath("child"))

	if len(child) != 1 {
		t.Fatalf("Incorrect amount of work returned for child: got %d, want 1", len(child))
	}

	other := m.GetWork(NewPath("other"))

	if len(other) != 1 {
		t.Fatalf("Incorrect amount of work returned for other: got %d, want 1", len(child))
	}

	par := m.GetWork(NewPath("par"))

	if len(par) != 0 {
		t.Fatalf("Incorrect amount of work returned for par: got %d, want 0", len(child))
	}
}

func TestWorkFromMatchingMultiRootScope(t *testing.T) {
	w1 := newWork("[cmd, api] Combined work.")
	w2 := newWork("cmd/api: Nested work.")

	m := NewManager().AddWork(&w1).AddWork(&w2)

	if m.ScopeCount() != 3 {
		t.Fatalf("Incorrect number of scopes created: got %d, want 3", m.ScopeCount())
	}

	if m.WorkCount() != 2 {
		t.Fatalf("Incorrect number of work created: got %d, want 2", m.WorkCount())
	}

	cmd := m.GetWork(NewPath("cmd"))

	if len(cmd) != 2 {
		t.Fatalf("Incorrect amount of work returned for cmd: got %d, want 2", len(cmd))
	}

	api := m.GetWork(NewPath("api"))

	if len(api) != 2 {
		t.Fatalf("Incorrect amount of work returned for api: got %d, want 2", len(api))
	}

	cmdApi := m.GetWork(NewPath("cmd/api"))

	if len(cmdApi) != 1 {
		t.Fatalf("Incorrect amount of work returned for cmd/api: got %d, want 1", len(api))
	}
}

func TestCaseSensitivity(t *testing.T) {
	w1 := newWork("api: Work.")
	w2 := newWork("API: WORK.")

	m := NewManager().AddWork(&w1).AddWork(&w2)

	if m.ScopeCount() != 1 {
		t.Fatalf("Incorrect number of scopes created: got %d, want 1", m.ScopeCount())
	}

	if m.WorkCount() != 2 {
		t.Fatalf("Incorrect number of work created: got %d, want 2", m.WorkCount())
	}

	api := m.GetWork(NewPath("api"))

	if len(api) != 2 {
		t.Fatalf("Incorrect amount of work returned for api: got %d, want 2", len(api))
	}

	API := m.GetWork(NewPath("API"))

	if len(API) != 2 {
		t.Fatalf("Incorrect amount of work returned for API: got %d, want 2", len(API))
	}
}

// TODO(hayden): Add tests for internal scope resolution
// Shared prefixes and scope reuse: Add a/b, a/c, and another a/b; verify the hierarchy is reused and lookups return the expected items.
// Repeated ResolveScopes: Call it more than once for the same work and verify it reuses existing scopes without adding duplicates to the manager’s tree.
