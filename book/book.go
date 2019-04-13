package book

type Book struct {
	ID     uint32 `query:"id"`
	Name   string `json:"name" query:"name"`
	Author string `json:"author" query:"author"`
	Qty    uint32 `json:"qty" query:"qty"`
}
