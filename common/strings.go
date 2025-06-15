package common

import "strings"

// 将字符串转换为驼峰命名
func ToCamelCase(s string) string {
	var result string
	capitalizeNext := false
	for i, char := range s {
		if char == '-' || char == '_' {
			capitalizeNext = true
		} else {
			if capitalizeNext || i == 0 {
				result += strings.ToUpper(string(char))
				capitalizeNext = false
			} else {
				result += strings.ToLower(string(char))
			}
		}
	}
	return result
}
