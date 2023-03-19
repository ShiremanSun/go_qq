package service

import (
	"QQ/models"
	"QQ/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/fatih/set"
	gin2 "github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"sync"
)

const PublishKey = "websocket"

var RedisSubscriber *redis.PubSub

// PublishMsg 发布消息到redis
func PublishMsg(ctx context.Context, channel string, msg []byte) error {
	fmt.Print("publish%s", msg)
	return utils.Redis.Publish(ctx, channel, msg).Err()
}

// Subscribe 订阅Redis
func Subscribe(ctx context.Context, channel string) {
	RedisSubscriber = utils.Redis.Subscribe(ctx, channel)
}

// SubscribeMsg 订阅消息
func SubscribeMsg(ctx context.Context) (string, error) {
	msg, error := RedisSubscriber.ReceiveMessage(ctx)
	if msg != nil {
		fmt.Print("Subscribe%s", msg.Payload)
		return msg.Payload, error
	}
	return "", error
}

// 防止跨域请求
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 客户端
var clientMap map[int64]*models.Node = make(map[int64]*models.Node, 0)

// 读写锁
var rwLocker = sync.RWMutex{}

// Connect user与服务端建立链接，用来接收消息
func Connect(c *gin2.Context) {
	con, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Print(err)
		return
	}
	userId, _ := strconv.ParseInt(c.Query("userId"), 10, 64)
	/*targetId := c.Query("targetId")
	messageType := c.Query("messageType")
	context := c.Query("context")*/
	node := &models.Node{
		Con:       con,
		DataQueue: make(chan []byte, 50),
		GroupSet:  set.New(set.ThreadSafe),
	}
	//3. 用户关系
	//4. userId 跟node绑定
	// 保证线程安全
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()
	//close ws
	//5. 开始订阅Redis
	Subscribe(c, PublishKey)
	//6. 开始读取消息
	go receiveMsg(node, c)
	//7. 开始订阅消息
	go subscribeMsg(node, c)
	go sendProc()
	con.WriteMessage(websocket.TextMessage, []byte("欢迎来到聊天室"))

}

// 使用ws接收消息，但是用户一般会使用http来发送消息
func receiveMsg(node *models.Node, ctx *gin2.Context) {
	for {
		// 接收消息
		_, p, err := node.Con.ReadMessage()
		//发布到redis中
		if err != nil {
			fmt.Println(err)
		}
		PublishMsg(ctx, PublishKey, p)
	}
}

func subscribeMsg(node *models.Node, ctx *gin2.Context) {
	// for循环不好，应该增加断开的逻辑，如果客户长期断开连接，可以尝试退出for循环
	for {
		msg, err := SubscribeMsg(ctx)
		if err != nil {
			fmt.Println(err)
			return
		}
		// 使用channel，将消息订阅阻塞住
		sendChan <- []byte(msg)
	}

}

// 定义一个channel
var sendChan chan []byte = make(chan []byte, 1024)

// 收到消息之后
func broadMsg(data []byte) {
	sendChan <- data
}

func init() {
	go udpSendProc()
	go udpRecvProc()
}
func udpRecvProc() {

}

// 完成udp发送服务
func udpSendProc() {

}

func dispatch(data []byte) {
	/*msg := &models.Message{}
	ere := json.Unmarshal(data, msg)*/
}

//给user发送消息
func sendProc() {
	for {
		select {
		//发送消息给user
		//协程的in
		case data := <-sendChan:
			//解析data
			msg := &models.Message{}
			err := json.Unmarshal(data, msg)
			fmt.Println(msg)
			// 解析到要发送给谁
			//err = node.Con.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

}
