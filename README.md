# Fresh 

[![Travis](https://img.shields.io/travis/tockins/fresh.svg?style=flat-square)](https://travis-ci.org/tockins/fresh)
[![Go Report Card](https://goreportcard.com/badge/github.com/tockins/fresh?style=flat-square)](https://goreportcard.com/report/github.com/tockins/fresh)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/tockins/fresh)

Fresh is a lightweight GoLang web framework for RESTful API management.

#### Wiki

- [Features list](#features)
- [Getting Started](#installation)

#### Features

- RESTful API with:
-- route group management
-- filters
-- before and after handlers
- Docker container ready
- Microservices architecture
- DDD (Domain Driven Design) example
- Cli commands to create project and logs

<p align="center">
<img src="https://i.imgur.com/mCCF2br.png" width="350px">
</p>


#### Installation

Run this to get/install:
```
$ go get github.com/tockins/fresh
```

##### Examples

To manage a simple GET API:
```
func main() {
	f := fresh.New()
	f.Config().SetPort(8080)

	// API definition with path and related controller
	g.GET("/todos/", list)

	//Start Fresh Server
	f.Run()
}

func list(){
    return f.Response().JSON(http.StatusOK, nil)
}
```
