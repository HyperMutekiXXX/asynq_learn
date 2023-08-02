package item

type Item struct {
	Id    int32   `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Num   int64   `json:"num"`
}
