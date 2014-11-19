package main

import "code.google.com/p/go-tour/tree"
import "fmt"

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	walk(t, ch)
	close(ch)
}

func walk(t *tree.Tree, ch chan int) {
	if t != nil {
		walk(t.Left, ch)
		ch <- t.Value
		walk(t.Right, ch)
	}
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go Walk(t1, ch1)
	go Walk(t2, ch2)

	for {
		v1, ok1 := <- ch1
		v2, ok2 := <- ch2
		if ok1 != ok2 {
			return false
		}
		if !ok1 {
			// all channels are closed - all values so far were equal
			return true
		}
		if v1 != v2 {
			return false
		}
	}
}

func main() {
	ch := make(chan int)
	go Walk(tree.New(1), ch)
	for i := 0; i < 10; i++ {
		fmt.Printf("%d\n", <-ch)
	}

	fmt.Printf("1 vs 1 %d\n", Same(tree.New(1), tree.New(1)))
	fmt.Printf("1 vs 2 %d\n", Same(tree.New(1), tree.New(2)))
}
