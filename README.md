ComplaintData service scrapes the phone number complaint data from  http://gsd-auth- callinfo.s3-website.us-east-2.amazonaws.com/ and provides parsed data via restful apis.
(Currenlty the website mentioned above is not working and so the service reads the data from ComplaintData.html which is part of the respository)

GETTING STARTED

To get started with the service ensure that GOPATH is configured to point to the src folder of your go repository.


To start the service use
./complaintData

The service listens for http requests on port 8080. It listens for requests on following endpoint
GET http://localhost:8080/complaint/[AreaCode]

GET http://localhost:8080/complaint/all

EXAMPLES

Fetch all:
GET http://localhost:8080/complaint/all

{
   "Complaints":[
      {
         "Area Code":"345",
         "PhoneNumEntries":[
            {
               "Phone Number":"3458882345",
               "NumberOfComments":4,
               "Comments":[
                  "Spam calls",
                  "Spam alerts",
                  "Spam",
                  "No one speaks from other end"
               ]
            }
         ]
      },
      {
         "Area Code":"456",
         "PhoneNumEntries":[
            {
               "Phone Number":"4567891234",
               "NumberOfComments":1,
               "Comments":[
                  "Unknown number"
               ]
            },
            {
               "Phone Number":"4568882345",
               "NumberOfComments":2,
               "Comments":[
                  "Fake IRS calls",
                  "Fake spam"
               ]
            }
         ]
      }
   ]
}

Fetch by area code:
GET http://localhost:8080/complaint/345
{
   "Complaints":[
      {
         "Area Code":"456",
         "PhoneNumEntries":[
            {
               "Phone Number":"4567891234",
               "NumberOfComments":1,
               "Comments":[
                  "Unknown number"
               ]
            },
            {
               "Phone Number":"4568882345",
               "NumberOfComments":2,
               "Comments":[
                  "Fake IRS calls",
                  "Fake spam"
               ]
            }
         ]
      }
   ]
}

If an area code is not found in the map it returns a message "No results for this area code"


MODULES
complaintMgr - Responsible for scraping the data from the website and maintaining that data in the memory.
It refreshes this data every 60 seconds. It provides apis so that the Webservice handlers can invoke them to get required information.

handler - It has the Webservice handler functions to send the requested complaint data to the user.

