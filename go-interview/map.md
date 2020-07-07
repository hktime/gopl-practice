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
哈希函数
* 便于查找
* 尽可能均匀

哈希冲突
* 开放定址法，装载因子
* 链地址法
* 再哈希法
* 建立公共溢出区

数据结构
bucket，oldbucket，正常桶，溢出桶

# Goroutinue调度
参考[goroutine调度](https://tiancaiamao.gitbooks.io/go-internals/content/zh/05.1.html)
Go调度的实现，涉及到几个重要的数据结构。分别是结构体G，结构体M，以及Sched结构体。简称为GMP调度。

