package utils

import (
	"reflect"
	"sort"
)

// Contains check a string
func Contains(slice []string, value string) bool {
	return ArrayContains(slice, value)
}

// ArrayContains checks if an array contains the element
func ArrayContains(l, elem interface{}) bool {
	vl := reflect.ValueOf(l)
	for i := 0; i < vl.Len(); i++ {
		if ObjectEqual(vl.Index(i).Interface(), elem) {
			return true
		}
	}
	return false
}

// ObjectEqual checks if objects are equal
func ObjectEqual(a, b interface{}) bool {
	if a == nil || b == nil {
		return a == b
	}
	return reflect.DeepEqual(a, b)
}

// StringArraysEqual checks if two string arrays are equal without the order
func StringArraysEqual(a, b []string) (ok bool) {
	if len(a) != len(b) {
		return
	}

	sort.Strings(a)
	sort.Strings(b)
	for i := range a {
		if a[i] != b[i] {
			return
		}
	}
	return true
}
