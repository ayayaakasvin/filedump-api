package models

import "fmt"

type Status struct {
	State string 	`json:"state"` // either error or success
	Message string 	`json:"message,omitempty"` // in case of error
}

func StatusError(msg string, args... any) Status {
	return Status{
		State: "error",
		Message: fmt.Sprintf(msg, args...),
	}
}

func StatusSuccess() Status {
	return Status{
		State: "success",
	}
}