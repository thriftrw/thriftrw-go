package lint

import "sort"

type byLine []Lint

func (ll byLine) Len() int {
	return len(ll)
}

func (ll byLine) Less(i, j int) bool {
	return ll[i].Line < ll[j].Line
}

func (ll byLine) Swap(i, j int) {
	ll[i], ll[j] = ll[j], ll[i]
}

// SortByLine sorts the given collection of lints by line number.
func SortByLine(ll []Lint) {
	sort.Sort(byLine(ll))
}
