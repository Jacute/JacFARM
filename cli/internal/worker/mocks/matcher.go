package mocks

import (
	"fmt"
	"reflect"

	"go.uber.org/mock/gomock"
)

type unorderedSliceMatcher struct {
	want interface{}
}

func (m unorderedSliceMatcher) Matches(x interface{}) bool {
	wantVal := reflect.ValueOf(m.want)
	gotVal := reflect.ValueOf(x)

	// оба должны быть слайсами
	if wantVal.Kind() != reflect.Slice || gotVal.Kind() != reflect.Slice {
		return false
	}

	// длина должна совпадать
	if wantVal.Len() != gotVal.Len() {
		return false
	}

	// карта счётчиков для элементов want
	used := make([]bool, wantVal.Len())

	// перебираем got и пытаемся найти deep-равный в want
	for i := 0; i < gotVal.Len(); i++ {
		match := false
		for j := 0; j < wantVal.Len(); j++ {
			if used[j] {
				continue
			}
			if reflect.DeepEqual(gotVal.Index(i).Interface(), wantVal.Index(j).Interface()) {
				used[j] = true
				match = true
				break
			}
		}
		if !match {
			return false
		}
	}

	return true
}

func (m unorderedSliceMatcher) String() string {
	return fmt.Sprintf("unordered slice %v", m.want)
}

func UnorderedSlice(want interface{}) gomock.Matcher {
	return unorderedSliceMatcher{want: want}
}
