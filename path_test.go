package scope

import "testing"

var tests = []struct {
	in  string
	out [][]string
}{
	{"cmd: First command.", [][]string{{"cmd"}}},
	{"cmd, api: Combined work.", [][]string{{"cmd"}, {"api"}}},
	{"[cmd, api] Combined work.", [][]string{{"cmd"}, {"api"}}},
	{"[cmd, api]: Combined work.", [][]string{{"cmd"}, {"api"}}},
	{"parent/child: Nested work.", [][]string{{"parent", "child"}}},
	{"parent/child, child: Nested and combined work.", [][]string{{"parent", "child"}, {"child"}}},
}

func TestNewPathInvalidInput(t *testing.T) {
	for _, input := range []string{"", " ", "/", "/a", "a/", "a//b", "a/ /b"} {
		t.Run(input, func(t *testing.T) {
			if path := NewPath(input); path != nil {
				t.Fatalf("NewPath(%q) = %v, want nil", input, path)
			}
		})
	}
}

func TestExtractPathsInvalidInput(t *testing.T) {
	for _, input := range []string{"", "no scope prefix", "good,,bad: message", "[unfinished"} {
		t.Run(input, func(t *testing.T) {
			if paths := ExtractPaths(input); paths != nil {
				t.Fatalf("ExtractPaths(%q) = %v, want nil", input, paths)
			}
		})
	}
}

func TestExtractPaths(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			paths := ExtractPaths(tt.in)
			if len(paths) != len(tt.out) {
				t.Fatalf("incorrect number of paths: got %d, want %d", len(paths), len(tt.out))
			}

			for i, root := range paths {
				if root == nil {
					t.Fatalf("path %d is nil", i)
				}
				if root.parent != nil {
					t.Errorf("path %d root has parent %q, want nil", i, root.parent.Value)
				}

				current := root
				var parent *Path
				for j, want := range tt.out[i] {
					if current == nil {
						t.Fatalf("path %d has only %d segments, want %d", i, j, len(tt.out[i]))
					}
					if current.Value != want {
						t.Errorf("path %d segment %d: got %q, want %q", i, j, current.Value, want)
					}
					if current.parent != parent {
						t.Errorf("path %d segment %d has incorrect parent", i, j)
					}

					parent = current
					current = current.child
				}

				if current != nil {
					t.Errorf("path %d has unexpected extra segment %q", i, current.Value)
				}
			}
		})
	}
}
