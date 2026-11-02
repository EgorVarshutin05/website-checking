package service

func GetStatusUD(statusCode int, err error) string {
	if statusCode > 399 || err != nil || statusCode < 200 {
		return "down"
	}
	return "up"
}
