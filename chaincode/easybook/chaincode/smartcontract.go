package chaincode

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing hotel rating
type SmartContract struct {
	contractapi.Contract
}

// Hotel stores rating of a particular hotel
type Hotel struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	IsActive bool    `json:"isActive"`
	Rating   float32 `json:"rating"`
}

// InitLedger adds a base set of hotels to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	hotels := []Hotel{
		{ID: "hotel1", Name: "Venice", IsActive: true, Rating: 5.0},
		{ID: "hotel2", Name: "Milan", IsActive: true, Rating: 4.5},
		{ID: "hotel3", Name: "Roma", IsActive: true, Rating: 4.2},
	}

	for _, h := range hotels {
		hotelJson, err := json.Marshal(h)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(h.ID, hotelJson)
		if err != nil {
			return fmt.Errorf("failed to put to world state: %v", err)
		}
	}

	return nil
}

// CreateHotel issues a new hotel to the world state with given details.
func (s *SmartContract) CreateHotel(ctx contractapi.TransactionContextInterface, id string, name string, isActive bool, rating float32) error {
	exist, err := s.HotelExists(ctx, id)
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("the hotel %s already exists", id)
	}

	hotel := Hotel{
		ID:       id,
		Name:     name,
		IsActive: isActive,
		Rating:   rating,
	}

	hotelJson, err := json.Marshal(hotel)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(id, hotelJson)
	if err != nil {
		return fmt.Errorf("can not create hotel %s", id)
	}
	return nil
}

// ReadHotel returns the hotel stored in the world state with given id.
func (s *SmartContract) ReadHotel(ctx contractapi.TransactionContextInterface, id string) (*Hotel, error) {
	hotelJson, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, err
	}
	if hotelJson == nil {
		return nil, fmt.Errorf("the hotel %s does not exist", id)
	}

	var hotel Hotel
	err = json.Unmarshal(hotelJson, &hotel)
	if err != nil {
		return nil, err
	}

	return &hotel, nil
}

// UpdateHotel updates an existing hotel in the world state with provided parameters.
func (s *SmartContract) UpdateHotel(ctx contractapi.TransactionContextInterface, id string, name string, isActive bool, rating float32) error {
	exist, err := s.HotelExists(ctx, id)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("the hotel %s does not exist", id)
	}

	hotel := &Hotel{
		ID:       id,
		Name:     name,
		IsActive: isActive,
		Rating:   rating,
	}

	hotelJson, err := json.Marshal(hotel)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(id, hotelJson)
	if err != nil {
		return fmt.Errorf("can not update information for the hotel %s", id)
	}

	return nil
}

// DeleteHotel deletes an given hotel from the world state.
func (s *SmartContract) DeleteHotel(ctx contractapi.TransactionContextInterface, id string) error {
	exist, err := s.HotelExists(ctx, id)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("the hotel %s does not exist", id)
	}

	err = ctx.GetStub().DelState(id)
	if err != nil {
		return fmt.Errorf("can not delete the hotel %s", id)
	}

	return nil
}

// HotelExists returns true when hotel with given ID exists in world state
func (s *SmartContract) HotelExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	exist, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return exist != nil, nil
}

// GetAllHotels returns all hotels found in world state
func (s *SmartContract) GetAllHotels(ctx contractapi.TransactionContextInterface) ([]*Hotel, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}

	defer resultsIterator.Close()

	var hotels []*Hotel
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var hotel Hotel
		err = json.Unmarshal(queryResponse.Value, &hotel)
		if err != nil {
			return nil, err
		}
		hotels = append(hotels, &hotel)
	}

	return hotels, nil
}
