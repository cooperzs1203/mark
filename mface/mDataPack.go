/**
* @Author: Cooper
* @Date: 2019/11/18 15:10
 */

package mface

type MDataPack interface {
	UnPack([]byte)
	Pack(MMessage) []byte
	CompleteMsgChan() chan MMessage
}