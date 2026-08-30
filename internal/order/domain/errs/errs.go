package errs

import "errors"

var UnknownStatus = errors.New("unknown status")
var ErrForbidden = errors.New("order forbidden")
