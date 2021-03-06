# Linux下的I/O模型

## 概念说明
Linux下的I/O模型的概念说明
1. 用户空间和内核空间
2. 进程切换
3. 进程的阻塞
4. 文件描述符
5. 缓存I/O

### 用户空间和内核空间
对于32位的操作系统而言，寻址空间位4G。操作系统的核心是内核，独立于普通的应用程序，可以访问受保护的内存空间，也有访问底层硬件设备的所有权限。

为了保护内核的安全，保证用户进程不能直接操作内核，操作系统将虚拟空间划分为两部分，一部分为内核空间，一部分为用户空间。针对Linux操作系统而言，将最高的1G字节，供内核使用，称为内核空间，而将较低的3G字节，供各个进程使用，称为用户空间。

### 进程切换
为了控制进程的执行，内核必须有能力挂起正在CPU上运行的进程，并恢复以前挂起的某个进程的执行。这种行为被称为进程切换。

从一个进程的运行转到另一个进程上运行，经过这些变化：
1. 保存处理机上下文，包括程序计数器和其他寄存器；
2. 更新PCB信息；
3. 把进程的PCB移入相应的队列，如就绪、某事件阻塞等队列；
4. 选择另一个进程执行，并更新其PCB；
5. 更新内存管理的数据结构；
6. 恢复处理机上下文。

很耗费资源。

### 进程阻塞
正在执行的进程，由于期待的某些事情未发生，如请求系统资源失败、等待某种操作的完成、新数据尚未到达或无新工作等，则由系统自动执行阻塞原语，使自己由运行态变为阻塞状态。可见，进程的阻塞是进程自身的一种主动行为，也只有处于运行态的进程，才可能将其转为阻塞状态。

### 文件描述符
文件描述符用于表述指向文件的引用的抽象化概念，在形式上是一个非负整数，也就是一个索引值，指向内核为每一个进程所维护的该进程打开文件的记录表。当程序打开一个现有文件或者创建一个新文件时，内核向进程返回一个文件描述符。

### 缓存I/O
缓存I/O又被称为称为标准I/O，也就是操作系统会将I/O的数据缓存在文件系统的页缓存，数据会先被拷贝到操作系统内核的缓冲区中，然后才会从操作系统内核的缓冲区拷贝到应用程序的地址空间。

优点：

1. 在一定程度上分离了内核空间和用户空间，保护系统本身的运行安全；

2. 可以减少读盘的次数，从而提高性能。

缺点：

数据在传输过程中需要在应用程序地址空间和内核进行多次数据拷贝操作，会带来很大的CPU和内存开销。

## I/O模式
对于一次I/O访问，如read，数据会先被拷贝到操作系统内核的缓冲区中，然后才会从操作系统内核的缓冲区拷贝到应用程序的地址空间。所以说，当一个read操作发生时，它会经历两个阶段
* 等待数据准备好
* 从内核向进程复制数据

对于一个套接字上的输入操作，第一步通常设计等待数据从网络中到达。当所等待数据到达时，它被复制到内核中的某个缓冲区。第二步就是把数据从内核缓冲区复制到应用进程缓冲区。因为这两个阶段，linux产生了这五种I/O模型：
* 阻塞式I/O
* 非阻塞式I/O
* I/O复用（select和poll）
* 信号驱动式I/O
* 异步I/O

### 阻塞式I/O
当用户进程调用了recvfrom这个系统调用，内核就开始了I/O的第一个阶段：准备数据，这个过程需要等待。而在用户进程这边，整个进程被阻塞，直到数据从内核缓冲区复制到应用进程缓冲区中才返回。

阻塞式I/O的特点就是在I/O执行的两个阶段都被block了，CPU利用率比较高。

### 非阻塞式I/O
应用进程执行系统调用之后，内核返回一个error。应用进程可以继续执行，但是需要不断的执行系统调用来获知I/O是否完成，这种方式称为轮询polling。

非阻塞式I/O的特点就是用户进程需要不断的主动询问，CPU要处理更多的系统调用，CPU利用率比较低。

### I/O多路复用
I/O多路复用指的就是select、poll和epoll，也被称为事件驱动型I/O。

当用户进程调用了select，那么整个进程会被block，同时kernel会“监视”所有select负责的socket，当任何一个socket中的数据准备好了，select就会返回，之后再使用recvfrom把数据从内核复制到进程中。

select和epoll的好处就在于单个process就可以同时处理多个网络连接的I/O。所以I/O多路复用的特点就是通过一种机制一个进程能同时等待多个文件描述符，而这些文件描述符任意一个进入读就绪状态，select函数就可以返回。select/epoll的优势并不是对于单个连接能处理的更快，而在于能处理更多的connection。

### 信号驱动式I/O
应用进程使用系统调用，内核立即返回，应用进程可以继续执行，等待数据阶段应用进程是非阻塞的。内核在数据到达时向应用进程发送SIGIO信号，应用进程收到之后在信号处理程序中调用recvfrom把数据从内核复制到进程中。

### 异步式I/O
应用进程执行系统调用会立即返回，应用进程可以继续执行，不会被阻塞，内核会在所有操作完成之后向应用进程发送信号。

异步 I/O 与信号驱动 I/O 的区别在于，异步 I/O 的信号是通知应用进程 I/O 完成，而信号驱动 I/O 的信号是通知应用进程可以开始 I/O。

同步 I/O 与异步 I/O 的区别在于第二阶段，从内核缓冲区向应用进程缓冲区复制数据的过程，应用进程是否会阻塞。

## I/O 复用
select/poll/epoll都是I/O多路复用的具体实现方式，出现顺序从早到晚。这三种I/O本质上都是同步I/O，都需要在读写事件就绪后自己负责进行读写，也就是说这个读写过程是阻塞的，而异步I/O则无须自己负责进行读写，异步I/O的实现会负责把数据从内核拷贝到用户空间。

### select
允许应用程序监听一组文件描述符，等待一个或者多个描述符成为就绪状态，从而完成I/O操作。

fdset使用数组实现，数组大小使用FD_SETSIZE定义，函数监听的文件描述符分3类，调用后select函数会阻塞，直到有描述符就绪或者超时，函数返回。当select函数返回后，可以通过遍历fdset，来找到就绪的描述符。

select目前几乎在所有的平台上支持，其良好跨平台支持也是其优点。缺点在于单个进程能够监视的文件描述符的数量存在最大限制，在linux上一般为1024。

### poll
poll没有最大数量的限制，和select一样，poll返回后，需要轮询pollfd来获取就绪的描述符。

### epoll
epoll涉及到的只有三个系统调用：

epoll_create创建一个epoll实例并返回epollfd，epoll_ctl注册file descriptor等待的I/O事件到epoll实例上，epoll_wait则是阻塞监听epoll实例上所有的file descriptor的I/O事件，接收一个用户空间上的一块内存地址，kernel会在有I/O事件发生的时候把文件描述符列表复制到这块内存地址上，然后epoll_wait解除阻塞并返回，最后用户空间上的程序就可以对相应的fd进行读写了。

在实现上epoll采用红黑树来存储所有监听的fd，而红黑树本身插入和删除性能比较稳定，时间复杂度O(logN)，通过epoll_ctl函数添加进来的fd都会被放在红黑树的某个节点内，所以重复添加是没有用的。当把fd添加进来的时候会完成关键的一步：该fd会与相应的设备驱动程序建立毁掉关系，也就是在内存中断处理程序为他注册一个回调函数，在fd相应的事件触发后，内核就会调用这个回调函数。

epoll对文件描述符的操作有两种模式，LT（level trigger）和ET（edge trigger），LT模式是默认模式，区别如下：

LT模式：当epoll_wait监测到描述符事件发生并将此事件通知应用程序，**应用程序可以不立刻处理该事件**。下次调用epoll_wait时，会再次响应应用程序并通知此事件。

ET模式：当epoll_wait监测到描述符事件发生并将此事件通知应用程序，**应用程序必须立刻处理该事件**。如果不处理，下次调用epoll_wait时，不会再次响应应用程序并通知此事件。



