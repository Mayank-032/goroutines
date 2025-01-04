package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
)

func (p *Payload) UploadToS3WithUnbufferedChannel(wg *sync.WaitGroup, result chan Resp) {
	defer wg.Done()

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(p)
	if err != nil {
		log.Println("err: ", err.Error())
		result <- Resp{
			Message:   "unable to upload to S3",
			IsSuccess: false,
			Error:     err.Error(),
		}
		return
	}

	result <- Resp{
		Message:   "successfully uploaded to S3",
		IsSuccess: true,
	}
}

func UploadToS3WithUnbufferedChannel(w http.ResponseWriter, r *http.Request) {
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

	var wg sync.WaitGroup
	var results chan Resp
	for _, payload := range data.Payloads {
		wg.Add(1)
		go payload.UploadToS3WithUnbufferedChannel(&wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var res = make([]Resp, 0)
	for result := range results {
		res = append(res, result)
	}

	resp["success"] = true
	resp["data"] = res

	respBytes, _ := json.Marshal(resp)
	w.Write(respBytes)
}
