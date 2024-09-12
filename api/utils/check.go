package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/mandelsoft/goutils/errors"
)

func CheckNonEmpty(s string, msg string, args ...interface{}) error {
	if s == "" {
		return fmt.Errorf("%s is not set or empty", fmt.Sprintf(msg, args...))
	}
	return nil
}

func CheckFlatName(s string, msg string, args ...interface{}) error {
	if s == "" {
		return fmt.Errorf("%s is not set or empty", fmt.Sprintf(msg, args...))
	}
	if strings.Index(s, "/") >= 0 {
		return fmt.Errorf("invalid hierachical name for %s: %s", fmt.Sprintf(msg, args...), s)
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

func CheckVersion(v string, msg string, args ...interface{}) error {
	_, err := semver.NewVersion(v)
	return errors.Wrapf(err, "%s: %s", fmt.Sprintf(msg, args...), v)
}

func CheckVersionConstraint(v string, msg string, args ...interface{}) error {
	_, err := semver.NewConstraint(v)
	return errors.Wrapf(err, "%s: %s", fmt.Sprintf(msg, args...), v)
}

var alnum = regexp.MustCompile("^[0-9a-zA-Z]*$")

func IsAlphaNumeric(s string) bool {
	return alnum.MatchString(s)
}
