package http_transport

import "fmt"

func GetV1Prefix(prefix string) string {
	return fmt.Sprintf("/api/v1/%s", prefix)
}