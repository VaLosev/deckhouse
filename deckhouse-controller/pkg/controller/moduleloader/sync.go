// Copyright 2024 Flant JSC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package moduleloader

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/deckhouse/deckhouse/deckhouse-controller/pkg/apis/deckhouse.io/v1alpha1"
	"github.com/deckhouse/deckhouse/deckhouse-controller/pkg/controller/module-controllers/downloader"
	"github.com/deckhouse/deckhouse/deckhouse-controller/pkg/controller/module-controllers/utils"
)

// // restoreAbsentModulesFromOverrides checks ModulePullOverrides and restore them on the FS
func (l *Loader) restoreAbsentModulesFromOverrides(ctx context.Context) error {
	currentNodeName := os.Getenv("DECKHOUSE_NODE_NAME")
	if len(currentNodeName) == 0 {
		return errors.New("determine the node name deckhouse pod is running on: missing or empty DECKHOUSE_NODE_NAME env")
	}

	mpos := new(v1alpha1.ModulePullOverrideList)
	if err := l.client.List(ctx, mpos); err != nil {
		return fmt.Errorf("list module pull overrides: %w", err)
	}

	for _, mpo := range mpos.Items {
		// ignore deleted mpo
		if !mpo.ObjectMeta.DeletionTimestamp.IsZero() {
			continue
		}

		// get relevant module source
		source := new(v1alpha1.ModuleSource)
		if err := l.client.Get(ctx, client.ObjectKey{Name: mpo.Spec.Source}, source); err != nil {
			return fmt.Errorf("get the '%s' module source for the '%s' module: %w", mpo.Spec.Source, mpo.Name, err)
		}

		// mpo's status.weight field isn't set - get it from the module's definition
		if mpo.Status.Weight == 0 {
			opts := utils.GenerateRegistryOptionsFromModuleSource(source, l.clusterUUID, l.log)
			md := downloader.NewModuleDownloader(l.dependencyContainer, l.downloadedModulesDir, source, opts)

			def, err := md.DownloadModuleDefinitionByVersion(mpo.Name, mpo.Spec.ImageTag)
			if err != nil {
				return fmt.Errorf("get the '%s' module definition from repository: %w", mpo.Name, err)
			}

			mpo.Status.Weight = def.Weight

			mpo.Status.UpdatedAt = metav1.NewTime(l.dependencyContainer.GetClock().Now().UTC())
			mpo.Status.Weight = def.Weight
			// we need not be bothered - even if the update fails, the weight will be set one way or another
			_ = l.client.Status().Update(ctx, &mpo)
		}

		// if deckhouseNodeName annotation isn't set or its value doesn't equal to current node name - overwrite the module from the repository
		if annotationNodeName, set := mpo.GetAnnotations()[v1alpha1.ModuleReleaseAnnotationDeckhouseNodeName]; !set || annotationNodeName != currentNodeName {
			l.log.Infof("reinitialize the '%s' module pull override due to stale/absent %s annotation", mpo.Name, v1alpha1.ModuleReleaseAnnotationDeckhouseNodeName)
			moduleDir := path.Join(l.downloadedModulesDir, mpo.Name, downloader.DefaultDevVersion)
			if err := os.RemoveAll(moduleDir); err != nil {
				return fmt.Errorf("delete the '%s' stale directory of the '%s' module: %w", moduleDir, mpo.Name, err)
			}

			if mpo.ObjectMeta.Annotations == nil {
				mpo.ObjectMeta.Annotations = make(map[string]string)
			}
			mpo.ObjectMeta.Annotations[v1alpha1.ModuleReleaseAnnotationDeckhouseNodeName] = currentNodeName

			if err := l.client.Update(ctx, &mpo); err != nil {
				l.log.Warnf("failed to annotate the '%s' module pull override: %v", mpo.Name, err)
			}
		}

		// if annotations is ok - we have to check that the file system is in sync
		moduleSymLink := filepath.Join(l.symlinksDir, fmt.Sprintf("%d-%s", mpo.Status.Weight, mpo.Name))
		if _, err := os.Stat(moduleSymLink); err != nil {
			// module symlink not found
			l.log.Infof("the '%s' module symlink is absent on file system, restore it", mpo.Name)
			if !os.IsNotExist(err) {
				return fmt.Errorf("check the '%s' module symlink: %w", mpo.Name, err)
			}
			if err = l.createModuleSymlink(mpo.Name, mpo.Spec.ImageTag, source, mpo.Status.Weight); err != nil {
				return fmt.Errorf("create the '%s' module symlink: %w", mpo.Name, err)
			}
		} else {
			dstDir, err := filepath.EvalSymlinks(moduleSymLink)
			if err != nil {
				return fmt.Errorf("evaluate the '%s' module symlink %s: %w", mpo.Name, moduleSymLink, err)
			}

			// check if module symlink leads to current version
			if filepath.Base(dstDir) != downloader.DefaultDevVersion {
				l.log.Infof("the '%s' module symlink is incorrect, restore it", mpo.Name)
				if err = l.createModuleSymlink(mpo.Name, mpo.Spec.ImageTag, source, mpo.Status.Weight); err != nil {
					return fmt.Errorf("create the '%s' module symlink: %w", mpo.Name, err)
				}
			}
		}

		// sync registry spec
		if err := utils.SyncModuleRegistrySpec(l.downloadedModulesDir, mpo.Name, downloader.DefaultDevVersion, source); err != nil {
			return fmt.Errorf("sync the '%s' module's registry settings with the '%s' module source: %w", mpo.Name, source.Name, err)
		}
		l.log.Infof("resynced the '%s' module's registry settings with the '%s' module source", mpo.Name, source.Name)
	}
	return nil
}

// restoreAbsentModulesFromReleases checks ModuleReleases with Deployed status and restore them on the FS
func (l *Loader) restoreAbsentModulesFromReleases(ctx context.Context) error {
	releases := new(v1alpha1.ModuleReleaseList)
	if err := l.client.List(ctx, releases); err != nil {
		return fmt.Errorf("list releases: %w", err)
	}

	// TODO: add labels to list only Deployed releases
	for _, release := range releases.Items {
		// ignore deleted release and not deployed
		if release.Status.Phase != v1alpha1.ModuleReleasePhaseDeployed || !release.ObjectMeta.DeletionTimestamp.IsZero() {
			continue
		}

		moduleWeight := release.Spec.Weight
		moduleVersion := "v" + release.Spec.Version.String()
		moduleName := release.Spec.ModuleName
		moduleSource := release.GetModuleSource()

		// if ModulePullOverride is set, don't check and restore overridden release
		exists, err := utils.ModulePullOverrideExists(ctx, l.client, moduleSource, moduleName)
		if err != nil {
			return fmt.Errorf("get module pull override for the '%s' module: %w", moduleName, err)
		}
		if exists {
			l.log.Infof("the '%s' module is overriden, skip release restoring", moduleName)
			continue
		}

		// get relevant module source
		source := new(v1alpha1.ModuleSource)
		if err = l.client.Get(ctx, client.ObjectKey{Name: moduleSource}, source); err != nil {
			return fmt.Errorf("get the '%s' module source for the '%s' module: %w", source.Name, moduleName, err)
		}

		moduleSymLink := filepath.Join(l.symlinksDir, fmt.Sprintf("%d-%s", release.Spec.Weight, release.Spec.ModuleName))
		if _, err = os.Stat(moduleSymLink); err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("check the '%s' module symlink: %w", moduleName, err)
			}
			l.log.Infof("the '%s' module symlink is absent on file system, restore it", moduleName)
			if err = l.createModuleSymlink(moduleName, moduleVersion, source, moduleWeight); err != nil {
				return fmt.Errorf("create module symlink: %w", err)
			}
		} else {
			dstDir, err := filepath.EvalSymlinks(moduleSymLink)
			if err != nil {
				return fmt.Errorf("evaluate the '%s' module symlink %s: %w", moduleName, moduleSymLink, err)
			}

			// check if module symlink leads to current version
			if filepath.Base(dstDir) != moduleVersion {
				l.log.Infof("the '%s' module symlink is incorrect, restore it", moduleName)
				if err = l.createModuleSymlink(moduleName, moduleVersion, source, moduleWeight); err != nil {
					return fmt.Errorf("create module symlink: %w", err)
				}
			}
		}

		// sync registry spec
		if err = utils.SyncModuleRegistrySpec(l.downloadedModulesDir, moduleName, moduleVersion, source); err != nil {
			return fmt.Errorf("sync the '%s' module's registry settings with the '%s' module source: %w", moduleName, source.Name, err)
		}
		l.log.Infof("resynced the '%s' module's registry settings with the '%s' module source", moduleName, source.Name)
	}
	return nil
}

// deleteModulesWithAbsentRelease deletes modules with absent releases
func (l *Loader) deleteModulesWithAbsentRelease(ctx context.Context) error {
	modulesLinks, err := l.readModulesFromFS()
	if err != nil {
		return fmt.Errorf("read source modules from the filesystem failed: %w", err)
	}

	releases := new(v1alpha1.ModuleReleaseList)
	if err = l.client.List(ctx, releases); err != nil {
		return fmt.Errorf("list releases: %w", err)
	}

	l.log.Debugf("found %d releases", len(releases.Items))

	// remove modules with release
	for _, release := range releases.Items {
		delete(modulesLinks, release.Spec.ModuleName)
	}

	for module, moduleLinkPath := range modulesLinks {
		mpo := new(v1alpha1.ModulePullOverride)
		if err = l.client.Get(ctx, client.ObjectKey{Name: module}, mpo); err != nil && apierrors.IsNotFound(err) {
			l.log.Warnf("the '%s' module has neither release nor override, purge it from fs", module)
			_ = os.RemoveAll(moduleLinkPath)
		}
	}

	return nil
}

func (l *Loader) readModulesFromFS() (map[string]string, error) {
	symlinksDir := filepath.Join(l.downloadedModulesDir, "modules")

	moduleLinks, err := os.ReadDir(symlinksDir)
	if err != nil {
		return nil, err
	}

	modules := make(map[string]string, len(moduleLinks))

	for _, moduleLink := range moduleLinks {
		index := strings.Index(moduleLink.Name(), "-")
		if index == -1 {
			continue
		}

		moduleName := moduleLink.Name()[index+1:]
		modules[moduleName] = path.Join(symlinksDir, moduleLink.Name())
	}

	return modules, nil
}

// createModuleSymlink checks if there are any other symlinks for a module in the symlink dir and deletes them before
// attempting to download current version of the module and creating correct symlink
func (l *Loader) createModuleSymlink(moduleName, moduleVersion string, moduleSource *v1alpha1.ModuleSource, moduleWeight uint32) error {
	l.log.Infof("the '%s' module is absent on filesystem, restore it from the '%s' source", moduleName, moduleSource.Name)

	// remove possible symlink doubles
	if err := wipeModuleSymlinks(l.symlinksDir, moduleName); err != nil {
		return err
	}

	// check if module's directory exists on fs
	info, err := os.Stat(path.Join(l.downloadedModulesDir, moduleName, moduleVersion))
	if err != nil || !info.IsDir() {
		l.log.Infof("downloading the '%s' module from the registry", moduleName)
		options := utils.GenerateRegistryOptionsFromModuleSource(moduleSource, l.clusterUUID, l.log)
		md := downloader.NewModuleDownloader(l.dependencyContainer, l.downloadedModulesDir, moduleSource, options)
		if _, err = md.DownloadByModuleVersion(moduleName, moduleVersion); err != nil {
			return fmt.Errorf("download the '%s' module with the %v version: %w", moduleName, moduleVersion, err)
		}
	}

	// restore symlink
	moduleRelativePath := filepath.Join("../", moduleName, moduleVersion)
	symlinkPath := filepath.Join(l.symlinksDir, fmt.Sprintf("%d-%s", moduleWeight, moduleName))
	if err = restoreModuleSymlink(l.downloadedModulesDir, symlinkPath, moduleRelativePath); err != nil {
		return fmt.Errorf("creating symlink for module %v failed: %w", moduleName, err)
	}
	l.log.Infof("the '%s/%s' module restored to %s", moduleName, moduleVersion, moduleRelativePath)

	return nil
}

func restoreModuleSymlink(downloadedModulesDir, symlinkPath, moduleRelativePath string) error {
	// make absolute path for versioned module
	moduleAbsPath := filepath.Join(downloadedModulesDir, strings.TrimPrefix(moduleRelativePath, "../"))
	// check that module exists on a disk
	if _, err := os.Stat(moduleAbsPath); os.IsNotExist(err) {
		return err
	}

	return os.Symlink(moduleRelativePath, symlinkPath)
}

// wipeModuleSymlinks checks if there are symlinks for the module with different weight in the symlink folder
func wipeModuleSymlinks(symlinksDir, moduleName string) error {
	// delete all module's symlinks in a loop
	for {
		anotherModuleSymlink, err := utils.FindExistingModuleSymlink(symlinksDir, moduleName)
		if err != nil {
			return fmt.Errorf("check if there are any other symlinks for the '%s' module: %w", moduleName, err)
		}

		if len(anotherModuleSymlink) > 0 {
			if err = os.Remove(anotherModuleSymlink); err != nil {
				return fmt.Errorf("delete stale symlink %v for the '%s' module: %w", anotherModuleSymlink, moduleName, err)
			}
			// go for another spin
			continue
		}

		// no more symlinks found
		break
	}
	return nil
}
