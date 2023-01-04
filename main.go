package main

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"os"
	"time"
)

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
	f, err := os.Open("sample.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	mkt := createMerkleTree(lines, 0, len(lines)-1)
	for {
		time.Sleep(time.Second)
		f.Seek(0, io.SeekStart)
		scanner = bufio.NewScanner(f)
		lines = []string{}
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		mkt2 := createMerkleTree(lines, 0, len(lines)-1)
		diff := mkt.Compare(mkt2)
		if diff != "" {
			fmt.Println(diff)
			mkt = mkt2
		}
	}
}

func createMerkleTree(arr []string, start int, end int) *MerkleNode {
	if len(arr) == 0 {
		return &MerkleNode{}
	}
	if len(arr) == 1 {
		return &MerkleNode{Hash: hash(arr[0]), Left: nil, Right: nil, lineNum: fmt.Sprintf("%d", start)}
	}
	n := len(arr)
	lNode := createMerkleTree(arr[:n/2], start, start+(end-start)/2)
	rNode := createMerkleTree(arr[n/2:], start+(end-start)/2+1, end)
	newHash := hash(fmt.Sprintf("%d", lNode.Hash+rNode.Hash))
	return &MerkleNode{Hash: newHash, Left: lNode, Right: rNode, lineNum: fmt.Sprintf("%d-%d", start, end)}
}

func (mkl *MerkleNode) Compare(mkl2 *MerkleNode) string {
	if mkl.lineNum != mkl2.lineNum {
		if mkl.lineNum < mkl2.lineNum {
			return "line(s) added"
		} else {
			return "line(s) deleted"
		}
	}

	if mkl.Hash == 0 && mkl2.Hash == 0 {
		return ""
	}
	if mkl.Hash == 0 || mkl2.Hash == 0 {
		return mkl.lineNum
	}

	if mkl.Hash == mkl2.Hash {
		return ""
	}
	// mkl.child == nil && mkl2.child == nil
	lMatch := ""
	rMatch := ""

	// mkl.child == nil && mkl2.child != nil
	if mkl2.Left != nil {
		lMatch = mkl.Left.lineNum
	}
	if mkl2.Right != nil {
		rMatch = mkl.Right.lineNum
	}

	// mkl.child != nil && mkl2.child != nil
	if mkl.Left != nil {
		lMatch = mkl.Left.Compare(mkl2.Left)
	}
	if mkl.Right != nil {
		rMatch = mkl.Right.Compare(mkl2.Right)
	}

	if lMatch == "" && rMatch == "" {
		return mkl.lineNum
	}
	return lMatch + rMatch
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
