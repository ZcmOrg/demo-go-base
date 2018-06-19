package main

import (
	"container/heap"
	"container/list"
	"container/ring"
	"fmt"
)

type IntHeap []int

//我们自定义一个堆需要实现5个接口
//Len(),Less(),Swap()这是继承自sort.Interface
//Push()和Pop()是堆自已的接口

//返回长度
func (h *IntHeap) Len() int {
	return len(*h)
}

//比较大小(实现最小堆)
func (h *IntHeap) Less(i, j int) bool {
	return (*h)[i] < (*h)[j]
}

//交换值
func (h *IntHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

//压入数据
func (h *IntHeap) Push(x interface{}) {
	//将数据追加到h中
	*h = append(*h, x.(int))
}

//弹出数据
func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	//让h指向新的slice
	*h = old[0 : n-1]
	//返回最后一个元素
	return x
}

func printRing(r *ring.Ring) {
	r.Do(func(v interface{}) {
		fmt.Print(v.(int), " ")
	})
	fmt.Println()
}

func printList(l *list.List) {
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Print(e.Value, " ")
	}
	fmt.Println()
}

func main() {
	{ //heap 的例子
		a := IntHeap{6, 2, 3, 1, 5, 4}
		//初始化堆
		heap.Init(&a)
		fmt.Println(a)
		//弹出数据，保证每次操作都是规范的堆结构
		fmt.Println(heap.Pop(&a))
		fmt.Println(a)
		fmt.Println(heap.Pop(&a))
		fmt.Println(a)
		heap.Push(&a, 0)
		heap.Push(&a, 8)
		fmt.Println("up", a)
	}
	{ //list 的例子
		//创建一个链表
		l := list.New()

		//链表最后插入元素
		a1 := l.PushBack(1)
		b2 := l.PushBack(2)

		//链表头部插入元素
		l.PushFront(3)
		l.PushFront(4)

		fmt.Println("list", l)
		printList(l)

		//取第一个元素
		f := l.Front()
		fmt.Println(f.Value)

		//取最后一个元素
		b := l.Back()
		fmt.Println(b.Value)

		//获取链表长度
		fmt.Println(l.Len())

		//在某元素之后插入
		l.InsertAfter(66, a1)

		//在某元素之前插入
		l.InsertBefore(88, a1)

		printList(l)

		l2 := list.New()
		l2.PushBack(11)
		l2.PushBack(22)
		//链表最后插入新链表
		l.PushBackList(l2)
		printList(l)

		//链表头部插入新链表
		l.PushFrontList(l2)
		printList(l)

		//移动元素到最后
		l.MoveToBack(a1)
		printList(l)

		//移动元素到头部
		l.MoveToFront(a1)
		printList(l)

		//移动元素在某元素之后
		l.MoveAfter(b2, a1)
		printList(l)

		//移动元素在某元素之前
		l.MoveBefore(b2, a1)
		printList(l)

		//删除某元素
		l.Remove(a1)
		printList(l)
	}
	{ //环的例子
		//创建环形链表
		r := ring.New(5)
		//循环赋值
		for i := 0; i < 5; i++ {
			r.Value = i
			//取得下一个元素
			r = r.Next()
		}
		printRing(r)
		//环的长度
		fmt.Println(r.Len())

		//移动环的指针
		r.Move(2)

		//从当前指针删除n个元素
		r.Unlink(2)
		printRing(r)

		//连接两个环
		r2 := ring.New(3)
		for i := 0; i < 3; i++ {
			r2.Value = i + 10
			//取得下一个元素
			r2 = r2.Next()
		}
		printRing(r2)

		r.Link(r2)
		printRing(r)
	}
}
