在 Go 中,Slice（切片）是抽象在 Array（数组）之上的特殊类型

# Array数组
[3]int{} 表示 3 个整数的数组

# Slice切片
```
type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}
```
Slice 的底层数据结构共分为三部分,如下：

- array：指向所引用的数组指针（unsafe.Pointer 可以表示任何可寻址的值的指针）
- len：长度,当前引用切片的元素个数
- cap：容量,当前引用切片的容量（底层数组的元素总数）
在实际使用中,cap 一定是大于或等于 len 的.否则会导致 panic

```
nums := [3]int{}
nums[0] = 1
dnums := nums[0:2]
# nums: [1 0 0] 
# dnums: [1 0], len: 2, cap: 3
```
Slice 是对 Array 的抽象,类型为 []T
dnums 变量通过 nums[:] 进行赋值
dnums 变量通过 nums[:] 进行赋值



## Slice 的创建方式：
- test := []int{2,3}
- test := make([]int, 5, 5)  // 创建一个类型为 int，长度为 5，容量为 5 的切片
- test1 := make([]int, 3)                        //如果不指定容量，默认容量等于初始时的长度
- test := make([]int,0)                              // 创建一个长度为0，容量为0 的数组
  test = append(test, 1)
//当数组的容量不够时，会重新申请一个两倍于当前长度的 slice，所以在使用过程中，尤其是频繁去往一个 slice 中 append 数据，需要尽可能给一个相对准确的容量， 减少分配过程的损耗。

