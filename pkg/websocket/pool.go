package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Client 表示一个WebSocket客户端连接
type Client struct {
	Conn     *websocket.Conn
	UserID   uint64
	ParentID uint64
	Mu       sync.Mutex
}

// Pool 管理所有WebSocket连接
type Pool struct {
	clients map[uint64]*Client // 用户ID到客户端的映射
	mu      sync.RWMutex
}

// NewPool 创建一个新的连接池
func NewPool() *Pool {
	return &Pool{
		clients: make(map[uint64]*Client),
	}
}

// Register 注册一个新的客户端连接
func (p *Pool) Register(client *Client) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.clients[client.UserID] = client
}

// Unregister 注销一个客户端连接
func (p *Pool) Unregister(userID uint64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.clients, userID)
}

// GetParentClient 获取指定孩子的家长客户端
func (p *Pool) GetParentClient(childID uint64) *Client {
	p.mu.RLock()
	defer p.mu.RUnlock()

	parentID, exists := p.parentMap[childID]
	if !exists {
		return nil
	}

	return p.clients[parentID]
}

// SendToParent 向指定孩子的家长发送消息
func (p *Pool) SendToParent(childID uint64, message []byte) error {
	parentClient := p.GetParentClient(childID)
	if parentClient == nil {
		return nil // 家长不在线，忽略消息
	}

	parentClient.Mu.Lock()
	defer parentClient.Mu.Unlock()

	return parentClient.Conn.WriteMessage(websocket.TextMessage, message)
}

// GetClient 获取指定用户的客户端连接
func (p *Pool) GetClient(userID uint64) *Client {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.clients[userID]
}
