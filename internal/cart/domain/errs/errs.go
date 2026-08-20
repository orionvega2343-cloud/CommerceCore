package errs

import "errors"

var FailedCheckedOut = errors.New("cart already checked out")
var InvalidCartItemCap = errors.New("invalid cart item cap")
var ProductNotFound = errors.New("product not found")
var CartNotFound = errors.New("cart not found")
var FailedGetUserIdFromContext = errors.New("failed to get userId from context")
