package utils

import (
	"backend/internal/model"
)

type GoodTransform struct{}
type UserTransform struct{}
type CommentTransform struct{}

func (gt *GoodTransform) FindGoodsByIdDb2ResponseModel(dbModel model.Good, theUser model.User) model.GetGoodsResponse {
	return model.GetGoodsResponse{
		GoodID:      dbModel.ID,
		Title:       dbModel.Title,
		Description: dbModel.Description,
		Images:      dbModel.Images,
		Price:       dbModel.Price,
		Views:       dbModel.Views,
		Favorites:   uint(len(dbModel.FavoUsers)),
		User: model.UserForGetGoodsResponse{
			UserID: theUser.ID,
			Name:   theUser.Name,
			Avatar: theUser.Avatar,
			Rating: theUser.Rating,
		},
	}
}

func (gt *GoodTransform) Post2DbModel(apiModel model.PostGoodsReceive) model.Good {
	return model.Good{
		Title:       apiModel.Title,
		Description: apiModel.Description,
		Images:      apiModel.Images,
		Price:       apiModel.Price,
		Tags:        apiModel.Tags,
		SellerID:    apiModel.SellerID,
	}
}
