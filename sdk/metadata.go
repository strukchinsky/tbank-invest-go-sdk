package investgo

import "google.golang.org/grpc/metadata"

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
