package complaintMgr

import (
	"time"
	"sync"

	"strings"
	"golang.org/x/net/html"

	"io/ioutil"
	"bytes"
)
type ComplaintList struct {
	Complaints []ComplaintEntry
}

type ComplaintEntry struct {
	AreaCode         string 	  `json:"Area Code"`
	PhoneNumEntries  []PhoneNumEntry
}

type PhoneNumEntry struct {
	PhoneNumber string 	`json:"Phone Number"`
	NComments   int 	`json:"NumberOfComments"`
	Comments    []string    `json:"Comments"`
}
/*
 * commentMap - It is used to store mapping of phone number as key and
 *             value as slice of comments for that phone number
 */
type commentMap map[string][]string
var ComplaintMapLock sync.RWMutex
var ComplaintMap = make(map[string]commentMap)

/*
 * insertIntoComplaintMap - Insert entry read from html data to the ComplaintMap
 */

func insertIntoComplaintMap(content []string) {
	//content[0] = AreaCode, content[1] = phone number and content[3]=Comment
	areaCode := content[0]
	phoneNum := content[1]
	comment := content[2]

        //New entry for this area code
	if ComplaintMap[areaCode] == nil {
		commentMap := make(map[string][]string)
		commentSlice := []string{comment}
		commentMap[phoneNum] = commentSlice
		ComplaintMap[areaCode] = commentMap
		return
	}
	//Existing entry for area code
	commentMap := ComplaintMap[areaCode]
	//New entry for this phone number
	if commentMap[phoneNum] == nil {
		commentSlice := []string{comment}
		commentMap[phoneNum] = commentSlice
		return
	}
	//Existing entry for phone num
	commentMap[phoneNum] = append(commentMap[phoneNum],comment)

}

/*
 * scrapeComplaintData - Function to get the website contents and store them in the complaint map
 */
func scrapeComplaintData() {

	/*resp, err := http.Get("http://gsd-auth- callinfo.s3-website.us-east-2.amazonaws.com/")


	if err != nil {
		fmt.Println("ERROR: Failed to get contents")
		return
	}

	b := resp.Body
	defer b.Close() // close Body when the function returns*/

	byteData, err := ioutil.ReadFile("complaintData.html")
	if err != nil {
		panic(err)
	}
	b := bytes.NewReader(byteData)
	z := html.NewTokenizer(b)

	ComplaintMapLock.Lock()
	// Making a new map provides an easy way to get rid of any entries
	// that were deleted .  We're going to read in and replace
	// anything that hasn't been deleted .
	ComplaintMap = make(map[string]commentMap)


	content := []string{}
	for z.Token().Data != "html" {
		tt := z.Next()
		if tt == html.StartTagToken {
			t := z.Token()
			if t.Data == "td" {
				inner := z.Next()
				if inner == html.TextToken {
					text := (string)(z.Text())
					t := strings.TrimSpace(text)
					content = append(content, t)
					//fmt.Println(t, len(content))
					if len(content) == 3 {
						insertIntoComplaintMap(content)
						content = []string{}
					}
				}
			}
		}
	}

	ComplaintMapLock.Unlock()
        //fmt.Println(ComplaintMap)
}


/*
 * launchContentScraper - Function to launch a go routine that will scrape the complaint website
 *                        every 60 seconds.
 */
func launchContentScraper() {
	// scrape the website during initialization.
	scrapeComplaintData()

	// Now we'll check the website every 60 seconds.
	ticker := time.NewTicker(time.Second * 60)
	go func() {
		for range ticker.C {
			scrapeComplaintData()
		}
	}()
}

/*
 * GetAllComplaints - Go over the ComplaintMap and return list of all complaints
 */
func GetAllComplaints() ComplaintList {

	list := ComplaintList{[]ComplaintEntry{}}
	ComplaintMapLock.Lock()

	for areaCode, commentMap := range ComplaintMap {
		complaintObj := ComplaintEntry{AreaCode:areaCode}
		complaintObj.PhoneNumEntries = []PhoneNumEntry{}
		for phone, comments := range commentMap {
			phoneObj := PhoneNumEntry{PhoneNumber:phone,NComments:len(comments),Comments:comments}
			complaintObj.PhoneNumEntries = append(complaintObj.PhoneNumEntries, phoneObj)

		}
		list.Complaints = append(list.Complaints, complaintObj)

	}

	ComplaintMapLock.Unlock()
	return list
}


/*
 * GetComplaintsByAreaCodeHandler - Go over the ComplaintMap and return list of all complaints for area code
 */
func GetComplaintsByAreaCodeHandler(areaCode string) ComplaintList {
	list := ComplaintList{[]ComplaintEntry{}}
	ComplaintMapLock.Lock()
        if ComplaintMap[areaCode] != nil {

		complaintObj := ComplaintEntry{AreaCode: areaCode}
		complaintObj.PhoneNumEntries = []PhoneNumEntry{}
		for phone, comments := range ComplaintMap[areaCode] {
			phoneObj := PhoneNumEntry{PhoneNumber: phone, NComments: len(comments), Comments: comments}
			complaintObj.PhoneNumEntries = append(complaintObj.PhoneNumEntries, phoneObj)

		}
		list.Complaints = append(list.Complaints, complaintObj)


	}
	ComplaintMapLock.Unlock()
	return list
}

func Init() {
	launchContentScraper()
}