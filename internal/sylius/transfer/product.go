package transfer

import (
	"encoding/json"
)

// ProductEntire is a representation of a product in Sylius.
type ProductEntire struct {
	Code         string                 `json:"code,omitempty"`
	Translations map[string]Translation `json:"translations,omitempty"`
	Images       []Image                `json:"images,omitempty"`
	Enabled      bool                   `json:"enabled"`
}

// Help structure to unmarshal Sylius api response.
type productEntireRaw struct {
	Code         string                 `json:"code,omitempty"`
	Translations map[string]Translation `json:"translations,omitempty"`
	Images       json.RawMessage        `json:"images,omitempty"`
	Enabled      bool                   `json:"enabled"`
}

// UnmarshalJSON helps to fix inconsistency in sylius api response.
// Sylius returns image as a slice or, sometimes, as a map.
func (p *ProductEntire) UnmarshalJSON(value []byte) error {
	raw := &productEntireRaw{}
	if err := json.Unmarshal(value, raw); err != nil {
		return err
	}

	p.Code = raw.Code
	p.Translations = raw.Translations
	p.Enabled = raw.Enabled

	var images []Image
	if err := json.Unmarshal(raw.Images, &images); err == nil {
		p.Images = images
	} else {
		var imageMap map[string]Image
		if err := json.Unmarshal(raw.Images, &imageMap); err == nil {
			images = make([]Image, len(imageMap))
			i := 0
			for _, im := range imageMap {
				images[i] = im
				i++
			}
			p.Images = images
		}
	}

	return nil
}

// Product is a structure to be used in product create/update requests.
type Product struct {
	ProductEntire

	MainTaxon     string   `json:"mainTaxon,omitempty"`
	ProductTaxons string   `json:"productTaxons,omitempty"` // String in which the codes of taxons was written down (separated by comma)
	Channels      []string `json:"channels,omitempty"`
}
