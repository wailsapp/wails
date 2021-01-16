package str

import (
	"fmt"
	"time"
)

func UnixNow() string {
	return fmt.Sprintf("%+v", time.Now().Unix())
}
