package zkwasm

import (
	"reflect"
	"testing"
)

func TestChunkSlice(t *testing.T) {
	type args[T any] struct {
		slice     []T
		chunkSize int
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want [][]T
	}
	tests := []testCase[int]{
		{
			name: "empty",
			args: args[int]{
				slice:     []int{},
				chunkSize: 10,
			},
			want: [][]int{},
		},
		{
			name: "less than chunk",
			args: args[int]{
				slice:     []int{1, 2, 3, 4, 5},
				chunkSize: 10,
			},
			want: [][]int{{1, 2, 3, 4, 5}},
		},
		{
			name: "equal chunkSize",
			args: args[int]{
				slice:     []int{1, 2, 3, 4, 5},
				chunkSize: 5,
			},
			want: [][]int{{1, 2, 3, 4, 5}},
		}, {
			name: "more than chunkSize",
			args: args[int]{
				slice:     []int{1, 2, 3, 4, 5},
				chunkSize: 3,
			},
			want: [][]int{{1, 2, 3}, {4, 5}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ChunkSlice(tt.args.slice, tt.args.chunkSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChunkSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
