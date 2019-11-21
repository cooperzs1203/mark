/**
* @Author: Cooper
* @Date: 2019/11/18 19:27
 */

package mnet

import (
	"errors"
	"mark/mface"
	"sync"
)

func NewConnManager() mface.MConnManager {
	connManager := &ConnManager{
		Conns:    make(map[string]mface.MConnection),
		ConnLock: sync.RWMutex{},
	}
	return connManager
}

type ConnManager struct {
	Conns map[string]mface.MConnection
	ConnLock sync.RWMutex
}

func (cm *ConnManager) Add(id string , newConn mface.MConnection) error {
	cm.ConnLock.Lock()
	defer cm.ConnLock.Unlock()

	_ , ok := cm.Conns[id]
	if ok {
		return errors.New("connection id repeat")
	}

	cm.Conns[id] = newConn
	return nil
}

func (cm *ConnManager) Remove(id string) {
	cm.ConnLock.Lock()
	defer cm.ConnLock.Unlock()

	delete(cm.Conns , id)
}

func (cm *ConnManager) CleanAll() {
	cm.ConnLock.Lock()
	defer cm.ConnLock.Unlock()

	for id := range cm.Conns {
		conn := cm.Conns[id]
		conn.StopCommunicate() // 主动清除，需要通知停止通讯
		delete(cm.Conns , id)
	}
}

func (cm *ConnManager) Get(id string) (mface.MConnection , error) {
	cm.ConnLock.RLock()
	defer cm.ConnLock.RUnlock()

	var err error
	conn , ok := cm.Conns[id]
	if !ok {
		err = errors.New(id+" not exists")
	}

	return conn , err
}

func (cm *ConnManager) Len() uint32 {
	cm.ConnLock.RLock()
	defer cm.ConnLock.RUnlock()

	count := len(cm.Conns)
	return uint32(count)
}