Go面试问的问题也很有限嘛，channel、协程调度模型、切片底层、读写锁，也可能让你介绍一下GC、三色标记算法

1、空struct{}是否使用过？会在什么情况下使用，举例说明一下。 2、在Go语言中，结构体是否能够比较？该如何比较两个结构体？如何比较两个接口？可以顺便查考一下代码实现。 3、使用Go语言编程实现堆栈和队列这两个数据结构，该如何实现。可以只说实现思路。 4、var a []int和a := []int{}是否有区别？如果有的话，具体有什么区别？在开发过程中使用哪个更好，为什么？ 5、Go中，如何复制切片内容？如何复制map内容？如何复制接口内容？编程时会如何操作实现。

什么是goroutine，他与process， thread有什么区别？2. 什么是channel，为什么它可以做到线程安全？3. 了解读写锁吗，原理是什么样的，为什么可以做到？4. 如何用channel实现一个令牌桶？5. 如何调试一个go程序？6. 如何写单元测试和基准测试？7. goroutine 的调度是怎样的？8. golang 的内存回收是如何做到的？9. cap和len分别获取的是什么？10. netgo，cgo有什么区别？11. 什么是interface？

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
# Map
无序的key/value对的集合。

底层引用一个哈希表，用链表来解决冲突。

map中所有的key都有相同的类型，所有的value也有着相同的类型，但是key和value可以是不同的数据类型。key的数据类型必须支持==比较运算符，通过测试key是否相等来判断是否已经存在。

区分map中已经存在的0，和不存在而返回零值的0
```
if age, ok := ages["bob"]; !ok{
    //"bob" is not a key in ages.
}
```
## 哈希表
想要实现一个性能优异的哈希表，需要注意两个关键点 -- 哈希函数和冲突解决方法。

### 哈希函数
理想情况下哈希函数能将不同键映射到不同的索引上，这要求哈希函数输出范围大于输入范围。但在实际情况中，键的数量会远远大于映射的范围，所以在实际使用时，这个理想的结果是不可能实现的。

比较实际的方式是让哈希函数的结果能够尽可能的均匀分布，然后通过工程上的手段解决哈希碰撞的问题。

* 便于查找
* 结果尽可能均匀

### 哈希冲突
由于哈希函数输入范围一定会远远大于输出的范围，所以在使用哈希表时，当输入的键足够多时一定会遇到冲突。这时就需要一些方法来解决哈希碰撞的问题，常见方法就是开放寻址法和链地址法。

1. 开放定址法

核心思想是**对数组中的元素依次探测和比较以判断目标键值对是否存在于哈希表中**，如果使用开放定址法来实现哈希表，那么支撑哈希表的数据结构就是数组，不过由于数组的长度有限，存储(key, value)这个键值对时会从如下的索引开始遍历：

`index := hash("key") % len(array)`

当我们向当前哈希表写入新的数据时发生冲突，就会将键值对写入到下一个不为空的位置。当需要查找某个键对应的值时，就会从索引的位置开始对数组进行线性探测，找到目标键值对或者空内存就意味着这一次查询操作的结束。

开放寻址法对性能影响最大的就是**装载因子**，它是数组中元素的数量与数组大小的比值，随着装载因子的增加，先行探测的平均用时就会逐渐增加，同时影响哈希表的读写性能。

只用数组存储，容易产生堆积问题，不适用于大规模的数据存储；插入时可能会出现多次冲突的现象，删除的元素是多个冲突元素中的一个时，需要对后面的元素做处理，实现较复杂；节点规模很大时会浪费很多空间，对内存的利用率低于链表法。

适合数据量比较小、装载因子小的时候。

2. 链地址法

常用的哈希表的实现方法，实现比开放寻址法稍微复杂一些，但是平均查找的长度也比较短，各个用于存储节点的内存都是动态申请的，可以节省比较多的存储空间。

实现链地址法一般会使用数组加上链表，也就是链表数组作为哈希底层的数据结构。当我们需要存储(key, value)这个键值对时，先经过一个哈希函数，返回的哈希值帮助我们查找在哪个链表中，如下所示：

`index := hash("key") % len(array)`

选择了x号链表之后，就可以遍历当前链表了，在遍历链表的过程中可能遇到这两种情况：
* 找到键相同的键值对 -- 更新键对应的值；
* 没有找到键相同的键值对 -- 在链表的末尾追加新键值对。

当需要通过某个键获取在其中映射的值时，过程也类似，如果遍历到链表的末尾也没有找到期望的值，则哈希表中没有该键对应的值。

内存利用率比开放寻址法更高，对大装载因子的容忍度更高；链表法比较适合存储大对象、大数据量的散列表；而且更加灵活，支持更多的优化策略，比如用红黑树代替链表。

3. 再哈希法
4. 建立公共溢出区

随着键值对数量的增加，哈希的装载因子也会逐渐升高，超过一定范围就会触发扩容，会将桶的数量翻倍，元素再分配的过程也是在调用写操作时增量进行的，不会造成性能的瞬时巨大抖动。

数据结构：
bucket，oldbucket，正常桶，溢出桶


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

# 常用关键字
## defer
参考[5.3 defer](https://draveness.me/golang/docs/part2-foundation/ch05-keyword/golang-defer/)

参考[Golang 中的 Defer 必掌握的 7 知识点](https://learnku.com/articles/42255)

1. defer关键字的调用时机，与return谁先谁后

defer在return之后执行。

2. 多次调用defer时的执行顺序；

多次调用defer时，是一个类似与“栈”的关系，也就是先进后出。先调用的defer会被更后执行。

3. defer关键字使用传值的方式传递参数时

Go语言中所有的函数调用都是传值的，defer虽然是关键字，但是也继承了这个特性。

4. defer遇见panic

能够触发defer的是遇见return（或函数体到末尾）和遇见panic。

遇见panic时，会遍历协程的defer链表，并执行defer。在执行defer的过程中，遇到recover则停止panic，返回recover处继续往下执行。如果没有遇到recover，遍历完本协程的defer链表后，向stderr抛出panic信息。

defer最大的作用就是panic后依然有效，可以保证资源一定会被关闭，从而避免一些异常出现的问题。

5. defer中包含panic

panic中仅有最后一个可以被recover捕获。

6. defer下的函数参数包含子函数

## panic
panic能够改变程序的控制流，函数调用panic时会立刻停止执行函数的其他代码，并在执行结束后在**当前Goroutine**中递归执行调用方的延迟函数调用defer。

向调用者报告错误的一般方式是将`error`作为额外的值返回，但如果错误是不可恢复的，或者有时程序就是不能继续运行的情况。

这时就可以使用内建的`panic`函数，会产生一个运行时错误并终止程序，该函数接受一个任意类型的实参，并在程序终止时打印。

当`panic`被调用后，程序将立刻终止当前函数的执行

## recover
recover可以中止panic造成的程序崩溃，它是一个只能在defer中发挥作用的函数，在其他作用域中调用不会发挥任何作用。

recover会使程序从panic中恢复，并返回panic value。导致panic异常的函数不会继续运行，但能正常返回。在未发生panic时调用recover，recover会返回nil。

## select
一个select语句用来选择哪个case中的发送或接收操作可以被立即执行。类似于switch语句，但是select的case涉及到与channel有关的I/O操作

# 接口
接口的值由两个部分组成，一个具体的类型和那个类型的值。被称为接口的动态类型和动态值。
在Go语言中，变量总是被一个定义明确的值初始化，即使接口类型也不例外。对于一个接口的零值,就是它的类型和值的部分都是nil。

## 接口的比较
接口值可以使用==或!=来进行比较。两个接口值相等仅当它们都是nil值或者它们的动态类型相同并且动态值也根据这个动态类型的==操作相等。因为接口值是可比较的，所以它们可以用做map的键或者switch语句的操作数。

然而，如果两个接口值的动态类型相同，但是该类型是不可比较的（比如切片），将他们进行比较就会失败并且panic。

# 变量
## 变量的声明
变量声明的一般语法如下：

```var 变量名字 类型 = 表达式```

其中“类型”或“= 表达式”两个部分可以省略其中的一个，如果省略的是类型信息，那么将根据初始化表达式来推导变量的类型信息。如果初始化表达式被省略，那么将用零值初始化该变量。

简短变量声明如下：

```名字 := 表达式```
被用于大部分局部变量的声明和初始化。

## var a []int和a := []int{}是否有区别
var形式的声明语句往往是用于需要显式指定变量类型的地方，或者因为变量稍后会被重新赋值而初始值无关紧要的地方。

简短声明的限制是不能提供数据类型，而且只能用在函数内部。

## make和new
1. make的作用是初始化内置的数据结构；
2. new的作用是根据传入的类型分配一片内存空间并返回指向这片内存空间的指针。


# 内存分配
编程语言的内存分配器一般包含两种分配方法，一种是线性分配器，另一种是空闲链表分配器。

## 线性分配器
一种高效的内存分配方式，但有较大的局限性。只需要在内存中维护一个指向内存特定位置的指针，当用户程序申请内存时，分配器只需要检查剩余的空闲内存、返回分配的内存区域，并修改指针在内存中的位置。

具有较快的执行速度，以及较低的实现复杂度，但是线性分配器无法在内存被释放时重用内存。需要配合具有拷贝特性的垃圾回收算法。

## 空闲链表分配器
可以重用已经被释放的内存，在内部会维护一个类似链表的数据结构，当用户程序申请内存时，空闲链表分配器会依次遍历空闲的内存块，找到足够大的内存，然后申请新的资源并修改链表。

## 分级分配
go语言的内存分配核心理念是利用多级缓存将对象根据大小分配，并按照类别实施不同的分配策略。
### 对象大小
类别|大小|
---|---|
微对象|0-16B|
小对象|16B-32KB|
大对象|32KB-+∞|

### 虚拟内存布局
1.10以前的版本，堆区的内存空间都是连续的，但是在1.11版本，使用稀疏的堆内存空间替代了连续的内存，解决了连续内存带来的限制以及在特殊场景下可能出现的问题。

# 垃圾回收
三色标记算法，将程序中的对象分为白色、黑色和灰色三类：
* 白色对象 - 潜在的垃圾，其内存可能会被垃圾收集器回收；
* 黑色对象 - 活跃的对象，包括不存在任何引用外部指针的对象以及从根对象可达的对象；
* 灰色对象 - 活跃的对象，存在指向白色对象的外部指针，垃圾收集器会扫描这些对象的子对象。

工作原理：
1. 从灰色对象的集合中选择一个灰色对象并将其标记成黑色；
2. 将黑色对象指向的所有对象都标记成灰色，保证该对象和被该对象引用的对象都不会被回收；
3. 重复上述步骤直到对象图中不存在灰色对象，只剩下黑色的存活对象和白色的垃圾对象，垃圾回收器可以回收白色的垃圾。

# Goroutinue 和 channel
Goroutinue 和 channel是go语言并发编程的两大基石，goroutine用于执行并发任务，channel用于goroutine之间的同步和通信。

## goroutine调度
参考[goroutine调度](https://tiancaiamao.gitbooks.io/go-internals/content/zh/05.1.html)

Go调度的实现，涉及到几个重要的数据结构。分别是结构体G，结构体M，结构体P，以及Sched结构体。简称为GMP调度。

### 结构体G
G是goroutine的缩写，相当于操作系统中的进程控制块，是对goroutine的抽象。其中包含了栈信息和运行的函数，只要得到CPU就可以运行。

goroutine切换时，上下文信息保存在结构体的sched域中，切换时不必陷入到操作系统内核中，所以保存过程很轻量。结构体中的Gobuf只保存了当前栈指针，程序计数器以及goroutine自身。

### 结构体M
M是machine的缩写，是对机器的抽象，每个m都是对应到一条操作系统的物理线程。M必须关联了P才可以执行Go代码，但是当它处理阻塞或者系统调用时，可以不需要关联P。

### 结构体P
P是processor的缩写，代表go代码执行时需要的资源。当M执行go代码时，需要关联一个P。有刚好GOMAXPROCS个P。

### 调度过程
首先创建一个G对象，G对象保存到P本地队列或者是全局队列。P此时去唤醒一个M，接下来M执行一个调度循环（调用G对象->执行->清理线程->继续找新的goroutine执行）。

## channel
channel是线程安全的，是“先进先出的”。通过加锁的方式来实现线程安全。

不要通过共享内存来通信，而要通过通信来实现内存共享。

### 单向或双向
```
chan T //声明一个双向channel
chan<- T //声明一个只能用来发送的channel
<-chan T // 声明一个只能用来接收的channel
```

### 无缓冲和有缓冲
```
ch := make(chan bool) //无缓冲的channel
ch := make(chan bool, 1) // 带缓冲的channel，且缓冲区大小为1
```
对不带缓冲的channel进行的操作实际上可以看作“同步模式”，发送方和接收方要同步就绪，只有在两者都ready的情况下，数据才能在两者间传输。否则，任意一方先行进行发送或接收操作，都会被挂起，等待另一方的出现才能被唤醒。

带缓冲的channel则称为“异步模式”，在缓冲区可用的情况下，发送和接收操作都可以顺利进行。否则，操作的一方同样会被挂起，直到出现相反操作才会被唤醒。


# go的编译执行流程
在golang中，build过程主要由go build执行，完成了源码的编译与可执行文件的生成。

go build接收参数为.go文件或目录，默认情况下编译当前目录下所有.go文件。在main包下执行会生成相应的可执行文件，在非main包下，会做一些检查，将生成的库文件放在缓存目录下，在工作目录下并无新文件生成。

```
go help build
-n 不执行地打印流程中用到的命令
-x 执行并打印流程中用到的命令，要注意下它与-n选项的区别
-work 打印编译时的临时目录路径，并在结束时保留。默认情况下，编译结束会删除该临时目录。
```

## 打印执行流程
使用-n选项在命令不执行的情况下，查看go build的执行流程，如下：
```
go build -n main.go

mkdir -p $WORK\b001\
cat >$WORK\b001\importcfg << 'EOF' # internal        
# import config
packagefile fmt=D:\Go\pkg\windows_amd64\fmt.a        
packagefile runtime=D:\Go\pkg\windows_amd64\runtime.a
EOF
cd D:\GoWorkspace
"D:\\Go\\pkg\\tool\\windows_amd64\\compile.exe" -o "$WORK\\b001\\_pkg_.a" -trimpath "$WORK\\b001=>" -p main -complete -buildid j0ngvybnShLGnMrOSUpT/j0ngvybnShLGnMrOSUpT -dwarf=false -goversion go1.14.1 -D _/D_/GoWorkspace -importcfg "$WORK\\b001\\importcfg" -pack -c=4 "D:\\GoWorkspace\\main.go"
"D:\\Go\\pkg\\tool\\windows_amd64\\buildid.exe" -w "$WORK\\b001\\_pkg_.a" # internal
cat >$WORK\b001\importcfg.link << 'EOF' # internal
packagefile command-line-arguments=$WORK\b001\_pkg_.a

...

packagefile internal/syscall/windows/registry=D:\Go\pkg\windows_amd64\internal\syscall\windows\registry.a
EOF
mkdir -p $WORK\b001\exe\
cd .
"D:\\Go\\pkg\\tool\\windows_amd64\\link.exe" -o "$WORK\\b001\\exe\\main.exe" -importcfg "$WORK\\b001\\importcfg.link" -s -w -buildmode=exe -buildid=BEvhBXJqCS0TDjrlOn42/j0ngvybnShLGnMrOSUpT/j0ngvybnShLGnMrOSUpT/BEvhBXJqCS0TDjrlOn42 -extld=gcc "$WORK\\b001\\_pkg_.a"
"D:\\Go\\pkg\\tool\\windows_amd64\\buildid.exe" -w "$WORK\\b001\\exe\\a.out.exe" # internal
mv $WORK\b001\exe\a.out.exe main.exe
```

主要由几部分组成，分别是：
* 创建临时目录，mkdir -p $WORK\b001\
* 查找依赖信息，cat >$WORK\b001\importcfg
* 执行源代码编译，"D:\\Go\\pkg\\tool\\windows_amd64\\compile.exe" -o
* 收集链接库文件，cat >$WORK\b001\importcfg.link
* 生成可执行文件，"D:\\Go\\pkg\\tool\\windows_amd64\\link.exe" -o
* 移动可执行文件，mv $WORK\b001\exe\a.out.exe main.exe


如何读输入

inputReader := bufio.NewReader(os.Stdin)
input, _ := inputReader.ReadString('\n')

如何对二维数组进行排序
```
var intervals [][]int
sort.Slice(intervals, func(i, j int)bool{
    return intervals[i][0] < intervals[j][0]
})
```


