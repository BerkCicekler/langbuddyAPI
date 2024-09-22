# langbuddyAPI
Its a restfull api that i am working on for my new flutter project. <br>
I am beginner on Golang so if you are a developer who look for example api code you probly should use this project. <br>
There will be many missing and bad coded parts in this project cause for now i just want to get used with language and packages


## Features
JWT Access Token<br>
JWT Refresh Token<br>
Account registration/login<br>
Freind Requests<br>
Random Friend Search <br>
Push Notification via Firebase messaging <br>
Chat system (websocket)

## Requests
### AUTH
/api/v1/users/login POST<br>
/api/v1/users/register POST<br>
/api/v1/users/refreshToken POST

### Friends
/api/v1/friends/sendRequest POST<br>
/api/v1/friends/data GET<br>
/api/v1/friends/accept POST<br>
/api/v1/friends/reject POST

### Friend Search
/api/v1/search/ POST

### User
/api/v1/user/language POST

### Chat 
/api/v1/chat/:roomid (websocket)
