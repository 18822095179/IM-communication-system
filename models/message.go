// 消息体

package models

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
	"net"
	"net/http"
	"strconv"
	"sync"
)

type Message struct {
	gorm.Model
	FormId   int64  // 发送者
	TargetId int64  // 接受者
	Type     int    // 发送类型，1私聊 2群聊 3广播
	Media    int    // 消息类型，1文字 2表情包 3图片 4音频
	Content  string // 消息内容
	Pic      string
	Url      string
	Desc     string
	Amout    int // 其他数字统计
}

func (table *Message) TableName() string {
	return "message"
}

// Node WebSocket连接节点，包括连接对象、数据队列、群组集合。
type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// 映射关系 键是用户ID，值是Node
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

// Chat 用于处理WebSocket连接请求，主要功能包括参数获取、WebSocket升级、用户关系处理、节点绑定以及发送和接收逻辑处理。
// 需要：发送者ID，接收者ID，消息类型，发送的内容，发送类型
func Chat(writer http.ResponseWriter, request *http.Request) {
	//1.获取参数并校验token等合法性
	//token := query.Get("token")
	query := request.URL.Query()
	Id := query.Get("userId")
	userId, _ := strconv.ParseInt(Id, 10, 64)
	//msgType := query.Get("msgType")
	//targetId := query.Get("targetId")
	//context := query.Get("context")
	isvalida := true // checkToken() token校验先写死（为true），之后再写

	conn, err := (&websocket.Upgrader{ // 使用websocket.Upgrader升级HTTP连接为WebSocket连接
		// token 校验
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	//2.获取conn(连接)
	node := &Node{
		Conn:      conn, // 将HTTP升级成WebSocket连接的conn给node.Conn
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}

	//3.用户关系

	//4.userid跟node绑定并加锁
	rwLocker.Lock()
	clientMap[userId] = node // 将发送者的userID与node关联起来
	rwLocker.Unlock()

	//5.完成发送逻辑
	go sendProc(node)

	//6.完成接收逻辑
	go recvProc(node)

	sendMsg(userId, []byte("欢迎进入聊天系统"))
}

// sendProc 是发送数据的协程逻辑，不断从节点的数据队列中取出数据并发送到对应的WebSocket连接
func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue: // 从node.DataQueue中取数据放入WebSocket连接中
			fmt.Println("[ws] sendProc >>> msg :", string(data))
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

// recvProc 是接收数据的协程逻辑，不断从WebSocket连接中读取消息数据，并广播给所有连接
func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage() // 持续监听WebSocket中的信息
		if err != nil {
			fmt.Println(err)
			return
		}
		dispatch(data) // 后端调度逻辑处理，其中包括私信、群发等方式。
		broadMsg(data) // 广播给所有节点
		fmt.Println("[ws] recvProc <<<<", string(data))
	}
}

var udpsendChan chan []byte = make(chan []byte, 1024) // 是一个缓冲通道，用于发送UDP数据包的通道

func broadMsg(data []byte) { // 用于广播消息，将数据发送到udpsendChan通道中
	udpsendChan <- data
}

func init() {
	go updSendProc()
	go udpRecvProc()
	fmt.Println("init goroutine ")
}

// 完成udp数据发送协程 updSendProc函数是处理UDP数据发送的协程逻辑，与指定的IP和端口建立UDP连接，并从udpsendChan通道中不断获取数据进行发送
func updSendProc() {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{ // 与指定IP和端口进行连接
		IP:   net.IPv4(192, 168, 0, 255), // cmd中输入ipconfig，这是其中IPV4的地址
		Port: 3000,
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
	}
	for {
		select {
		case data := <-udpsendChan: // 接收udpsendChan中的数据并发送给已经建立的UDP连接（即con中）
			fmt.Println("updSendProc data:", string(data))
			_, err := con.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

// 完成udp数据接收协程 udpRecvProc()函数是处理UDP数据接收的协程逻辑，通过指定IP和端口进行UDP监听，不断接收数据并进行处理
func udpRecvProc() {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	if err != nil {
		fmt.Println(err)
	}
	defer con.Close()

	for {
		var buf [512]byte
		n, err := con.Read(buf[0:]) //UDP连接中（即con中）连续读取数据
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("udpRecvProc data :", string(buf[0:n]))
		dispatch(buf[0:n]) // 调用dispatch函数处理收到的数据
	}
}

// 后端调度逻辑处理 dispatch()函数是后端调度逻辑的处理函数，它从接收到的数据中解析出Message对象，并根据消息类型进行相应的处理
func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg) // 将接收到的data切片解析成Message对象
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case 1: // 私信
		fmt.Println("dispatch data:", string(data))
		sendMsg(msg.TargetId, data)
		//case 2: // 群发
		//	sendGroupMsg()
		//case 3: // 广播
		//	sendAllMsg()
		//case 4:
		//

	}
}

// sendMsg()函数用于向指定用户发送消息，通过获取用户ID对应的节点对象，将消息发送到该节点的数据队列中
func sendMsg(userId int64, msg []byte) {
	fmt.Println("sendMsg >>>> userID: ", userId, " msg:", string(msg))
	rwLocker.RLock()
	node, ok := clientMap[userId] // 将信息放入指定UserID对应node的DataQueue中
	rwLocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}
