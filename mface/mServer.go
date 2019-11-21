/**
* @Author: Cooper
* @Date: 2019/11/18 14:18
 */

package mface

import "net"

type MServer interface {
	Start()
	Stop()
	MaxConnLimitFilter(net.Conn) bool
	ConnManager() MConnManager
	MsgManager() MMsgManager
	RouteManager() MRouteManager
}
