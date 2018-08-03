package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"wxbot/src/utils"
	"wxbot/src/wx"
	"io/ioutil"
	"wxbot/src/models"
	"encoding/json"
)

func main() {
	if len(os.Args) == 1 {
		utils.LoadConfig("")
	} else if len(os.Args) == 2 {
		utils.LoadConfig(os.Args[1])
	} else {
		panic(errors.New("params error"))
	}
	//mirbase.InitClient()
	go wx.WxClient.Start()
	go func() {
		http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			result, _:= ioutil.ReadAll(r.Body)
			r.Body.Close()
			 messageSt:= models.MessageSt{}
			json.Unmarshal(result, &messageSt)
			wx.WxClient.SendMessage(messageSt.Message, messageSt.User)
			fmt.Fprintf(w, "send successfully!\nuser:%s,message:%s\n", messageSt.User, messageSt.Message)
			return
		})
		http.ListenAndServe(":"+utils.Conf.HttpConf.RestAPIPort, nil)
	}()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-signalChan
	log.Println("Shutdown signal received, exiting...")
}
