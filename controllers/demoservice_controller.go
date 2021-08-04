/*
Copyright 2021 zyx.

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

	"github.com/zyxpaomian/opdemo/resources"

	"github.com/go-logr/logr"
	appv1beta1 "github.com/zyxpaomian/opdemo/api/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DemoServiceReconciler reconciles a DemoService object
type DemoServiceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=app.zyx.io,resources=demoservices,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=app.zyx.io,resources=demoservices/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=app.zyx.io,resources=demoservices/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DemoService object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *DemoServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("demoservice", req.NamespacedName)

	// 获取crd的对象
	var demoService appv1beta1.DemoService
	err := r.Get(ctx, req.NamespacedName, &demoService)
	if err != nil {
		// crd对象 被删除的时候，忽略
		if client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	log.Info("发现crd对象", "demoservice", demoService)
	// Reconcile successful - don't requeue
	// return ctrl.Result{}, nil
	// Reconcile failed due to error - requeue
	// return ctrl.Result{}, err
	// Requeue for any reason other than an error
	// return ctrl.Result{Requeue: true}, nil

	// your logic here
	deploy := &appsv1.Deployment{}
	if err := r.Get(ctx, req.NamespacedName, deploy); err != nil && errors.IsNotFound(err) {
		log.Info("该service 下没有deployment")
		deploy := resources.NewDeploy(&demoService)
		if err := r.Client.Create(ctx, deploy); err != nil {
			return ctrl.Result{}, err
		}
		log.Info("创建deployment完成")
	}

	return ctrl.Result{}, nil

	// return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DemoServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1beta1.DemoService{}).
		Complete(r)
}
