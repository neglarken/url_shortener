# url_shortener
To start app
```
go build ./cmd/main.go
./cmd/cmd.exe
```
or
```
go run cmd/main.go
```
or start file cmd.exe in cmd folder

# POST localhost:8080/url
Add short URL in redirecting list
```
body:
{
"url":"example.com"
"alias":"example"
}
```
# DELETE localhost:8080/url
Delete short URL from redirecting list
```
body:
{
"alias":"example"
}
```
# GET localhost:8080/{alias}
Redirect to full URL by short URL
