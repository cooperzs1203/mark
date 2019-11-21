/**
* @Author: Cooper
* @Date: 2019/11/18 15:14
 */

package mface

type RouteHandleFunc func(MMessage)MMessage

type MRouteManager interface {
	Add(string , RouteHandleFunc) error
	AddRoutes(map[string]RouteHandleFunc) error
	Handle(string , MMessage) error
	GetHandle(string) (RouteHandleFunc , error)
	SetServer(MServer)
}
