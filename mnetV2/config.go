/**
* @Author: Cooper
* @Date: 2019/11/21 20:23
 */

package mnetV2

type DefaultConfig struct {
	name         string
	netType      string
	host         string
	port         string
	maxConnCount uint32
	mMReqCS      uint32 // MManager request chan space
	mMRspCS      uint32 // MManager response chan space
	connRepCS    uint32 // Connection reply chan space
}

type defaultConfig struct {
	Name         string `json:"name"`
	NetType      string `json:"netType"`
	Host         string `json:"host"`
	Port         string `json:"port"`
	MaxConnCount uint32 `json:"maxConnCount"`
	MMReqCS      uint32 `json:"mMReqCS"`   // MManager request chan space
	MMRspCS      uint32 `json:"mMRspCS"`   // MManager response chan space
	ConnRepCS    uint32 `json:"connRepCS"` // Connection reply chan space
}

func (dc *DefaultConfig) Load() {

}

func (dc *DefaultConfig) Name() string {
	return dc.name
}

func (dc *DefaultConfig) NetType() string {
	return dc.netType
}

func (dc *DefaultConfig) Host() string {
	return dc.host
}

func (dc *DefaultConfig) Port() string {
	return dc.port
}

func (dc *DefaultConfig) MaxConnCount() uint32 {
	return dc.maxConnCount
}

func (dc *DefaultConfig) MMReqCS() uint32 {
	return dc.mMReqCS
}

func (dc *DefaultConfig) MMRspCS() uint32 {
	return dc.mMRspCS
}

func (dc *DefaultConfig) ConnReqCS() uint32 {
	return dc.connRepCS
}