/**
* @Author: Cooper
* @Date: 2019/11/19 19:32
 */

package mnet

import (
	"errors"
	"mark/mface"
	"sync"
)

func NewRouteManager() mface.MRouteManager {
	routeManager := &RouteManager{
		Server:    nil,
		Routes:    make(map[string]mface.RouteHandleFunc),
		RouteLock: sync.RWMutex{},
	}
	return routeManager
}

type RouteManager struct {
	Server    mface.MServer
	Routes    map[string]mface.RouteHandleFunc
	RouteLock sync.RWMutex
}

func (rm *RouteManager) Add(msgID string, handleFunc mface.RouteHandleFunc) error {
	rm.RouteLock.Lock()
	defer rm.RouteLock.Unlock()

	_, ok := rm.Routes[msgID]
	if ok {
		return errors.New(msgID + " exists")
	}

	rm.Routes[msgID] = handleFunc

	return nil
}

func (rm *RouteManager) AddRoutes(routes map[string]mface.RouteHandleFunc) error {
	rm.RouteLock.Lock()
	defer rm.RouteLock.Unlock()

	for msgID , handleFunc := range routes {
		_, ok := rm.Routes[msgID]
		if ok {
			return errors.New(msgID + " exists")
		}

		rm.Routes[msgID] = handleFunc
	}

	return nil
}

func (rm *RouteManager) Handle(msgID string, msg mface.MMessage) error {
	rm.RouteLock.RLock()
	defer rm.RouteLock.RUnlock()

	handleFunc, ok := rm.Routes[msgID]
	if !ok {
		return errors.New(msgID + " not exists")
	}

	go func(msgManager mface.MMsgManager) {
		replyMsg := handleFunc(msg)
		if replyMsg != nil && !msgManager.InClosed() { // nil 表示不返回消息
			msgManager.ReplyChannel() <- replyMsg
		}
	}(rm.Server.MsgManager())

	return nil
}

func (rm *RouteManager) GetHandle(msgID string) (mface.RouteHandleFunc, error) {
	rm.RouteLock.RLock()
	defer rm.RouteLock.RUnlock()

	handleFunc, ok := rm.Routes[msgID]
	if !ok {
		return nil, errors.New(msgID + " not exists")
	}

	return handleFunc, nil
}

func (rm *RouteManager)	SetServer(server mface.MServer)  {
	rm.Server = server
}