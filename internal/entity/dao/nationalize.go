package dao

type Nationalize struct {
	Count   int    `json:"count"`
	Name    string `json:"name"`
	Country []struct {
		CountryId   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	} `json:"country"`
}
