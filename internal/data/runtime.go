package data


import (
	"fmt"
	"strconv"
)

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error){

	js := fmt.Sprintf("%d mins", r)
	quotedjs := strconv.Quote(js)

	return []byte(quotedjs), nil

}
