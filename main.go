package main

import (
	"fmt"
	"hash/fnv"
	"strings"
)

const textContent = `aaaaaaaaaaaaaaaaaa
bbbbbbbbbbbbbbbbbbbb
ccccccccccccccccccccc
ddddddddddddddddddddd`

const textContent2 = `aaaaaaaaaaaaaaaaaa
bbbbbbbbbbbbbbbbbbb
ccccccccccccccccccccc
ddddddddddddddddddddd`

type MerkleNode struct {
	Hash    uint32
	lineNum string
	Left    *MerkleNode
	Right   *MerkleNode
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func main() {
	mkt := createMerkleTree(strings.Split(textContent, "\n"), 0, 3)
	mkt.Print()
	mkt2 := createMerkleTree(strings.Split(textContent2, "\n"), 0, 3)
	mkt2.Print()
}

func createMerkleTree(arr []string, start int, end int) MerkleNode {
	if len(arr) == 0 {
		return MerkleNode{}
	}
	if len(arr) == 1 {
		return MerkleNode{Hash: hash(arr[0]), Left: nil, Right: nil, lineNum: fmt.Sprintf("%d", start)}
	}
	n := len(arr)
	lNode := createMerkleTree(arr[:n/2], start, start+(end-start)/2)
	rNode := createMerkleTree(arr[n/2:], start+(end-start)/2+1, end)
	newHash := hash(fmt.Sprintf("%d", lNode.Hash+rNode.Hash))
	return MerkleNode{Hash: newHash, Left: &lNode, Right: &rNode, lineNum: fmt.Sprintf("%d-%d", start, end)}
}

func (mkl1 *MerkleNode) Compare(mkl2 MerkleNode) bool {
	return true
}

func (mkl *MerkleNode) Print() {
	if mkl.Hash == 0 {
		return
	}
	fmt.Printf("(%s) - %d\n", mkl.lineNum, mkl.Hash)
	if mkl.Left != nil {
		mkl.Left.Print()
	}
	if mkl.Right != nil {
		mkl.Right.Print()
	}
}
