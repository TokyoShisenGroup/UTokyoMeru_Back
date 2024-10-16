package router

import (
	"backend/internal/db"
	"backend/internal/logic"
	mw "backend/internal/middlewares"
	"backend/internal/model"
	"backend/internal/utils"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func VerificationHandler(c *gin.Context) {
	//获取参数与参数校验
	var p model.ParamVerify
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误
		fmt.Println(err)
		c.JSON(400, gin.H{"message": "Invalid param", "error": err})
		return
	}
	to := p.MailAddress                 // 收件人邮箱
	code := utils.GenerateRandomCode(6) // 验证码

	err := utils.SendEmail(to, code)
	if err != nil {
		c.JSON(500, gin.H{"message": "sending code failed", "error": err})
	}

	//将验证码存入redis
	err = db.SetVerificationCode(p.MailAddress+p.VerificationCodeType, code)
	if err != nil {
		c.JSON(500, gin.H{"message": "storaging code failed"})
		return
	}

	//返回响应
	c.JSON(200, gin.H{"message": "sending code success"})
}

func SignUpHandler(c *gin.Context) {
	//获取参数与参数校验
	var p model.ParamSignup
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误
		c.JSON(400, gin.H{"message": "Invalid param", "error": err})
		return
	}
	// 从Redis获取存储的验证码并比对获取的验证码
	storedCode, err := db.GetVerificationCode(p.MailAddress + p.VerificationCodeType)
	if err == redis.Nil {
		c.JSON(400, gin.H{"error": "Verificationcode expired or not exist."})
		return
	} else if err != nil {
		c.JSON(400, gin.H{"error": "getting verificationcode failed"})
		return
	}

	// 验证码比对
	if p.VerificationCode != storedCode {
		c.JSON(400, gin.H{"error": "VerificationCode error"})
		return
	}

	//业务处理
	user, err := logic.SignUp(&p)
	if err != nil {
		if strings.Contains(err.Error(), "23505") {
			c.JSON(400, gin.H{"message": "User exist", "error": err})
			return
		}
		c.JSON(500, gin.H{"message": "Server busy", "error": err})
		return
	}
	address := strings.Split(user.Address, " ")
	response := model.UserInfoResponse{
		ID:          user.ID,
		UserName:    user.Name,
		MailAddress: user.MailAddress,
		Gender:      user.Gender,
		Birthday:    user.Birthday,
		PhoneNumber: user.PhoneNumber,
	}

	if len(address) == 4 {
		response.Address = model.Address{
			PostalCode:    address[0],
			Prefecture:    address[1],
			City:          address[2],
			AddressDetail: address[3],
		}
	}
	//返回响应
	c.JSON(200, response)
}

func LoginHandler(c *gin.Context) {
	//1.获取参数与参数校验
	var p model.ParamLogin
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误
		c.JSON(400, gin.H{"message": "Invalid param", "error": err})
		return
	}
	fmt.Println(p)

	//2.业务处理
	user, err := logic.Login(&p)
	if err != nil {
		if errors.Is(err, errors.New("user not exist")) {
			c.JSON(404, gin.H{"message": "User not exist", "error": err})
			return
		}
		c.JSON(500, gin.H{"message": "Invalid password", "error": err})
		return
	}

	address := strings.Split(user.Address, " ")
	response := model.UserInfoResponse{
		ID:          user.ID,
		UserName:    user.Name,
		MailAddress: user.MailAddress,
		Gender:      user.Gender,
		Birthday:    user.Birthday,
		PhoneNumber: user.PhoneNumber,
		Token:       user.Token,
	}

	if len(address) == 4 {
		response.Address = model.Address{
			PostalCode:    address[0],
			Prefecture:    address[1],
			City:          address[2],
			AddressDetail: address[3],
		}
	}

	//返回响应
	c.JSON(200, response)
}

func LoginByCodeHandler(c *gin.Context) {
	//获取参数与参数校验
	var p model.ParamLoginByCode
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误
		c.JSON(400, gin.H{"message": "Invalid param", "error": err})
		return
	}
	fmt.Println(p)

	// 从Redis获取存储的验证码并比对获取的验证码
	storedCode, err := db.GetVerificationCode(p.MailAddress + p.VerificationCodeType)
	if err == redis.Nil {
		c.JSON(400, gin.H{"message": "Verificationcode expired or not exist.", "error": err})
		return
	} else if err != nil {
		c.JSON(500, gin.H{"message": "Failed to get verification code", "error": err})
		return
	}

	// 验证码比对
	if p.VerificationCode != storedCode {
		c.JSON(400, gin.H{"message": "Verification code error", "error": "VerificationCode error"})
		return
	}

	//业务处理
	user, err := logic.LoginByCode(&p)
	if err != nil {
		if errors.Is(err, errors.New("user not exist")) {
			c.JSON(400, gin.H{"message": "User not exist", "error": err})
			return
		}
		c.JSON(400, gin.H{"message": "Invalid password", "error": err})
		return
	}

	address := strings.Split(user.Address, " ")
	response := model.UserInfoResponse{
		ID:          user.ID,
		UserName:    user.Name,
		MailAddress: user.MailAddress,
		Gender:      user.Gender,
		Birthday:    user.Birthday,
		PhoneNumber: user.PhoneNumber,
		Token:       user.Token,
	}

	if len(address) == 4 {
		response.Address = model.Address{
			PostalCode:    address[0],
			Prefecture:    address[1],
			City:          address[2],
			AddressDetail: address[3],
		}
	}

	//返回响应
	c.JSON(200, response)
}

func ResetPasswordHandler(c *gin.Context) {
	//获取参数与参数校验
	var p model.ParamResetPassword
	if err := c.ShouldBindJSON(&p); err != nil {
		//请求参数有误
		c.JSON(400, gin.H{"message": "Invalid param", "error": err})
		return
	}
	fmt.Println(p)

	// 从Redis获取存储的验证码并比对获取的验证码
	storedCode, err := db.GetVerificationCode(p.MailAddress + "reset")
	if err == redis.Nil {
		c.JSON(404, gin.H{"message": "Verificationcode expired or not exist.", "error": err})
		return
	} else if err != nil {
		c.JSON(500, gin.H{"message": "Failed to get verification code", "error": err})
		return
	}

	// 验证码比对
	if p.VerificationCode != storedCode {
		c.JSON(400, gin.H{"message": "VerificationCode error", "error": errors.New("VerificationCode error")})
		return
	}

	//业务处理
	if err := logic.ResetPassword(&p); err != nil {
		if errors.Is(err, errors.New("user not exist")) {
			c.JSON(404, gin.H{"message": "User not exist", "error": err})
			return
		}
		c.JSON(400, gin.H{"message": "Invalid password", "error": err})
		return
	}

	//3.返回响应
	c.JSON(200, gin.H{"message": "reset password"})
}

//需要userid因为可能需要看其他用户的出售商品
func GetAllSalesGoodsHandler(c *gin.Context) {
	crud := db.UsersCRUD{}
	gt := utils.GoodTransform{}
	//获取用户ID
	userID, err := strconv.ParseUint(c.Query("user_id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"message": "Failed to get params", "error": err})
		return
	}
	goods, err := crud.FindGoodsByFK(uint(userID), "Sales")
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to get all sales goods", "error": err})
		return
	}
	response := []model.GetGoodsResponse{}
	for _, good := range goods {
		response = append(response, gt.FindGoodsByIdDb2ResponseModel(good, good.Seller))
	}

	c.JSON(200, response)
}

//需要userid因为可能需要看其他用户的
func GetAllSellingGoodsHandler(c *gin.Context) {
	crud := db.UsersCRUD{}
	gt := utils.GoodTransform{}
	//获取用户ID
	userID, err := strconv.ParseUint(c.Query("user_id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"message": "Failed to get params", "error": err})
		return
	}
	goods, err := crud.FindGoodsByFK(uint(userID), "Sales")
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to get all sales goods", "error": err})
		return
	}
	response := []model.GetGoodsResponse{}
	for _, good := range goods {
		if !good.IsBought {
			response = append(response, gt.FindGoodsByIdDb2ResponseModel(good, good.Seller))
		}
	}

	c.JSON(200, response)
}

func GetAllLikedGoodsHandler(c *gin.Context) {
	crud := db.UsersCRUD{}
	gt := utils.GoodTransform{}
	//获取用户ID
	mailAddressInterface, _ := c.Get(mw.ContextUserIDKey)
	mailAddress, _ := mailAddressInterface.(string)

	user, err := crud.FindOneByUniqueField("mail_address", mailAddress)
	if err != nil {
		c.JSON(400, gin.H{"message": "User not exist", "error": err})
		return
	}
	goods := user.FavoList
	response := []model.GetGoodsResponse{}
	for _, good := range goods {
		response = append(response, gt.FindGoodsByIdDb2ResponseModel(good, good.Seller))
	}

	c.JSON(200, response)
}

func GetAllBoughtGoodsHandler(c *gin.Context) {
	crud := db.UsersCRUD{}
	gt := utils.GoodTransform{}
	//获取用户ID
	mailAddressInterface, _ := c.Get(mw.ContextUserIDKey)
	mailAddress, _ := mailAddressInterface.(string)

	user, err := crud.FindOneByUniqueField("mail_address", mailAddress)
	if err != nil {
		c.JSON(400, gin.H{"message": "User not exist", "error": err})
		return
	}
	goods := user.Buys
	response := []model.GetGoodsResponse{}
	for _, good := range goods {
		response = append(response, gt.FindGoodsByIdDb2ResponseModel(good, good.Seller))
	}

	c.JSON(200, response)
}

func GetAllSoldGoodsHandler(c *gin.Context) {
	crud := db.UsersCRUD{}
	gt := utils.GoodTransform{}
	//获取用户ID
	mailAddressInterface, _ := c.Get(mw.ContextUserIDKey)
	mailAddress, _ := mailAddressInterface.(string)

	user, err := crud.FindOneByUniqueField("mail_address", mailAddress)
	if err != nil {
		c.JSON(400, gin.H{"message": "User not exist", "error": err})
		return
	}
	goods := user.Sales
	response := []model.GetGoodsResponse{}
	for _, good := range goods {
		if good.IsBought {
			response = append(response, gt.FindGoodsByIdDb2ResponseModel(good, good.Seller))
		}
	}

	c.JSON(200, response)
}


func GetAllGoodsStatsHandler(c *gin.Context) {
	crud := db.UsersCRUD{}
	//获取用户ID
	mailAddressInterface, _ := c.Get(mw.ContextUserIDKey)
	mailAddress, _ := mailAddressInterface.(string)

	user, err := crud.FindOneByUniqueField("mail_address", mailAddress)
	if err != nil {
		c.JSON(400, gin.H{"message": "User not exist", "error": err})
		return
	}

	goods, err := crud.FindAllGoodsFK(user.ID)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to get all goods", "error": err})
		return
	}

	var lengths [4]int
	for i := 0; i < 3; i++ {
		lengths[i] = len(goods[i])
	}

	isBoughtCount := 0
	for _, sale := range goods[0] {
		if sale.IsBought {
			isBoughtCount++
		}
	}
	lengths[3] = isBoughtCount

	c.JSON(200, gin.H{
		"sale_number":  lengths[0] - lengths[3],
		"sold_number":  lengths[3],
		"buy_number":   lengths[1],
		"favor_number": lengths[2],
	})
}
