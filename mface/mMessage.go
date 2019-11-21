/**
* @Author: Cooper
* @Date: 2019/11/18 14:40
 */

package mface

type MMessage interface {
	GetHead() string
	GetID() string
	GetLen() uint32
	GetData() []byte
	GetCConnID() string
	Marshal() []byte
}
