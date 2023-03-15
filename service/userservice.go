package service

import (
	"QQ/models"
	"QQ/utils"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	gin2 "github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// GetUsers
// @Tags getUsers
// @Success 200 {string} json
// @Router /user/getUsers [get]
func GetUsers(context *gin2.Context) {
	context.JSON(200, gin2.H{"message": models.GetUserList()})
}

// CreateUser
// @Tags CreateUser
// @Param name formData string false "用户名"
// @Param pwd formData string false "密码"
// @Success 200 {string} json{"code", "message"}
// @Router /user/createUser [post]
func CreateUser(ct *gin2.Context) {
	name := ct.PostForm("name")
	pwd := ct.PostForm("pwd")
	user := models.UserBasic{
		Name: name,
	}
	if name == "" || pwd == "" {
		ct.JSON(200, gin2.H{
			"code":    -1,
			"message": "用户名或密码为空",
		})
		return
	}
	// 生成一个随机数
	user.Salt = fmt.Sprintf("%d", rand.Int())
	// 生成密码，设置密码
	user.PWD = utils.MakePWS(pwd, user.Salt)
	db := models.CreateUser(&user)
	if db.Error != nil {
		log.Print(db.Error)
		ct.JSON(200, gin2.H{
			"message": "新增用户失败",
		})
		return
	}
	ct.JSON(200, gin2.H{
		"message": "新增用户成功",
	})

}

// DeleteUser
// @Tags DeleteUser
// @Param id query string false "用户ID"
// @Success 200 {string} json{"code", "message"}
// @Router /user/deleteUser [get]
func DeleteUser(ct *gin2.Context) {

	user := models.UserBasic{}

	id, _ := strconv.Atoi(ct.Query("id"))
	fmt.Print(id)
	user.ID = uint(id)
	db := models.DeleteUser(&user)
	if db.Error != nil {
		ct.JSON(200, gin2.H{
			"message": "删除用户失败",
		})
		return
	}
	ct.JSON(200, gin2.H{
		"message": "删除用户成功",
	})
}

// UpdateUser
// @Tags UpdateUser
// @Param id formData string false "用户ID"
// @Param name formData string false "用户名"
// @Success 200 {string} json{"code", "message"}
// @Router /user/updateUser [post]
func UpdateUser(ct *gin2.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(ct.PostForm("id"))
	fmt.Print(id)
	user.ID = uint(id)
	db := models.UpdateUser(&user)
	if db.Error != nil {
		ct.JSON(200, gin2.H{
			"message": "更新用户失败",
		})
		return
	}
	ct.JSON(200, gin2.H{
		"message": "更新用户成功",
	})
}

func Login(name string, pwd string) *models.UserBasic {
	user := models.FindUserByName(name)
	if user.ID == 0 {
		return nil
	}
	if !utils.ValidatePwd(pwd, user.Salt, user.PWD) {
		return nil
	} else {
		return user
	}
}

var AuthMiddleware *jwt.GinJWTMiddleware

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

/**
鉴权工具
*/
func InitAuth() {
	authMiddleware, error := jwt.New(&jwt.GinJWTMiddleware{
		Key:        []byte(viper.GetString("auth.secret")),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		//generate token when call LoginHandler
		Authenticator: func(c *gin2.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password
			// find user
			if Login(userID, password) != nil {

				return nil, nil
			}
			// will call Unauthorized
			return nil, jwt.ErrFailedAuthentication
		},
		// 将token 保存起来
		LoginResponse: func(c *gin2.Context, code int, token string, expire time.Time) {
			var loginVals login
			c.ShouldBind(&loginVals)
			userID := loginVals.Username
			models.UpdateToken(userID, token)
			c.JSON(http.StatusOK, gin2.H{
				"code":   http.StatusOK,
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		},
		//provide an identity for Authorizator
		IdentityHandler: func(c *gin2.Context) interface{} {
			token := c.GetHeader("Token")
			return &models.UserBasic{
				Token: token,
			}
		},
		//verify token
		Authorizator: func(data interface{}, c *gin2.Context) bool {
			if v, ok := data.(*models.UserBasic); ok && models.VerifyToken(v.Token).Name != "" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin2.Context, code int, message string) {
			c.JSON(code, gin2.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Token",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	if error != nil {
		log.Fatal("JWT Error:" + error.Error())
	}

	errInit := authMiddleware.MiddlewareInit()
	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	AuthMiddleware = authMiddleware
}

// 防止跨域请求
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// SendMsg 发送消息，会启动一个websocket
func SendMsg(c *gin2.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Print(err)
		return
	}
	//close ws
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	// 开始订阅Redis
	utils.Subscribe(c, utils.PublishKey)
	//开始读取消息
	go receiveMsg(ws, c)
	//开始订阅消息
	subscribeMsg(ws, c)

}
func receiveMsg(ws *websocket.Conn, ctx *gin2.Context) {
	for {
		// 接收消息
		_, p, err := ws.ReadMessage()
		//发布到redis中
		utils.PublishMsg(ctx, utils.PublishKey, p)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func subscribeMsg(ws *websocket.Conn, ctx *gin2.Context) {
	// for循环不好，应该增加断开的逻辑，如果客户长期断开连接，可以尝试退出for循环
	for {
		msg, err := utils.SubscribeMsg(ctx)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(msg)
		tm := time.Now()
		m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
		//回复消息
		err = ws.WriteMessage(1, []byte(m))
		if err != nil {
			fmt.Println(err)
		}
	}
}
