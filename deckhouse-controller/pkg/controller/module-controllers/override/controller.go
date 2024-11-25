// Copyright 2023 Flant JSC
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

package override

import (
	"context"
	"fmt"
	addonmodules "github.com/flant/addon-operator/pkg/module_manager/models/modules"
	"k8s.io/apimachinery/pkg/util/wait"
	"os"
	"path"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	addonutils "github.com/flant/addon-operator/pkg/utils"
	cp "github.com/otiai10/copy"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/deckhouse/deckhouse/deckhouse-controller/pkg/apis/deckhouse.io/v1alpha1"
	"github.com/deckhouse/deckhouse/deckhouse-controller/pkg/controller/module-controllers/downloader"
	"github.com/deckhouse/deckhouse/deckhouse-controller/pkg/controller/module-controllers/utils"
	"github.com/deckhouse/deckhouse/go_lib/d8env"
	"github.com/deckhouse/deckhouse/go_lib/dependency"
	"github.com/deckhouse/deckhouse/pkg/log"
)

const (
	controllerName = "d8-module-override-controller"

	maxConcurrentReconciles = 1
	cacheSyncTimeout        = 3 * time.Minute
)

func RegisterController(runtimeManager manager.Manager, mm moduleManager, dc dependency.Container, logger *log.Logger) error {
	r := &reconciler{
		init:                 new(sync.WaitGroup),
		client:               runtimeManager.GetClient(),
		log:                  logger,
		moduleManager:        mm,
		downloadedModulesDir: d8env.GetDownloadedModulesDir(),
		symlinksDir:          filepath.Join(d8env.GetDownloadedModulesDir(), "modules"),
		dependencyContainer:  dc,
	}

	r.init.Add(1)

	// add preflight Check
	if err := runtimeManager.Add(manager.RunnableFunc(r.preflight)); err != nil {
		return err
	}

	pullOverrideController, err := controller.New(controllerName, runtimeManager, controller.Options{
		MaxConcurrentReconciles: maxConcurrentReconciles,
		CacheSyncTimeout:        cacheSyncTimeout,
		NeedLeaderElection:      ptr.To(false),
		Reconciler:              r,
	})
	if err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(runtimeManager).
		For(&v1alpha1.ModulePullOverride{}).
		WithEventFilter(predicate.Or(predicate.GenerationChangedPredicate{}, predicate.AnnotationChangedPredicate{})).
		Complete(pullOverrideController)
}

type reconciler struct {
	init                 *sync.WaitGroup
	client               client.Client
	log                  *log.Logger
	dependencyContainer  dependency.Container
	moduleManager        moduleManager
	downloadedModulesDir string
	symlinksDir          string
	clusterUUID          string
}

type moduleManager interface {
	DisableModuleHooks(moduleName string)
	GetModule(moduleName string) *addonmodules.BasicModule
	RunModuleWithNewOpenAPISchema(moduleName, moduleSource, modulePath string) error
	AreModulesInited() bool
}

func (r *reconciler) preflight(ctx context.Context) error {
	defer r.init.Done()

	// wait until module manager init
	r.log.Debug("wait until module manager is inited")
	if err := wait.PollUntilContextCancel(ctx, 2*time.Second, true, func(_ context.Context) (bool, error) {
		return r.moduleManager.AreModulesInited(), nil
	}); err != nil {
		r.log.Errorf("failed to init module manager: %v", err)
		return err
	}

	r.clusterUUID = utils.GetClusterUUID(ctx, r.client)

	//err = c.restoreAbsentModulesFromOverrides(ctx)
	//if err != nil {
	//	return fmt.Errorf("modules restoration from overrides failed: %w", err)
	//}

	return nil
}

func (r *reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// wait for init
	r.init.Wait()

	mpo := new(v1alpha1.ModulePullOverride)
	if err := r.client.Get(ctx, client.ObjectKey{Name: req.Name}, mpo); err != nil {
		if apierrors.IsNotFound(err) {
			r.log.Warnf("the '%s' module pull override not found", req.Name)
			return ctrl.Result{}, nil
		}
		r.log.Errorf("failed to get the '%s'module pull override: %v", req.Name, err)
		return ctrl.Result{Requeue: true}, nil
	}

	return r.handleModuleOverride(ctx, mpo)
}

func (r *reconciler) handleModuleOverride(ctx context.Context, mo *v1alpha1.ModulePullOverride) (ctrl.Result, error) {
	var metaUpdateRequired bool

	// check if RegistrySpecChanged annotation is set and process it
	if _, set := mo.GetAnnotations()[v1alpha1.ModuleReleaseAnnotationRegistrySpecChanged]; set {
		// if module is enabled - push runModule task in the main queue
		r.log.Infof("apply new registry settings to the '%s' module", mo.Name)
		modulePath := filepath.Join(r.downloadedModulesDir, mo.Name, downloader.DefaultDevVersion)
		source := mo.ObjectMeta.Labels[v1alpha1.ModuleReleaseLabelSource]
		if err := r.moduleManager.RunModuleWithNewOpenAPISchema(mo.Name, source, modulePath); err != nil {
			r.log.Errorf("failed to run the '%s' module with new OpenAPI schema: %v", mo.Name, err)
			return ctrl.Result{Requeue: true}, nil
		}
		// delete annotation and requeue
		delete(mo.ObjectMeta.Annotations, v1alpha1.ModuleReleaseAnnotationRegistrySpecChanged)
		metaUpdateRequired = true
	}

	// add labels if empty, source and release controllers are looking for this labels
	if _, ok := mo.Labels[v1alpha1.ModuleReleaseLabelModule]; !ok {
		if len(mo.Labels) == 0 {
			mo.Labels = map[string]string{}
		}
		mo.Labels[v1alpha1.ModuleReleaseLabelModule] = mo.Name
		mo.Labels[v1alpha1.ModuleReleaseLabelSource] = mo.Spec.Source
		metaUpdateRequired = true
	}

	if metaUpdateRequired {
		return ctrl.Result{RequeueAfter: 500 * time.Millisecond}, r.client.Update(ctx, mo)
	}

	source := new(v1alpha1.ModuleSource)
	if err := r.client.Get(ctx, client.ObjectKey{Name: mo.Spec.Source}, source); err != nil {
		if apierrors.IsNotFound(err) {
			mo.Status.Message = fmt.Sprintf("the '%s' module source not found", mo.Spec.Source)
			if uerr := r.updateModulePullOverrideStatus(ctx, mo); uerr != nil {
				r.log.Errorf("failed to update module pull override status: %v", uerr)
				return ctrl.Result{Requeue: true}, nil
			}
			return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
		}
		r.log.Errorf("failed to get the '%s' module source: %v", mo.Spec.Source, err)
		return ctrl.Result{Requeue: true}, nil
	}

	tmpDir, err := os.MkdirTemp("", "module*")
	if err != nil {
		r.log.Errorf("failed to create temporary directory: %v", err)
		return ctrl.Result{}, nil
	}

	// clear temp dir
	defer func() {
		if err = os.RemoveAll(tmpDir); err != nil {
			r.log.Errorf("failed to remove the '%s' old module dir: %v", tmpDir, err)
		}
	}()

	options := utils.GenerateRegistryOptionsFromModuleSource(source, r.clusterUUID, r.log)
	md := downloader.NewModuleDownloader(r.dependencyContainer, tmpDir, source, options)
	newChecksum, moduleDef, err := md.DownloadDevImageTag(mo.Name, mo.Spec.ImageTag, mo.Status.ImageDigest)
	if err != nil {
		mo.Status.Message = err.Error()
		if uerr := r.updateModulePullOverrideStatus(ctx, mo); uerr != nil {
			r.log.Errorf("failed to update module pull override status: %v", uerr)
			return ctrl.Result{Requeue: true}, nil
		}
		r.log.Errorf("failed to download dev image tag: %v", err)
		return ctrl.Result{RequeueAfter: mo.Spec.ScanInterval.Duration}, nil
	}

	if newChecksum == "" {
		// module is up-to-date
		if mo.Status.Message != "" {
			// drop error message, if exists
			mo.Status.Message = ""
			if uerr := r.updateModulePullOverrideStatus(ctx, mo); uerr != nil {
				r.log.Errorf("failed to update module pull override status: %v", uerr)
				return ctrl.Result{Requeue: true}, nil
			}
		}
		return ctrl.Result{RequeueAfter: mo.Spec.ScanInterval.Duration}, nil
	}

	if moduleDef == nil {
		return ctrl.Result{RequeueAfter: mo.Spec.ScanInterval.Duration}, fmt.Errorf("got an empty module definition for the '%s' module pull override", mo.Name)
	}

	var values = make(addonutils.Values)
	if module := r.moduleManager.GetModule(moduleDef.Name); module != nil {
		values = module.GetConfigValues(false)
	}

	if err = utils.ValidateModule(*moduleDef, values, r.log); err != nil {
		mo.Status.Message = fmt.Sprintf("validation failed: %v", err)
		if uerr := r.updateModulePullOverrideStatus(ctx, mo); uerr != nil {
			r.log.Errorf("failed to update module pull override status: %v", uerr)
			return ctrl.Result{Requeue: true}, nil
		}
		r.log.Errorf("failed to validate the '%s' module pull override: %v", mo.Name, err)
		return ctrl.Result{}, nil
	}

	moduleStorePath := path.Join(r.downloadedModulesDir, moduleDef.Name, downloader.DefaultDevVersion)
	if err = os.RemoveAll(moduleStorePath); err != nil {
		return ctrl.Result{}, fmt.Errorf("remove the '%s' old module dir: %w", r.downloadedModulesDir, err)
	}

	if err = cp.Copy(tmpDir, r.downloadedModulesDir); err != nil {
		r.log.Errorf("failed to copy the module from the downloaded module dir: %v", err)
		return ctrl.Result{}, nil
	}

	symlinkPath := filepath.Join(r.symlinksDir, fmt.Sprintf("%d-%s", moduleDef.Weight, mo.Name))
	if err = r.enableModule(mo.Name, symlinkPath); err != nil {
		mo.Status.Message = err.Error()
		if uerr := r.updateModulePullOverrideStatus(ctx, mo); uerr != nil {
			r.log.Errorf("failed to update module pull override status: %v", uerr)
			return ctrl.Result{Requeue: true}, nil
		}
		r.log.Errorf("failed to enable the '%s' module: %v", mo.Name, err)
		return ctrl.Result{Requeue: true}, nil
	}

	// disable target module hooks so as not to invoke them before restart
	if r.moduleManager.GetModule(mo.Name) != nil {
		r.moduleManager.DisableModuleHooks(mo.Name)
	}

	defer func() {
		r.log.Infof("restart Deckhouse because %q ModulePullOverride image was updated", mo.Name)
		if err = syscall.Kill(1, syscall.SIGUSR2); err != nil {
			r.log.Fatalf("failed to send SIGUSR2 signal: %v", err)
		}
	}()

	mo.Status.Message = ""
	mo.Status.ImageDigest = newChecksum
	mo.Status.Weight = moduleDef.Weight

	if err = r.updateModulePullOverrideStatus(ctx, mo); err != nil {
		r.log.Errorf("failed to update the '%s' module pull override status: %v", mo.Name, err)
		return ctrl.Result{Requeue: true}, nil
	}

	if _, ok := mo.Annotations["renew"]; ok {
		delete(mo.Annotations, "renew")
		_ = r.client.Update(ctx, mo)
	}

	if err = r.ensureModuleDocumentation(ctx, mo); err != nil {
		r.log.Errorf("failed to ensure module documentation: %v", err)
		return ctrl.Result{Requeue: true}, nil
	}

	return ctrl.Result{RequeueAfter: mo.Spec.ScanInterval.Duration}, nil
}

// ensureModuleDocumentation creates or updates module documentation
func (r *reconciler) ensureModuleDocumentation(ctx context.Context, mo *v1alpha1.ModulePullOverride) error {
	modulePath := fmt.Sprintf("/%s/dev", mo.GetModuleName())
	moduleVersion := mo.Spec.ImageTag
	moduleChecksum := mo.Status.ImageDigest
	moduleName := mo.Name
	moduleSource := mo.Spec.Source
	ownerRef := metav1.OwnerReference{
		APIVersion: v1alpha1.ModulePullOverrideGVK.GroupVersion().String(),
		Kind:       v1alpha1.ModulePullOverrideGVK.Kind,
		Name:       mo.GetName(),
		UID:        mo.GetUID(),
		Controller: ptr.To(true),
	}

	md := new(v1alpha1.ModuleDocumentation)
	if err := r.client.Get(ctx, client.ObjectKey{Name: moduleName}, md); err != nil {
		if !apierrors.IsNotFound(err) {
			return fmt.Errorf("get the '%s' module documentation: %w", moduleName, err)
		}
		md = &v1alpha1.ModuleDocumentation{
			TypeMeta: metav1.TypeMeta{
				Kind:       v1alpha1.ModuleDocumentationGVK.Kind,
				APIVersion: v1alpha1.ModuleDocumentationGVK.GroupVersion().String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: moduleName,
				Labels: map[string]string{
					v1alpha1.ModuleReleaseLabelModule: moduleName,
					v1alpha1.ModuleReleaseLabelSource: moduleSource,
				},
				OwnerReferences: []metav1.OwnerReference{ownerRef},
			},
			Spec: v1alpha1.ModuleDocumentationSpec{
				Version:  moduleVersion,
				Path:     modulePath,
				Checksum: moduleChecksum,
			},
		}

		if err = r.client.Create(ctx, md); err != nil {
			return fmt.Errorf("create module documentation: %w", err)
		}
	}

	if md.Spec.Version != moduleVersion || md.Spec.Checksum != moduleChecksum {
		// update module documentation
		md.Spec.Path = modulePath
		md.Spec.Version = moduleVersion
		md.Spec.Checksum = moduleChecksum
		md.SetOwnerReferences([]metav1.OwnerReference{ownerRef})

		if err := r.client.Update(ctx, md); err != nil {
			return fmt.Errorf("update the '%s' module documentation: %w", md.Name, err)
		}
	}

	return nil
}

func (r *reconciler) enableModule(moduleName, symlinkPath string) error {
	currentModuleSymlink, err := utils.FindExistingModuleSymlink(r.symlinksDir, moduleName)
	if err != nil {
		currentModuleSymlink = "900-" + moduleName // fallback
	}

	modulePath := path.Join("../", moduleName, downloader.DefaultDevVersion)
	return utils.EnableModule(r.downloadedModulesDir, currentModuleSymlink, symlinkPath, modulePath)
}

func (r *reconciler) updateModulePullOverrideStatus(ctx context.Context, mo *v1alpha1.ModulePullOverride) error {
	mo.Status.UpdatedAt = metav1.NewTime(r.dependencyContainer.GetClock().Now().UTC())
	return r.client.Status().Update(ctx, mo)
}
