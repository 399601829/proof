package main

import (
	"testing"
	"crypto/sha256"
)
func genMerkleList(howMany int) []StringMerkle{
	var listOfMerkleStrings []StringMerkle

	for i := 0; i < howMany; i ++ {
		sm := StringMerkle{Data: "Test"}
		listOfMerkleStrings = append(listOfMerkleStrings, sm)
	}

	return listOfMerkleStrings
}

func TestBuildTree (t *testing.T){
	merkleDataList          := genMerkleList(10)
	baseTree                := createBaseTree(convertMerkleDataToMerkable(merkleDataList))

	resultTree              := buildTree(baseTree)
	resultTreeHeight        := len(resultTree)
	rootLevel               := resultTree[resultTreeHeight - 1]
	root 	                := rootLevel[0]
	levelUnderRoot          := resultTree[resultTreeHeight - 2]
	rootLeftChild           := levelUnderRoot[0]
	rootRightChild          := levelUnderRoot[1]

	if len(rootLevel) != 1 {
		t.Error("root does not have one node")
	}

	if len(levelUnderRoot) != 2 {
		t.Error("level under root does not have two nodes")
	}

	if root.LeftChild != rootLeftChild {
		t.Error("root leftchild does not match")
	}

	if root.RightChild != rootRightChild {
		t.Error("root right child does not match")
	}

	h := sha256.New()
	h.Write(rootLeftChild.Hash)
	h.Write(rootRightChild.Hash)

	if string(root.Hash) != string(h.Sum(nil)){
		t.Error("hash of the root does not match the hashes of its children")
	}
}

func TestCreateBaseTree(t *testing.T){
	var baseTree [][]*MNode

	merkleDataList := genMerkleList(1)
	baseTree        = createBaseTree(convertMerkleDataToMerkable(merkleDataList))

	//calculate hash of single tree node
	h := sha256.New()
	h.Write([]byte(merkleDataList[0].Data))
	// h.Write(merkleDataList.Sig)

	// test base tree only has one level
	if len(baseTree) != 1{
		t.Error("Base tree should only have one layer")
	}

	//test the transaction hashes match
	if string(baseTree[0][0].Hash) == string(h.Sum(nil)) {
		t.Error("Hash of the transaction in the baseTree is incorrect")
	}

	// try again with a transaction list of three
	merkleDataList = genMerkleList(3)
	baseTree = createBaseTree(convertMerkleDataToMerkable(merkleDataList))

	// test base tree only has one level
	if len(baseTree) != 1{
		t.Error("Base tree should only have one layer")
	}

	// test there are three transactions in base layer
	if len(baseTree[0]) != 3 {
		t.Error("Base layer should have exactly three transactions")
	}
}

func TestBuildNextLevel(t *testing.T) {

	merkleDataList          := genMerkleList(10)
	baseTree    := createBaseTree(convertMerkleDataToMerkable(merkleDataList))
	baseLevel   := baseTree[0]
	secondLevel := buildNextLevel(baseLevel)
	thirdLevel  := buildNextLevel(secondLevel)
	fourthLevel := buildNextLevel(thirdLevel) //should have two nodes
	fifthLevel  := buildNextLevel(fourthLevel) //should just be the root

	if len(secondLevel) != 5 {
		t.Error("Second level has incorrect number of nodes")
	}
	if len(thirdLevel) != 3 {
		t.Error("Third level has incorrect number of nodes")
	}
	if len(fifthLevel) != 1 {
		t.Error("Fifth level has incorrect number of nodes")
	}

	root := fifthLevel[0]
	rootLeftChild  := fourthLevel[0]
	rootRightChild := fourthLevel[1]

	if root.LeftChild != rootLeftChild{
		t.Error("root's left child is incorrect")
	}
	if root.RightChild != rootRightChild{
		t.Error("root's right child is incorrect")
	}

	h := sha256.New()
	h.Write(rootLeftChild.Hash)
	h.Write(rootRightChild.Hash)

	if string(h.Sum(nil)) != string(root.Hash){
		t.Error("roots children do not add to hash of the root")
	}
}


func TestHashNodes(t *testing.T){
	merkleDataList  := genMerkleList(2)
	baseTree        := createBaseTree(convertMerkleDataToMerkable(merkleDataList))
	baseLevel       := baseTree[0]

	leftChild       := baseLevel[0]
	rightChild      := baseLevel[1]

	parent          := hashNodes(leftChild, rightChild)

	h := sha256.New()
	h.Write(leftChild.Hash)
	h.Write(rightChild.Hash)

	if parent.LeftChild != leftChild {
		t.Error("Left child of parent does not correspond")
	}
	if parent.RightChild != rightChild {
		t.Error("Right child of parent does not correspond")
	}

	if string(parent.Hash) !=  string(h.Sum(nil)){
		t.Error("Hash of parent does not correspond to the hashes of its children")
	}
}