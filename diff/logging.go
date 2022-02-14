package diff

import (
	"fmt"
	"reflect"

	"github.com/streamingfast/logging"
)

var zlog, tracer = logging.PackageLogger("solana-go", "github.com/streamingfast/solana-go/diff")

type reflectType struct {
	in interface{}
}

func (r reflectType) String() string {
	if r.in == nil {
		return "<nil>"
	}

	valueOf := reflect.ValueOf(r.in)
	return fmt.Sprintf("%s (zero? %t, value %s)", valueOf.Type(), valueOf.IsZero(), r.in)
}
