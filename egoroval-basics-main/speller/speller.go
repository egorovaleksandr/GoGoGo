//go:build !solution

package speller

func thouthConv(n int64) string {
	var result string
	m := n % 100
	d := n / 100
	dozens := map[int64]string{2: "twenty", 3: "thirty", 4: "forty", 5: "fifty", 6: "sixty", 7: "seventy", 8: "eighty", 9: "ninety"}
	uptonineteen := []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen"}
	if d > 0 {
		result += uptonineteen[d]
		result += " hundred"
	}
	if m < 20 && m != 0 {
		if d > 0 {
			result += " "
		}
		result += uptonineteen[m]
		return result
	}
	if m > 0 && d > 0 {
		result += " "
	}
	result += dozens[m/10]
	if m%10 > 0 {
		result += "-"
		result += uptonineteen[m%10]
	}
	return result
}

func Spell(n int64) string {
	var result string
	isNeg := false
	s := []string{}
	c := map[int]string{1: "thousand", 2: "million", 3: "billion"}
	var cap string
	if n == 0 {
		return "zero"
	}
	if n < 0 {
		isNeg = true
		n = -n
	}
	for n > 0 {
		s = append(s, thouthConv(n%1000))
		n /= 1000
	}
	if isNeg {
		result += "minus"
		result += " "
	}
	l := len(s)
	for i := l - 1; i >= 0; i-- {
		cap = ""
		if v, ok := c[i]; ok {
			cap = v
		}
		if i != l-1 && s[i] != "" {
			result += " "
		}
		result += s[i]
		if l > 1 && s[i] != "" && cap != "" {
			result += " "
			result += cap
		}
	}
	return result
}
