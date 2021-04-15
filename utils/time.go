package utils

import "time"

func GetDuration(t int64, v string) time.Duration {
	switch v {
	case "h":
		return time.Hour * time.Duration(t)
	case "m":
		return time.Minute * time.Duration(t)
	case "s":
		return time.Second * time.Duration(t)
	default:
		return time.Minute * time.Duration(t)
	}
}
