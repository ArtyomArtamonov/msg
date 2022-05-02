package utils

// Unwrap
//
// Safely casts object to error.
//
// We need this, because there's no way to check if [T] is [nil] prematurely,
// otherwise we fall into error trying to cast [nil] to [T].
func Unwrap[T any](object interface{}) T {
	if object == nil {
		var result T
		return result
	}
	return object.(T)
}
