package errors


var (
	ErrGroupNotFound = &SherlockErrTemplate{
		ErrorCode:   5,
		ErrorReason: "group not found",
	}

	ErrAccountNotFound = &SherlockErrTemplate{
		ErrorCode:   5,
		ErrorReason: "account not found",
	}
)
