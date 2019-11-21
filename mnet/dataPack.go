/**
* @Author: Cooper
* @Date: 2019/11/18 20:58
 */

package mnet

import (
	"bytes"
	"encoding/binary"
	"mark/mface"
)

const (
	HEAD               = "HEAD"
	HEAD_LENGTH        = 4
	MSG_ID_LENGTH      = 10
	LENGTH_DATA_LENGTH = 4
)

func NewDataPack(id string) mface.MDataPack {
	dataPack := &DataPack{
		ConnID:  id,
		Buffer:  make([]byte, 0),
		CompMsg: make(chan mface.MMessage),
	}
	return dataPack
}

type DataPack struct {
	ConnID  string
	Buffer  []byte
	CompMsg chan mface.MMessage
}

func (dp *DataPack) UnPack(data []byte) {
	dp.Buffer = append(dp.Buffer, data...)
	length := len(dp.Buffer)

	var i int
	for i = 0; i < length; i++ {
		if i+TotalHeaderLen() > length { // 剩余长度不足以获取消息完整头部
			break
		}
		if string(dp.Buffer[i:i+HeadLen()]) == HEAD {
			msgLength := BytesToInt(dp.Buffer[i+HeadMsgIDLen() : i+TotalHeaderLen()]) // 获取数据长度
			if i+TotalHeaderLen()+msgLength > length { // 长度不足以取出完整数据
				break
			}

			msg := dp.Buffer[i : i+TotalHeaderLen()+msgLength]
			dp.CompMsg <- NewMessage(msg , dp.ConnID)

			i += TotalHeaderLen()+msgLength - 1
		}
	}

	if i != length {
		leftData := dp.Buffer[i:]
		dp.Buffer = make([]byte, 0)
		dp.Buffer = append(dp.Buffer, leftData...)
	} else {
		dp.Buffer = make([]byte, 0)
	}
}

func (dp *DataPack) Pack(msg mface.MMessage) []byte {
	data := make([]byte , 0)
	data = append(data , []byte(msg.GetHead())...)
	data = append(data , []byte(msg.GetID())...)
	data = append(data , IntToBytes(int(msg.GetLen()))...)
	data = append(data , msg.GetData()...)

	return data
}

func (dp *DataPack) CompleteMsgChan() chan mface.MMessage {
	return dp.CompMsg
}

func IntToBytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	_ = binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}

func HeadLen() int {
	return HEAD_LENGTH
}

func HeadMsgIDLen() int {
	return HEAD_LENGTH + MSG_ID_LENGTH
}

func TotalHeaderLen() int {
	return HEAD_LENGTH + MSG_ID_LENGTH + LENGTH_DATA_LENGTH
}