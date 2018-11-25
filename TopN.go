package main

import (
	"fmt"
	"math/rand"
	"time"
)

//最大堆实现的topN算法
/*
算法思路:
1.当节点数量小于topN数量时直接加入slice
2.当节点数量等于topN时,和堆顶节点进行比较,当大于堆顶节点则丢弃,否则替代堆顶节点并调整堆
假设数据量为M,返回排序结果数为N.则算法复杂度为M(log2N)
*/

type MaxHeap struct {
    heapItemSlice     []Item // for adjust heap
    //不要存放指针,指针操作时间大约是值的操作的两倍
    sortItemSlice     []Item  // sorted
    len               int
    cap               int
    sorted            bool    //sortItemSlice是否已经排好序
    heaped            bool   // 是否为最大堆
}

type TOPN struct {
	mhp     *MaxHeap
	topN    int
}

type Item interface {
	Less (other Item) bool   //for user define their struct
}

func NewMaxHeapP(cap int) *MaxHeap{
	if cap <= 0 {
		fmt.Println("cap must a positive num.")
		return nil
	}
	heapItemSlice := make([]Item , cap)
	sortItemSlice := make([]Item , cap)
	return &MaxHeap{heapItemSlice,sortItemSlice,0,cap,false,false}
}


func NewTOPNP(topN int) *TOPN{
	if topN <= 0 {
		fmt.Println("topN must a positive num.")
		return nil
	}
	return &TOPN{NewMaxHeapP(topN),topN}
}

//调整堆得某个位置和堆得大小
func AdjustMaxHeapOnPosition(itemSlice []Item,position int,end int) {
    if itemSlice == nil || position >= end {
		return
	}
	j  := position *2 + 1 // 子节点位置
	//fmt.Println("position:",position)
	positionValue := itemSlice[position] // 先将positionValue保存起来
	for j < end {
		//在两个子节点中找一个最大的值
		if j+1 <end && (itemSlice[j]).Less(itemSlice[j+1]){
			j = j+1
		}//if
		//和最大值进行比较,如果大于最大值则直接停止
		if !(positionValue).Less(itemSlice[j]){ //use cached positionValue
			break
		}//if
		//向上上升
		itemSlice[position] = itemSlice[j]
		position               = j
		j = position * 2 + 1
	}//for
	itemSlice[position]     = positionValue // 最后放入position的值
	//调整完毕
}//ad

//调整整个堆
func (mh *MaxHeap)adjustFullMaxHeap(){
	var i int
	var len = mh.len
	var heapItemSlice = mh.heapItemSlice
	if mh.len > 2 {
		for i= int(mh.len / 2 -1);i >=0 ;i--{
			AdjustMaxHeapOnPosition(heapItemSlice,i,len)
		}
	}//if
	mh.heaped  = true
}

//大根堆不一定是完全有序堆,需要再次调整才能成为有序堆
//方法为不断将堆
func (mh *MaxHeap) sortFullMaxHeap(){
    var temp Item
    var i int
    copy(mh.sortItemSlice[0:mh.len],mh.heapItemSlice[0:mh.len])
    sortItemSlice := mh.sortItemSlice
    for i = int( mh.len-1);i>=0;i--{
    	temp = sortItemSlice[0]
		sortItemSlice[0] = sortItemSlice[i]
		sortItemSlice[i] = temp
    	AdjustMaxHeapOnPosition(sortItemSlice,0,i)
	}//for
	mh.sorted   = true
}

func (mh *MaxHeap) insert(item Item) bool{
	if mh == nil {
		return false
	}//if
    if mh.len+1 <= mh.cap{
    	mh.heapItemSlice[mh.len] = item
		mh.len += 1
    	mh.heaped = false   // 顺序打乱了
    	mh.sorted = false   //两个都为false
    	return true
	}else{
		if mh.heaped == false{
			mh.adjustFullMaxHeap()
		}//
		if (mh.heapItemSlice[0]).Less(item){
			return true
		}else{
			mh.heapItemSlice[0] = item
			AdjustMaxHeapOnPosition(mh.heapItemSlice,0,mh.len) // 仍然是大根堆
			mh.sorted  = false  //增加新元素sortItem又无序了
			return true
		}//
	}
}

func (mh *MaxHeap) topNOrdered(topN int) []Item{
    if 	mh.sorted == false{
    	mh.sortFullMaxHeap()
	}//
	if topN > mh.len{
		topN = mh.len
	}
	return mh.sortItemSlice[0:topN]  //at most len
}//

func (mh *MaxHeap) topN(topN int)[]Item{
	if topN > mh.len{
		topN = mh.len
	}
	return mh.heapItemSlice[0:topN]  //at most len
}

func (tn *TOPN) Insert(item Item){
	tn.mhp.insert(item)
}

func (tn *TOPN) TopNOrdered(topN int) []Item{
	if topN <= 0{
		return nil
	}
    return tn.mhp.topNOrdered(topN)
}//

func (tn *TOPN) TopN(topN int)[]Item{
	if topN <= 0{
		return nil
	}
	return tn.mhp.topN(topN)
}
func (tn *TOPN) Len() int{
	return (*tn.mhp).len
}

func (tn *TOPN) Cap() int{
	return tn.topN
}

type Int int

func (i Int) Less(other Item) bool{
    return i < other.(Int)
}

type String string

func (s String) Less(other Item) bool{
	return s < other.(String)
}

type Float64 float64

func (f64 Float64) Less(other Item) bool{
	return f64 < other.(Float64)
}

///example
func main() {
	rand.Seed(time.Now().Unix())
	tn := NewTOPNP(1000)
	intSlice := genRandInt(500000)
	t1 := time.Now().Nanosecond()
	for i := 0; i < 500000; i++ {
		tn.Insert(Int(intSlice[i]))
	}
	t2 := time.Now().Nanosecond()
	fmt.Println((t2-t1)/1e6, "ms")
	fmt.Println(tn.TopNOrdered(1000))
}

//about 20 ms
func genRandInt(num int) []int{
	intSlice  := make([]int,num)
	for i:=0;i<num;i++{
		intSlice[i]  = rand.Intn(1000000)
	}
	fmt.Println(intSlice)
	return intSlice
}


