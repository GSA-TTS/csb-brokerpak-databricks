package apps

import (
	"csbbrokerpakdatabricks/acceptance-tests/helpers/testpath"
)

type ManifestCode string

func (a ManifestCode) Path() string {
	return testpath.BrokerpakFile("acceptance-tests", "apps", string(a))
}

func WithTestAppManifest(manifest ManifestCode) Option {
	return func(a *App) {
		a.manifest = manifest.Path()
	}
}
