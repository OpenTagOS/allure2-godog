package allure

type Status string

const (
	Broken  Status = "broken"
	Passed  Status = "passed"
	Failed  Status = "failed"
	Skipped Status = "skipped"
	Unknown Status = "unknown"
)
