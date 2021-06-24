package allure

import (
	uuid5 "github.com/satori/go.uuid"
)

type Container struct {
	UUID     string   `json:"uuid"`
	Name     string   `json:"name"`
	Children []string `json:"children"`
	Start    int64    `json:"start"`
	Stop     int64    `json:"stop"`
}

func NewContainer() *Container {
	return &Container{
		UUID:  uuid5.NewV4().String(),
		Start: timestampMs(),
	}
}

func (c *Container) AddChildren(testCase *TestCase) {
	c.Children = append(c.Children, testCase.UUID)
}

func (c *Container) Finish() {
	c.Stop = timestampMs()
}
