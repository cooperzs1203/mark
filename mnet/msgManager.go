/**
* @Author: Cooper
* @Date: 2019/11/19 16:13
 */

package mnet

import (
	"log"
	"mark/mface"
)

func NewMsgManager(mMReqCS , mMRepCS uint32) mface.MMsgManager {
	msgManager := &MsgManager{
		Server:      nil,
		IsClosed:    false,
		RequestChan: make(chan mface.MMessage, mMReqCS),
		ReplyChan:   make(chan mface.MMessage, mMRepCS),
	}

	return msgManager
}

type MsgManager struct {
	Server      mface.MServer
	IsClosed    bool
	RequestChan chan mface.MMessage
	ReplyChan   chan mface.MMessage
}

func (mm *MsgManager) ReplyChannel() chan mface.MMessage {
	return mm.ReplyChan
}

func (mm *MsgManager) RequestChannel() chan mface.MMessage {
	return mm.RequestChan
}

func (mm *MsgManager) Start() {
	go mm.startAcceptRequest()
	go mm.startAcceptReply()
}

func (mm *MsgManager) Stop() {
	if mm.IsClosed {
		return
	}
	mm.IsClosed = true

	close(mm.RequestChan)
	close(mm.ReplyChan)

	// TODO:将数据持久化
}

func (mm *MsgManager) InClosed() bool {
	return mm.IsClosed
}

func (mm *MsgManager) SetServer(server mface.MServer) {
	mm.Server = server
}

func (mm *MsgManager) startAcceptRequest() {
	for {
		msg, ok := <-mm.RequestChan
		if !ok && mm.IsClosed {
			log.Printf("MsgManager accept request stop")
			break
		}
		// TODO:可以设计一个中间件插件，功能用来过滤消息

		err := mm.Server.RouteManager().Handle(msg.GetID() , msg)
		if err != nil {
			log.Printf("[ERROR] handle %s message error : %v" , msg.GetID() , err)
		}
	}
}

func (mm *MsgManager) startAcceptReply() {
	for {
		msg, ok := <-mm.ReplyChan
		if !ok && mm.IsClosed {
			log.Printf("MsgManager accept reply stop")
			break
		}

		conn , err := mm.Server.ConnManager().Get(msg.GetCConnID())
		if err != nil {
			log.Printf("[REPLY] message to %s error : %v" , msg.GetCConnID() , err)
			continue
		}
		if conn.InClosed() {
			log.Printf("[%s] is closed , can't reply message" , conn.ID())
			continue
		}

		conn.WriteChan() <- msg.Marshal()
	}
}
