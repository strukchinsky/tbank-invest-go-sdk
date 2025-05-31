package investgo

import (
	"strconv"

	"google.golang.org/grpc/metadata"
)

var messageHeaderName string = "message"

func MessageFromMetadata(md metadata.MD) string {
	msgs := md.Get(messageHeaderName)
	if len(msgs) > 0 {
		return msgs[0]
	}
	return ""
}

var trackingIdHeaderName string = "x-tracking-id"

func TrackingIdFromMetadata(md metadata.MD) string {
	msgs := md.Get(trackingIdHeaderName)
	if len(msgs) > 0 {
		return msgs[0]
	}
	return ""
}

var ratelimitRemainingHeaderName string = "x-ratelimit-remaining"

func RemainingLimitFromMetadata(md metadata.MD) int {
	limits := md.Get(ratelimitRemainingHeaderName)
	if len(limits) > 0 {
		lim := limits[0]
		limAsNum, err := strconv.Atoi(lim)
		if err != nil {
			return -1
		}
		return limAsNum
	}
	return -1
}
