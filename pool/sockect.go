package pool

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	// 读写缓冲区大小
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// 在线用户管理

type userPool struct {
	users map[uint]*websocket.Conn
	mutex sync.RWMutex
}

func (p *userPool) AddUser(userID uint, conn *websocket.Conn) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	oldConn, k := p.users[userID]

	// 移除之前的链接
	if k {
		oldConn.Close()
	}

	p.users[userID] = conn
}

func (p *userPool) RemoveUser(userID uint) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if conn, ok := p.users[userID]; ok {
		conn.Close()
		delete(p.users, userID)
	}
}

func (p *userPool) GetUserConn(userID uint) *websocket.Conn {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	conn, _ := p.users[userID]
	return conn
}

func (p *userPool) IsOnline(uesrID uint) *websocket.Conn {
	p.mutex.RLocker()
	defer p.mutex.RUnlock()

	return p.users[uesrID]
}

var UserPool = &userPool{
	users: make(map[uint]*websocket.Conn),
}
