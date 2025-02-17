package render

type MissingKeyPolicy string

const (
	MissingKeyPolicyZero    MissingKeyPolicy = "zero"
	MissingKeyPolicyDefault MissingKeyPolicy = "default"
	MissingKeyPolicyError   MissingKeyPolicy = "error"
)

func (msp MissingKeyPolicy) String() string {
	return "missingkey=" + string(msp)
}
