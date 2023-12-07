package slice

func Prepend[T any](dest []T, value T, bufferCap int) []T {
	ogLength := len(dest)

	// Available capacity, shift content right and prepend
	if cap(dest) > ogLength {
		dest = dest[:ogLength+1]
		copy(dest[1:], dest)
		dest[0] = value

		return dest
	}

	// Capped slice, create new and pre-allocate some extra space for the future
	res := make([]T, ogLength+1, len(dest)+bufferCap)
	res[0] = value
	copy(res[1:], dest)

	return res
}
