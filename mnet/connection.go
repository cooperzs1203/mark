/**
* @Author: Cooper
* @Date: 2019/11/18 18:30
 */

package mnet

import (
	"errors"
	"log"
	"mark/mface"
	"net"
	"sync"
)

func NewConnection(conn *net.Conn, id string, server mface.MServer, connRepCS uint32) mface.MConnection {
	connection := &Connection{
		Server:       server,
		Conn:         conn,
		ConnID:       id,
		ReplyChan:    make(chan []byte, connRepCS),
		IsClosed:     false,
		DP:           NewDataPack(id),
		Property:     make(map[string]interface{}),
		PropertyLock: sync.RWMutex{},
	}
	return connection
}

type Connection struct {
	Server    mface.MServer
	Conn      *net.Conn
	ConnID    string
	ReplyChan chan []byte
	IsClosed  bool
	DP        mface.MDataPack

	Property     map[string]interface{}
	PropertyLock sync.RWMutex
}

func (c *Connection) StartCommunicate() {
	log.Printf("[%s] start communication", c.ConnID)

	go c.startReadData()           // 开始监听读取数据，并将数据放入分包工具中
	go c.startMonitorCompleteMsg() // 开始监听分包工具的完整消息通道
	go c.startMonitorReply()       // 开始监听分发到的回复信息并写给连接
}

func (c *Connection) StopCommunicate() {
	log.Printf("[%s] stop communication", c.ConnID)

	if c.IsClosed {
		return
	}
	c.IsClosed = true
	if err := (*c.Conn).Close(); err != nil {
		log.Printf("[%s] close error : %v", c.ConnID, err)
	}
	close(c.ReplyChan)

	c.Server.ConnManager().Remove(c.ConnID) //从connManager中移除
}

func (c *Connection) GetConn() *net.Conn {
	return c.Conn
}

func (c *Connection) ID() string {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return (*c.Conn).RemoteAddr()
}

func (c *Connection) WriteChan() chan []byte {
	return c.ReplyChan
}

func (c *Connection) InClosed() bool {
	return c.IsClosed
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.PropertyLock.Lock()
	defer c.PropertyLock.Unlock()

	c.Property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.PropertyLock.RLock()
	defer c.PropertyLock.RUnlock()

	value , ok := c.Property[key]
	if !ok {
		return nil , errors.New(key+" not exists")
	}

	return value , nil
}

func (c *Connection) CleanProperty() {
	c.PropertyLock.Lock()
	defer c.PropertyLock.Unlock()

	for key := range c.Property {
		delete(c.Property , key)
	}

	c.Property = make(map[string]interface{})
}

func (c *Connection) startReadData() {
	for {
		buffer := make([]byte, 512)
		cnt, err := (*c.Conn).Read(buffer)
		if err != nil {
			if err.Error() == "EOF" {
				c.StopCommunicate() // 客户端主动断开连接
				break
			} else {
				log.Printf("[%s] read request error : %v", c.ConnID, err)
				continue
			}
		}

		// 分包粘包中
		c.DP.UnPack(buffer[:cnt])
	}
}

func (c *Connection) startMonitorCompleteMsg() {
	for {
		msg := <-c.DP.CompleteMsgChan()

		if !c.Server.MsgManager().InClosed() {
			c.Server.MsgManager().RequestChannel() <- msg
		}
	}
}

func (c *Connection) startMonitorReply() {
	for {
		msg, ok := <-c.ReplyChan
		if !ok && c.IsClosed {
			break
		}
		if _, err := (*c.Conn).Write(msg); err != nil {
			log.Printf("[%s] write reply %s --- error : %v", c.ConnID, string(msg), err)
		}
	}
}
