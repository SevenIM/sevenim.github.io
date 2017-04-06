package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"public"
	"sql"
)

type Message struct {
	SendUser    uint64 `json:"send_user"`
	RecvUser    uint64 `json:"recv_user"`
	Message     string `json:"message"`
	ExternMsg   string `json:"extern_msg"`
	MessageType int    `json:"msg_type"`
	UniqMsgId   string `json:"uniq_msgid"`
}

const TableNum = 1

var sqlObj sql.SqlAlchemy

func Init(dbInfo *sql.DbInfo) int {
	_, err := sqlObj.Init(dbInfo)
	if err != nil {
		log.Fatal(err)
		return -1
	}
	return 0
}

func StoreMessage(writer http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	fmt.Println(req.URL.Path)
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(string(body))
	var msg Message
	err = json.Unmarshal(body, &msg)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
	out, err := json.Marshal(msg)
	fmt.Println(string(out))

	sessionId, err := public.GenSessionId(msg.SendUser, msg.RecvUser)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	msgTbNo := public.APHash(sessionId) % TableNum
	sqlObj.Insert("message_"+string(msgTbNo)).V("session_id",
		sessionId).V("sender", msg.SendUser).V("recver",
		msg.RecvUser).V("message", msg.Message).V("extern_msg",
		msg.ExternMsg).V("message_type", msg.MessageType).V("uniq_msg_id", msg.UniqMsgId).Execute()
}

func main() {
	dbInfo := sql.DbInfo{
		"mysql",
		"192.168.72.128:3306",
		"im",
		"root",
		"root",
		"utf8",
		1000,
		100,
	}

	if Init(&dbInfo) != 0 {
		return
	}

	http.HandleFunc("/message/store/", StoreMessage)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
