package setx_test

import (
	"github.com/shijl0925/go-toolkits/setx"
	"reflect"
	"testing"
)

// TestNewSet tests the New function.
func TestNewSet(t *testing.T) {
	t.Run("empty set", func(t *testing.T) {
		s := setx.New[int]()
		if len(s) != 0 {
			t.Errorf("expected empty set, got length %d", len(s))
		}
	})

	t.Run("set with initial elements", func(t *testing.T) {
		s := setx.New(3, 1, 2, 3)
		expected := setx.Set[int]{1: struct{}{}, 2: struct{}{}, 3: struct{}{}}
		if !reflect.DeepEqual(s, expected) {
			t.Errorf("expected %v, got %v", expected, s)
		}
	})
}

// TestAdd tests the Add method.
func TestAdd(t *testing.T) {
	s := setx.New[int]()

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
	s := setx.New([]int{1, 2}...)

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
	s := setx.New[int]()
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
	s := setx.New(1)

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
			src:  setx.New(1, 2, 3),
			want: []int{3, 2, 1},
		},
		{
			name: "TC02 - Src is empty",
			src:  setx.New[int](),
			want: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.src.Keys()
			if !setx.NewFromSlice(got).Equal(setx.NewFromSlice(tt.want)) {
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
			src:  setx.New(1, 2, 3),
			dst:  setx.New(2, 3, 4),
			want: []int{1},
		},
		{
			name: "TC02 - Src is empty",
			src:  setx.New[int](),
			dst:  setx.New(1, 2),
			want: []int{},
		},
		{
			name: "TC03 - Dst is empty",
			src:  setx.New[int](5, 6),
			dst:  setx.New[int](),
			want: []int{5, 6},
		},
		{
			name: "TC04 - All elements overlap",
			src:  setx.New[int](7, 8),
			dst:  setx.New[int](7, 8),
			want: []int{},
		},
		{
			name: "TC05 - No overlap",
			src:  setx.New[int](9, 10),
			dst:  setx.New[int](11, 12),
			want: []int{9, 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := setx.Difference(tt.src, tt.dst)
			if !setx.NewFromSlice(got).Equal(setx.NewFromSlice(tt.want)) {
				t.Errorf("Difference() = %v, want %v", got, tt.want)
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
			src:  setx.New(1, 2, 3),
			dst:  setx.New(2, 3, 4),
			want: []int{2, 3},
		},
		{
			name: "TC02 - Src is empty",
			src:  setx.New[int](),
			dst:  setx.New(1, 2),
			want: []int{},
		},
		{
			name: "TC03 - Dst is empty",
			src:  setx.New[int](5, 6),
			dst:  setx.New[int](),
			want: []int{},
		},
		{
			name: "TC04 - All elements overlap",
			src:  setx.New[int](7, 8),
			dst:  setx.New[int](7, 8),
			want: []int{7, 8},
		},
		{
			name: "TC05 - No overlap",
			src:  setx.New[int](9, 10),
			dst:  setx.New[int](11, 12),
			want: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := setx.Intersect(tt.src, tt.dst)
			if !setx.NewFromSlice(got).Equal(setx.NewFromSlice(tt.want)) {
				t.Errorf("Intersect() = %v, want %v", got, tt.want)
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
			src:  setx.New(1, 2, 3),
			dst:  setx.New(2, 3, 4),
			want: []int{1, 2, 3, 4},
		},
		{
			name: "TC02 - Src is empty",
			src:  setx.New[int](),
			dst:  setx.New(1, 2),
			want: []int{1, 2},
		},
		{
			name: "TC03 - Dst is empty",
			src:  setx.New[int](5, 6),
			dst:  setx.New[int](),
			want: []int{5, 6},
		},
		{
			name: "TC04 - All elements overlap",
			src:  setx.New[int](7, 8),
			dst:  setx.New[int](7, 8),
			want: []int{7, 8},
		},
		{
			name: "TC05 - No overlap",
			src:  setx.New[int](9, 10),
			dst:  setx.New[int](11, 12),
			want: []int{9, 10, 11, 12},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := setx.Union(tt.src, tt.dst)
			if !setx.NewFromSlice(got).Equal(setx.NewFromSlice(tt.want)) {
				t.Errorf("Union() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestClear_EmptySet 测试清空一个空集合
func TestClear_EmptySet(t *testing.T) {
	set := make(setx.Set[int])
	set.Clear()

	if len(set) != 0 {
		t.Errorf("Expected empty set after Clear(), got %d elements", len(set))
	}
}

// TestClear_NonEmptySet 测试清空一个包含元素的集合
func TestClear_NonEmptySet(t *testing.T) {
	set := make(setx.Set[string])
	set.Add("a")
	set.Add("b")

	set.Clear()

	if len(set) != 0 {
		t.Errorf("Expected empty set after Clear(), got %d elements", len(set))
	}
}

// TestClear_MultipleElements 测试多个不同类型的元素是否都能被清除
func TestClear_MultipleElements(t *testing.T) {
	set := make(setx.Set[interface{}])
	set.Add(1)
	set.Add("hello")
	set.Add(true)

	set.Clear()

	if len(set) != 0 {
		t.Errorf("Expected empty set after Clear(), got %d elements", len(set))
	}
}

// TestSet_Equal 测试 Equal 方法的各种情况
func TestSet_Equal(t *testing.T) {
	tests := []struct {
		name string
		s    setx.Set[int]
		dst  setx.Set[int]
		want bool
	}{
		{
			name: "TC01 - Both empty sets",
			s:    setx.Set[int]{},
			dst:  setx.Set[int]{},
			want: true,
		},
		{
			name: "TC02 - Sets with same elements",
			s:    setx.Set[int]{1: {}, 2: {}, 3: {}},
			dst:  setx.Set[int]{1: {}, 2: {}, 3: {}},
			want: true,
		},
		{
			name: "TC03 - Different lengths",
			s:    setx.Set[int]{1: {}, 2: {}},
			dst:  setx.Set[int]{1: {}, 2: {}, 3: {}},
			want: false,
		},
		{
			name: "TC04 - s has extra element",
			s:    setx.Set[int]{1: {}, 2: {}, 4: {}},
			dst:  setx.Set[int]{1: {}, 2: {}, 3: {}},
			want: false,
		},
		{
			name: "TC05 - dst has extra element",
			s:    setx.Set[int]{1: {}, 2: {}, 3: {}, 4: {}},
			dst:  setx.Set[int]{1: {}, 2: {}, 3: {}},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.Equal(tt.dst)
			if got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 辅助函数：比较两个 Set 是否相等
func assertSetEqual[T comparable](t *testing.T, expected, actual setx.Set[T]) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected set %v but got %v", expected, actual)
	}
}

// TC01: s 为空，dst 包含多个元素
func TestUpdate_EmptySource(t *testing.T) {
	dst := setx.New[int]([]int{1, 2, 3}...)
	s := make(setx.Set[int])

	expected := setx.New[int]([]int{1, 2, 3}...)
	result := s.Update(dst)

	assertSetEqual(t, expected, result)
}

// TC02: s 已有部分元素，dst 包含新旧混合元素
func TestUpdate_OverlapElements(t *testing.T) {
	s := setx.New[int]([]int{1, 2}...)
	dst := setx.New[int]([]int{2, 3, 4}...)

	expected := setx.New[int]([]int{1, 2, 3, 4}...)
	result := s.Update(dst)

	assertSetEqual(t, expected, result)
}

// TC03: dst 为空，s 不变
func TestUpdate_EmptyDestination(t *testing.T) {
	s := setx.New[string]([]string{"a", "b"}...)
	dst := make(setx.Set[string])

	expected := s // 原始集合不变
	result := s.Update(dst)

	assertSetEqual(t, expected, result)
}

// TC04: dst 中所有元素都已存在于 s 中
func TestUpdate_AllElementsExist(t *testing.T) {
	s := setx.New[int]([]int{10, 20, 30}...)
	dst := setx.New[int]([]int{10, 20}...)

	expected := s // 不应改变
	result := s.Update(dst)

	assertSetEqual(t, expected, result)
}
