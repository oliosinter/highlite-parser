package transfer

import (
	"highlite2-import/internal/math"
)

// VariantEntire is a representation of a product variant in Sylius.
type VariantEntire struct {
	Code          string                  `json:"code,omitempty"`
	Translations  map[string]Translation  `json:"translations,omitempty"`
	ChannelPrices map[string]ChannelPrice `json:"channelPricings,omitempty"`
}

// Variant is a structure to be used in new product variant request.
type Variant struct {
	VariantEntire
}

// VariantsEqual checks if variants are equal
func VariantsEqual(e VariantEntire, v Variant) bool {
	if e.Code != v.Code {
		return false
	}

	// checking prices
	if len(e.ChannelPrices) != len(v.ChannelPrices) {
		return false
	}

	for key, price1 := range e.ChannelPrices {
		price2 := v.ChannelPrices[key]
		if int(price1.Price) != int(math.Round(price2.Price*100)) {
			return false
		}
	}

	return true
}
