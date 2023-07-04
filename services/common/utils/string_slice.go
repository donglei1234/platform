package utils

type StringSlice []string

// Returns whether the given string is contained in this StringSlice.
func (s StringSlice) Contains(str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// Returns whether the given StringSlice is the same length and contains
// the same strings as this StringSlice.
func (s StringSlice) Equals(o StringSlice) bool {
	if len(s) != len(o) {
		return false
	}

	for i, v := range s {
		if v != o[i] {
			return false
		}
	}

	return true
}

// Filters out the given string from this StringSlice.
// Returns the resulting StringSlice and whether any changes were made from the source.
func (s StringSlice) FilterExclude(str string) (result StringSlice, found bool) {
	if !s.Contains(str) {
		return s, false
	}

	// Initialize result to len 0, but enough cap to accommodate only 1 instance of str being removed.
	// We know at least 1 will be removed because of the Contains test above.
	result = make(StringSlice, 0, len(s)-1)

	for _, v := range s {
		if v != str {
			result = append(result, v)
		}
	}

	return result, true
}

// Searches for a string and does a swap remove if found.
// Returns the modified StringSlice and whether the string was found.
func (s StringSlice) FindRemoveSwap(toRemove string) (StringSlice, bool) {
	for i, str := range s {
		if str == toRemove {
			s[i] = s[len(s)-1]
			return s[:len(s)-1], true
		}
	}

	return s, false
}
