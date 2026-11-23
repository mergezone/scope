// Package scope extracts scope paths from work titles and organizes work in a
// searchable hierarchy. Implementations of Work provide their scope paths,
// typically by passing a title to ExtractPaths, and Manager stores and queries
// the resulting scope tree.
package scope
