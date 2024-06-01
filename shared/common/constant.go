package common

// $(*%S$FDd!3)96|12AP&LR
const (
	STATUS_CREATED               = 201
	STATUS_DUPLICATE             = 409
	STATUS_OK                    = 200
	STATUS_FOUND                 = 302
	STATUS_BAD_REQUEST           = 400
	STATUS_UNAUTHORIZED          = 401
	STATUS_INTERNAL_SERVER_ERROR = 500
)

const (
	ONE_TO_ONE_COMMUNICATION = "one-to-one"
	PRIVATE_COMMUNICATION    = "private"
)

const (
	META_SUCCESS = 1
	META_FAILED  = 0
)

const (
	PAGE        = "1"
	PAGE_OFFSET = "20"
)

const (
	POST_NOTIFICATION_EVENT = "post-notification"
	GET_NOTIFICATION_EVENT  = "get-notification"
	JOIN_ROOM_EVENT         = "join-room"
	SEND_NEW_MESSAGE_EVENT  = "send-new-message"
	READ_NEW_MESSAGE_EVENT  = "read-new-message"
	STATUS_AWAY_EVENT       = "status-away"
	STATUS_ONLINE_EVENT     = "status-online"
)

const (
	USER_STATUS_ONLINE         = 1
	USER_STATUS_OFFLINE        = 2
	USER_STATUS_AWAY           = 3
	USER_STATUS_DO_NOT_DISTURB = 4
)
