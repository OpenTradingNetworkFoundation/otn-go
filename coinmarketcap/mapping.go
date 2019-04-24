package coinmarketcap

func MapListingsBySymbol(listings []*Listing) map[string]*Listing {
	result := make(map[string]*Listing)
	for _, l := range listings {
		result[l.Symbol] = l
	}
	return result
}

func MapListingsByID(listings []*Listing) map[int]*Listing {
	result := make(map[int]*Listing)
	for _, l := range listings {
		result[l.ID] = l
	}
	return result
}
