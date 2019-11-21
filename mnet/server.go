/**
* @Author: Cooper
* @Date: 2019/11/18 15:48
 */
/*
1.Config , 入口只要
*/

package mnet

import (
	"log"
	"mark/mface"
	"mark/utils"
	"net"
)

var (
	introduce = `
┌───────────────────────────────────────────────────┐
│                    Mark V0.8                      │
└───────────────────────────────────────────────────┘`
)

func init() {
	log.Println(introduce)
}

func NewServer(path string) (mface.MServer , error) {
	err := utils.LoadConfigFile(path)
	if err != nil {
		return nil , err
	}

	server := &Server{
		Name:         utils.GlobalConfig.Name,
		NetType:      utils.GlobalConfig.NetType,
		Host:         utils.GlobalConfig.Host,
		Port:         utils.GlobalConfig.Port,
		Listener:     nil,
		IsClosed:     false,
		MaxConnCount: utils.GlobalConfig.MaxConnCount,
		CManager:     NewConnManager(),
		MManager:     NewMsgManager(utils.GlobalConfig.MMReqCS , utils.GlobalConfig.MMRepCS),
		RManager:     NewRouteManager(),
	}

	// 设置路由管理器
	server.RManager.SetServer(server)

	// 设置并开启消息管理器
	server.MManager.SetServer(server)
	server.MManager.Start()

	return server , nil
}

// 服务器
type Server struct {
	Name     string
	NetType  string
	Host     string
	Port     string
	Listener *net.Listener
	IsClosed bool

	MaxConnCount uint32
	CManager     mface.MConnManager
	MManager     mface.MMsgManager
	RManager     mface.MRouteManager
}

func (s *Server) Start() {
	log.Printf("[START] Start %s server on %s with %s", s.Name, s.monitorAddr(), s.NetType)

	listener, err := net.Listen(s.NetType, s.monitorAddr())
	if err != nil {
		panic(err)
	}

	s.Listener = &listener

	s.startAcceptConnection()
}

func (s *Server) Stop() {
	if s.IsClosed {
		return
	}
	s.IsClosed = true
	if err := (*s.Listener).Close(); err != nil {
		log.Printf("[ERROR] Close listener error : %v", err)
	}

	// TODO:清理其它相关信息，例如日志未打印完、消息未转发完等待
	s.CManager.CleanAll() // 清除所有连接
	s.MManager.Stop()	  // 关闭消息管理器运作
}

func (s *Server) monitorAddr() string {
	if s.Port == "" {
		return ""
	}

	return net.JoinHostPort(s.Host, s.Port)
}

func (s *Server) startAcceptConnection() {
	log.Printf("[RUN] %s server start accept connection on %s", s.Name, s.monitorAddr())

	for {
		conn, err := (*s.Listener).Accept()
		if err != nil {
			if s.IsClosed {
				break
			} else {
				log.Printf("[RUN] %s server start accept connection on %s", s.Name, s.monitorAddr())
				continue
			}
		}

		// 限制最大连接数
		if pass := s.MaxConnLimitFilter(conn); !pass {
			log.Printf("[LIMIT] %s server is over limit running , close [%s]", s.Name, conn)
			continue
		}

		connection := NewConnection(&conn, utils.GetRandString(32), s , utils.GlobalConfig.ConnRepCS)
		connection.StartCommunicate()
		if err := s.CManager.Add(connection.ID(), connection); err != nil { // 如果添加错误，直接关闭连接
			connection.StopCommunicate()
		}
	}
}

func (s *Server) ConnManager() mface.MConnManager {
	return s.CManager
}

func (s *Server) MsgManager() mface.MMsgManager {
	return s.MManager
}

func (s *Server) RouteManager() mface.MRouteManager {
	return s.RManager
}

func (s *Server) MaxConnLimitFilter(conn net.Conn) bool {
	if s.CManager.Len() >= s.MaxConnCount {
		msg := []byte("Sorry , the server is over limit running")
		_ , _ = conn.Write(msg)
		_ = conn.Close()
		return false
	}

	return true
}
