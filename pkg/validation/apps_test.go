package validation_test

import (
	"os"
	"testing"

	"github.com/codeready-toolchain/argocd-checker/pkg/validation"

	charmlog "github.com/charmbracelet/log"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckApplications(t *testing.T) {

	t.Run("success", func(t *testing.T) {

		t.Run("empty apps", func(t *testing.T) {
			// given
			logger := NewTestLogger(os.Stdout, charmlog.Options{
				Level: charmlog.InfoLevel,
			})
			afs := afero.Afero{
				Fs: afero.NewMemMapFs(),
			}
			err := afs.Mkdir("/path/to/apps", os.ModeDir)
			require.NoError(t, err)

			// when
			err = validation.CheckApplications(logger, afs, "/path/to", "apps")

			// then
			require.NoError(t, err)
			assert.Empty(t, logger.Errors())
			assert.Empty(t, logger.Warnings())
		})

		t.Run("empty kustomization", func(t *testing.T) {
			// given
			logger := NewTestLogger(os.Stdout, charmlog.Options{
				Level: charmlog.DebugLevel,
			})

			afs := afero.Afero{
				Fs: afero.NewMemMapFs(),
			}

			err := afs.MkdirAll("/path/to/apps", 0755)
			require.NoError(t, err)
			err = addFile(afs, "/path/to/apps/kustomization.yaml", `kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1`)
			require.NoError(t, err)

			// when
			err = validation.CheckApplications(logger, afs, "/path/to", "apps")

			// then
			require.Error(t, err, "kustomization.yaml is empty")
			assert.Empty(t, logger.Errors())
			assert.Empty(t, logger.Warnings())
		})

		t.Run("kustomization with valid apps", func(t *testing.T) {
			// given
			logger := NewTestLogger(os.Stdout, charmlog.Options{
				Level: charmlog.InfoLevel,
			})

			afs := afero.Afero{
				Fs: afero.NewMemMapFs(),
			}

			err := afs.MkdirAll("/path/to/apps", 0755)
			require.NoError(t, err)
			err = addFile(afs, "/path/to/apps/kustomization.yaml", `kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1
resources:
- app-cookie.yaml
- appset-pasta.yaml`)
			require.NoError(t, err)
			err = addFile(afs, "/path/to/apps/app-cookie.yaml", `apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: app-cookie
spec:
  destination:
    server: https://kubernetes.default.svc
  project: default
  source:
    path: components/cookie`)
			require.NoError(t, err)
			err = addFile(afs, "/path/to/components/cookie/kustomization.yaml", `kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1`)
			require.NoError(t, err)

			err = addFile(afs, "/path/to/apps/appset-pasta.yaml", `apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: appset-pasta
spec:
  template:
    spec:
      destination:
        server: https://kubernetes.default.svc
      project: default
      source:
        path: components/pasta`)
			require.NoError(t, err)
			err = addFile(afs, "/path/to/components/pasta/kustomization.yaml", `kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1`)
			require.NoError(t, err)

			// when
			err = validation.CheckApplications(logger, afs, "/path/to", "apps")

			// then
			require.NoError(t, err)
			assert.Empty(t, logger.Errors())
			assert.Empty(t, logger.Warnings())
		})

		t.Run("kustomization with invalid app", func(t *testing.T) {
			t.Run("unknown source path", func(t *testing.T) {

				// given
				logger := NewTestLogger(os.Stdout, charmlog.Options{
					Level: charmlog.InfoLevel,
				})

				afs := afero.Afero{
					Fs: afero.NewMemMapFs(),
				}

				err := afs.MkdirAll("/path/to/apps", 0755)
				require.NoError(t, err)
				err = addFile(afs, "/path/to/apps/kustomization.yaml", `kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1
resources:
- app-cookie.yaml`)
				require.NoError(t, err)

				err = addFile(afs, "/path/to/apps/app-cookie.yaml", `apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: app-cookie
spec:
  destination:
    server: https://kubernetes.default.svc
  project: default
  source:
    path: components/cookie # path does not exist/does not contain a 'kustomization.yaml' file`)
				require.NoError(t, err)

				// when
				err = validation.CheckApplications(logger, afs, "/path/to", "apps")

				// then
				require.EqualError(t, err, "components/cookie is not valid")
				assert.Empty(t, logger.Errors())
				assert.Empty(t, logger.Warnings())
			})

			t.Run("missing component kustomization.yaml", func(t *testing.T) {

				// given
				logger := NewTestLogger(os.Stdout, charmlog.Options{
					Level: charmlog.InfoLevel,
				})

				afs := afero.Afero{
					Fs: afero.NewMemMapFs(),
				}

				err := afs.MkdirAll("/path/to/apps", 0755)
				require.NoError(t, err)
				err = addFile(afs, "/path/to/apps/kustomization.yaml", `kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1
resources:
- app-cookie.yaml`)
				require.NoError(t, err)

				err = addFile(afs, "/path/to/apps/app-cookie.yaml", `apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: app-cookie
spec:
  destination:
    server: https://kubernetes.default.svc
  project: default
  source:
    path: components/cookie # path does not exist/does not contain a 'kustomization.yaml' file`)
				require.NoError(t, err)
				// empty dir
				err = addDir(afs, "/path/to/components/cookie")
				require.NoError(t, err)

				// when
				err = validation.CheckApplications(logger, afs, "/path/to", "apps")

				// then
				require.EqualError(t, err, "components/cookie does not contain a 'kustomization.yaml' file")
				assert.Empty(t, logger.Errors())
				assert.Empty(t, logger.Warnings())
			})
		})

		t.Run("kustomization with invalid appset", func(t *testing.T) {
			t.Run("unknown source path", func(t *testing.T) {

				// given
				logger := NewTestLogger(os.Stdout, charmlog.Options{
					Level: charmlog.InfoLevel,
				})

				afs := afero.Afero{
					Fs: afero.NewMemMapFs(),
				}

				err := afs.MkdirAll("/path/to/apps", 0755)
				require.NoError(t, err)
				err = addFile(afs, "/path/to/apps/kustomization.yaml", `kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1
resources:
- appset-cookie.yaml`)
				require.NoError(t, err)

				err = addFile(afs, "/path/to/apps/appset-cookie.yaml", `apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: app-cookie
spec:
  template:
    spec:
      destination:
        server: https://kubernetes.default.svc
      project: default
      source:
        path: components/cookie # path does not exist/does not contain a 'kustomization.yaml' file`)
				require.NoError(t, err)

				// when
				err = validation.CheckApplications(logger, afs, "/path/to", "apps")

				// then
				require.EqualError(t, err, "components/cookie is not valid")
				assert.Empty(t, logger.Errors())
				assert.Empty(t, logger.Warnings())
			})

			t.Run("missing component kustomization.yaml", func(t *testing.T) {

				// given
				logger := NewTestLogger(os.Stdout, charmlog.Options{
					Level: charmlog.InfoLevel,
				})

				afs := afero.Afero{
					Fs: afero.NewMemMapFs(),
				}

				err := afs.MkdirAll("/path/to/apps", 0755)
				require.NoError(t, err)
				err = addFile(afs, "/path/to/apps/kustomization.yaml", `kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1
resources:
- appset-cookie.yaml`)
				require.NoError(t, err)

				err = addFile(afs, "/path/to/apps/appset-cookie.yaml", `apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: appset-cookie
spec:
  template:
    spec:
      destination:
        server: https://kubernetes.default.svc
      project: default
      source:
        path: components/cookie # path does not exist/does not contain a 'kustomization.yaml' file`)
				require.NoError(t, err)
				// empty dir
				err = addDir(afs, "/path/to/components/cookie")
				require.NoError(t, err)

				// when
				err = validation.CheckApplications(logger, afs, "/path/to", "apps")

				// then
				require.EqualError(t, err, "components/cookie does not contain a 'kustomization.yaml' file")
				assert.Empty(t, logger.Errors())
				assert.Empty(t, logger.Warnings())
			})
		})
	})

}

func addFile(afs afero.Afero, path string, data string) error {
	return afs.WriteFile(path, []byte(data), 0755)
}

func addDir(afs afero.Afero, path string) error {
	return afs.Mkdir(path, 0755)
}
