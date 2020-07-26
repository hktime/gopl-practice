## 一面
### 项目经历
### 操作系统
1. 进程与线程的区别
2. 线程独有的资源，线程间通信
3. 操作系统内存管理，段页式内存管理，缺页中断
### 计算机网络
1. TCP和UDP区别
2. TCP包最大大小，切包
3. 三次握手，为什么三次
### 数据结构和算法
1. 常用数据结构
2. 堆和栈的区别
3. map 哈希表， 键值冲突，解决方法优缺点
4. 树，b+树和b树，二叉查找树和二叉平衡查找树
### 多态和继承
1. 协程与线程的区别
2. GMP调度器
3. 两个线程、四个协程，协程能否并发

## 一面
### 项目相关
1. 云原生cicd
2. 具体工作，在cicd迁移过程中。

### 网络相关
1. 浏览器输入https链接到显示页面全过程；
2. ssl认证过程，对称加密和非对称加密；
3. timewait状态；
4. socket无状态，如何保证信息完整性，客户端与服务器通信；
5. 抓包工具，https，中间人攻击

### 操作系统相关
1. 进程与线程区别
2. fork()，父进程和子进程，在内存上的区别
3. 进程间通信，信号量，与互斥锁、条件变量的区别，使用场景
4. I/O模型


### 算法
1. 翻转链表
2. 最长不重复连续子串

### 测试相关
1. 对如流客户端进行测试

## 二面
### 自我介绍
### 项目经验
1. 传统CICD和云原生CICD区别

### Python
1. 常用python标准库
os，与操作系统相关联的函数；
shutil，针对日常的文件和目录管理任务；
sys，命令行参数；
re，正则匹配；
math，数学运算相关；
random，生成随机数的工具；
urllib，访问互联网；
unittest，测试。

2. python常用数据结构，列表和元组的区别
列表、元组、字典、集合，列表与元组类似，不同之处在于元组的元素不能修改。

### 算法
1. 二分查找
循环解法：
```
func binarySearch(nums []int, target int) int {
    left, right := 0, len(nums)-1
    for left <= right {
        mid := left + (right-left)/2
        if nums[mid] == target {
            return mid
        } else if nums[mid] > target {
            right = mid - 1
        } else {
            left = mid + 1
        }
    }
    return -1
}
```

## 查找左边界
```
func binarySearchLeftBorder(nums []int, target int) int {
    left, right := 0, len(nums)
    for left < right {
        mid := left + (right-left)/2
        if nums[mid] == target {
            right = mid
        } else if nums[mid] > target {
            right = mid
        } else {
            left = mid + 1
        }
    }
    if left >= len(nums) {
        return -1
    }
    if nums[left] != target {
        return -1
    }
    return left
}
```

## 查找右边界
```
func binarySearchRightBorder(nums []int, target int) int {
    left, right := 0, len(nums)
    for left < right {
        mid := left + (right-left)/2
        if nums[mid] == target {
            left = mid + 1
        } else if nums[mid] > target {
            right = mid
        } else {
            left = mid + 1
        }
    }
    if left == 0 {
        return -1
    }
    if nums[left-1] != target {
        return -1
    }
    return left - 1
}
```

### 操作系统
1. python多线程
threading模块

2. 进程与线程的区别

### 区别
线程依赖于进程而存在，一个进程至少有一个线程。
1. 拥有资源
* 进程有自己的独立地址空间，线程共享所属进程的地址空间；
* 进程是拥有系统资源的一个独立单位，而线程自己基本上不拥有系统资源，
* 只拥有一点在运行中必不可少的资源(如程序计数器,一组寄存器和栈)，和其他线程共享本进程的相关资源如内存、I/O、cpu等；

2. 调度
* 在进程切换时，涉及到整个当前进程CPU环境的保存环境的设置以及新被调度运行的CPU环境的设置；
* 而线程切换只需保存和设置少量的寄存器的内容，并不涉及存储器管理方面的操作，可见，进程切换的开销远大于线程切换的开销；
* 线程是独立调度的基本单位，在同一进程中，线程的切换不会引起进程切换，而不同进程中的线程切换，则会引起进程的切换。

3. 系统开销
由于创建或撤销进程时，系统都要为他分配或回收资源，所付出的开销远大于创建或撤销线程时的开销。线程切换时只需保存或设置少量寄存器的内容，开销远小于进程切换。

4. 通信
* 线程之间的通信更方便，同一进程下的线程共享全局变量等数据，
* 而进程之间的通信需要以进程间通信(IPC)的方式进行；

5. 健壮性
* 多线程程序只要有一个线程崩溃，整个程序就崩溃了，
* 但多进程程序中一个进程崩溃并不会对其它进程造成影响，因为进程有自己的独立地址空间，因此多进程更加健壮。

### 同一进程中的线程可以共享哪些数据？
堆、地址空间、全局变量等。

### 线程独占哪些资源？
栈和程序计数器，用来保存线程的执行历史和状态。

3. 进程间的通信方式
进程通信是一种手段，而进程同步是一种目的。可以通过进程通信的方法来达到进程同步的目的。

1. 管道
通过调用pipe函数创建的。具有以下限制：
* 只支持半双工通信（单向交替传输）；
* 只能在父子进程或者兄弟进程中使用。

2. 命名管道（FIFO）
去除了管道只能用在父子进程中使用的限制。

3. 消息队列
可以独立于读写进程存在，避免了FIFO的同步阻塞问题，不需要进程自己提供同步方法。

读进程可以根据消息类型有选择地接受信息，而不像FIFO只能默认的接收。

克服了信号承载信息量少，管道只能承载无格式字节流以及缓冲区大小受限等缺点。

4. 信号量
计数器，用于为多个进程提供对共享数据对象的访问。

5. 共享内存
允许多个进程共享一个给定的内存，数据不需要在进程之间复制，所以是最快的一种IPC。

6. 套接字
用于不同机器间的进程通信。

4. Linux常用命令
ping, kill, ps, netstat, df, top, chmod, cd, mkdir, cat, rm -rf.

5. 服务器上某些端口是否被占用
参考[linux下查看端口占用](https://www.cnblogs.com/zjfjava/p/10513399.html)
lsof -i:端口号，lsof（list open files）；
netstat -tunlp | grep 端口号，用于查看指定端口号的进程情况

6. 批量杀死僵尸进程
``` ps -A -ostat,ppid,pid,cmd |grep -e '^[Zz]' | grep -v grep | cut -c 5-10 | xargs kill -9```

### 安全相关
1. TCP flood 
TCP flood攻击是利用了TCP三次握手的设定，伪造大量的IP地址给服务器发送SYN报文，服务器会维护一个庞大的等待列表，同时占用着大量的资源无法释放。
2. UDP flood
UDP flood是利用大量UDP小包冲击DNS服务器或认证服务器，通过将线路上的骨干设备打瘫，可以造成整个网段的瘫痪。正常情况下，当服务器在特定端口接收到UDP数据包时，会经过两个步骤：服务器首先检查是否有正在运行，侦听指定端口请求的程序；如果没有则使用ICMP数据包进行响应，以通知发送方目的地不可达。