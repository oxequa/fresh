# Fresh 

[![Build Status](https://travis-ci.org/tockins/fresh.svg?branch=master)](https://travis-ci.org/tockins/fresh) [![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

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
