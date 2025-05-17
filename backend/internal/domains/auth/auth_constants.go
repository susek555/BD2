package auth

import "time"

var AccessTokenExpirationTime = 30 * time.Minute
var RefreshTokenExpirationTime = 30 * 24 * time.Hour
