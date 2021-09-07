package alluregodog

import (
	"io"
	"regexp"

	"github.com/OpenTagOS/allure2-godog/allure"
	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v16"
)

var tagRegexp = regexp.MustCompile(`^@(.*):(.*)$`)

type allureFmt struct {
	writer          Writer
	suiteName       string
	currentFeature  *messages.Feature
	currentScenario *allure.TestCase
	currentStep     *allure.Step
	container       *allure.Container
	tagLabelMapping map[string]string
}

type Writer interface {
	Init() error
	WriteTestCaseResults(testCase *allure.TestCase) error
	WriteContainerResults(container *allure.Container) error
}

type Option func(f *allureFmt)

func WithTagLabelMapping(mapping map[string]string) Option {
	return func(f *allureFmt) {
		f.tagLabelMapping = mapping
	}
}

func NewFormatter(writer Writer, options ...Option) func(suite string, out io.Writer) godog.Formatter {
	return func(suite string, out io.Writer) godog.Formatter {
		f := &allureFmt{suiteName: suite, writer: writer, tagLabelMapping: map[string]string{}}

		for _, option := range options {
			option(f)
		}

		return f
	}
}

func (f *allureFmt) TestRunStarted() {
	if err := f.writer.Init(); err != nil {
		panic(err)
	}

	f.container = allure.NewContainer()
}

func (f *allureFmt) Feature(document *messages.GherkinDocument, s string, bytes []byte) {
	f.currentFeature = document.Feature
}

func (f *allureFmt) Pickle(scenario *godog.Scenario) {
	fullName := scenario.Uri + ":" + scenario.Name
	testCase := allure.NewTestCase(scenario.Name, fullName, f.currentFeature.Description)

	testCase.AddLabel("feature", f.currentFeature.Name)
	testCase.AddLabel("suite", f.suiteName)

	testCase.AddLabels(f.tagsToLabels(scenario.Tags)...)

	f.currentScenario = testCase
	f.container.AddChildren(testCase)
}

func (f *allureFmt) Defined(scenario *godog.Scenario, step *godog.Step, d *godog.StepDefinition) {
	allureStep := allure.NewStep(step.Text)
	allureStep.AddParam(stepArgumentToParam(step.Argument))

	f.currentStep = allureStep
	f.currentScenario.AddStep(*f.currentStep)
}

func stepArgumentToParam(argument *messages.PickleStepArgument) *allure.Parameter {
	if argument == nil {
		return nil
	}

	isString := defArgumentType(argument.DocString)
	if isString != nil {
		return isString
	} else {
		return defArgumentType(argument.DataTable)
	}
}

func defArgumentType(i interface{}) *allure.Parameter {
	switch arg := i.(type) {
	case messages.PickleDocString:
		return &allure.Parameter{Name: "Message", Value: arg.Content}
	case messages.PickleTable:
		for key, value := range firstTableRow(&arg) {
			return &allure.Parameter{Name: key, Value: value}
		}
	default:
		return nil
	}

	return nil
}

func (f *allureFmt) Passed(scenario *godog.Scenario, step *godog.Step, d *godog.StepDefinition) {
	f.currentStep.Finish(allure.Passed)

	if isLastStep(scenario, step) {
		f.currentScenario.Finish(allure.Passed)

		if err := f.writer.WriteTestCaseResults(f.currentScenario); err != nil {
			panic(err)
		}
	}
}

func (f *allureFmt) Undefined(*godog.Scenario, *godog.Step, *godog.StepDefinition) {
	f.currentStep.Finish(allure.Unknown)
	f.currentScenario.Finish(allure.Unknown)

	if err := f.writer.WriteTestCaseResults(f.currentScenario); err != nil {
		panic(err)
	}
}

func (f *allureFmt) Failed(scenario *godog.Scenario, step *godog.Step, d *godog.StepDefinition, err error) {
	f.currentStep.Finish(allure.Failed)
	f.currentScenario.Finish(allure.Failed)
	f.currentScenario.Error(err)

	if err := f.writer.WriteTestCaseResults(f.currentScenario); err != nil {
		panic(err)
	}
}

func (f *allureFmt) Pending(scenario *godog.Scenario, step *godog.Step, d *godog.StepDefinition) {
	f.currentStep.Finish(allure.Skipped)

	// Need to check if it's the last step because after it we can receive Skipped.
	if isLastStep(scenario, step) {
		f.currentScenario.Finish(allure.Skipped)

		if err := f.writer.WriteTestCaseResults(f.currentScenario); err != nil {
			panic(err)
		}
	}
}

// All other steps after Pending will be Skipped.
func (f *allureFmt) Skipped(scenario *godog.Scenario, step *godog.Step, d *godog.StepDefinition) {
	f.currentStep.Finish(allure.Skipped)

	if isLastStep(scenario, step) {
		f.currentScenario.Finish(allure.Skipped)

		if err := f.writer.WriteTestCaseResults(f.currentScenario); err != nil {
			panic(err)
		}
	}
}

func (f *allureFmt) Summary() {
	f.container.Finish()

	if err := f.writer.WriteContainerResults(f.container); err != nil {
		panic(err)
	}
}

func (f *allureFmt) tagsToLabels(tags []*messages.PickleTag) []allure.Label {
	var labels []allure.Label

	for _, tag := range tags {
		matches := tagRegexp.FindStringSubmatch(tag.Name)
		if len(matches) != 3 {
			continue
		}

		tagName, tagValue := matches[1], matches[2]

		labelName, ok := f.tagLabelMapping[tagName]
		if ok {
			labels = append(labels, allure.Label{
				Name:  labelName,
				Value: tagValue,
			})
		}
	}

	return labels
}

func isLastStep(pickle *messages.Pickle, step *messages.PickleStep) bool {
	return pickle.Steps[len(pickle.Steps)-1].Id == step.Id
}

func firstTableRow(tableData *messages.PickleTable) map[string]string {
	mapData := make(map[string]string)

	if len(tableData.Rows) < 2 {
		return mapData
	}

	for i, header := range tableData.Rows[0].Cells {
		mapData[header.Value] = tableData.Rows[1].Cells[i].Value
	}

	return mapData
}
