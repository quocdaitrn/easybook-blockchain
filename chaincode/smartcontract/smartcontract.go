package smartcontract

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing hotel rating
type SmartContract struct {
	contractapi.Contract
}

// Hotel stores information of a hotel
type Hotel struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	IsActive      bool            `json:"isActive"`
	Rating        float32         `json:"rating"`
	ServiceLevels []*ServiceLevel `json:"serviceLevels"`
}

// ServiceLevel stores information of a service level in a hotel
type ServiceLevel struct {
	ID               string       `json:"id"`
	Name             string       `json:"name"`
	IsUsed           bool         `json:"isUsed"`
	SatisfactionRate float32      `json:"satisfactionRate"`
	RuleAbidingRate  float32      `json:"ruleAbidingRate"`
	HotelID          string       `json:"hotelId"`
	Agreements       []*Agreement `json:"agreements"`
}

// Agreement stores information of a agreement of a level in a hotel
type Agreement struct {
	ID                          string `json:"id"`
	IsApplied                   bool   `json:"isApplied"`
	TotalFeedbacks              uint   `json:"totalFeedbacks"`
	TotalUnfulfilledCommitments uint   `json:"totalUnfulfilledCommitments"`
	IsAppliedPenalty            bool   `json:"isAppliedPenalty"`
	TotalCompensations          uint   `json:"totalCompensations"`
	TotalNoCompensations        uint   `json:"totalNoCompensations"`
	ServiceLevelID              string `json:"serviceLevelId"`
	HotelID                     string `json:"hotelId"`
}

// InitLedger adds a base set of hotels to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	hotels := []Hotel{
		{
			ID:       "1",
			Name:     "Rex Hotel",
			IsActive: true,
			Rating:   8.4,
			ServiceLevels: []*ServiceLevel{
				{
					ID:               "1",
					Name:             "Standard",
					IsUsed:           true,
					SatisfactionRate: 0.84,
					RuleAbidingRate:  0.0,
					HotelID:          "1",
					Agreements: []*Agreement{
						{
							ID:                          "1",
							IsApplied:                   true,
							TotalFeedbacks:              100,
							TotalUnfulfilledCommitments: 16,
							IsAppliedPenalty:            false,
						},
						{
							ID:                          "2",
							IsApplied:                   true,
							TotalFeedbacks:              1000,
							TotalUnfulfilledCommitments: 100,
							IsAppliedPenalty:            false,
						},
					},
				},
			},
		},
		// {ID: "2", Name: "Mia Saigon", IsActive: true, Rating: 9.2},
		// {ID: "3", Name: "The Myst Dong Khoi", IsActive: true, Rating: 8.8},
		// {ID: "4", Name: "Nikko Saigon", IsActive: true, Rating: 9.1},
	}

	for _, h := range hotels {
		hotelJSON, err := json.Marshal(h)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(h.ID, hotelJSON)
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

	hotelJSON, err := json.Marshal(hotel)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(id, hotelJSON)
	if err != nil {
		return fmt.Errorf("can not create hotel %s", id)
	}
	return nil
}

// ReadHotel returns the hotel stored in the world state with given id.
func (s *SmartContract) ReadHotel(ctx contractapi.TransactionContextInterface, id string) (*Hotel, error) {
	hotelJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, err
	}
	if hotelJSON == nil {
		return nil, fmt.Errorf("the hotel %s does not exist", id)
	}

	var hotel Hotel
	err = json.Unmarshal(hotelJSON, &hotel)
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

	hotelJSON, err := json.Marshal(hotel)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(id, hotelJSON)
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
