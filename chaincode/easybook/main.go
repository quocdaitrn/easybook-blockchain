package main

import (
	"log"

	"bitbucket.org/quocdaitrn/hotel-rating/chaincode/hotel-rating/chaincode"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	hotelChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating hotel-rating chaincode: %v", err)
	}

	if err := hotelChaincode.Start(); err != nil {
		log.Panicf("Error starting hotel-rating chaincode: %v", err)
	}
}
