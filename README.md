# Fresh 

[![Travis](https://img.shields.io/travis/tockins/fresh.svg?style=flat-square)](https://travis-ci.org/tockins/fresh)
[![Go Report Card](https://goreportcard.com/badge/github.com/tockins/fresh?style=flat-square)](https://goreportcard.com/report/github.com/tockins/fresh)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/tockins/fresh)
[![AUR](https://img.shields.io/aur/license/yaourt.svg?style=flat-square)](https://raw.githubusercontent.com/tockins/fresh/v1/LICENSE)
[![](https://img.shields.io/badge/fresh-examples-yellow.svg?style=flat-square)](https://github.com/tockins/fresh-examples)
[![Gitter](https://img.shields.io/gitter/room/tockins/fresh.svg?style=flat-square)](https://gitter.im/tockins/fresh?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)


<p align="center">
    <img src="https://i.imgur.com/ogzKrn7.png" width="300px">
</p>

### About Fresh

Fresh lightweight and high performance RESTful framework written in GoLang. Fresh aim to be ready to use for environments where scaling is needed.

### Content

- [Features list](#features)
- [Get Started](#get-started)
- [Contribute](#contribute)

### Features

- RESTful API with:
  - route group management 
  - filters
  - before and after handlers
- Microservices architecture
- DDD (Domain Driven Design) example
- Cli commands to create project and logs

### Get Started

Run this to get and install fresh:
```
$ go get github.com/tockins/fresh
```

### Examples

To manage a simple GET API:

```
func main() {
    f := fresh.New()
    f.Config().Port(8080)

    // API definition with path and related controller
    f.GET("/todo/", func(c fresh.Context) error{
	    return c.Response().JSON(http.StatusOK, nil)
	})
    f.GET("/todo/:uuid", func(c fresh.Context) error{
        todoUuid := c.Request().URLParam("uuid")
        res := map[string]interface{}{ "uuid": todoUuid}
        return c.Response().JSON(http.StatusOK, res)
     })
    //Start Fresh Server
    f.Run()
}
```

### Contribute

[See our guidelines](https://github.com/tockins/fresh/blob/master/CONTRIBUTING.md)
