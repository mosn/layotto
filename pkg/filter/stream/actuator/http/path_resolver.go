package http

import "strings"

// PathResolver is a util to extract elements in http path.
type PathResolver struct {
	// e.g. /a/b/c/d
	rawPath string
	//unsolved path,which starts by "/". e.g.  /b/c/d
	unresolved string
}

func NewPathResolver(path string) *PathResolver {
	return &PathResolver{
		rawPath:    path,
		unresolved: path,
	}
}

func (p *PathResolver) HasNext() bool {
	path := p.UnresolvedPath()
	return path != "" && path != "/" && path != "\\"
}

func (p *PathResolver) Next() string {
	if !p.HasNext() {
		return ""
	}
	// /a/b/c
	// remove first /
	path := p.UnresolvedPath()[1:]
	// a/b/c
	// find first /
	idx := strings.Index(path, "/")
	if idx < 0 {
		idx = len(path)
	}
	// a
	tmp := path[:idx]
	// /b/c
	p.unresolved = path[idx:]
	return tmp
}

func (p *PathResolver) UnresolvedPath() string {
	return p.unresolved
}
