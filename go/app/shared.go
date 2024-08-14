package app

type OrderInput struct {
	OrderId string
	Address string
}

type OrderOutput struct {
	TrackingId string `json:"trackingId"`
	Address    string `json:"address"`
}

type Items []Item

type Item struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}

// Item sort methods
func (p Items) Len() int {
	return len(p)
}

func (p Items) Less(i, j int) bool {
	return p[i].Id < p[j].Id
}

func (p Items) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
