<A name="toc0_1" title="Thali Price Index"/>
##  Thali Price Index ##
A cost of living index for cities across India using user contributed / owned data. TPI v1 focuses on the price of a thali (meal) while v2 will focus on apartment rentals. In theory the platform can be extended to cover all sorts of price data which is otherwise difficult to obtain. 

##Contents     
**<a href="toc1_1">Methodology</a>**  
**<a href="toc1_2">Data</a>**  
**<a href="toc1_3">JSON API</a>**  
**<a href="toc1_4">App</a>**  
**<a href="toc1_5">Incentives</a>**  
**<a href="toc1_6">References</a>**  


<A name="toc1_1" title="Methodology" />
## Methodology ##
We're a highly stratified society so our thalis are also stratified. A thali can be broadly classified based on the target customer:

1. Blue Collar (Unorganised Labour Workers)
2. Yellow Collar (Retail & Organized Labour Workers)
3. White Collar (Office Workers)
4. Leisure

In addition the thali can also be classified on it's characteristics:

1. Limited
2. Unlimited

Or

1. South Indian
2. North Indian
3. Other Regional

After filtering outliers, the price index will be based on a weighted average of the data collected with a Yellow Collar Limited South Indian Thali in Bengaluru in 2016 being the benchmark of 100.  

<A name="toc1_2" title="Data" />
## Data ##
In v1 there's three data structures of interest:

+ User
    + Name string
    + Email string
    + Confirmed bool
    + Thalis []Thali // thalis contributed
    + Venues []int64 // venues contributed - []int64 due to datastore restriction of no nested slices
    + Rep int
    + Submitted time.Time

+ Venue
    + Name string
    + Latitude float64 // can be replaced with Location appengine.GeoPoint
    + Longitude float64 // can be replaced with Location appengine.GeoPoint
    + Thalis []int64
    + Submitted time.Time

+ Thali
    + Name string
    + Target int // 1-4 target customer profile
    + Limited bool
    + Region int // 1-3 target cuisine
    + Price float64 //
    + Photo string // filename in GCS
    + Venueid int64  // available at venue with id
    + Userid int64 // contributing by user with id
    + Verified bool
    + Accepted bool
    + Submitted time.Time

User -> Thali = One-to-many

We need a appengine datastore access structure and also a Postgres and/or Mongo access structure for deployment in case of move away from Appengine. All in Go.

In the appengine datastore version, Thali is slightly modified to include Id of Venue rather than a Venue (see appengine datastore reference). 


<A name="toc1_3" title="JSON API" />
## Endpoints ##
Data contribution, edit & retrieval is done via a simple HTTP/S REST JSON API. 

##VALIDATION##
https://thalipriceindex.appspot.com/token_auth
https://thalipriceindex.appspot.com/refresh_token_auth
https://thalipriceindex.appspot.com/logout

##CREATE (POST ONLY)##
https://thalipriceindex.appspot.com/create/user
Request body must consist of json with Name, Email, Password

https://thalipriceindex.appspot.com/create/venue
https://thalipriceindex.appspot.com/create/thali
https://thalipriceindex.appspot.com/create/data

Post JSON data, receive 200 OK if user/venue/thali/data created successfully

##RETRIEVE (POST ONLY)##


Post JSON data

##UPDATE (POST ONLY)##


##DELETE##


##HTML TEMPLATES##

HTML templates for logs and list of users/venues/thalis/data are available at:

https://thalipriceindex.appspot.com/logs
https://thalipriceindex.appspot.com/list/users
https://thalipriceindex.appspot.com/list/venues
https://thalipriceindex.appspot.com/list/thalis
https://thalipriceindex.appspot.com/list/datas


<A name="toc1_4" title="App & Browser" />
## App  ##

Mobile app needs to be very simple. Should use location and camera for posting data. 

Preferable to avoid any and all javascript in browser version. Need to consider data contributors with older phones/computers so app needs to be very basic. Should have some basic user validationinterface for Google/Facebook oAuth2, some basic data input functionality and ability to post that data to a server. Ability to get data and display tpi at user's location is secondary. 

Responsiveness of app is of primary importance rather than bells & whistles.

+Android
+ IOS app 

<A name="toc1_2" title="Incentives" />
## Incentives ##
The starting group of users will be a small number - 30 colleagues, friends, families, willing acquaintances. So no real need to have a super scalable backend. Other users and spammers will hopefully contribute. 

As soon as a user contributes 10 verified/accepted data points, they get access to the entire data set via the JSON API. 

Spammers should gain negative reputation for every unverified/unaccepted data point and after 10 such points unable to contribute.

<A name="toc1_5" title="Gotchas" />
## Gotchas ##



<A name="toc1_6" title="References" />
## References ##
+ [Writing images to templates](http://www.sanarias.com/blog/1214PlayingwithimagesinHTTPresponseingolang)
+ [Appengine datastore api](https://godoc.org/google.golang.org/appengine/datastore)
+ [GCP Appengine Console](https://console.cloud.google.com/appengine?project=tpi)
+ [Method: apps.repair](https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1/apps/repair)
+ [Google Cloud Platform Datastore Reference](https://cloud.google.com/appengine/docs/go/datastore/reference)


