/*
 * Copyright 2021 Layotto Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
