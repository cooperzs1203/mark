/**
* @Author: Cooper
* @Date: 2019/11/21 20:14
 */

package mfaceV2

type MConfig interface {
	Name() string
	NetType() string
	Host() string
	Port() string
	MaxConnCount() uint32
	MMReqCS() uint32
	MMRspCS() uint32
	ConnReqCS() uint32
}
