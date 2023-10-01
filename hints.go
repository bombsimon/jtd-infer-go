package jtdinfer

const Wildcard = "-"

type Hints struct {
	DefaultNumType NumType
	Enums          HintSet
	Values         HintSet
	Discriminator  HintSet
}

func New(defaultNumType NumType, enums, values, discriminator HintSet) Hints {
	return Hints{
		DefaultNumType: defaultNumType,
		Enums:          enums,
		Values:         values,
		Discriminator:  discriminator,
	}
}

func (h Hints) SubHints(key string) Hints {
	return Hints{
		DefaultNumType: h.DefaultNumType,
		Enums:          h.Enums.SubHints(key),
		Values:         h.Values.SubHints(key),
		Discriminator:  h.Discriminator.SubHints(key),
	}
}

func (h Hints) IsEnumActive() bool {
	return h.Enums.IsActive()
}

func (h Hints) IsValuesActive() bool {
	return h.Values.IsActive()
}

func (h Hints) PeekActiveDiscriminator() (string, bool) {
	return h.Discriminator.PeekActive()
}

type HintSet struct {
	Values [][]string
}

func NewHintSet(values [][]string) HintSet {
	return HintSet{
		Values: values,
	}
}

func (h HintSet) SubHints(key string) HintSet {
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

	return HintSet{
		Values: filteredValues,
	}
}

func (h HintSet) IsActive() bool {
	for _, valueList := range h.Values {
		if len(valueList) == 0 {
			return true
		}
	}

	return false
}

func (h HintSet) PeekActive() (string, bool) {
	for _, values := range h.Values {
		if len(values) != 1 {
			continue
		}

		return values[0], true
	}

	return "", false
}
