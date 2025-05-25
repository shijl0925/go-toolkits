# go-toolkits

[![DeepSource](https://app.deepsource.com/gh/shijl0925/go-toolkits.svg/?label=active+issues&show_trend=true&token=4DmceHPfSvsFfGwQ3g4OPPvT)](https://app.deepsource.com/gh/shijl0925/go-toolkits/)

go-toolkits 是一个基于 Go 泛型的工具库, 它涵盖了在开发过程中经常会用到的一些辅助方法。


## 安装
ekit 目前要求 Go >= 1.23。

使用 go get 命令安装：
```shell
go get github.com/shijl0925/go-toolkits@latest
```

## Documentation

#### Slice
[Slice](https://pkg.go.dev/github.com/shijl0925/go-toolkits/slicex)

#### Function list:
- **<big>Sum</big>**:
    * Sum returns the sum of all elements in the slice.
- **<big>Avg</big>**:
- **<big>Max</big>**: 
    * returns the maximum value in the slice.
- **<big>Min</big>**:
    * returns the minimum value in the slice.
- **<big>RemoveHead</big>**: 
    * removes the first element from a slice and returns the remaining elements.
    * 返回一个新的切片，其中去除了原切片的第一个元素, 原始 slice 不会被改变。
- **<big>RemoveTail</big>**: 
    * removes the last element from a slice and returns the remaining elements.
    * 返回一个新的切片, 其中去除了原切片的最后一个元素, 原始 slice 不会被改变。
- **<big>Pop</big>**: 
    * removes and returns the last element from a slice.
    * 从切片 s 中弹出最后一个元素，并返回该元素和弹出状态。原始 slice 会被改变。
- **<big>Drop</big>**: 
    * removes the n elements from the slice and returns the remaining elements.
    * 删除并返回切片的前 n 个元素和剩余元素组成的切片, 原始 slice 不会被改变。
- **<big>DropLeft</big>**:
    * DropLeft removes the n elements from the slice and returns the remaining elements.
    * 删除并返回切片的左侧 n 个元素和剩余元素组成的切片, 原始 slice 不会被改变。
- **<big>Remove</big>**: 
    * returns a slice with the all occurrence of element removed.
    * 返回一个切片，该切片移除了所有等于 element 的元素。该函数不会修改原始切片内容，而是返回一个新的切片。
- **<big>Delete</big>**: 
    * removes and returns the element at a specific index from a slice.
    * 删除位于索引 index 的元素, 原始 slice 会被改变。
- **<big>DeleteAt</big>**: 
    * removes and returns the element at a specific index from a slice.
    * 删除位于索引 index 的元素, 原始 slice 会被改变。(性能比Delete好)
- **<big>DeleteMany</big>**:
    * removes and returns the elements between start and end from a slice.
    * 删除位于索引 start 和 end 之间的元素, 原始 slice 会被改变。
- **<big>Append</big>**:
    * appends one or more elements to a slice.
    * 追加一个或多个元素到 slice。
- **<big>Extend</big>**:
    * extends a slice by appending one or more elements.
    * 将切片 elements 中的所有元素追加到 src 切片后并返回新切片。
- **<big>InsertAtV2</big>**:
    * inserts an element at a specific index in a slice.
    * 在 slice 的指定位置插入一个元素, 原始 slice 不会被改变。
- **<big>InsertMany</big>**:
    * inserts multiple elements at a specific index in a slice.
    * 在 slice 的指定位置插入新切片, 原始 slice 不会被改变。
- **<big>AddMany</big>**:
    * adds multiple elements to a slice at a specific index.
    * 在 slice 的指定位置添加新切片, 原始 slice 不会被改变。
- **<big>Reverse</big>**:
    * returns a new slice with the elements in reverse order.
    * 返回一个新的 slice，元素顺序相反, 原始 slice 不会被改变。
- **<big>ReverseSelf</big>**:
    * reverses the elements in a slice in place.
    * 反转 slice 本身
- **<big>Repeat</big>**:
    * returns a new slice with n copies of the given element.
    * 返回一个新 slice，其中包含 n 个给定元素的副本, 原始 slice 不会被改变。
- **<big>Product</big>**:
    * returns a new slice of tuples, where each tuple contains the corresponding elements from two slices.
- **<big>Contains</big>**:
    * 判断 src 里面是否存在 dst
- **<big>ContainsAny</big>**:
    * 判断 src 里面是否存在 dst 中的任何一个元素
- **<big>ContainsAll</big>**:
    * 判断 src 里面是否存在 dst 中的所有元素
- **<big>Unique</big>**:
    * removes duplicate elements from a slice.
    * 删除 slice 中的重复元素, 原始 slice 不会被改变 。
- **<big>Difference</big>**:
    * returns the difference between two slices.
    * 返回 src 和 dst 的差集, 即 src 中存在，dst 中不存在的元素。
- **<big>Intersect</big>**:
    * returns the intersection between two slices.
    * 返回 src 和 dst 的交集, 即 src 中存在，dst 中也存在元素。
- **<big>Union</big>**:
    * returns the union of two slices.
    * 返回 src 和 dst 的并集, 即 src 和 dst中所有元素的集合。
- **<big>Index</big>**: 
    * 返回和 dst 相等的第一个元素下标，-1 表示没找到
- **<big>LastIndex</big>**:
    * 返回和 dst 相等的最后一个元素下标，-1 表示没找到
- **<big>IndexAll</big>**:
    * 返回和 dst 相等的所有元素的下标
- **<big>JoinSlice</big>**:
    * join []any slice to string.
- **<big>Chunk</big>**:
    * creates a slice of elements split into groups the length of size.
    * 将切片按照给定的大小进行分组，并返回一个二维切片。
- **<big>Compact</big>**:
    * returns a new slice with all falsey values removed. The values false, nil, 0, and "" are falsey.
    * 返回一个新的 slice，其中包含原始 slice 中所有非零值的元素。注意：本函数会改变原始 slice。
- **<big>Concat</big>**:
    * returns a new slice with all slices concatenated.
- **<big>Equal</big>**:
    * returns true if the two slices are equal.
    * 判断两个 slice 是否相等。
- **<big>EqualUnordered</big>**:
    * checks if two slices are equal: the same length and all elements' value are equal (unordered).
    * 判断两个 slice 是否相等：长度相同且所有元素值相等（不考虑顺序）。
- **<big>Counter</big>**:
    * returns a map of the counts of each value in the slice.
    * 返回一个 map，其中键为 slice 中的每个值，值为该值在 slice 中出现的次数。
- **<big>Replace</big>**:
    * returns a new slice with all occurrences of old replaced by new.
    * 返回一个新的 slice，其中包含原始 slice 中所有 old 值被 new 替换后的元素。n 表示替换的次数，-1 表示替换所有。
- **<big>IsAscending</big>**:
    * checks if a slice is ascending order.
    * 判断一个 slice 是否是升序排列的。
- **<big>IsDescending</big>**:
    * checks if a slice is descending order.
    * 判断一个 slice 是否是降序排列的。
- **<big>RightPadding</big>**:
    * adds padding to the right end of a slice.
- **<big>LeftPadding</big>**:
    * adds padding to the left begin of a slice.
- **<big>FlattenDeep</big>**:
    * flattens slice recursive.
    * 递归地展平切片
- **<big>Combine</big>**:
  * combines two slices into a map.
  * 将两个切片组合成一个 map
- **<big>SortSlice</big>**:
  * 对任意类型的切片进行排序，使用提供的 less 函数定义排序规则。









