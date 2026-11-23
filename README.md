# scope

`scope` parses scope prefixes from item titles and organizes items in a
hierarchical scope tree. Use it to group and query items by paths such as
`api/http`.

An item implements `scope.Pathable` by returning its scope paths from `Paths()`.
`ExtractPaths` provides the package's built-in parser for title prefixes such as
`api/http: Add an endpoint` and `[api, cli] Add a feature`.

## Example

```go
package main

import (
	"fmt"

	"github.com/mergezone/scope"
)

type item struct {
	title string
}

func (t item) Paths() []*scope.Path {
	return scope.ExtractPaths(t.title)
}

func main() {
	manager := scope.NewManager()
	manager.AddItem(item{title: "api/http: Add an endpoint"})

	fmt.Printf("scopes: %d, items: %d\n", manager.ScopeCount(), manager.ItemCount())
	fmt.Printf("items under api: %d\n", len(manager.GetItems(scope.NewPath("api"))))
}
```

The example prints:

```text
scopes: 2, items: 1
items under api: 1
```
