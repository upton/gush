package gush

import (
	"bufio"
	"net"
	"strings"
	"time"
)

const (
	AUTH       = "_A:"
	HB         = "_H:"
	END        = "__END"
	MSG_PREFIX = "_M:"
	NEW_LINE   = "\n"
)

var userMap = make(map[string]*UserChannel)

type UserChannel struct {
	uid string
	msg chan string
}

func Run() {
	ln, err := net.Listen("tcp", ":"+config.Port_tcp)
	if err != nil {
		Logger.Error(err)
	}

	go regNotifyApi()

	Logger.Info("start success")
	for {
		conn, err := ln.Accept()
		if err != nil {
			Logger.Error(err)
			continue
		}

		channel := make(chan string, 5)
		uc := &UserChannel{"", channel}

		go readConn(conn, uc)
		go wirteConn(conn, uc)
	}
}

func readConn(c net.Conn, uc *UserChannel) {
	defer func() {
		c.Close()
		Logger.Warn("conn closed on readConn.")
	}()

	reader := bufio.NewReader(c)

	for {
		c.SetReadDeadline(time.Now().Add(time.Duration(config.Read_timeout) * time.Second))
		bb, _, err := reader.ReadLine()

		if err != nil {
			Logger.Error("read error: ", err)
			break
		}

		if len(bb) > 0 {
			m := string(bb)
			routeMsg(m, uc)
		}
	}
}

func wirteConn(c net.Conn, uc *UserChannel) {
	defer func() {
		c.Close()
		delete(userMap, uc.uid)
		Logger.Warn("conn closed on writeConn.")
	}()

	for {
		select {
		case msg := <-uc.msg:
			c.SetWriteDeadline(time.Now().Add(time.Duration(config.Write_timeout) * time.Second))
			_, err := c.Write([]byte(msg + NEW_LINE))

			if err != nil {
				Logger.Error("Write error: ", err)
				break
			}

		case <-time.After(time.Duration(config.Read_timeout) * time.Second):
			break
		}
	}
}

func routeMsg(request string, uc *UserChannel) {
	if strings.HasPrefix(request, AUTH) {
		auth(strings.Replace(request, AUTH, "", 1), uc)
	} else if request == HB {
		heartbeat(uc)
	} else {
		if uc.uid != "" {
			//do something
		}
	}
}

func auth(request string, uc *UserChannel) {
	//TODO: do auth check
	//if fail {
	//	uc.msg <- AUTH + "FAIL"
	//	uc.msg <- END
	//}

	uc.uid = request
	userMap[request] = uc

	uc.msg <- AUTH + "OK[uid=" + request + "]"
}

func heartbeat(uc *UserChannel) {
	uc.msg <- HB + "OK[uid=" + uc.uid + "]"
}
