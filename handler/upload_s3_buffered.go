package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
)

func UploadToS3WithBufferedChannel(workerId int, wg *sync.WaitGroup, payloads chan Payload, result chan Resp) {
	defer wg.Done()
	for payload := range payloads {
		b := new(bytes.Buffer)
		err := json.NewEncoder(b).Encode(payload)
		if err != nil {
			log.Println("err: ", err.Error())
			result <- Resp{
				PId:       payload.Id,
				WorkerId:  workerId,
				Message:   "unable to upload to S3",
				IsSuccess: false,
				Error:     err.Error(),
			}
			return
		}

		result <- Resp{
			PId:       payload.Id,
			WorkerId:  workerId,
			Message:   "successfully uploaded to S3",
			IsSuccess: true,
		}
	}
}

func UploadToS3WithBufferedChannelAndWorkedPool(w http.ResponseWriter, r *http.Request) {
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

	var (
		maxWorkers = 10
		wg         sync.WaitGroup
		payloads   = make(chan Payload, len(data.Payloads))
		results    = make(chan Resp, len(data.Payloads))
	)

	// init workers waiting to exec payload received via channel
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go UploadToS3WithBufferedChannel(i, &wg, payloads, results)
	}

	// send payload to channels to assign them to workers
	for _, payload := range data.Payloads {
		payloads <- payload
	}
	close(payloads)

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
