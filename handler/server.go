package handler

import (
	"net/http"
	"github.com/gorilla/mux"
	"complaintData/complaintMgr"
	"encoding/json"
	"complaintData/complaintMgrR"
)

func GetAllComplaintsHandler(list complaintMgr.ComplaintList, w http.ResponseWriter) {

	if len(list.Complaints) == 0 {
		out := []byte("No results for this area code")
		w.Write(out)
		return
	}
	out, err := json.Marshal(list)
	w.Header().Add("Content-type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		out = []byte("Unknown app error")
		w.Write(out)

	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(out)
	}

}


func GetAllComplaintsHandlerR(list complaintMgrR.ComplaintListR, w http.ResponseWriter) {

	if len(list.Complaints) == 0 {
		out := []byte("No results for this area code")
		w.Write(out)
		return
	}
	out, err := json.Marshal(list)
	w.Header().Add("Content-type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		out = []byte("Unknown app error")
		w.Write(out)

	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(out)
	}

}


func ComplaintHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var list complaintMgr.ComplaintList
	if vars["area"] == "all" {
		list = complaintMgr.GetAllComplaints()

	} else {
		list = complaintMgr.GetComplaintsByAreaCodeHandler(vars["area"])
	}
	GetAllComplaintsHandler(list, w)
}

func ComplaintHandlerR(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var list complaintMgrR.ComplaintListR
	if vars["area"] == "all" {
		list = complaintMgrR.GetAllComplaintsR()

	} else {
		list = complaintMgrR.GetComplaintsByAreaCodeHandlerR(vars["area"])
	}
	GetAllComplaintsHandlerR(list, w)
}