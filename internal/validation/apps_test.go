package validation_test

import (
	"os"
	"testing"

	"github.com/codeready-toolchain/argocd-checker/internal/validation"

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
			require.NoError(t, err)
			assert.Empty(t, logger.Errors())
			assert.Empty(t, logger.Warnings())
		})

		t.Run("kustomization with resources", func(t *testing.T) {
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
- configmap-1.yaml
- configmap-2.yaml`)
			require.NoError(t, err)
			err = addFile(afs, "/path/to/apps/configmap-1.yaml", `apiVersion: v1
kind: ConfigMap
metadata:
  namespace: test
  name: cm-1
data:
  cookie: yummy`)
			require.NoError(t, err)
			err = addFile(afs, "/path/to/apps/configmap-2.yaml", `apiVersion: v1
kind: ConfigMap
metadata:
  namespace: test
  name: cm-2
data:
  pasta: yummy`)
			require.NoError(t, err)

			// when
			err = validation.CheckApplications(logger, afs, "/path/to", "apps")

			// then
			require.NoError(t, err)
			assert.Empty(t, logger.Errors())
			assert.Empty(t, logger.Warnings())
		})

		t.Run("kustomization with overlays", func(t *testing.T) {
			// given
			logger := NewTestLogger(os.Stdout, charmlog.Options{
				Level: charmlog.InfoLevel,
			})

			afs := afero.Afero{
				Fs: afero.NewMemMapFs(),
			}

			err := afs.MkdirAll("/path/to/apps", 0755)
			require.NoError(t, err)
			err = addFile(afs, "/path/to/apps/dev/kustomization.yaml", `kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1
resources:
- ../base`)
			require.NoError(t, err)
			err = addFile(afs, "/path/to/apps/base/kustomization.yaml", `kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1
resources:
- configmap-1.yaml`)
			require.NoError(t, err)
			err = addFile(afs, "/path/to/apps/base/configmap-1.yaml", `apiVersion: v1
kind: ConfigMap
metadata:
  namespace: test
  name: cm-1
data:
  cookie: yummy`)
			require.NoError(t, err)

			// when
			err = validation.CheckApplications(logger, afs, "/path/to", "apps")

			// then
			require.NoError(t, err)
			assert.Empty(t, logger.Errors())
			assert.Empty(t, logger.Warnings())
		})
	})

}

func addFile(afs afero.Afero, path string, data string) error {
	return afs.WriteFile(path, []byte(data), 0755)
}
