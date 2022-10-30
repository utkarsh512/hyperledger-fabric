// CS61065 - Assignment 4 - Part B
//
// Authors:
// Utkarsh Patel (18EC35034)
// Saransh Patel (18CS30039)

package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

var seed = rand.NewSource(42)
var rng = rand.New(seed)

type SmartContract struct {
	contractapi.Contract
}

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

type MyBST struct {
	PrimaryKey string
	Root       *TreeNode
}

// MyBSTExists checks whether key is already present in the ledger
func (s *SmartContract) MyBSTExists(ctx contractapi.TransactionContextInterface, key string) (bool, error) {
	fmt.Println("MyBSTExists: Started..")
	bst, err := ctx.GetStub().GetState(key)
	if err != nil {
		fmt.Println("MyBSTExists: Error while reading state")
		return false, err
	}
	if bst != nil {
		fmt.Println("MyBSTExists: Found BST for key %s", key)
		return true, nil
	}
	fmt.Println("MyBSTExists: BST for key %s doesn't exist", key)
	return false, nil
}

// ReadMyBST returns the BST present in the ledger as a structure
func (s *SmartContract) ReadMyBST(ctx contractapi.TransactionContextInterface) (*MyBST, error) {
	fmt.Println("ReadMyBST: Started..")
	bstIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		fmt.Println("ReadMyBST: Error in creating iterator")
		return nil, err
	}
	defer bstIterator.Close()

	fmt.Println("ReadMyBST: Iteration started")

	// Iterate over the contract
	var result []*MyBST
	for bstIterator.HasNext() {
		res, err := bstIterator.Next()
		if err != nil {
			fmt.Println("ReadMyBST: Error in getting next element from iterator")
			return nil, err
		}

		var bst MyBST
		err = json.Unmarshal(res.Value, &bst)
		if err != nil {
			fmt.Println("ReadMyBST: Error in unmarshalling BST JSON")
			return nil, err
		}

		result = append(result, &bst)
	}

	n := len(result)
	fmt.Println("ReadMyBST: Found %d BST(s) in the ledger", n)
	if n == 0 {
		return nil, nil
	}

	return result[0], nil
}

// UpdateMyBST inserts/deletes values into/from the BST
func (s *SmartContract) UpdateMyBST(ctx contractapi.TransactionContextInterface, val int, bst *MyBST, operation int) error {
	fmt.Println("UpdateMyBST: Started..")
	if operation == 0 {
		// Insert value in BST
		fmt.Println("UpdateMyBST: Performing insertion for %d", val)
		newroot, err := s.InsertValue(ctx, bst.Root, val)
		if err != nil {
			fmt.Println("UpdateMyBST: Error thrown by InsertValue")
			return err
		}
		fmt.Println("UpdateMyBST: Insertion completed")
		bst.Root = newroot
	} else {
		// Delete value in BST
		fmt.Println("UpdateMyBST: Performing deletion for %d", val)
		newroot, err := s.DeleteValue(ctx, bst.Root, val)
		if err != nil {
			fmt.Println("UpdateMyBST: Error thrown by DeleteValue")
			return err
		}
		fmt.Println("UpdateMyBST: Deletion completed")
		bst.Root = newroot
	}

	bstJSON, err := json.Marshal(bst)
	if err != nil {
		fmt.Println("UpdateMyBST: Error in marshalling BST object")
		return err
	}
	fmt.Println("UpdateMyBST: Commiting update")
	return ctx.GetStub().PutState(bst.PrimaryKey, bstJSON)
}

// InsertValue insert the value as per BST rule
func (s *SmartContract) InsertValue(ctx contractapi.TransactionContextInterface, node *TreeNode, val int) (*TreeNode, error) {
	if node == nil {
		return &TreeNode{Val: val}, nil
	}

	if val < node.Val {
		newnode, err := s.InsertValue(ctx, node.Left, val)
		if err != nil {
			return node, err
		}
		node.Left = newnode
	} else if val > node.Val {
		newnode, err := s.InsertValue(ctx, node.Right, val)
		if err != nil {
			return node, err
		}
		node.Right = newnode
	}
	return node, nil
}

// DeleteValue removes the node with passed value
func (s *SmartContract) DeleteValue(ctx contractapi.TransactionContextInterface, node *TreeNode, val int) (*TreeNode, error) {
	if node == nil {
		return node, fmt.Errorf("value to be deleted doesn't exist")
	}

	if val < node.Val {
		newnode, err := s.DeleteValue(ctx, node.Left, val)
		if err != nil {
			return node, err
		}
		node.Left = newnode

	} else if val > node.Val {
		newnode, err := s.DeleteValue(ctx, node.Right, val)
		if err != nil {
			return node, err
		}
		node.Right = newnode

	} else {
		if node.Left == nil {
			temp := node.Right
			node = nil
			return temp, nil

		} else if node.Right == nil {
			temp := node.Left
			node = nil
			return temp, nil
		}

		temp := node.Right
		for temp.Left != nil {
			temp = temp.Left
		}
		node.Val = temp.Val

		newnode, err := s.DeleteValue(ctx, node.Right, temp.Val)
		if err != nil {
			return node, err
		}
		node.Right = newnode
	}

	return node, nil
}

// Insert takes a value to be inserted into the BST.
// If the BST doesn't exist in the ledger, it creates one.
func (s *SmartContract) Insert(ctx contractapi.TransactionContextInterface, val int) error {
	fmt.Println("Insert: Started..")
	bst, err := s.ReadMyBST(ctx)

	if err != nil {
		fmt.Println("Insert: Error in reading BST")
		return err
	}

	if bst != nil {
		// BST exists in the ledger
		fmt.Println("Insert: BST exists, calling UpdateMyBST..")
		return s.UpdateMyBST(ctx, val, bst, 0)
	}

	fmt.Println("Insert: BST not found, creating a new BST")
	// Create a new BST
	root := &TreeNode{Val: val}

	// Generate unique key
	var key string = strconv.FormatUint(rng.Uint64(), 10)
	for {
		fmt.Println("Insert: Key generated: %s", key)
		exists, err := s.MyBSTExists(ctx, key)
		if err != nil {
			fmt.Println("Insert: Error thrown by MyBSTExists")
			return err
		}

		if exists {
			fmt.Println("Insert: Key already exists, generating a new key")
			key = strconv.FormatUint(rng.Uint64(), 10)
		} else {
			fmt.Println("Insert: Unique key generated")
			break
		}
	}

	newbst := MyBST{PrimaryKey: key, Root: root}
	bstJSON, err := json.Marshal(newbst)
	if err != nil {
		fmt.Println("Insert: Error in marshalling MyBST object")
		return err
	}
	fmt.Println("Insert: Returning..")
	return ctx.GetStub().PutState(key, bstJSON)
}

// Delete takes a value to be deleted from BST.
// Returns error if the value doesn't exist in the BST or there is no BST.
func (s *SmartContract) Delete(ctx contractapi.TransactionContextInterface, val int) error {
	fmt.Println("Delete: Started..")
	bst, err := s.ReadMyBST(ctx)
	if err != nil {
		fmt.Println("Delete: Error in reading BST")
		return err
	}
	if bst == nil {
		// No BST in the ledger
		fmt.Println("Delete: No BST created")
		return fmt.Errorf("tree not found")
	}
	fmt.Println("Delete: Calling UpdateMyBST..")
	return s.UpdateMyBST(ctx, val, bst, 1)
}

// Preorder returns the preorder traversal of BST
func (s *SmartContract) Preorder(ctx contractapi.TransactionContextInterface) (string, error) {
	fmt.Println("Preorder: Started...")
	bst, err := s.ReadMyBST(ctx)
	if err != nil {
		fmt.Println("Preorder: Error in reading BST")
		return "", err
	}
	if bst == nil {
		fmt.Println("Preorder: No BST created")
		return "", fmt.Errorf("tree not found")
	}

	var elements []string
	err = s.preorderTraversal(ctx, bst.Root, &elements)
	if err != nil {
		fmt.Println("Preorder: Error thrown by inorderTraversal")
		return "", err
	}
	fmt.Println("Preorder: Traversal (as array): %+q", elements)
	traversal := strings.Join(elements, ",")
	fmt.Println("Preorder: Traversal (as string): %s", traversal)
	return traversal, nil
}

func (s *SmartContract) preorderTraversal(ctx contractapi.TransactionContextInterface, node *TreeNode, elements *[]string) error {
	if node == nil {
		return nil
	}
	(*elements) = append((*elements), strconv.Itoa(node.Val))
	fmt.Println("preorderTraversal: State of traversal till now: %+q", (*elements))

	err := s.preorderTraversal(ctx, node.Left, elements)
	if err != nil {
		return err
	}

	err = s.preorderTraversal(ctx, node.Right, elements)
	if err != nil {
		return err
	}

	return nil
}

// Inorder returns the inorder traversal of BST
func (s *SmartContract) Inorder(ctx contractapi.TransactionContextInterface) (string, error) {
	fmt.Println("Inorder: Started...")
	bst, err := s.ReadMyBST(ctx)
	if err != nil {
		fmt.Println("Inorder: Error in reading BST")
		return "", err
	}
	if bst == nil {
		fmt.Println("Inorder: No BST created")
		return "", fmt.Errorf("tree not found")
	}

	var elements []string
	err = s.inorderTraversal(ctx, bst.Root, &elements)
	if err != nil {
		fmt.Println("Inorder: Error thrown by inorderTraversal")
		return "", err
	}
	fmt.Println("Inorder: Traversal (as array): %+q", elements)
	traversal := strings.Join(elements, ",")
	fmt.Println("Inorder: Traversal (as string): %s", traversal)
	return traversal, nil
}

func (s *SmartContract) inorderTraversal(ctx contractapi.TransactionContextInterface, node *TreeNode, elements *[]string) error {
	if node == nil {
		return nil
	}

	err := s.inorderTraversal(ctx, node.Left, elements)
	if err != nil {
		return err
	}

	(*elements) = append((*elements), strconv.Itoa(node.Val))
	fmt.Println("inorderTraversal: State of traversal till now: %+q", (*elements))

	err = s.inorderTraversal(ctx, node.Right, elements)
	if err != nil {
		return err
	}

	return nil
}

// TreeHeight returns the height of the BST
// Height of empty BST is "0", height of tree with only one node is "1"
func (s *SmartContract) TreeHeight(ctx contractapi.TransactionContextInterface) (string, error) {
	fmt.Println("TreeHeight: Started..")
	bst, err := s.ReadMyBST(ctx)
	if err != nil {
		fmt.Println("TreeHeight: Error in reading BST")
		return "", err
	}
	if bst == nil {
		fmt.Println("TreeHeight: No BST created")
		return "0", fmt.Errorf("tree not found")
	}

	height, err := s.heightOfTree(ctx, bst.Root)
	if err != nil {
		fmt.Println("TreeHeight: Error thrown by heightOfTree")
		return "", err
	}
	fmt.Println("TreeHeight: height of BST is %d", height)
	return strconv.Itoa(height), nil
}

func (s *SmartContract) heightOfTree(ctx contractapi.TransactionContextInterface, node *TreeNode) (int, error) {
	if node == nil {
		return 0, nil
	} else {
		lHeight, err := s.heightOfTree(ctx, node.Left)
		if err != nil {
			return 0, err
		}

		rHeight, err := s.heightOfTree(ctx, node.Right)
		if err != nil {
			return 0, err
		}

		fmt.Println("heightOfTree: Val=%d, lHeight=%d, rHeight=%d", node.Val, lHeight, rHeight)
		if lHeight > rHeight {
			return lHeight + 1, nil
		} else {
			return rHeight + 1, nil
		}
	}
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Println("failed to create chaincode. %v", err)
		return
	}

	err = chaincode.Start()

	if err != nil {
		fmt.Println("failed to start chaincode. %v", err)
		return
	}
}
