package utils

import (
	"fmt"
	"strings"
)

func CheckNonEmpty(s string, msg string) error {
	if s == "" {
		return fmt.Errorf("%s is not set or empty", msg)
	}
	return nil
}

func CheckFlatName(s string, msg string) error {
	if s == "" {
		return fmt.Errorf("%s is not set or empty", msg)
	}
	if strings.Index(s, "/") >= 0 {
		return fmt.Errorf("invalid hierachical name for %s: %s", msg, s)
	}
	return nil
}

func CheckValues(s string, msg string, values ...string) error {
	if len(values) == 0 {
		return nil
	}
	for _, v := range values {
		if s == v {
			return nil
		}
	}
	if s == "" {
		return fmt.Errorf("%s is not set or empty", msg)
	}
	return fmt.Errorf("invalid %s: %s (expected one ov %s)", msg, s, strings.Join(values, ", "))
}
