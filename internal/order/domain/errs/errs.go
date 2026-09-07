package errs

import "errors"

var UnknownStatus = errors.New("unknown status")
var ErrForbidden = errors.New("order forbidden")
var OrderNotPayable = errors.New("order is not in a payable state")
