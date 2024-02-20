package entity

type WishList struct {
	WishListId string `json:"wish_list_id"`
	Name       string `json:"name"`
	UID        int    `json:"uid"`
	Alias      string `json:"alias"`
}

type GiftList struct {
	GiftId       int    `json:"gift_id"`
	WishListId   int    `json:"wish_list_id"`
	WishListName string `json:"wish_list_name"`
	Name         string `json:"name"`
	Url          string `json:"url"`
}
