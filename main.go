package main

import (
	"net/http"
	"complaintData/handler"
	"complaintData/complaintMgr"
	"github.com/gorilla/mux"
	"complaintData/complaintMgrR"
)

func main() {
	complaintMgr.Init()
	complaintMgrR.Init()
	r := mux.NewRouter()
	r.HandleFunc("/complaint/{area}", handler.ComplaintHandlerR)
	http.ListenAndServe(":8080", r)

}


