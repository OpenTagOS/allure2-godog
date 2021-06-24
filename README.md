Allure 2 godog integration
---

This library includes formatter that allows you to generate [Allure 2](https://github.com/allure-framework/allure2)
report from [godog](https://github.com/cucumber/godog) test results.

Installation
---

```
go get https://github.com/OpenTagOS/allure2-godog
```

Usage
---

```go
"github.com/OpenTagOS/allure2-godog/allure"
"github.com/OpenTagOS/allure2-godog/alluregodog"

allureWriter := allure.NewReportWriter("/tmp/report/")
godog.Format("allure", "Allure 2 formatter", alluregodog.NewFormatter(allureWriter))

opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Paths:  []string{"."},
    /// other options...
	Format: "allure",
}

status := godog.TestSuite{
	Name: "godogs",
	// other params...
	Options: &opts,
}.Run()
```

Allure report zip archive will be generated in the /tmp/report/ dir after execution.

`WithTagLabelMapping` option is used to map scenario tag to Allure test case label:

```go
tagLabelMapping := map[string]string{
	"issueId": "issue",
}
godog.Format("allure", "Allure formatter", alluregodog.NewFormatter(allureWriter, alluregodog.WithTagLabelMapping(tagLabelMapping)))
```

```
Feature: eat godogs
  In order to be happy
  As a hungry gopher
  I need to be able to eat godogs

  @issueId:PROJECT-2257
  Scenario: Eat 5 out of 12
    Given there are 12 godogs
    When I eat 5
    Then there should be 7 remaining
```
