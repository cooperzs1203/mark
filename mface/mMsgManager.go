/**
* @Author: Cooper
* @Date: 2019/11/18 14:48
 */

package mface

type MMsgManager interface {
	ReplyChannel() chan MMessage
	RequestChannel() chan MMessage
	Start()
	Stop()
	InClosed() bool
	SetServer(MServer)
}
