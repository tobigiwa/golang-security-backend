package http

import "encoding/gob"

func init() {
	gob.Register(&UserResponseModel{})
}
