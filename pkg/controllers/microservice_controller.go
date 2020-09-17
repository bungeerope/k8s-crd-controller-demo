/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"encoding/json"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"my.domain/example/resources"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	log "github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	devopsv1 "my.domain/example/pkg/api/v1"
)

// MicroServiceReconciler reconciles a MicroService object
type MicroServiceReconciler struct {
	client.Client
	Log    log.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=devops.my.domain,resources=microservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=devops.my.domain,resources=microservices/status,verbs=get;update;patch

func (r *MicroServiceReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {

	r.Log.Info("Request.Namespace", "Request.Name", req.Name)
	r.Log.Info("Reconciling AppService")

	// Fetch the AppService instance
	instance := &devopsv1.MicroService{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile req.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the req.
		return reconcile.Result{}, err
	}

	if instance.DeletionTimestamp != nil {
		return reconcile.Result{}, err
	}

	r.Log.Info("走到这里意味着 DemoMicroService resource 被找到，即该 resource 被成功创建，进入到了可根据该 resource 来执行逻辑的主流程")

	// 如果不存在，则创建关联资源
	// 如果存在，判断是否需要更新
	//   如果需要更新，则直接更新
	//   如果不需要更新，则正常返回

	deploy := &appsv1.Deployment{}
	if err := r.Client.Get(context.TODO(), req.NamespacedName, deploy); err != nil && errors.IsNotFound(err) {
		// 创建关联资源
		// 1. 创建 Deploy
		deploy := resources.NewDeployment(instance)
		if err := r.Client.Create(context.TODO(), deploy); err != nil {
			return reconcile.Result{}, err
		}
		// 2. 创建 Service
		service := resources.NewService(instance)
		if err := r.Client.Create(context.TODO(), service); err != nil {
			return reconcile.Result{}, err
		}
		// 3. 关联 Annotations
		data, _ := json.Marshal(instance.Spec)
		if instance.Annotations != nil {
			instance.Annotations["spec"] = string(data)
		} else {
			instance.Annotations = map[string]string{"spec": string(data)}
		}

		if err := r.Client.Update(context.TODO(), instance); err != nil {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, nil
	}

	oldspec := devopsv1.MicroServiceSpec{}
	if err := json.Unmarshal([]byte(instance.Annotations["spec"]), &oldspec); err != nil {
		return reconcile.Result{}, err
	}

	if !reflect.DeepEqual(instance.Spec, oldspec) {
		// 更新关联资源
		newDeploy := resources.NewDeployment(instance)
		oldDeploy := &appsv1.Deployment{}
		if err := r.Client.Get(context.TODO(), req.NamespacedName, oldDeploy); err != nil {
			return reconcile.Result{}, err
		}
		oldDeploy.Spec = newDeploy.Spec
		if err := r.Client.Update(context.TODO(), oldDeploy); err != nil {
			return reconcile.Result{}, err
		}

		newService := resources.NewService(instance)
		oldService := &corev1.Service{}
		if err := r.Client.Get(context.TODO(), req.NamespacedName, oldService); err != nil {
			return reconcile.Result{}, err
		}
		oldService.Spec = newService.Spec
		if err := r.Client.Update(context.TODO(), oldService); err != nil {
			return reconcile.Result{}, err
		}

		return reconcile.Result{}, nil
	}

	return reconcile.Result{}, nil

}

func (r *MicroServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&devopsv1.MicroService{}).
		Complete(r)
}
