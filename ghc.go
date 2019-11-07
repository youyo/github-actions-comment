package ghc

import (
	"bytes"
	"html/template"
	"net/http"
	"os"

	pipeline "github.com/mattn/go-pipeline"
)

type Ghc struct {
	Workflow    string
	Action      string
	Token       string
	Url         string
	EventPath   string
	RequestBody string
}

// Template, use .Title .Body .Workflow .Action
const Template string = "#### `{{.Title}}` `{{.Status}}`" + `

<details>
<summary>Show Output</summary>

` + "```" + `
{{.Body}}
` + "```" + `

</details>

` + "*Workflow: `{{.Workflow}}`, Action: `{{.Action}}*"

const (
	Success bool = true
	Failure bool = false
)

// New
func New() *Ghc {
	workflow := os.Getenv("GITHUB_WORKFLOW")
	action := os.Getenv("GITHUB_ACTION")
	token := os.Getenv("GITHUB_TOKEN")
	eventPath := os.Getenv("GITHUB_EVENT_PATH")

	return &Ghc{
		Workflow:  workflow,
		Action:    action,
		Token:     token,
		EventPath: eventPath,
	}
}

func (g *Ghc) GetCommentUrl() error {
	url, err := getCommentUrl(g.EventPath)
	if err != nil {
		return err
	}

	g.Url = url

	return nil
}

func (g *Ghc) GenerateComment(title, body string, status bool) (*bytes.Buffer, error) {
	tmpl, err := template.New("test").Parse(Template)
	if err != nil {
		return nil, err
	}

	var result string
	if status {
		result = "Success"
	} else {
		result = "Failure"
	}

	params := map[string]string{
		"Title":    title,
		"Body":     body,
		"Status":   result,
		"Workflow": g.Workflow,
		"Action":   g.Action,
	}

	var comment *bytes.Buffer
	if err := tmpl.Execute(comment, params); err != nil {
		return nil, err
	}

	return comment, nil
}

func (g *Ghc) PostComment(comment *bytes.Buffer) error {
	req, err := http.NewRequest("POST", g.Url, comment)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "token "+g.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return err
}

func getCommentUrl(filepath string) (string, error) {
	out, err := pipeline.Output(
		[]string{"cat", filepath},
		[]string{"jq", "-r", ".pull_request.comments_url"},
	)
	if err != nil {
		return "", err
	}

	output := string(out)

	return output, nil
}
