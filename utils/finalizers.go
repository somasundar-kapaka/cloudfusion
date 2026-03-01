package utils

const (
	FinalizerKey = "cloudfusion.com/compute"
)

func AddFinalizers(f *[]string, groups ...string) {
	if f == nil || len(*f) == 0 {
		*f = make([]string, 0)
	}

	for _, key := range groups {
		*f = append(*f, key)
	}

}

func DeleteFinalizers(f *[]string, groups ...string) {
	if f == nil || len(*f) == 0 {
		return
	}

	// Convert groups to a map for O(1) lookup
	removeMap := make(map[string]struct{})
	for _, g := range groups {
		removeMap[g] = struct{}{}
	}

	// Rebuild slice
	newSlice := (*f)[:0] // reuse underlying array

	for _, v := range *f {
		if _, shouldRemove := removeMap[v]; !shouldRemove {
			newSlice = append(newSlice, v)
		}
	}

	*f = newSlice
}
