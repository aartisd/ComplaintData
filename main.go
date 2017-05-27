package main

import (
	"net/http"
	"complaintData/handler"
	"complaintData/complaintMgr"
	"github.com/gorilla/mux"
)

func main() {
	complaintMgr.Init()
	//http.HandleFunc("/complaint/all", handler.RegisterHandler)
	//http.HandleFunc("/complaint/.*", handler.RegisterHandler)
	r := mux.NewRouter()
	r.HandleFunc("/complaint/{area}", handler.ComplaintHandler)
	http.ListenAndServe(":8080", r)

}