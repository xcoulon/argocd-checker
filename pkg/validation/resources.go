package validation

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"sigs.k8s.io/kustomize/api/types"
)

// Compares the entries of `resources` in the Kustomize file with the contents in the current directory to see if
// any local file is missing (not referenced as a resource).
// Files starting with an underscore character (`_`) are ignored
func checkKustomizeResources(logger Logger, afs afero.Afero, basedir, kpath string) error {
	logger.Debug("checking kustomization resource", "path", kpath)
	data, err := afs.ReadFile(kpath)
	if err != nil {
		return err
	}
	var kobj types.Kustomization
	if err := kobj.Unmarshal(data); err != nil {
		return err
	}

	// list resources
	logger.Debug("checking kustomization resources", "dir", filepath.Dir(kpath))
	entries, err := afs.ReadDir(filepath.Dir(kpath))
	if err != nil {
		return err
	}
entries:
	for _, e := range entries {
		switch {
		case e.IsDir():
			fallthrough
		case e.Name() == filepath.Base(kpath):
			fallthrough
		case strings.HasPrefix(e.Name(), "_"):
			fallthrough
		case !(filepath.Ext(e.Name()) == ".yaml" || filepath.Ext(e.Name()) == ".yml"):
			fallthrough
		case filepath.Base(e.Name()) == "kustomization.yaml":
			logger.Debug("ignoring file", "path", kpath)
			continue entries
		}
		for _, r := range kobj.Resources {
			if r == e.Name() {
				continue entries
			}
		}
		for _, sg := range kobj.ConfigMapGenerator {
			for _, f := range sg.FileSources {
				if i := strings.LastIndex(f, "="); i > 0 {
					if f[i+1:] == e.Name() {
						continue entries
					}
				} else if f == e.Name() {
					continue entries
				}
			}
		}
		for _, sg := range kobj.SecretGenerator {
			for _, f := range sg.FileSources {
				if i := strings.LastIndex(f, "="); i > 0 {
					if f[i+1:] == e.Name() {
						continue entries
					}
				} else if f == e.Name() {
					continue entries
				}
			}
		}
		for _, m := range kobj.PatchesStrategicMerge { //nolint:staticcheck
			if string(m) == e.Name() {
				continue entries
			}
		}
		for _, m := range kobj.Patches {
			if string(m.Path) == e.Name() {
				continue entries
			}
		}
		for _, m := range kobj.Transformers {
			if string(m) == e.Name() {
				continue entries
			}
		}
		rkpath, _ := filepath.Rel(basedir, kpath)

		return fmt.Errorf("resource is not referenced in %s: %s", rkpath, e.Name())
	}
	return nil
}
