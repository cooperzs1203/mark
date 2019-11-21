/**
* @Author: Cooper
* @Date: 2019/11/18 13:23
 */

package main

import (
	"fmt"
	"log"
	"mark/mface"
	"mark/mnet"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	server mface.MServer
	n      int
)

func main() {
	Server()
}

func Client() {
	time.Sleep(time.Second * time.Duration(3))

	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		log.Println("Client ----- ", err)
		return
	}

	go func() {
		n := 0
		for {
			time.Sleep(time.Second * time.Duration(3))

			a := fmt.Sprintf("HEAD100000000%d", (n%3)+1)
			data := "ABCDE" + strconv.Itoa((n%3)+1)

			aBytes := []byte(a)
			aBytes = append(aBytes, mnet.IntToBytes(len(data))...)
			aBytes = append(aBytes, []byte(data)...)
			if _, err := conn.Write(aBytes); err != nil {
				log.Println("Client ----- ", err)
			}
			n++
		}
	}()

	for {
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			log.Println("Client ----- ", err)
			return
		}

		log.Println("Client ----- " + string(buf[:cnt]))
	}
}

func Server() {
	path, _ := os.Getwd()
	path = filepath.Join(path, "conf", "defaultConfig.json")
	log.Println(path)

	payServer, err := mnet.NewServer(path)
	if err != nil {
		log.Printf("New server error : %v", err)
		return
	}

	server = payServer

	routes := map[string]mface.RouteHandleFunc{
		"1000000001": routeOne,
		"1000000002": routeTwo,
		"1000000003": heartCheck,
	}

	err = server.RouteManager().AddRoutes(routes)
	if err != nil {
		log.Printf("add routes error : %v", err)
		return
	}

	//err = server.RouteManager().Add("1000000001" , routeOne)
	//err = server.RouteManager().Add("1000000002" , routeTwo)
	//err = server.RouteManager().Add("1000000003" , heartCheck)

	go Client()

	server.Start()
}

func heartCheck(message mface.MMessage) mface.MMessage {
	log.Printf("[%s] health check", message.GetCConnID())
	return nil
}

func routeOne(message mface.MMessage) mface.MMessage {
	conn, err := server.ConnManager().Get(message.GetCConnID())
	if err != nil {
		log.Printf("[%s] connection get error : %v", message.GetCConnID(), err)
		return nil
	}

	n++
	conn.SetProperty("UserInfo", n)

	if n == 6 {
		conn.CleanProperty()
	}

	log.Printf("Route 1 : %+v", message)
	log.Printf("Route 1 : %s", string(message.GetData()))
	replyMsg := mnet.NewMessage(append([]byte("Route Got You : "), message.GetData()...), message.GetCConnID())
	return replyMsg
}

func routeTwo(message mface.MMessage) mface.MMessage {
	conn, err := server.ConnManager().Get(message.GetCConnID())
	if err != nil {
		log.Printf("[%s] connection get error : %v", message.GetCConnID(), err)
		return nil
	}

	userInfo, err := conn.GetProperty("UserInfo")
	if err != nil {
		log.Printf("[%s] connection get property error : %v", message.GetCConnID(), err)
		return nil
	}

	log.Printf("[%s] UserInfo : %v", message.GetCConnID(), userInfo)

	log.Printf("Route 2 : %+v", message)
	log.Printf("Route 2 : %s", string(message.GetData()))
	replyMsg := mnet.NewMessage(append([]byte("Route Got You : "), message.GetData()...), message.GetCConnID())
	return replyMsg
}
