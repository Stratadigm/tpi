<A name="toc0_1" title="Thali Price Index"/>
# Thali Price Index is a cost of living index for cities across India using user contributed / owned data. TPI v1 focuses on the price of a thali while v2 will focus on apartment rentals. In theory the platform can be extended to cover all sorts of price data which is otherwise difficult to obtain. 

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
        Name string
        Email string
        Confirmed bool
        Points []Data // data points contributed
        Rep int
        JDte time.Time

+ Thali
        Target int // 1-4 target customer profile
        Limited bool
        Region int // 1-3 target cuisine
        Price float64 //
        Photo image

+ Venue
        Name string
        Latitude float64
        Longitude float64
        Thalis []Thali

+ Data
        Tha Thali
        Ven Venue
        Time time.Time
        Cntrbtr Contributor 
        Verfied bool
        Accepted bool

User -> Data = One-to-many
Venue -> Thali = One-to-many

We need a appengine datastore access structure and also a Postgres or Mongo access structure for deployment in case of move away from Appengine. All in Go.

<A name="toc1_3" title="JSON API" />
## Endpoints ##
Data contribution, edit & retrieval is done via a simple HTTP/S REST JSON API. 



<A name="toc1_4" title="App & Browser" />
## App  ##

Mobile app needs to be very simple. Should use location and camera for submission of the data. 

Preferable to avoid any and all javascript in browser version. Need to consider data contributors with older phones/computers so app needs to be very basic. Should have some basic user validationinterface for Google/Facebook oAuth2, some basic data input functionality and ability to post that data to a server. Ability to get data and display tpi at user's location is secondary. 

Responsiveness of app is of primary importance rather than bells & whistles.

+Android
+ IOS app 

<A name="toc1_2" title="Incentives" />
## Incentives ##
The starting group of users will be a small number - 30 colleagues, friends, families, willing acquaintances. So no real need to have a super scalable backend. Other users and spammers will hopefully contribute. 

As soon as a user contributes 10 verified/accepted data points, they get access to the entire data set via the JSON API. 

Spammers should gain negative reputation for every unverified/unaccepted data point and after 10 such points unable to contribute.

<A name="toc1_5" title="References" />
## References ##
[Writing images to templates](http://www.sanarias.com/blog/1214PlayingwithimagesinHTTPresponseingolang)
[Appengine datastore api](https://godoc.org/google.golang.org/appengine/datastore)
[GCP Appengine Console](https://console.cloud.google.com/appengine?project=tpi)
[Method: apps.repair](https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1/apps/repair)


