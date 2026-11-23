package scope

import "reflect"

// unique returns the unique elements of elems.
func unique[T any](elems []T) []T {
	result := make([]T, 0, len(elems))
	for _, elem := range elems {
		duplicate := false
		for _, existing := range result {
			if same(elem, existing) {
				duplicate = true
				break
			}
		}
		if !duplicate {
			result = append(result, elem)
		}
	}
	return result
}

func isNilValue(value any) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}

// same compares comparable values normally and falls back to deep equality for
// non-comparable values, which can be stored in interface-typed collections.
func same[T any](a, b T) bool {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)
	if !av.IsValid() || !bv.IsValid() {
		return !av.IsValid() && !bv.IsValid()
	}
	if av.Type() != bv.Type() {
		return false
	}
	if av.Comparable() {
		return av.Interface() == bv.Interface()
	}
	return reflect.DeepEqual(a, b)
}
