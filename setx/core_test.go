package setx_test

import (
	"github.com/shijl0925/go-toolkits/setx"
	"github.com/shijl0925/go-toolkits/slice"
	"reflect"
	"testing"
)

// TestNewSet tests the NewSet function.
func TestNewSet(t *testing.T) {
	t.Run("empty set", func(t *testing.T) {
		s := setx.NewSet[int](0)
		if len(s) != 0 {
			t.Errorf("expected empty set, got length %d", len(s))
		}
	})

	t.Run("set with initial elements", func(t *testing.T) {
		s := setx.NewSet(3, 1, 2, 3)
		expected := setx.Set[int]{1: struct{}{}, 2: struct{}{}, 3: struct{}{}}
		if !reflect.DeepEqual(s, expected) {
			t.Errorf("expected %v, got %v", expected, s)
		}
	})
}

// TestAdd tests the Add method.
func TestAdd(t *testing.T) {
	s := setx.NewSet[int](0)

	t.Run("add new element", func(t *testing.T) {
		s.Add(1)
		if !s.Exists(1) {
			t.Errorf("element 1 should exist")
		}
		if s.Len() != 1 {
			t.Errorf("expected length 1, got %d", s.Len())
		}
	})

	t.Run("add duplicate element", func(t *testing.T) {
		s.Add(1)
		if s.Len() != 1 {
			t.Errorf("duplicate add should not increase length, got %d", s.Len())
		}
	})
}

// TestRemove tests the Remove method.
func TestRemove(t *testing.T) {
	s := setx.NewSet(0, 1, 2)

	t.Run("remove existing element", func(t *testing.T) {
		s.Remove(1)
		if s.Exists(1) {
			t.Errorf("element 1 should be removed")
		}
		if s.Len() != 1 {
			t.Errorf("expected length 1 after remove, got %d", s.Len())
		}
	})

	t.Run("remove non-existing element", func(t *testing.T) {
		s.Remove(3)
		if s.Len() != 1 {
			t.Errorf("removing non-existing element should not affect length, got %d", s.Len())
		}
	})
}

// TestLen tests the Len method.
func TestLen(t *testing.T) {
	s := setx.NewSet[int](0)
	if s.Len() != 0 {
		t.Errorf("expected length 0, got %d", s.Len())
	}

	s.Add(1)
	s.Add(2)
	if s.Len() != 2 {
		t.Errorf("expected length 2, got %d", s.Len())
	}
}

// TestExists tests the Exists method.
func TestExists(t *testing.T) {
	s := setx.NewSet(0, 1)

	t.Run("existing element", func(t *testing.T) {
		if !s.Exists(1) {
			t.Errorf("element 1 should exist")
		}
	})

	t.Run("non-existing element", func(t *testing.T) {
		if s.Exists(2) {
			t.Errorf("element 2 should not exist")
		}
	})
}

func Test_Keys(t *testing.T) {
	tests := []struct {
		name string
		src  setx.Set[int]
		want []int
	}{
		{
			name: "TC01 - Src is not empty",
			src:  setx.NewSet(3, 1, 2, 3),
			want: []int{3, 2, 1},
		},
		{
			name: "TC02 - Src is empty",
			src:  setx.NewSet[int](0),
			want: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.src.Keys()
			if len(slice.DiffSet(got, tt.want)) != 0 || len(slice.DiffSet(tt.want, got)) != 0 {
				t.Errorf("Keys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_DiffSet(t *testing.T) {
	tests := []struct {
		name string
		src  setx.Set[int]
		dst  setx.Set[int]
		want []int
	}{
		{
			name: "TC01 - Some elements in src not in dst",
			src:  setx.NewSet(3, 1, 2, 3),
			dst:  setx.NewSet(3, 2, 3, 4),
			want: []int{1},
		},
		{
			name: "TC02 - Src is empty",
			src:  setx.NewSet[int](0),
			dst:  setx.NewSet(2, 1, 2),
			want: []int{},
		},
		{
			name: "TC03 - Dst is empty",
			src:  setx.NewSet[int](2, 5, 6),
			dst:  setx.NewSet[int](0),
			want: []int{5, 6},
		},
		{
			name: "TC04 - All elements overlap",
			src:  setx.NewSet[int](2, 7, 8),
			dst:  setx.NewSet[int](2, 7, 8),
			want: []int{},
		},
		{
			name: "TC05 - No overlap",
			src:  setx.NewSet[int](2, 9, 10),
			dst:  setx.NewSet[int](2, 11, 12),
			want: []int{9, 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := setx.DiffSet(tt.src, tt.dst)
			if len(slice.DiffSet(got, tt.want)) != 0 || len(slice.DiffSet(tt.want, got)) != 0 {
				t.Errorf("DiffSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IntersectSet(t *testing.T) {
	tests := []struct {
		name string
		src  setx.Set[int]
		dst  setx.Set[int]
		want []int
	}{
		{
			name: "TC01 - Some elements in src not in dst",
			src:  setx.NewSet(3, 1, 2, 3),
			dst:  setx.NewSet(3, 2, 3, 4),
			want: []int{2, 3},
		},
		{
			name: "TC02 - Src is empty",
			src:  setx.NewSet[int](0),
			dst:  setx.NewSet(2, 1, 2),
			want: []int{},
		},
		{
			name: "TC03 - Dst is empty",
			src:  setx.NewSet[int](2, 5, 6),
			dst:  setx.NewSet[int](0),
			want: []int{},
		},
		{
			name: "TC04 - All elements overlap",
			src:  setx.NewSet[int](2, 7, 8),
			dst:  setx.NewSet[int](2, 7, 8),
			want: []int{7, 8},
		},
		{
			name: "TC05 - No overlap",
			src:  setx.NewSet[int](2, 9, 10),
			dst:  setx.NewSet[int](2, 11, 12),
			want: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := setx.IntersectSet(tt.src, tt.dst)
			if len(slice.DiffSet(got, tt.want)) != 0 || len(slice.DiffSet(tt.want, got)) != 0 {
				t.Errorf("IntersectSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_UnionSet(t *testing.T) {
	tests := []struct {
		name string
		src  setx.Set[int]
		dst  setx.Set[int]
		want []int
	}{
		{
			name: "TC01 - Some elements in src not in dst",
			src:  setx.NewSet(3, 1, 2, 3),
			dst:  setx.NewSet(3, 2, 3, 4),
			want: []int{1, 2, 3, 4},
		},
		{
			name: "TC02 - Src is empty",
			src:  setx.NewSet[int](0),
			dst:  setx.NewSet(2, 1, 2),
			want: []int{1, 2},
		},
		{
			name: "TC03 - Dst is empty",
			src:  setx.NewSet[int](2, 5, 6),
			dst:  setx.NewSet[int](0),
			want: []int{5, 6},
		},
		{
			name: "TC04 - All elements overlap",
			src:  setx.NewSet[int](2, 7, 8),
			dst:  setx.NewSet[int](2, 7, 8),
			want: []int{7, 8},
		},
		{
			name: "TC05 - No overlap",
			src:  setx.NewSet[int](2, 9, 10),
			dst:  setx.NewSet[int](2, 11, 12),
			want: []int{9, 10, 11, 12},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := setx.UnionSet(tt.src, tt.dst)
			if len(slice.DiffSet(got, tt.want)) != 0 || len(slice.DiffSet(tt.want, got)) != 0 {
				t.Errorf("UnionSet() = %v, want %v", got, tt.want)
			}
		})
	}
}
