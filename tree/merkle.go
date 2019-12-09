package main

import (
	"crypto/sha256"
)


type Merkable interface {
	// 叶子节点
	Hash() []byte
}

type MNode struct {
	// 中间节点
	Hash       []byte
	LeftChild  *MNode
	RightChild *MNode
}

type MTree struct {
	// merkle tree
	RootHash   []byte // root
	ResultTree [][]*MNode
	Data       []Merkable // 叶子节点
}

func EmptyMTree() *MTree { // 创建一棵空树
	emptyDataList := []Merkable{}
	emptyBlockData := CreateMTree(emptyDataList)
	return emptyBlockData
}

func CreateMTree(data []Merkable) (*MTree) { // 根据叶子节点创建树
	if len(data) == 0 {
		return &MTree{RootHash: []byte{}, Data: data}
	}

	baseTree := createBaseTree(data)
	resultTree := buildTree(baseTree)
	root := resultTree[len(resultTree)-1][0]
	mt := MTree{RootHash: root.Hash, ResultTree: resultTree, Data: data}

	return &mt
}

func createBaseTree(data []Merkable) (baseTree [][]*MNode) { // 计算叶子节点
	var baseLevel []*MNode // 第一层

	for _, d := range data {
		h := sha256.New()
		h.Write(d.Hash())
		m := &MNode{Hash: h.Sum(nil), LeftChild: nil, RightChild: nil}

		baseLevel = append(baseLevel, m)
	}

	baseTree = append(baseTree, baseLevel)
	return baseTree
}

func buildTree(inputTree [][]*MNode) (resultTree [][]*MNode) { //  递归构建树
	numLevels := len(inputTree)
	highestLevel := inputTree[numLevels-1]

	if len(highestLevel) > 1 {
		nextLevel := buildNextLevel(highestLevel) //create the net level
		inputTree = append(inputTree, nextLevel)  // add to tree

		return buildTree(inputTree)
	} else { // we have the top level
		return inputTree
	}
}

func buildNextLevel(level []*MNode) (nextLevel []*MNode) { // 使用当前level的节点 构建上一层level

	if len(level)%2 != 0 {
		level = append(level, level[len(level)-1])
	}

	for i := 0; i < len(level); i = i + 2 {
		nextLevel = append(nextLevel, hashNodes(level[i], level[i+1]))
	}
	return nextLevel
}

func (mt *MTree) hasData(m Merkable) bool { //是否存在叶子节点
	seenData := mt.Data
	for _, seenDatum := range seenData {
		if string(seenDatum.Hash()) == string(m.Hash()) {
			return true
		}
	}
	return false
}

func AddDataToTree(mt MTree, datum Merkable) MTree { // 向树中添加数据，并更新树上的节点
	newData := append(mt.Data, datum)

	baseTree := createBaseTree(newData)
	resultTree := buildTree(baseTree)
	resultTreeRoot := resultTree[len(resultTree)-1][0]

	finalTree := MTree{RootHash: resultTreeRoot.Hash, ResultTree: resultTree, Data: newData}

	return finalTree
}

func hashNodes(leftChild, rightChild *MNode) (*MNode) { // 计算左右子树的hash，sha256
	h := sha256.New()
	h.Write(leftChild.Hash)
	h.Write(rightChild.Hash)

	parent := MNode{Hash: h.Sum(nil), LeftChild: leftChild, RightChild: rightChild}

	return &parent
}



func (mt *MTree) GetRoot() []byte {
	return mt.RootHash
}

func (mt *MTree) GetLeaf(index uint64) (hash []byte) {
	return mt.Data[index].Hash()
}

func (mt *MTree) GetLeafIndex(hash []byte) (uint64) {
	for i := uint64(0); i < uint64(len(mt.Data)); i++ {
		if string(mt.GetLeaf(i)) == string(hash) {
			return i
		}
	}
	return -1
}

// proof 单独生成
func (mt *MTree) MerkleProof(index uint64) (proof [][]byte, proof_index []uint64) {
	tree_depth := uint64(len(mt.ResultTree))
	address_bits := make([]uint64, tree_depth)
	merkleProof := make([][]byte, tree_depth)
	for i := uint64(0); i < tree_depth; i++ {
		if index%2 == 0 {
			address_bits[i] = 1;
			merkleProof[i] = mt.ResultTree[i][index+1].Hash;
		} else {
			address_bits[i] = 0;
			merkleProof[i] = mt.ResultTree[i][index-1].Hash;
		}
		index = index / 2;
	}
	return merkleProof, address_bits
}

func verifyMTree(proof [][]byte,proof_index []uint64,root []byte) bool { // 校验两棵树是否相等

	return true
}

// StringMerkle implements Merkable, with just a string
//type StringMerkle struct {
//	Data string
//}
//
//func (s StringMerkle) Hash() []byte{
//	h := sha256.New()
//	h.Write([]byte(s.Data))
//	return h.Sum(nil)
//}

/* Since []StringMerkables does not implement []Merkables using
 * Go, you have to do this converion manually
 */
//func convertMerkleDataToMerkable(data []StringMerkle) []Merkable{
//	merkables := make([]Merkable, len(data))
//	for i, v := range data {
//		merkables[i] = Merkable(v)
//	}
//	return merkables
//}
