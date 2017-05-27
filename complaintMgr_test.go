package main

import (
	"testing"
	"complaintData/complaintMgr"
)


func TestScrapeComplaintData(t *testing.T) {
	complaintMgr.Init()
	if complaintMgr.ComplaintMap["456"] == nil || complaintMgr.ComplaintMap["345"] == nil {
		t.Errorf("Expected keys not present in ComplaintMap")
	}
}


func TestGetAllComplaintData(t *testing.T) {
	complaintMgr.Init()
	list := complaintMgr.GetAllComplaints()
	if len(list.Complaints) == 0  {
		t.Errorf("List is empty")
	}
}

func TestGetComplaintByAreaCode(t *testing.T) {
	complaintMgr.Init()
	list := complaintMgr.GetComplaintsByAreaCodeHandler("345")
	if len(list.Complaints) == 0  {
		t.Errorf("List is empty")
	}
}
