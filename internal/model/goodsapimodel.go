package model

// send this when frontend requests goods
type GetGoodsResponse struct {
	GoodID      uint                    `json:"good_id"`
	Title       string                  `json:"title"`
	Images      []string                `json:"images"`
	Price       float64                 `json:"price"`
	Views       uint                    `json:"views"`
	Favorites   uint                    `json:"favorites"`
	Description string                  `json:"description"`
	User        UserForGetGoodsResponse `json:"user"`
}
type UserForGetGoodsResponse struct {
	UserID   uint    `json:"user_id"`
	Name     string  `json:"name"`
	Avatar   string  `json:"avatar"`
	Rating   float32 `json:"rating"`
}

// receive this when frontend posts or updates goods
type PostGoodsReceive struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
	Price       float64  `json:"price"`
	Tags        []string `json:"tags"`
	SellerID    uint     `json:"seller_id"`
	IsInvisible bool     `json:"is_invisible"`
	IsDeleted   bool     `json:"is_deleted"`
	IsBought    bool     `json:"is_bought"`
}

// send this when frontend posts or updates goods
type PostGoodsResponse struct {
	Message  string `json:"message"`
	GoodInfo Good   `json:"good_info"`
}

// receive this when frontend deletes goods
type DeleteGoodsReceive struct {
	ID       uint `json:"id"`
	SellerID uint `json:"seller_id"`
}

// send this when frontend deletes goods
type DeleteGoodsResponse struct {
	Message string `json:"message"`
}
