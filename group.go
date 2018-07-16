package fresh

import (
	"path/filepath"
)

type (
	Group interface {
		Rest
		Group(string) Group
		After(...HandlerFunc) Group
		Before(...HandlerFunc) Group
	}

	group struct {
		parent *fresh
		route  *route
	}
)


// Group registration
func (g *group) Group(path string) Group {
	return g.parent.Group(filepath.Join(g.route.path, path)).After(g.route.after...).Before(g.route.before...)
}

// WS api registration
func (g *group) WS(path string, handler HandlerFunc) Handler {
	return g.parent.WS(filepath.Join(g.route.path, path), handler).After(g.route.after...).Before(g.route.before...)
}

// Register a resource (get, post, put, delete)
func (g *group) CRUD(path string, h ...HandlerFunc) Resource {
	return g.parent.CRUD(filepath.Join(g.route.path, path), h...).After(g.route.after...).Before(g.route.before...)
}

// GET api registration
func (g *group) GET(path string, handler HandlerFunc) Handler {
	return g.parent.GET(filepath.Join(g.route.path, path), handler).After(g.route.after...).Before(g.route.before...)
}

// PUT api registration
func (g *group) PUT(path string, handler HandlerFunc) Handler {
	return g.parent.PUT(filepath.Join(g.route.path, path), handler).After(g.route.after...).Before(g.route.before...)
}

// POST api registration
func (g *group) POST(path string, handler HandlerFunc) Handler {
	return g.parent.POST(filepath.Join(g.route.path, path), handler).After(g.route.after...).Before(g.route.before...)
}

// TRACE api registration
func (g *group) TRACE(path string, handler HandlerFunc) Handler {
	return g.parent.TRACE(filepath.Join(g.route.path, path), handler).After(g.route.after...).Before(g.route.before...)
}

// PATCH api registration
func (g *group) PATCH(path string, handler HandlerFunc) Handler {
	return g.parent.PATCH(filepath.Join(g.route.path, path), handler).After(g.route.after...).Before(g.route.before...)
}

// DELETE api registration
func (g *group) DELETE(path string, handler HandlerFunc) Handler {
	return g.parent.DELETE(filepath.Join(g.route.path, path), handler).After(g.route.after...).Before(g.route.before...)
}

// OPTIONS api registration
func (g *group) OPTIONS(path string, handler HandlerFunc) Handler {
	return g.parent.OPTIONS(filepath.Join(g.route.path, path), handler).After(g.route.after...).Before(g.route.before...)
}

// ASSETS serve a list of static files. Array of files or directories TODO write logic
func (g *group) STATIC(static map[string]string) {
	g.parent.STATIC(static)
}

// After middleware
func (g *group) After(middleware ...HandlerFunc) Group {
	g.route.after = append(g.route.after, middleware...)
	return g
}

// Before middleware
func (g *group) Before(middleware ...HandlerFunc) Group {
	g.route.before = append(g.route.before, middleware...)
	return g
}
