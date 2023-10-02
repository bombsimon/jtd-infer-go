package jtdinfer

// Wildcard represents the character that matches any value for hints.
const Wildcard = "-"

// Hints contains the default number type to use and all the hints for enums,
// values and discriminators.
type Hints struct {
	DefaultNumType NumType
	Enums          *HintSet
	Values         *HintSet
	Discriminator  *HintSet
}

// NewHints creates a new empty non nil `Hints`.
func NewHints() *Hints {
	return &Hints{
		Enums:         NewHintSet(),
		Values:        NewHintSet(),
		Discriminator: NewHintSet(),
	}
}

// SubHints will return the sub hints for all hint sets for the passed key.
func (h Hints) SubHints(key string) *Hints {
	return &Hints{
		DefaultNumType: h.DefaultNumType,
		Enums:          h.Enums.SubHints(key),
		Values:         h.Values.SubHints(key),
		Discriminator:  h.Discriminator.SubHints(key),
	}
}

// IsEnumActive checks if the enum hint set is active.
func (h *Hints) IsEnumActive() bool {
	return h.Enums.IsActive()
}

// IsValuesActive checks if the values hint set is active.
func (h *Hints) IsValuesActive() bool {
	return h.Values.IsActive()
}

// PeekActiveDiscriminator will peek the currently active discriminator, if any.
// The returned boolean tells if there is an active discriminator.
func (h *Hints) PeekActiveDiscriminator() (string, bool) {
	return h.Discriminator.PeekActive()
}

// HintSet represents a list of paths (lists) to match for hints.
type HintSet struct {
	Values [][]string
}

// NewHintSet creates a new empty `HintSet`.
func NewHintSet() *HintSet {
	return &HintSet{
		Values: [][]string{},
	}
}

// Add will add a path (slice) to the `HintSet`.
func (h *HintSet) Add(v []string) *HintSet {
	h.Values = append(h.Values, v)
	return h
}

// SubHints will filter all the current sets and keep those who's first element
// matches the passed key or wildcard.
func (h *HintSet) SubHints(key string) *HintSet {
	filteredValues := [][]string{}

	for _, values := range h.Values {
		if len(values) == 0 {
			continue
		}

		first := values[0]
		if first == Wildcard || first == key {
			filteredValues = append(filteredValues, values[1:])
		}
	}

	return &HintSet{
		Values: filteredValues,
	}
}

// IsActive returns true if any set in the hint set his active.
func (h *HintSet) IsActive() bool {
	for _, valueList := range h.Values {
		if len(valueList) == 0 {
			return true
		}
	}

	return false
}

// PeekActive returns the currently active value if any. The returned boolean
// tells if a value was found.
func (h *HintSet) PeekActive() (string, bool) {
	for _, values := range h.Values {
		if len(values) != 1 {
			continue
		}

		return values[0], true
	}

	return "", false
}
