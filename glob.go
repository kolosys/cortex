package cortex

import "strings"

// matchGlob reports whether name matches pattern.
// * matches within a single '/' segment; ** matches across segments.
func matchGlob(pattern, name string) bool {
	return matchSegs(splitSegs(pattern), splitSegs(name))
}

func splitSegs(p string) []string {
	if p == "" {
		return nil
	}
	parts := strings.Split(p, "/")
	n := 0
	for _, s := range parts {
		if s != "" {
			parts[n] = s
			n++
		}
	}
	return parts[:n]
}

func matchSegs(pat, name []string) bool {
	for len(pat) > 0 {
		if pat[0] == "**" {
			for len(pat) > 0 && pat[0] == "**" {
				pat = pat[1:]
			}
			if len(pat) == 0 {
				return true
			}
			for i := 0; i <= len(name); i++ {
				if matchSegs(pat, name[i:]) {
					return true
				}
			}
			return false
		}
		if len(name) == 0 {
			return false
		}
		if !matchSeg(pat[0], name[0]) {
			return false
		}
		pat = pat[1:]
		name = name[1:]
	}
	return len(name) == 0
}

func matchSeg(pat, s string) bool {
	pi, si := 0, 0
	star, match := -1, 0
	for si < len(s) {
		if pi < len(pat) && pat[pi] == s[si] {
			pi++
			si++
			continue
		}
		if pi < len(pat) && pat[pi] == '*' {
			star = pi
			match = si
			pi++
			continue
		}
		if star >= 0 {
			pi = star + 1
			match++
			si = match
			continue
		}
		return false
	}
	for pi < len(pat) && pat[pi] == '*' {
		pi++
	}
	return pi == len(pat)
}
