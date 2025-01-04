package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type PayloadCollection struct {
	Token    string    `json:"token"`
	Payloads []Payload `json:"payloads"`
}

type Payload struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Resp struct {
	Message   string `json:"message"`
	IsSuccess bool   `json:"is_succes"`
	Error     string `json:"error"`
}

func (p *Payload) UploadToS3(wg *sync.WaitGroup) {
	defer wg.Done()

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(p)
	if err != nil {
		log.Println("err: ", err.Error())
		return
	}

	time.Sleep(3 * time.Second)
}

func UploadToS3WithSimpleGoRoutine(w http.ResponseWriter, r *http.Request) {
	var resp = map[string]interface{}{
		"success": false,
	}

	file, err := os.Open("./examples/data.json")
	if err != nil {
		log.Println("unable to open json file, err: ", err.Error())
		respByte, _ := json.Marshal(resp)
		w.Write(respByte)
		return
	}

	var data = PayloadCollection{}
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		log.Println("unable to decode json file to struc, err: ", err.Error())
		respByte, _ := json.Marshal(resp)
		w.Write(respByte)
		return
	}

	/*
		-> The issue with below approach is we never know how many goroutines we can spawn like this
		-> Resources are always limited and we are using huge amount of resources thus huge load on our memory
	*/
	var wg sync.WaitGroup
	for _, payload := range data.Payloads {
		wg.Add(1)
		go payload.UploadToS3(&wg)
	}
	log.Printf("success uploading payload")
	wg.Wait()
	fmt.Println()

	resp["success"] = true
	respByte, _ := json.Marshal(resp)
	w.Write(respByte)
}
