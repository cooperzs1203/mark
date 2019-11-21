/**
* @Author: Cooper
* @Date: 2019/11/18 21:00
 */

package mnet

import (
	"mark/mface"
)

func NewMessage(data []byte , connID string) mface.MMessage {
	head , id , dataLength , msgData := analysisData(data)
	message := &Message{
		ID:      id,
		CConnID: connID,
		Head:    head,
		Len:     dataLength,
		Data:    msgData,
	}
	return message
}

// 解析完整的数据
func analysisData(data []byte) (string, string, uint32, []byte) {
	var head , id string
	var dataLength int
	var msgData []byte

	head = string(data[:HeadLen()])
	id = string(data[HeadLen():HeadMsgIDLen()])
	dataLength = BytesToInt(data[HeadMsgIDLen():TotalHeaderLen()])
	msgData = data[TotalHeaderLen():]

	return head , id , uint32(dataLength) , msgData
}

type Message struct {
	ID string
	CConnID string
	Head string
	Len uint32
	Data []byte
}

func (m *Message) GetHead() string {
	return m.Head
}

func (m *Message) GetID() string {
	return m.ID
}

func (m *Message) GetLen() uint32 {
	return m.Len
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) GetCConnID() string {
	return m.CConnID
}

func (m *Message) Marshal() []byte {
	data := make([]byte , 0)
	data = append(data , []byte(m.Head)...)
	data = append(data , []byte(m.ID)...)
	data = append(data , IntToBytes(int(m.Len))...)
	data = append(data , m.Data...)

	return data
}