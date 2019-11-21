/**
* @Author: Cooper
* @Date: 2019/11/18 14:37
 */

package mface

import "net"

type MConnection interface {
	StartCommunicate()
	StopCommunicate()
	GetConn() *net.Conn
	ID() string
	RemoteAddr() net.Addr
	InClosed() bool

	WriteChan() chan []byte
	SetProperty(string , interface{})
	GetProperty(string) (interface{} , error)
	CleanProperty()
}
