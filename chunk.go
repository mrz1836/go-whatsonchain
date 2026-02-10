package whatsonchain

// chunkSlice breaks a slice into chunks of the given size.
// The last chunk may be smaller than size.
func chunkSlice[T any](items []T, size int) [][]T {
	if size <= 0 || len(items) == 0 {
		return nil
	}
	numBatches := (len(items) + size - 1) / size
	batches := make([][]T, 0, numBatches)
	for i := 0; i < len(items); i += size {
		end := i + size
		if end > len(items) {
			end = len(items)
		}
		batches = append(batches, items[i:end])
	}
	return batches
}
