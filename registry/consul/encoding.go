package consul

import "strings"

func decodeVersion(tags []string) (string, bool) {
	for _, tag := range tags {
		if len(tag) < 2 || tag[0] != 'v' {
			continue
		}

		if tag[1] == '-' {
			return tag[2:], true
		}
	}
	return "", false
}


func encodeMetadata(md map[string]string) []string {
	result := make([]string, 0, len(md))
	for k, v := range md {
		result = append(result, k+"-"+v)
	}
	return result
}

func decodeMetadata(tags []string) map[string]string {

	result :=make(map[string]string,16)
	for _, tag := range tags {

		parts := strings.Split(tag, "-")
		if len(parts) == 1 {
			result[parts[0]] = ""
		} else {

			if parts[0]=="version" || parts[0]=="endpoints" {
				continue
			}
			result[parts[0]] = parts[1]
		}
	}
	return result
}