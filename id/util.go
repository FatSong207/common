package id

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const base = 62

// Int64ToA 將int64轉換為Base62編碼的字符串
func Int64ToA(num int64) string {
	if num == 0 {
		return "0"
	}

	var result []byte
	n := num

	for n > 0 {
		result = append([]byte{alphabet[n%base]}, result...)
		n /= base
	}

	return string(result)
}

// AToInt64 將Base62編碼的字符串轉換為int64
func AToInt64(s string) int64 {
	var result int64

	for _, c := range s {
		var digit int64
		if c >= '0' && c <= '9' {
			digit = int64(c - '0')
		} else if c >= 'A' && c <= 'Z' {
			digit = int64(c - 'A' + 10)
		} else if c >= 'a' && c <= 'z' {
			digit = int64(c - 'a' + 36)
		}
		result = result*base + digit
	}

	return result
}
