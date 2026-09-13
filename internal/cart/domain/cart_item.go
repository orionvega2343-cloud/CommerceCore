package domain

type CartItem struct {
	Id            int `db:"id"`
	CartId        int `db:"cart_id"`
	ProductId     int `db:"product_id"`
	Quantity      int `db:"quantity"`
	PriceSnapshot int `db:"price_snapshot"`
}
