package allure

import (
	"runtime"

	uuid5 "github.com/satori/go.uuid"
)

type TestCase struct {
	UUID          string         `json:"uuid,omitempty"`
	Name          string         `json:"name,omitempty"`
	Description   string         `json:"description,omitempty"`
	Status        Status         `json:"status,omitempty"`
	StatusDetails *StatusDetails `json:"statusDetails,omitempty"`
	Stage         string         `json:"stage,omitempty"`
	Steps         []Step         `json:"steps,omitempty"`
	Parameters    []Parameter    `json:"parameters,omitempty"`
	Start         int64          `json:"start,omitempty"`
	Stop          int64          `json:"stop,omitempty"`
	Children      []string       `json:"children,omitempty"`
	FullName      string         `json:"fullName,omitempty"`
	Labels        []Label        `json:"labels,omitempty"`
	Links         []Link         `json:"links,omitempty"`
}

type StatusDetails struct {
	Known   bool   `json:"known,omitempty"`
	Muted   bool   `json:"muted,omitempty"`
	Flaky   bool   `json:"flaky,omitempty"`
	Message string `json:"message,omitempty"`
	Trace   string `json:"trace,omitempty"`
}

type Label struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

func NewTestCase(name, fullName, description string) *TestCase {
	return &TestCase{
		UUID:        uuid5.NewV4().String(),
		Name:        name,
		Description: description,
		FullName:    fullName,
		Start:       timestampMs(),
	}
}

func (t *TestCase) AddStep(step Step) {
	t.Steps = append(t.Steps, step)
}

func (t *TestCase) Finish(status Status) {
	t.Status = status
	t.Stage = "finished"
	t.Stop = timestampMs()
}

func (t *TestCase) Error(err error) {
	b := make([]byte, 2048)
	n := runtime.Stack(b, false)
	s := string(b[:n])

	t.StatusDetails = &StatusDetails{
		Message: err.Error(),
		Trace:   s,
	}
}

func (t *TestCase) AddLabel(name string, value string) {
	t.Labels = append(t.Labels, Label{Name: name, Value: value})
}

func (t *TestCase) AddLabels(labels ...Label) {
	t.Labels = append(t.Labels, labels...)
}
