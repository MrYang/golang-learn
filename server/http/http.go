package http

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"fmt"
	"os"

	"github.com/MrYang/golang-learn/conf"

	gws "github.com/gorilla/websocket"
	"golang.org/x/net/websocket"
)

func Start() {
	log.Println("init controller")
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		query := r.Form.Get("query")
		log.Println("query:", query)
		RenderMsgJson(w, "test")
	})

	http.Handle("/ws", websocket.Handler(func(ws *websocket.Conn) {
		writeMsg := make(chan string, 1)
		exitChan := make(chan struct{})
		go func() {
			for {
				var receive string
				if err := websocket.Message.Receive(ws, &receive); err != nil {
					exitChan <- struct{}{}
					break
				}
				writeMsg <- receive
			}
		}()
		for {
			select {
			case receive := <-writeMsg:
				msg := "reply " + receive
				if err := websocket.Message.Send(ws, msg); err != nil {
					break
				}
			case <-exitChan:
				break
			}

		}
	}))

	http.HandleFunc("/gws", func(w http.ResponseWriter, r *http.Request) {
		upgrader := gws.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		writeMsg := make(chan string, 1)
		exitChan := make(chan struct{})

		go func() {
			for {
				msgType, msg, err := conn.ReadMessage()
				log.Println("msgType", msgType)
				if err != nil {
					log.Println("read error", err)
					exitChan <- struct{}{}
					break
				}

				writeMsg <- string(msg)
			}

		}()

		for {
			select {
			case msg := <-writeMsg:
				if err := conn.WriteMessage(gws.TextMessage, []byte("reply "+string(msg))); err != nil {
					log.Println("write error", err)
					break
				}
			case <-exitChan:
				break
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		indexFile, err := os.Open("http/ws.html")
		if err != nil {
			http.NotFound(w, r)
			return
		}

		index, err := io.ReadAll(indexFile)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, string(index))
	})

	port := conf.Config().Server.Http.Port
	log.Println("start http server on port", port)
	http.ListenAndServe(port, nil)
}

func RenderJson(w http.ResponseWriter, v interface{}) {
	bs, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bs)
}

type Dto struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func RenderDataJson(w http.ResponseWriter, data interface{}) {
	RenderJson(w, Dto{Msg: "success", Data: data})
}

func RenderMsgJson(w http.ResponseWriter, msg string) {
	RenderJson(w, map[string]string{"msg": msg})
}

func AutoRender(w http.ResponseWriter, data interface{}, err error) {
	if err != nil {
		RenderMsgJson(w, err.Error())
		return
	}

	RenderDataJson(w, data)
}
