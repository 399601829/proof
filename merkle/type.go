package merkle

import "io"

type Merkle interface {
	GetMerkleProof(index uint64) (path [][]byte, pathIndex []uint64)                                    // 获取某个叶子节点的Merkle proof
	VerifyMerkleProof(leafData []byte, merklePath [][]byte, merklePathIndex []uint64, root []byte) bool // 校验Merkle proof 是否正确

	GetLeaf(index uint64) (hash []byte) // 获取某个叶子节点
	GetRoot() (hash []byte)             // 获取当前Merkle tree root
	Insert(hash []byte) (index uint64)  // 插入叶子节点
	UpdateTree() (root []byte)          // 插入叶子节点后，更新merkle tree

	Save(w io.Writer) error // 保存
	Load()
}

type Hasher interface {
	Hash() []byte
}



type MerkleTree struct {
	data [][]byte

}