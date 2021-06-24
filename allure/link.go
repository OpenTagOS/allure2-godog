package allure

type LinkType string

const (
	Issue  LinkType = "issue"
	TMS    LinkType = "tms"
	Custom LinkType = "custom"
)

type Link struct {
	Name string   `json:"name,omitempty"`
	Type LinkType `json:"type,omitempty"`
	URL  string   `json:"url,omitempty"`
}
