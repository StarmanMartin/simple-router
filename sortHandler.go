package router

import (
	"sort"
)

// ByAge implements sort.Interface for []Person based on
// the Age field.
type ByIndex []*finalRouteElement

func (a ByIndex) Len() int           { return len(a) }
func (a ByIndex) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByIndex) Less(i, j int) bool { return a[i].index < a[j].index }

func sortRoutingList(list []*finalRouteElement) {
	sort.Sort(ByIndex(list))
}
