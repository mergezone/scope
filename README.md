# scope

`scope` parses scope prefixes from work titles and organizes work items in a
hierarchical scope tree. Use it to group and query work by paths such as
`api/http`.

A work item implements `scope.Work` by returning its scope paths from `Paths()`.
`ExtractPaths` provides the package's built-in parser for title prefixes such as
`api/http: Add an endpoint` and `[api, cli] Add a feature`.

## Example

```go
package main

import (
	"fmt"

	"github.com/mergezone/scope"
)

type task struct {
	title string
}

func (t task) Paths() []*scope.Path {
	return scope.ExtractPaths(t.title)
}

func main() {
	manager := scope.NewManager()
	manager.AddWork(task{title: "api/http: Add an endpoint"})

	fmt.Printf("scopes: %d, work: %d\n", manager.ScopeCount(), manager.WorkCount())
	fmt.Printf("work under api: %d\n", len(manager.GetWork(scope.NewPath("api"))))
}
```

The example prints:

```text
scopes: 2, work: 1
work under api: 1
```
