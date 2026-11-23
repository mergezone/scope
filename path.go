package scope

import "strings"

// Path represents a slash-separated scope path as a linked sequence of
// segments. For example, "parent/child" contains the segments "parent" and
// "child".
type Path struct {
	parent *Path
	child  *Path

	// Value is the name of this path segment.
	Value string
}

// NewPath creates a scope path from str. It returns nil if str is blank or
// contains an empty or whitespace-only path segment.
func NewPath(str string) *Path {
	if strings.TrimSpace(str) == "" {
		return nil
	}

	str = strings.ToLower(str)
	parts := strings.Split(str, pathDelim)
	for _, part := range parts {
		if strings.TrimSpace(part) == "" {
			return nil
		}
	}

	root := &Path{Value: parts[0]}
	curr := root

	// Create children for remaining parts
	for _, part := range parts[1:] {
		next := &Path{
			parent: curr,
			Value:  part,
		}
		curr.child = next
		curr = next
	}

	return root
}

// String returns the slash-separated path starting at p.
func (p *Path) String() string {
	var parts []string
	for current := p; current != nil; current = current.child {
		parts = append(parts, current.Value)
	}

	return strings.Join(parts, pathDelim)
}

// ExtractPaths extracts scope paths from a title prefix. It supports simple
// prefixing schemes and conventional commits, and returns nil if no valid
// scope prefix is present.
func ExtractPaths(str string) []*Path {
	paths := extractSimplePrefixPaths(str)
	if len(paths) == 0 {
		return nil
	}

	result := make([]*Path, 0, len(paths))
	for _, path := range paths {
		parsed := NewPath(path)
		if parsed == nil {
			return nil
		}
		result = append(result, parsed)
	}

	return unique(result)
}

// extractSimplePrefixPaths extracts strings that may be converted to Path.
// Strings must be of one of the following forms, with scopes in the prefix:
// <scope>: <message>
// <scope>, <scope>: <message>
// [<scope>, <scope>] <message>
// [<scope>, <scope>]: <message>
func extractSimplePrefixPaths(str string) []string {
	str = strings.TrimSpace(str)
	if str == "" {
		return []string{}
	}

	var prefix string

	if strings.HasPrefix(str, "[") {
		end := strings.Index(str, "]")
		if end == -1 {
			return []string{}
		}
		prefix = str[1:end]
	} else if idx := strings.Index(str, ":"); idx != -1 {
		prefix = str[:idx]
	} else {
		return []string{}
	}

	prefix = strings.Join(strings.Fields(prefix), "")
	if prefix == "" {
		return []string{}
	}

	return strings.Split(prefix, ",")
}
