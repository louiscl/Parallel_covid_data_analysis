package main

import (
	"fmt"
	"sync/atomic"
)

// Double Ended LinkedList implementation ------------------------------
// Could be simplified

type node struct {
	content int
	task    Runnable
	next    *node
	prev    *node
}

type LinkedList struct {
	// NOTE
	// In a list of length 1, the single node is always
	// considered the first node, and the last node is nil
	firstNode *node
	lastNode  *node
	length    int
}

// LinkedList methods

func NewLinkedList() *LinkedList {
	newList := LinkedList{length: 0}
	return &newList
}

func (list *LinkedList) addToBeginning(node *node) {
	if list.firstNode == nil {
		list.firstNode = node
		list.length = 1
	} else {
		prevFirst := list.firstNode
		prevFirst.prev = node
		list.firstNode = node
		node.next = prevFirst
		list.length++
	}
	// could add the bound; to not overfill [Lecture; not in slides]
	// But not necessary, since tasks are pre-allocated
}

func (list *LinkedList) addToTail(node *node) {
	if list.length == 0 {
		list.firstNode = node
		list.length = 1
	} else if list.length == 1 {
		list.firstNode.next = node
		list.lastNode = node
		node.prev = list.firstNode
		list.length = 2
	} else {
		list.lastNode.next = node
		temp := list.lastNode
		list.lastNode = node
		node.prev = temp
		list.length++
	}
	// could add the bound; to not overfill
	// But not necessary, since tasks are pre-allocated
}

func (list *LinkedList) popHead() *node {
	if list.length == 0 {
		return nil
	} else if list.length == 1 {
		temp := list.firstNode
		list.firstNode = nil
		list.length = 0
		return temp
	} else if list.length == 2 {
		temp := list.firstNode
		list.firstNode = list.lastNode
		list.firstNode.prev = nil
		list.lastNode = nil
		list.length = 1
		return temp
	} else {
		temp := list.firstNode
		list.firstNode = list.firstNode.next
		list.firstNode.prev = nil
		list.length--
		return temp
	}
}

func (list *LinkedList) popTail() *node {
	if list.length == 0 {
		return nil
	} else if list.length == 1 {
		temp := list.firstNode
		list.firstNode = nil
		list.length = 0
		return temp
	} else if list.length == 2 {
		temp := list.lastNode
		list.firstNode.next = nil
		list.lastNode = nil
		list.length = 1
		return temp
	} else {
		temp := list.lastNode
		list.lastNode = list.lastNode.prev
		list.lastNode.next = nil
		list.length--
		return temp
	}
}

func (list *LinkedList) displayList() {
	fmt.Println("linkedList Content with index:")
	nextNode := list.firstNode
	for i := 0; i < list.length; i++ {
		fmt.Println("   ", i, nextNode.content, "task:", nextNode.task)
		nextNode = nextNode.next
	}
	fmt.Println("End of LinkedList")
}

// Runnable

type RunnableStruct struct {
	number int
}

type Runnable interface {
	ReturnFileNum() int
}

func NewRunnable(num int) Runnable {
	NewRunnable := &RunnableStruct{number: num}
	return NewRunnable
}

func (task *RunnableStruct) ReturnFileNum() int {
	return task.number
}

// DEQueue

type BoundedDEQueue struct {
	grabbedStamp int32
	linkedList   LinkedList
}

func NewBoundedDEQueue() DEQueue {
	newList := NewLinkedList()
	newDEQ := &BoundedDEQueue{grabbedStamp: 0, linkedList: *newList}
	return newDEQ
}

type DEQueue interface {
	PushBottom(task Runnable)
	PopTop() Runnable
	PopBottom() Runnable
	Length() int
	DisplayList()
}

func (BQ *BoundedDEQueue) PushBottom(task Runnable) {
	taskNode := node{task: task}
	BQ.linkedList.addToTail(&taskNode)
}

func (BQ *BoundedDEQueue) PopBottom() Runnable {
	if BQ.linkedList.length == 0 {
		return nil
	}

	// check the top to see whether there's a synchronization conflict
	oldStamp := BQ.grabbedStamp
	newStamp := oldStamp + 1

	// if head and tail are more than one apart; no conflict
	if BQ.linkedList.length > 1 {
		bottomNode := BQ.linkedList.popTail()
		bottomRun := bottomNode.task
		return bottomRun
	}

	// if only one item left, we need to make sure it's not being stolen simultaneously
	if BQ.linkedList.length == 1 {
		// no need of reset bottom to 0 because done by linkedList automatically
		swapBoolStamp := atomic.CompareAndSwapInt32(&BQ.grabbedStamp, oldStamp, newStamp)
		if swapBoolStamp {
			bottomNode := BQ.linkedList.popTail()
			bottomRun := bottomNode.task
			return bottomRun
		}
	}
	return nil
}

func (BQ *BoundedDEQueue) PopTop() Runnable {
	// == steal work

	// Given that in this exercise we first fill up the dequeues with tasks,
	// and then work on them, it's not necessary to compare oldTopNode and newTopNode.
	// The top node stays the same, unless taken by another thief; but we control for this
	// by comparing and swapping the stamp int

	oldStamp := BQ.grabbedStamp
	newStamp := oldStamp + 1

	// is there anything to steal? Not stealing if just one task
	if BQ.linkedList.length <= 1 {
		return nil
	}

	swapBoolStamp := atomic.CompareAndSwapInt32(&BQ.grabbedStamp, oldStamp, newStamp)

	if swapBoolStamp {
		return BQ.linkedList.popHead().task
	}

	return nil // => potentially find other deque to steal from
}

// helpers

func (BQ *BoundedDEQueue) Length() int {
	return BQ.linkedList.length
}

func (BQ *BoundedDEQueue) DisplayList() {
	BQ.linkedList.displayList()
}
