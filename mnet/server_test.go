/**
* @Author: Cooper
* @Date: 2019/11/18 16:05
 */

package mnet

import (
	"fmt"
	"log"
	"mark/mface"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func Client() {
	time.Sleep(time.Second*time.Duration(3))

	conn , err := net.Dial("tcp" , "127.0.0.1:8888")
	if err != nil {
		log.Println("Client ----- " , err)
		return
	}

	go func() {
		n := 0
		for {
			a := fmt.Sprintf("HEAD100000000%d",(n%2)+1)
			data := "ABCDE"+strconv.Itoa((n%2)+1)

			aBytes := []byte(a)
			aBytes = append(aBytes , intToBytes(len(data))...)
			aBytes = append(aBytes , []byte(data)...)
			if _, err := conn.Write(aBytes); err != nil {
				log.Println("Client ----- " , err)
			}
			time.Sleep(time.Second*time.Duration(3))
			n++
		}
	}()

	for {
		buf := make([]byte , 512)
		cnt , err := conn.Read(buf)
		if err != nil {
			log.Println("Client ----- " , err)
			return
		}

		log.Println("Client ----- "+string(buf[:cnt]))
	}
}

func TestServer(t *testing.T) {
	path , _ := os.Getwd()
	path = strings.TrimRight(path , "\\mnet")
	path = filepath.Join(path , "conf" , "defaultConfig.json")
	log.Println(path)

	server , err := NewServer(path)
	if err != nil {
		log.Printf("New server error : %v" , err)
		return
	}

	n := 0
	_ = server.RouteManager().Add("1000000001" , func(message mface.MMessage) mface.MMessage {
		conn , err := server.ConnManager().Get(message.GetCConnID())
		if err != nil {
			log.Printf("[%s] connection get error : %v" , message.GetCConnID() , err)
			return nil
		}

		n++
		conn.SetProperty("UserInfo" , n)

		if n == 6 {
			conn.CleanProperty()
		}

		log.Printf("Route 1 : %+v" , message)
		log.Printf("Route 1 : %s" , string(message.GetData()))
		replyMsg := NewMessage(append([]byte("Route Got You : ") , message.GetData()...) , message.GetCConnID())
		return replyMsg
	})

	_ = server.RouteManager().Add("1000000002" , func(message mface.MMessage) mface.MMessage {
		conn , err := server.ConnManager().Get(message.GetCConnID())
		if err != nil {
			log.Printf("[%s] connection get error : %v" , message.GetCConnID() , err)
			return nil
		}

		userInfo , err := conn.GetProperty("UserInfo")
		if err != nil {
			log.Printf("[%s] connection get property error : %v" , message.GetCConnID() , err)
			return nil
		}

		log.Printf("[%s] UserInfo : %v" , message.GetCConnID() , userInfo)

		log.Printf("Route 2 : %+v" , message)
		log.Printf("Route 2 : %s" , string(message.GetData()))
		replyMsg := NewMessage(append([]byte("Route Got You : ") , message.GetData()...) , message.GetCConnID())
		return replyMsg
	})

	go Client()

	server.Start()
}

