package ec2i

import "errors"

var (
	InvalidRegion = errors.New("region cannot be empty")
	InvalidInstanceType = errors.New("invalid instance type %s")
)
