package helpers

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindLeftSequence(t *testing.T) {
	type args struct {
		data       []int
		finderFunc func(item int) bool
	}

	numbers := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "0 - return empty list",
			args: args{data: numbers, finderFunc: func(item int) bool { return false }},
			want: []int{},
		},
		{
			name: "1 - grab first element",
			args: args{data: numbers, finderFunc: func(item int) bool { return item == 0 }},
			want: []int{0},
		},
		{
			name: "2 - grab first two elements",
			args: args{data: numbers, finderFunc: func(item int) bool { return item == 1 }},
			want: []int{0, 1},
		},
		{
			name: "3 - grab first elements up to the middle of list",
			args: args{data: numbers, finderFunc: func(item int) bool { return item == 4 }},
			want: []int{0, 1, 2, 3, 4},
		},
		{
			name: "4 - grab all elements",
			args: args{data: numbers, finderFunc: func(item int) bool { return item == 9 }},
			want: numbers,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindLeftSequence(tt.args.data, tt.args.finderFunc)
			assert.Equal(t, tt.want, got)

			// if got := FindLeftSequence(tt.args.data, tt.args.finderFunc); !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("FindLeftSequence() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func TestFindRightSequence(t *testing.T) {
	type args struct {
		data       []int
		finderFunc func(item int) bool
	}

	numbers := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "0 - return empty list",
			args: args{data: numbers, finderFunc: func(item int) bool { return false }},
			want: []int{},
		},
		{
			name: "1 - grab last element",
			args: args{data: numbers, finderFunc: func(item int) bool { return item == 9 }},
			want: []int{9},
		},
		{
			name: "2 - grab last two elements",
			args: args{data: numbers, finderFunc: func(item int) bool { return item == 8 }},
			want: []int{8, 9},
		},
		{
			name: "3 - grab last elements from the middle of list",
			args: args{data: numbers, finderFunc: func(item int) bool { return item == 4 }},
			want: []int{4, 5, 6, 7, 8, 9},
		},
		{
			name: "4 - grab all elements",
			args: args{data: numbers, finderFunc: func(item int) bool { return item == 0 }},
			want: numbers,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindRightSequence(tt.args.data, tt.args.finderFunc)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Map(t *testing.T) {
	type args struct {
		collection    []int
		transformFunc func(int, int) int
	}

	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "multiply",
			args: args{
				collection: []int{3, 5, 7},
				transformFunc: func(i int, n int) int {
					return n * n
				},
			},
			want: []int{9, 25, 49},
		},
		{
			name: "add",
			args: args{
				collection: []int{3, 5, 7},
				transformFunc: func(i int, n int) int {
					return n + 1
				},
			},
			want: []int{4, 6, 8},
		},
		{
			name: "subtract",
			args: args{
				collection: []int{3, 5, 7},
				transformFunc: func(i int, n int) int {
					return n - 1
				},
			},
			want: []int{2, 4, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Map(tt.args.collection, tt.args.transformFunc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}
