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

func TestNewInMemory(t *testing.T) {

	// given
	logger := NewTestLogger(os.Stdout, charmlog.Options{
		Level: charmlog.InfoLevel,
	})
	afs := afero.Afero{
		Fs: afero.NewMemMapFs(),
	}
	err := afs.MkdirAll("/basedir/apps", 0755)
	require.NoError(t, err)
	data := []byte("cookies are yummy")
	err = afs.WriteFile("/basedir/apps/kustomization.yaml", data, 0755)
	require.NoError(t, err)

	// when
	fsys, err := validation.NewInMemoryFS(logger, afs, "/basedir")

	// then
	require.NoError(t, err)

	assert.True(t, fsys.Exists("/basedir/apps"))
	assert.True(t, fsys.Exists("/basedir/apps/kustomization.yaml"))
	actual, err := fsys.ReadFile("/basedir/apps/kustomization.yaml")
	require.NoError(t, err)
	assert.Equal(t, data, actual)
}
