package apps

import (
	"csbbrokerpakdatabricks/acceptance-tests/helpers/testpath"
)

type AppCode string

const (
	Databricks AppCode = "databricksapp"
)

func (a AppCode) Dir() string {
	return testpath.BrokerpakFile("acceptance-tests", "apps", string(a))
}

func WithApp(app AppCode) Option {
	return WithGoPreBuild(app.Dir())
}
