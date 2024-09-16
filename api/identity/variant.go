package identity

import (
	"fmt"
	"maps"
	"strings"

	"github.com/mandelsoft/goutils/jsonutils"
	"github.com/mandelsoft/goutils/maputils"
	"github.com/open-component-model/service-model/api/utils"
)

type Variant map[string]string

var (
	_ jsonutils.MapKeyMarshaler   = (Variant)(nil)
	_ jsonutils.MapKeyUnmarshaler = (*Variant)(nil)
)

func (m Variant) Copy() Variant {
	return maps.Clone(m)
}

func (m Variant) Equals(o Variant) bool {
	if len(m) != len(o) {
		return false
	}
	for k, v := range m {
		if ov, ok := o[k]; !ok || ov != v {
			return false
		}
	}
	return true
}

func (m Variant) Validate() error {
	for k, v := range m {
		if !utils.IsAlphaNumeric(k) {
			return fmt.Errorf("variant key %q is not alpha numeric", k)
		}
		if !utils.IsAlphaNumeric(v) {
			return fmt.Errorf("variant value %q[%s] is not alpha numeric", v, k)
		}
	}
	return nil
}

func (m Variant) String() string {
	if len(m) == 0 {
		return ""
	}
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
