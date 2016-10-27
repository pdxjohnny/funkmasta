package url

import (
	"net/url"
	"strings"
)

func EnvArray(v url.Values) []string {
	if v == nil {
		return nil
	}

	a := make([]string, 0)
	for k, val := range v {
		switch len(val) {
		case 1:
			a = append(a, k+"="+val[0])
			break
		default:
			a = append(a, k+"="+strings.Join(val, " "))
			break
		}
	}
	return a
}
