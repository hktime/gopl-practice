Go面试问的问题也很有限嘛，channel、协程调度模型、切片底层、读写锁，也可能让你介绍一下GC、三色标记算法

11、go defer 的执行顺序。

12、go interface.

1、空struct{}是否使用过？会在什么情况下使用，举例说明一下。 2、在Go语言中，结构体是否能够比较？该如何比较两个结构体？如何比较两个接口？可以顺便查考一下代码实现。 3、使用Go语言编程实现堆栈和队列这两个数据结构，该如何实现。可以只说实现思路。 4、var a []int和a := []int{}是否有区别？如果有的话，具体有什么区别？在开发过程中使用哪个更好，为什么？ 5、Go中，如何复制切片内容？如何复制map内容？如何复制接口内容？编程时会如何操作实现。

什么是goroutine，他与process， thread有什么区别？2. 什么是channel，为什么它可以做到线程安全？3. 了解读写锁吗，原理是什么样的，为什么可以做到？4. 如何用channel实现一个令牌桶？5. 如何调试一个go程序？6. 如何写单元测试和基准测试？7. goroutine 的调度是怎样的？8. golang 的内存回收是如何做到的？9. cap和len分别获取的是什么？10. netgo，cgo有什么区别？11. 什么是interface？

作者：「已注销」
链接：https://www.zhihu.com/question/60952598/answer/456838652
来源：知乎
著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。

# Array和slice 或者 数组与切片 的联系与区别

## 数组
数组是一个由**固定**长度的特定类型元素组成的序列。
* 数组是不可变的。
* 数组的长度是**数组类型**的一个组成部分，因此[3]int和[4]int是两种不同的数组类型。
* 数组的长度必须是**常量**表达式，因为数组的长度需要在编译阶段确定。

## 切片
切片代表变长的序列，序列中每个元素都有相同的类型，写作[]T。
* 切片是可变的，没有固定长度。
* slice由三部分构成：指针、长度和容量，指针指向第一个slice元素对应的底层数组元素的地址。
* 和数组不同的是，slice之间不能比较，不能用==操作符判断两个slice是否含有全部相等元素。

切片不支持==操作的原因：
1. 一个slice的元素是间接引用的，一个slice甚至可以包含自身。缺乏一个简单有效的方法处理这种情形。举例：s := []interface{}{"one", nil}; s[1] = s;
2. 因为slice的元素是间接引用的，一个固定的slice值在不同的时刻可能包含不同的元素，因为底层数组的元素可能会被修改。安全起见。slice可以与nil比较。

### append操作
先判断slice底层数组是否有足够的容量来保存新添加的元素。
1. 如果有，直接拓展slice（依然在原有的底层数组之上），将新元素添加到新拓展的空间，并返回slice。输入的x和输出的z共享相同的底层数组。

2. 如果没有，先**分配**一个足够大的slice用于保存新的结果，先将输入的x**复制**到新的空间，然后添加y元素。结果z和输入的x引用的是不同的底层数组。

拓展数组时长度直接翻倍。

### 比较两个slice是否相等
参考[比较两个slice是否相等](https://www.jianshu.com/p/80f5f5173fca)
参考[golang中判断两个slice是否相等](https://www.cnblogs.com/apocelipes/p/11116725.html)

* 判断两个[]byte是否相等
使用系统函数bytes.Equal()
```
return bytes.Equal(a, b)
```
* reflect
使用reflect.DeepEqual()函数：
```
func StringSliceReflectEqual(a, b []string) bool {
    return reflect.DeepEqual(a, b)
}
```
* 循环遍历比较
手动实现，比reflect快
```
func StringSliceEqual(a, b []string) bool {
    if len(a) != len(b) {
        return false
    }

    if (a == nil) != (b == nil) {
        return false
    }

    for i, v := range a {
        if v != b[i] {
            return false
        }
    }

    return true
}
```

# 结构体 struct

## 空结构体的作用
参考[空结构体struct{}解析](https://cloud.tencent.com/developer/article/1068400)
空结构体写作struct{}，大小为0，也不包含任何信息，有时候依然是有价值的。在用map来模拟set数据结构时，用它来代替map中布尔类型的value，只是强调key的重要性。但是因为节约的空间有限，而且语法比较复杂，所以通常避免这样的用法。
```
seen := make(map[string]struct{})
if _, ok := seen[s]; !ok{
    seen[s] = struct{}{}
}
```
## 结构体的比较
如果结构体的全部成员都是可以比较的，那么结构体也是可以比较的，可以使用==或!=运算符进行比较。相等比较运算符==将比较两个结构体的每个成员。
```
type Point struct{ X, Y int }

p := Point{1, 2}
q := Point{2, 1}
// 两种比较方式是等价的：
fmt.Println(p.X == q.X && p.Y == q.Y) // "false"
fmt.Println(p == q)                   // "false"
```

# 接口
接口的值由两个部分组成，一个具体的类型和那个类型的值。被称为接口的动态类型和动态值。
在Go语言中，变量总是被一个定义明确的值初始化，即使接口类型也不例外。对于一个接口的零值,就是它的类型和值的部分都是nil。

## 接口的比较
接口值可以使用==或!=来进行比较。两个接口值相等仅当它们都是nil值或者它们的动态类型相同并且动态值也根据这个动态类型的==操作相等。因为接口值是可比较的，所以它们可以用做map的键或者switch语句的操作数。

然而，如果两个接口值的动态类型相同，但是该类型是不可比较的（比如切片），将他们进行比较就会失败并且panic。

# 变量
## var a []int和a := []int{}是否有区别

