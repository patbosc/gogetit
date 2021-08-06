# go get it

a crawler wirtten in go because uneidel tells me all day long how awesome go is :P

Using colly which lives 
[on Github ](https://github.com/gocolly/colly) thanks to the [docs no one finds](http://go-colly.org/docs/) because of the [friendly crawler documentation on gopkg](https://pkg.go.dev/github.com/gocolly/colly#section-readme) 

First  installing Colly like this:

```
go get -u github.com/gocolly/colly/...
```

but you also need the go modules file to be like

```
 github.com/gocolly/colly/v2 v2.1.0
```

Okay the little service is a httpListener 

Supported Commands:
- Ping 
```
curl 'http://127.0.0.1:7272/ping
```
- Search
```
curl 'http://127.0.0.1:7272/search?url=http://www.google.com
```

Check out this stuff as well: 
--
- https://github.com/el10savio/GoCrawler
- http://go-colly.org/

I found this during my research and simply love it. 
- https://goclone.imthaghost.dev/ 