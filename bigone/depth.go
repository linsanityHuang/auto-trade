package bigone

import (
	"context"
	"fmt"
)

// Depth Depth
type Depth struct {
	AssetPairName string       `json:"asset_pair_name,omitempty"`
	Bids          []PriceLevel `json:"bids,omitempty"`
	Asks          []PriceLevel `json:"asks,omitempty"`
}

// ReadDepth Order Book is the ask orders and bid orders collection of a asset pair
func ReadDepth(assetPairName string) (*Depth, error) {

	resp, err := HTTPRequest(context.Background()).SetQueryParam("limit", "50").Get(fmt.Sprintf("/asset_pairs/%v/depth", assetPairName))

	if err != nil {
		return nil, err
	}

	depths := &Depth{}

	if err := UnmarshalResponse(resp, depths); err != nil {
		return nil, err
	}

	return depths, nil
}
