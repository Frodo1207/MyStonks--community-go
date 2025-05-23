package utils

import "strings"

func CallerFormater(i interface{}) string {
	s, ok := i.(string)
	if !ok {
		return ""
	}
	s1 := strings.Split(s, ":")
	s = ""
	for i := 0; i < len(s1)-1; i++ {
		s += s1[i] + ":"
	}
	s2 := strings.Split(s, "/")
	s = ""
	p := len(s2) - 3
	if p < 0 {
		p = 0
	}
	for i := p; i < len(s2); i++ {
		s += "/" + s2[i]
	}
	return s + s1[len(s1)-1]
}
