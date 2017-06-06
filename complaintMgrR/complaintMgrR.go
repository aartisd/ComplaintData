package complaintMgrR

import (
	"time"
	"sync"

	"strings"
	"github.com/PuerkitoBio/goquery"
	"strconv"
)
type ComplaintListR struct {
	Complaints []ComplaintEntryR
}

type ComplaintEntryR struct {
	AreaCode         string 	  `json:"Area Code"`
	PhoneNumEntries  []PhoneNumEntryR
}

type PhoneNumEntryR struct {
	PhoneNumber string 	`json:"Phone Number"`
	NComments   int 	`json:"NumberOfComments"`
	Comment     string      `json:"Comment"`
}

var ComplaintMapLockR sync.RWMutex
var ComplaintMapR = make(map[string][]PhoneNumEntryR)



func insertIntoComplaintMapR(areaCode string, phoneNum string, comment string, nComments int) {
	phoneObj := PhoneNumEntryR{PhoneNumber:phoneNum,NComments:nComments,Comment:comment}
	//New entry for this area code
	if ComplaintMapR[areaCode] == nil {
		ComplaintMapR[areaCode] = []PhoneNumEntryR{phoneObj}
		return
	}
	//Existing entry for area code
	entries := ComplaintMapR[areaCode]
	ComplaintMapR[areaCode] = append(entries,phoneObj)

}

func scrapeComplaintDataR() {
	doc, err := goquery.NewDocument("http://gsd-auth-callinfo.s3-website.us-east-2.amazonaws.com/")
	if err != nil {
		panic(err)
		return
	}
	ComplaintMapLockR.Lock()
	// Making a new map provides an easy way to get rid of any entries
	// that were deleted .  We're going to read in and replace
	// anything that hasn't been deleted .
	ComplaintMapR = make(map[string][]PhoneNumEntryR)

	doc.Find(".oos_noUserSelect.oos_list.oos_l4 .oos_listItem").Each(func(index int, item *goquery.Selection) {
		var phone,areaC string
		postCount := item.Find(".postCount")
		comment := item.Find(".oos_previewBody")
		if(postCount == nil || comment == nil) {
			return
		}
		//title := item.Text()
		item.Find("a").Each(func(index int, item *goquery.Selection) {
			link, _ := item.Attr("href")
			s := strings.Split(link,"/")
			if(s[3] == "Phone.aspx") {
				phone = s[4]
			} else if(s[3] == "AreaCode.aspx") {
				areaC = s[4]
			}
		})
		if(phone == "" || areaC == "") {
			return
		}
		//fmt.Printf("#%d: %s - %s - %s - %s\n", index, postCount.Text(), comment.Text(), phone, areaC)
		nComment,_ := strconv.Atoi(postCount.Text())
		insertIntoComplaintMapR(areaC, phone, comment.Text(),nComment)
	})
	ComplaintMapLockR.Unlock()
}

/*
 * launchContentScraper - Function to launch a go routine that will scrape the complaint website
 *                        every 60 seconds.
 */
func launchContentScraperR() {
	// scrape the website during initialization.
	scrapeComplaintDataR()

	// Now we'll check the website every 60 seconds.
	ticker := time.NewTicker(time.Second * 60)
	go func() {
		for range ticker.C {
			scrapeComplaintDataR()
		}
	}()
}

/*
 * GetAllComplaints - Go over the ComplaintMap and return list of all complaints
 */
func GetAllComplaintsR() ComplaintListR {

	list := ComplaintListR{[]ComplaintEntryR{}}
	ComplaintMapLockR.Lock()

	for areaCode, entries := range ComplaintMapR {
		complaintObj := ComplaintEntryR{AreaCode:areaCode}
		complaintObj.PhoneNumEntries = entries
		list.Complaints = append(list.Complaints, complaintObj)

	}

	ComplaintMapLockR.Unlock()
	return list
}


/*
 * GetComplaintsByAreaCodeHandler - Go over the ComplaintMap and return list of all complaints for area code
 */
func GetComplaintsByAreaCodeHandlerR(areaCode string) ComplaintListR {
	list := ComplaintListR{[]ComplaintEntryR{}}
	ComplaintMapLockR.Lock()
        if ComplaintMapR[areaCode] != nil {

		complaintObj := ComplaintEntryR{AreaCode: areaCode}
		complaintObj.PhoneNumEntries = ComplaintMapR[areaCode]
		list.Complaints = append(list.Complaints, complaintObj)

	}
	ComplaintMapLockR.Unlock()
	return list
}

func Init() {
	launchContentScraperR()
}