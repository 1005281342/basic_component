package rbtree

import (
	"fmt"
	"github.com/1005281342/test_tools/survey"
	"sort"
)

var x []int

func InOrder(rb *RBTree) []int {
	inOrder(rb.root)
	var t = x
	x = make([]int, 0)
	return t
}

func inOrder(node *Node) {
	if node == nil {
		return
	}
	if node.left != nil {
		inOrder(node.left)
	}
	//fmt.Printf("(%v, %v) ", node.value, node.color)
	x = append(x, node.key.(int))
	if node.right != nil {
		inOrder(node.right)
	}
}

func Check(n int) bool {
	return checkNums(n)
}

func checkNums(n int) bool {
	var rb = New()
	var nums = survey.RandShuffle(n)
	for _, num := range nums {
		rb.Add(num, num)
	}
	sort.Ints(nums)
	var nu = InOrder(rb)

	if len(nums) != len(nu) {
		fmt.Printf("check failed, n: %d, want: %v, got: %v",
			n, nums, nu)
		return false
	}

	for i := 0; i < len(nums); i++ {
		if nums[i] != nu[i] {
			fmt.Printf("check failed, n: %d, want: %v, got: %v",
				n, nums, nu)
			return false
		}
	}

	if !CheckColor(rb) {
		return false
	}

	if !checkBlackHigh(rb) {
		return false
	}

	return true
}

// 根节点是黑色的
// 每个红色节点的的两个子节点一定是黑色的
func CheckColor(tree *RBTree) bool {
	if tree == nil {
		return true
	}

	if tree.root.color != BLACK {
		return false
	}

	return checkColor(tree.root)
}

func checkColor(node *Node) bool {
	if node == nil {
		return true
	}

	if node.color == RED {
		if node.left != nil && node.left.color == RED {
			return false
		}
		if node.right != nil && node.right.color == RED {
			return false
		}
	}
	return checkColor(node.left) && checkColor(node.right)
}

// 任意节点到其子节点的路径上所包含的黑色节点数量相等（从这个结论又可以得出：如果一个节点存在黑色子节点，那么该节点一定有两个子节点）
func checkBlackHigh(tree *RBTree) bool {
	if tree == nil {
		return true
	}
	if tree.root == nil {
		return true
	}

	return checkBlackNum(tree.root)
}

func checkBlackNum(node *Node) bool {
	if node == nil {
		return true
	}
	if !checkBlackNumsSub(node) {
		return false
	}
	return checkBlackNumsSub(node.left) && checkBlackNumsSub(node.right)
}

type sNode struct {
	n *Node
	c int
}

func checkBlackNumsSub(node *Node) bool {
	if node == nil {
		return true
	}
	var t = -1
	var dq = []*sNode{{n: node, c: 0}}
	for len(dq) > 0 {
		var x = dq[0]
		dq = dq[1:]
		if x.n.color == BLACK {
			x.c++
		}

		if x.n.left == nil && x.n.right == nil {
			if t < 0 {
				t = x.c
			} else {
				if t != x.c {
					return false
				}
			}
		}

		if x.n.left != nil {
			dq = append(dq, &sNode{n: x.n.left, c: x.c})
		}
		if x.n.right != nil {
			dq = append(dq, &sNode{n: x.n.right, c: x.c})
		}
	}
	return true
}
