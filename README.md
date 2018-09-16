# mbc-api
a simple REST API built on scrapping mbc.net
>## DISCLAIMER
>#### this is mere practice i can not or any body use this for commerical purposes as the site scrapped stated clearly
##### What will you get?
##### You get the whole catelogue of shows displayed in the channels networks ie. mbc3,mbc2 and so on
##### a request can be made to get the currently being displayed show from any channel (/{channel}/currently) -> add channels' urls to the global map  
##### a request can be made to get a specific day's catelogue of any channel (/{channel}/{day})
##### you can request the whole week catelogue too (/{channel}/week)
##### every show has the following properties :
- Title
- Description (IN Arabic)
- Thumbnail Image
- Show Times -> EGY,GMT,KSA according to each channel
##### all Times and dates were parsed so they can be parsed back to golang time.Time objects
##### response is obviously in JSon..

