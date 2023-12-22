package validation

import (
	"fmt"
	iofs "io/fs"
	"path/filepath"

	"github.com/spf13/afero"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// Look for all YAML files in the given paths and when the contents if an Argo CD Application or ApplicationSet,
// verify that the `spec.source.path` matches an existing component
func CheckApplications(logger Logger, afs afero.Afero, baseDir string, apps ...string) error {

	for _, path := range apps {
		p := filepath.Join(baseDir, path)
		logger.Info("👀 checking Applications and ApplicationSets", "path", p)
		fsys, err := NewInMemoryFS(logger, afs, p)
		if err != nil {
			return err
		}
		if err := afs.Walk(p, func(path string, info iofs.FileInfo, err error) error {
			if err != nil {
				logger.Error("prevent panic by handling failure", "path", path)
				return err
			}
			if info.IsDir() {
				logger.Debug("👀 checking contents", "path", path)
				if kp, found := lookupKustomizationFile(logger, afs, path); found {
					if err := checkKustomizeResources(logger, afs, kp); err != nil {
						return err
					}
					if info.Name() != "base" {
						if err := checkBuild(logger, fsys, path); err != nil {
							return err
						}
					}
				}
				return nil
			}
			data, err := afs.ReadFile(path)
			if err != nil {
				return err
			}
			if filepath.Ext(info.Name()) == ".yaml" {
				logger.Debug("checking contents", "path", path)
				obj, err := yaml.Parse(string(data))
				if err != nil {
					return err
				}
				switch {
				case obj.GetKind() == "Application":
					path, err := obj.GetString("spec.source.path")
					if err != nil {
						return err
					}
					return checkPath(logger, afs, baseDir, path)
				case obj.GetKind() == "ApplicationSet":
					path, err := obj.GetString("spec.template.spec.source.path")
					if err != nil {
						return err
					}
					return checkPath(logger, afs, baseDir, path)
				default:
					// logger.Debug("ignoring %s (%T)", path, obj)
				}
			}
			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}

func checkPath(_ Logger, afs afero.Afero, repoURL, path string) error {
	p := filepath.Join(repoURL, path)
	if _, err := afs.ReadDir(p); err != nil {
		return fmt.Errorf("%s is not valid", path)
	}
	// logger.Debugf("%s is valid", path)
	return nil
}
