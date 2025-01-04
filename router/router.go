package router

import (
	"lld/handler"
	"net/http"
)

func InitRoutes(router *http.ServeMux) {
	router.HandleFunc("/test", handler.UploadToS3WithSimpleGoRoutine)
}