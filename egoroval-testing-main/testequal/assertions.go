//go:build !solution

package testequal

import "fmt"

func erf(t T, f interface{}, s interface{}, ma ...interface{}) {
	t.Helper()
	var msg string
	format :=
		`
		expected: %v
        actual  : %v
        message : %v`
	if len(ma) > 0 {
		if len(ma) == 1 {
			msg = ma[0].(string)
		} else {
			msg = fmt.Sprintf(ma[0].(string), ma[1:]...)
		}
	} else {
		msg = ""
	}
	t.Errorf(format, f, s, msg)

}

func isEq(f interface{}, s interface{}) bool {
	if f == nil {
		return s == nil
	}
	if s == nil {
		return f == nil
	}
	switch t := f.(type) {
	case struct{}:
		return false
	case string:
		sc, ok := s.(string)
		if !ok {
			return false
		}
		if sc != t {
			return false
		}
	case int, uint, uint8, int8, uint16, int32, uint32, int16, uint64, int64:
		return f == s
	case map[string]string:
		sm, ok := s.(map[string]string)
		if !ok {
			return false
		}
		if len(t) != len(sm) {
			return false
		}
		if t == nil {
			return sm == nil
		}
		if sm == nil {
			return t == nil
		}
		if len(sm) == 0 {
			return false
		}
		for k, v := range t {
			vv, ok := sm[k]
			if !ok {
				return false
			}
			if v == vv {
				continue
			}
			return false
		}
	case []int:
		si, ok := s.([]int)
		if !ok {
			return false
		}
		if t == nil {
			return si == nil
		}
		if si == nil {
			return t == nil
		}
		if len(t) != len(si) {
			return false
		}
		for k, v := range t {
			vv := si[k]
			if v == vv {
				continue
			}
			return false
		}
	case []byte:
		sb, ok := s.([]byte)
		if !ok {
			return false
		}
		if t == nil {
			return sb == nil
		}
		if sb == nil {
			return t == nil
		}
		if len(t) != len(sb) {
			return false
		}
		for k, v := range t {
			vv := sb[k]
			if v == vv {
				continue
			}
			return false
		}
	default:
		return false
	}
	return true
}

func AssertEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	res := isEq(expected, actual)
	if !res {
		erf(t, expected, actual, msgAndArgs...)
	}
	return res
}

func AssertNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	res := !isEq(expected, actual)
	if !res {
		erf(t, expected, actual, msgAndArgs...)
	}
	return res
}

func RequireEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	res := isEq(expected, actual)
	if !res {
		erf(t, expected, actual, msgAndArgs...)
		t.FailNow()
	}
}

func RequireNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	res := !isEq(expected, actual)
	if !res {
		erf(t, expected, actual, msgAndArgs...)
		t.FailNow()
	}
}
