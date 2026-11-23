// Package scope extracts scope paths from item titles and organizes items in a
// searchable hierarchy. Implementations of Pathable provide their scope paths,
// typically by passing a title to ExtractPaths, and Manager stores and queries
// the resulting scope tree.
package scope
