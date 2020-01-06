package encoding


import "strings"

func EncodeMetadata(md map[string]string) []string {
	result := make([]string, 0, len(md))
	for k, v := range md {
		result = append(result, k+"-"+v)
	}
	return result
}

func DecodeMetadata(tags []string) map[string]string {

	result :=make(map[string]string,16)
	for _, tag := range tags {

		parts := strings.Split(tag, "-")
		if len(parts) == 1 {
			result[parts[0]] = ""
		} else {
			result[parts[0]] = parts[1]
		}
	}
	return result
}

