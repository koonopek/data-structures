package main

import (
	"fmt"
	"reflect"

	"golang.org/x/crypto/sha3"
)

type MerkleTree struct {
	root *Node
}

type Node struct {
	hash   []byte
	value  *[]byte
	left   *Node
	right  *Node
	parent *Node
}

type Proof struct {
	hash   []byte
	isLeft bool
}

func (node *Node) GenProof(proofPath []Proof) []Proof {
	if node.parent == nil {
		return proofPath
	}

	if node.parent.left != nil && node.parent.left != node {
		proofPath = append(proofPath, Proof{hash: node.parent.left.hash, isLeft: true})
	}

	if node.parent.right != nil && node.parent.right != node {
		proofPath = append(proofPath, Proof{hash: node.parent.right.hash, isLeft: false})
	}

	return node.parent.GenProof(proofPath)
}

func GenerateProof(root *Node, value *[]byte) (error, []Proof) {
	hash := keccak(*value)
	nodeToProof := root.FindByHash(hash)
	if nodeToProof == nil {
		return fmt.Errorf("hash=%x not found in tree", hash), nil
	}
	return nil, nodeToProof.GenProof([]Proof{})
}

func VerifyProof(rootHash []byte, proofPath []Proof, valueToProof *[]byte) bool {
	lastHash := keccak(*valueToProof)

	for _, proof := range proofPath {
		if proof.isLeft {
			lastHash = keccak(append(proof.hash, lastHash...))
		} else {
			lastHash = keccak(append(lastHash, proof.hash...))
		}
	}

	return reflect.DeepEqual(rootHash, lastHash)
}

func (node *Node) Depth() int {
	currentNode := node
	depth := 1
	for {
		depth++
		if currentNode.parent == nil {
			return depth
		}
		currentNode = currentNode.parent
	}
}

func (node *Node) FindByHash(hash []byte) *Node {
	var foundNode *Node

	if reflect.DeepEqual(node.hash, hash) {
		return node
	}

	if node.left != nil {
		foundNode = node.left.FindByHash(hash)
		if foundNode != nil {
			return foundNode
		}
	}

	if node.right != nil {
		foundNode = node.right.FindByHash(hash)
		if foundNode != nil {
			return foundNode
		}
	}

	return nil
}

// func VerifyProof

func BuildTree(values []*[]byte) Node {
	if len(values)%2 == 1 {
		panic("Build tree require event values count")
	}
	leafs := make([]Node, len(values))

	for i := 0; i < len(values); i++ {
		leafs[i].hash = keccak(*values[i])
		leafs[i].value = values[i]
	}

	root := buildParentNodes(leafs)

	return root
}

func buildParentNodes(leafs []Node) Node {
	if len(leafs) == 1 {
		return leafs[0]
	}

	parents := make([]Node, len(leafs)/2)

	i := 0
	for offset := 1; offset < len(leafs); offset += 2 {
		parents[i] = Node{
			hash:  keccak(append(leafs[offset-1].hash, leafs[offset].hash...)),
			left:  &leafs[offset-1],
			right: &leafs[offset],
		}
		leafs[offset-1].parent = &parents[i]
		leafs[offset].parent = &parents[i]
		i++
	}

	return buildParentNodes(parents)
}

func keccak(value []byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(value)
	return hash.Sum(nil)
}

func main() {
	milenka, michas, mia, kazik := []byte("milenka"),
		[]byte("michas"),
		[]byte("mia"),
		[]byte("kazik")

	values := []*[]byte{
		&milenka, &michas, &mia, &kazik,
	}

	root := BuildTree(values)
	fmt.Printf("root hash %x tree depth %d \n", root.hash, root.Depth())

	fmt.Println("generating proof for kazik")
	err, proof := GenerateProof(&root, &kazik)

	if err != nil {
		panic(err.Error())
	}

	fmt.Print("checking for george ")
	george := []byte("george_the_liar")
	if VerifyProof(root.hash, proof, &george) == false {
		fmt.Println("proof invalid")
	} else {
		fmt.Println("proof valid")
	}

	fmt.Print("checking for kazik ")
	if VerifyProof(root.hash, proof, &kazik) == false {
		fmt.Println("proof invalid")
	} else {
		fmt.Println("proof valid")
	}
}
