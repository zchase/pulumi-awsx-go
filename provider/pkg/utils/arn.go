package utils

import (
	"fmt"
	"strings"
)

type ARN struct {
	ResourceType string
	ResourceID   string
	Partition    string
	Service      string
	Region       string
	AccountID    string
}

func ParseARN(arnString string) (*ARN, error) {
	parts := strings.Split(arnString, ":")

	if parts[0] != "arn" {
		return nil, fmt.Errorf("Invalid ARN: must start with \"arn:\"")
	}

	partsLength := len(parts)

	if (partsLength != 6) && (partsLength != 7) {
		return nil, fmt.Errorf("Invalid ARN: must be between 6 or 7 parts")
	}

	arn := &ARN{
		Partition: parts[1],
		Service:   parts[2],
		Region:    parts[3],
		AccountID: parts[4],
	}

	if partsLength == 7 {
		arn.ResourceType = parts[5]
		arn.ResourceID = parts[6]
		return arn, nil
	}

	slashIndex := strings.Index(parts[5], "/")
	if slashIndex > -1 {
		arn.ResourceType = parts[5][:slashIndex]
		arn.ResourceID = parts[5][slashIndex+1:]
		return arn, nil
	}

	arn.ResourceID = parts[5]
	return arn, nil
}
