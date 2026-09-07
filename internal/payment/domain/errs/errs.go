package errs

import "errors"

var InvalidStatus = errors.New("invalid payment status")
var InvalidUserOrAdmin = errors.New("invalid user or admin")
