package bigone

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/shopspring/decimal"
)

// SpotAccount represents the state of one spot account.
type SpotAccount struct {
	AssetSymbol   string `json:"asset_symbol,omitempty"`
	Balance       string `json:"balance,omitempty"`
	LockedBalance string `json:"locked_balance,omitempty"`
}

// ReadSpotAccounts Balance of all assets
func ReadSpotAccounts() ([]*SpotAccount, error) {
	resp, err := HTTPRequest(context.Background()).Get("/viewer/accounts")

	if err != nil {
		return nil, err
	}

	accounts := []*SpotAccount{}

	if err := UnmarshalResponse(resp, &accounts); err != nil {
		return nil, err
	}

	return accounts, nil
}

// ReadSpotAccount Balance of one asset
func ReadSpotAccount(assetSymbol string) (*SpotAccount, error) {
	path := fmt.Sprintf("/viewer/accounts/%v", strings.ToUpper(assetSymbol))
	resp, err := HTTPRequest(context.Background()).Get(path)

	if err != nil {
		return nil, err
	}

	account := &SpotAccount{}

	if err := UnmarshalResponse(resp, &account); err != nil {
		return nil, err
	}

	return account, nil
}

// FundAccount represents the state of one fund account.
type FundAccount struct {
	AssetSymbol   string `json:"asset_symbol"`
	Balance       string `json:"balance"`
	LockedBalance string `json:"locked_balance"`
}

// ReadFundAccounts Balance of all assets
func ReadFundAccounts() ([]*FundAccount, error) {
	resp, err := HTTPRequest(context.Background()).Get("/viewer/fund/accounts")

	if err != nil {
		return nil, err
	}

	accounts := []*FundAccount{}

	if err := UnmarshalResponse(resp, &accounts); err != nil {
		return nil, err
	}

	return accounts, nil
}

// ReadFundAccount Balance of one asset
func ReadFundAccount(assetSymbol string) (*FundAccount, error) {
	path := fmt.Sprintf("/viewer/fund/accounts/%v", strings.ToUpper(assetSymbol))
	resp, err := HTTPRequest(context.Background()).Get(path)

	if err != nil {
		return nil, err
	}

	account := &FundAccount{}

	if err := UnmarshalResponse(resp, &account); err != nil {
		return nil, err
	}

	return account, nil
}

func SpotBalance(asset string) float32 {
	spotAccount, err := ReadSpotAccount(strings.ToUpper(asset))
	if err != nil {
		log.Fatalf("read spot account err: %v\n", err)
	}

	balance, err := decimal.NewFromString(spotAccount.Balance)
	if err != nil {
		log.Fatalf("balance to decimal err: %v\n", err)
	}

	b, _ := balance.BigFloat().Float32()

	return b
}
