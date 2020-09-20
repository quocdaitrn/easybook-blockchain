package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/quocdaitrn/easybook-blockchain/chaincode/smartcontract"
)

func main() {
	slaChaincode, err := contractapi.NewChaincode(&smartcontract.SmartContract{})
	if err != nil {
		log.Panicf("Error creating SLA chaincode: %v", err)
	}

	if err := slaChaincode.Start(); err != nil {
		log.Panicf("Error starting SLA chaincode: %v", err)
	}
}
