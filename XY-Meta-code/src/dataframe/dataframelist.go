//List Data struct module
//This package includes linked list and data struct of parameter
package dataframe

import (
	"fmt"
)

type Node struct {
	data interface{}
	next *Node
}

type LinkList struct {
	head *Node
	tail *Node
	size int
}

func CreateLinkList() LinkList {
	l := LinkList{}
	l.head = nil
	l.tail = nil
	l.size = 0
	return l
}

func (l *LinkList) EmptyIs() bool {
	return l.size == 0
}

func (l *LinkList) GetLength() int {
	return l.size
}

func (l *LinkList) Exist(node *Node) bool {
	var p *Node = l.head
	for p != nil {
		if p == node {
			return true
		} else {
			p = p.next
		}
	}
	return false
}

func (l *LinkList) GetNode(e interface{}) *Node {
	var p *Node = l.head
	for p != nil {
		if e == p.data {
			return p
		} else {
			p = p.next
		}
	}
	return nil
}

func (l *LinkList) Append(e interface{}) {
	newNode := Node{}
	newNode.data = e
	newNode.next = nil

	if l.size == 0 {
		l.head = &newNode
		l.tail = &newNode
	} else {
		l.tail.next = &newNode
		l.tail = &newNode
	}
	l.size++
}

func (l *LinkList) InsertHead(e interface{}) {
	newNode := Node{}
	newNode.data = e
	newNode.next = l.head
	l.head = &newNode
	if l.size == 0 {
		l.tail = &newNode
	}
	l.size++
}

func (l *LinkList) InsertAfterNode(pre *Node, e interface{}) {
	if l.Exist(pre) {
		newNode := Node{}
		newNode.data = e
		if pre.next == nil {
			l.Append(e)
		} else {
			newNode.next = pre.next
			pre.next = &newNode
		}
		l.size++
	} else {
		fmt.Println("链表中不存在该结点")
	}
}

func (l *LinkList) InsertAfterData(preData interface{}, e interface{}) bool {
	var p *Node = l.head
	for p != nil {
		if p.data == preData {
			l.InsertAfterNode(p, e)
			return true
		} else {
			p = p.next
		}
	}
	fmt.Println("链表中没有该数据，插入失败")
	return false
}

func (l *LinkList) Insert(position int, e interface{}) bool {
	if position < 0 {
		fmt.Println("指定下标不合法")
		return false
	} else if position == 0 {
		l.InsertHead(e)
		return true
	} else if position == l.size {
		l.Append(e)
		return true
	} else if position > l.size {
		fmt.Println("指定下标超出链表长度")
		return false
	} else {
		var index int
		var p *Node = l.head
		for index = 0; index < position-1; index++ {
			p = p.next
		}
		l.InsertAfterNode(p, e)
		return true
	}

}

func (l *LinkList) DeleteNode(node *Node) {
	if l.Exist(node) {
		if node == l.head {
			l.head = l.head.next
		} else if node == l.tail {
			var p *Node = l.head
			for p.next != l.tail {
				p = p.next
			}
			p.next = nil
			l.tail = p
		} else {
			var p *Node = l.head
			for p.next != node {
				p = p.next
			}
			p.next = node.next
		}
		l.size--
	}
}

func (l *LinkList) Delete(e interface{}) {
	p := l.GetNode(e)
	if p == nil {
		fmt.Println("链表中无该数据，删除失败")
	} else {
		l.DeleteNode(p)
	}
}

func (l *LinkList) Traverse() {
	var p *Node = l.head
	if l.EmptyIs() {
		fmt.Println("LinkList is empty")
	} else {
		for p != nil {
			fmt.Print(p.data, " ")
			p = p.next
		}
		fmt.Println()
	}
}

func (l *LinkList) Search(index int) interface{} {
	var p *Node = l.head
	var step int
	if l.EmptyIs() {
		fmt.Println("LinkList is empty")
	} else {
		for p != nil {
			if step == index {
				return p.data
			} else {
				p = p.next
				step += 1
			}
		}
	}
	return "out of index"
}

func (l *LinkList) PrintInfo() {
	fmt.Println("***********************************************")
	fmt.Println("链表长度为：", l.GetLength())
	fmt.Println("链表是否为空:", l.EmptyIs())
	fmt.Print("遍历链表：")
	l.Traverse()
	fmt.Println("***********************************************")
}

type Grade struct {
	Queryindex     int
	Index          int
	Score          float64
	TSNR           float64
	ESNR           float64
	Dot_product    float64
	Cosine_sum     float64
	FDR            float64
	Cosine_similar float64
}

type Parameters struct {
	Current_path       string
	Input              string
	Reference          string
	Output             string
	Decoyinput         string
	Search_pattern     int
	Decoy_pattern      int
	Electric_pattern   int
	HPLC_pattern       int
	Adduct_path        string
	Clear              bool
	Min_mass           int
	Max_mass           int
	Merge_tolerance    float64
	Merge_type         int
	Threshold_peaks    float64
	Precur             int
	Tolerance_precur   float64
	Isotope            int
	Tolerance_isotopic float64
	Threads            int
	Decoy_similitude   float64
	Database_filter    int
	Match_model        int
	MMI                int
}

type Adduct_isotope struct {
	Isotope_mass_list     []float64
	Isotope_precusor_nums []float64
	Isotope_type_list     []string
}
