package v1

import (
	"fmt"
	"maps"
	"strings"

	"github.com/mandelsoft/goutils/jsonutils"
	"github.com/mandelsoft/goutils/maputils"
)

type Variant map[string]string

var (
	_ jsonutils.MapKeyMarshaler   = (Variant)(nil)
	_ jsonutils.MapKeyUnmarshaler = (*Variant)(nil)
)

func (m Variant) Copy() Variant {
	return maps.Clone(m)
}

func (m Variant) String() string {
	key := "{"
	for i, k := range maputils.OrderedKeys(m) {
		if i > 0 {
			key += ","
		}
		v := m[k]
		key += escapeString(k) + "=" + escapeString(v)
	}
	return key + "}"
}

func (m *Variant) Parse(s string) error {
	if len(s) < 2 || s[0] != '{' || s[len(s)-1] != '}' {
		return fmt.Errorf("invalid key %q", s)
	}
	comps := splitSep(s[1 : len(s)-1])
	*m = Variant{}
	for i := 0; i < len(comps); i += 2 {
		(*m)[comps[i]] = comps[i+1]
	}
	return nil
}

func (m *Variant) UnmarshalMapKey(s string) error {
	return m.Parse(s)
}

func (m Variant) MarshalMapKey() (string, error) {
	return m.String(), nil
}

func escapeString(s string) string {
	s = strings.Replace(s, "\\", "\\\\", -1)
	s = strings.Replace(s, "}", "\\}", -1)
	s = strings.Replace(s, "=", "\\=", -1)
	return strings.Replace(s, ",", "\\,", -1)
}

func unescapeString(s string) string {
	s = strings.Replace(s, "\\\\", "\\", -1)
	s = strings.Replace(s, "\\}", "}", -1)
	s = strings.Replace(s, "\\=", "=", -1)
	return strings.Replace(s, "\\,", ",", -1)
}

func splitSep(in string) []string {
	var result []string

	escape := false
	s := ""
	for _, c := range in {
		switch c {
		case '\\':
			if !escape {
				escape = true
				continue
			} else {
				s += string(c)
			}
		case '=', ',':
			if escape {
				s += string(c)
			} else {
				result = append(result, unescapeString(s))
				s = ""
			}
		default:
			s += unescapeString(string(c))
		}
		escape = false
	}
	if len(result)%2 == 1 {
		result = append(result, unescapeString(s))
	}
	return result
}
