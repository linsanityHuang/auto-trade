package bigone

import (
	"context"
	"fmt"
	"strings"
)

// Depth Depth
type Depth struct {
	AssetPairName string       `json:"asset_pair_name,omitempty"`
	Bids          []PriceLevel `json:"bids,omitempty"`
	Asks          []PriceLevel `json:"asks,omitempty"`
}

// ReadDepth Order Book is the ask orders and bid orders collection of a asset pair
func ReadDepth(assetPairName string) (*Depth, error) {
	path := fmt.Sprintf("/asset_pairs/%v/depth", strings.ToUpper(assetPairName))

	resp, err := HTTPRequest(context.Background()).SetQueryParam("limit", "50").Get(path)

	if err != nil {
		return nil, err
	}

	depths := &Depth{}

	if err := UnmarshalResponse(resp, depths); err != nil {
		return nil, err
	}

	return depths, nil
}
