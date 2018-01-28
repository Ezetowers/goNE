package processing

import (
	"github.com/op/go-logging"
)

var Log *logging.Logger

func init() {
	Log = logging.MustGetLogger("")
}
