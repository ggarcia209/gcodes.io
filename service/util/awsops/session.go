package awsops

import (
	"github.com/ggarcia209/go-aws/goaws"
)

// AWSSession wraps the goaws.Session object for use in higher level packages.
type AWSSession struct {
	Session goaws.Session
}

// NewDefaultAWSSession creates a new AWS Session using the defualt credentials and settings.
func NewDefaultAWSSession() AWSSession {
	s := goaws.NewDefaultSession()
	session := AWSSession{Session: s}
	return session
}

// NewAWSSessionWithProfile creates a new AWS Session using the given AWS_PROFILE name to authorize
// with non-default credentials.
func NewAWSSessionWithProfile(profile string) AWSSession {
	s := goaws.NewSessionWithProfile(profile)
	session := AWSSession{Session: s}
	return session
}
