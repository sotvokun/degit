package version

import "time"

const (
	FallbackValue = "unknown"
)

var (
	version string
	commit  string
	date    string
)

type VersionService struct{}

func NewVersionService() *VersionService {
	return &VersionService{}
}

func (v *VersionService) Version() string {
	if version == "" {
		return FallbackValue
	}
	return version
}

func (v *VersionService) Commit() string {
	if commit == "" {
		return FallbackValue
	}
	return commit
}

func (v *VersionService) Date() time.Time {
	if date == "" {
		return time.Unix(0, 0).UTC()
	}

	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return time.Unix(0, 0).UTC()
	}
	return t
}
