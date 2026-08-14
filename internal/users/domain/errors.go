package domain

import "errors"

var InvalidPassword = errors.New("invalid password")
var FailedCreatedUser = errors.New("failed to create user")
var FailedHashingPassword = errors.New("failed to hash password")
var FailedToGetUser = errors.New("failed to get user")
var FailedToUpdateUser = errors.New("failed to update user")
var InvalidUserRole = errors.New("invalid user role")
