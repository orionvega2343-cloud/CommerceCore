package errs

import "errors"

// Validate
var InvalidProductName = errors.New("invalid product name")
var InvalidPrice = errors.New("price is negative or zero")
var InvalidQuantity = errors.New("quantity is negative or zero")

// Redis
var ProductNotFound = errors.New("product not found")

// Service
var InvalidRole = errors.New("invalid role")

var InsufficientStock = errors.New("insufficient stock")
