package booking

import "time"

type Booking struct {
	Id             int64
	Name           string
	Phone          string
	Email          string
	PackageType    string
	SessionType    string
	Date           time.Time
	Time           time.Time
	Subjects       string
	AdditionalInfo string
	Referral       string

	SubmissionTime time.Time
}
