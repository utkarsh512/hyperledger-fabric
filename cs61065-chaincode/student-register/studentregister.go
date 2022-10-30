// CS61065 - Assignment 4 - Part A
//
// Authors:
// Utkarsh Patel (18EC35034)
// Saransh Patel (18CS30039)

package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

// ReadStudent returns the Student's name corresponding to roll
func (s *SmartContract) ReadStudent(ctx contractapi.TransactionContextInterface, roll string) (string, error) {
	val, err := ctx.GetStub().GetState(roll)

	if err != nil {
		return "", fmt.Errorf("failed to query world state. %v", err)
	}

	return string(val), nil
}

// StudentExists return true if passed roll exists in the contract, otherwise false
func (s *SmartContract) StudentExists(ctx contractapi.TransactionContextInterface, roll string) (bool, error) {
	val, err := ctx.GetStub().GetState(roll)

	if err != nil {
		return true, fmt.Errorf("failed to query world state. %v", err)
	}

	if val != nil {
		return true, nil
	}

	return false, nil
}

// CreateStudent inserts roll-name pair to contract if roll doesn't exist in the contract
func (s *SmartContract) CreateStudent(ctx contractapi.TransactionContextInterface, roll string, name string) error {
	exists, err := s.StudentExists(ctx, roll)

	if err != nil {
		return fmt.Errorf("failed to query world state. %v", err)
	}

	if exists {
		return fmt.Errorf("given roll number already exists in the contract")
	} else {
		return ctx.GetStub().PutState(roll, []byte(name))
	}
}

// ReadAllStudents iterates over the ledger and returns all roll-name pair
func (s *SmartContract) ReadAllStudents(ctx contractapi.TransactionContextInterface) (string, error) {
	// Get an iterator over <roll, name> pair
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return "", err
	}
	defer resultsIterator.Close()

	// Iterate over the contract and store <roll, name> pair
	var result []string
	for resultsIterator.HasNext() {
		res, err := resultsIterator.Next()
		if err != nil {
			return "", err
		}

		var kvpair string = string(res.Key) + ":" + string(res.Value)
		result = append(result, kvpair)
	}

	// Marshal the map into a JSON object and convert that into string for return value
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(resultJSON), nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("failed to create chaincode. %v", err)
		return
	}

	err = chaincode.Start()

	if err != nil {
		fmt.Printf("failed to start chaincode. %v", err)
		return
	}
}
