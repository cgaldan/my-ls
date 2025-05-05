package utils

import "os"

// • Base("foo")   == "foo"
func Base(path string) string {
	if path == "" {
		return "."
	}

	// strip trailing separators
	end := len(path)
	for end > 0 && isSep(path[end-1]) {
		end--
	}
	// if it was all separators, return the first one
	if end == 0 {
		return string(path[0])
	}

	// find the start of the base name
	start := end - 1
	for start >= 0 && !isSep(path[start]) {
		start--
	}

	return path[start+1 : end]
}

func Ext(path string) string {
	for i := len(path) - 1; i >= 0 && !isSep(path[i]); i-- {
		if path[i] == '.' {
			return path[i:]
		}
	}
	return ""
}

// Join joins any number of path elements into a single path, adding
// OS-specific separators as needed, ignoring empties, and then
// cleaning the result (removing “.” and “..”).  Just like filepath.Join
// but without importing path/filepath.
func Join(elem ...string) string {
	if len(elem) == 0 {
		return "."
	}
	var path string
	for _, e := range elem {
		if e == "" {
			continue
		}
		// if first non-empty or e is absolute, reset
		if path == "" || isAbs(e) {
			path = e
		} else {
			// ensure exactly one separator
			if isSep(path[len(path)-1]) {
				path += e
			} else {
				path += string(os.PathSeparator) + e
			}
		}
	}
	if path == "" {
		return "."
	}
	return Clean(path)
}

// Clean returns the shortest path name equivalent to path by purely
// lexical processing: no filesystem calls.  It collapses “.”, “..” and
// multiple separators.
func Clean(path string) string {
	if path == "" {
		return "."
	}
	sep := os.PathSeparator
	isSep := func(c byte) bool {
		return c == byte(sep) || (sep == '\\' && c == '/')
	}
	rooted := isSep(path[0])

	// 1) split into components
	var comps []string
	for i := 0; i < len(path); {
		// skip separators
		for i < len(path) && isSep(path[i]) {
			i++
		}
		if i >= len(path) {
			break
		}
		j := i
		for j < len(path) && !isSep(path[j]) {
			j++
		}
		comps = append(comps, path[i:j])
		i = j
	}

	// 2) process “.” and “..”
	var out []string
	for _, comp := range comps {
		switch comp {
		case ".":
			// no-op
		case "..":
			if len(out) > 0 && out[len(out)-1] != ".." {
				out = out[:len(out)-1]
			} else if !rooted {
				out = append(out, "..")
			}
		default:
			out = append(out, comp)
		}
	}

	// 3) rebuild
	result := ""
	if rooted {
		result = string(sep)
	}
	for i, comp := range out {
		if i > 0 {
			result += string(sep)
		}
		result += comp
	}
	if result == "" {
		if rooted {
			return string(sep)
		}
		return "."
	}
	return result
}

// isSep reports whether c is a path separator on this OS.
func isSep(c byte) bool {
	sep := os.PathSeparator
	return c == byte(sep) || (sep == '\\' && c == '/')
}

// isAbs reports whether the provided path starts with a separator.
func isAbs(path string) bool {
	return len(path) > 0 && isSep(path[0])
}

func Dir(path string) string {
	if path == "" {
		return "."
	}

	// is it rooted (first byte is a separator)?
	rooted := isSep(path[0])

	// 1) strip trailing separators
	end := len(path) - 1
	for end >= 0 && isSep(path[end]) {
		end--
	}
	// if nothing left, return "/" or "."
	if end < 0 {
		if rooted {
			return string(path[0])
		}
		return "."
	}

	// 2) find the last separator
	for end >= 0 && !isSep(path[end]) {
		end--
	}
	// no separator found
	if end < 0 {
		if rooted {
			return string(path[0])
		}
		return "."
	}

	// 3) trim any duplicate separators at the end of the result
	dir := path[:end]
	for len(dir) > 1 && isSep(dir[len(dir)-1]) {
		dir = dir[:len(dir)-1]
	}
	return dir
}
