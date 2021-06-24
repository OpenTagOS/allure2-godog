package allure

type Step struct {
	Name       string      `json:"name,omitempty"`
	Status     Status      `json:"status,omitempty"`
	Stage      string      `json:"stage"`
	Steps      []Step      `json:"steps"`
	Parameters []Parameter `json:"parameters"`
	Start      int64       `json:"start"`
	Stop       int64       `json:"stop"`
}

func NewStep(name string) *Step {
	return &Step{
		Name:  name,
		Start: timestampMs(),
	}
}

func (s *Step) AddParam(param *Parameter) {
	if param == nil {
		return
	}

	s.Parameters = append(s.Parameters, *param)
}

func (s *Step) Finish(status Status) {
	s.Status = status
	s.Stage = "finished"
	s.Stop = timestampMs()
}
