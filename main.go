package main

import (
	"net/http"
	"complaintData/handler"
	"complaintData/complaintMgr"
	"github.com/gorilla/mux"
)

func main() {
	complaintMgr.Init()

	r := mux.NewRouter()
	r.HandleFunc("/complaint/{area}", handler.ComplaintHandler)
	http.ListenAndServe(":8080", r)

}