/*
Copyright 2022 fangfenghuang.

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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	hffappv1 "github.com/fangfenghuang/kubebuilder-demo/api/v1"
)

// HffDemoReconciler reconciles a HffDemo object
type HffDemoReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=hffapp.fangfenghuang.io,resources=hffdemoes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=hffapp.fangfenghuang.io,resources=hffdemoes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=hffapp.fangfenghuang.io,resources=hffdemoes/finalizers,verbs=update
//+kubebuilder:rbac:groups=hffapp.fangfenghuang.io,resources=pods,verbs=get;list;watch;create;update;patch;delete

//上面增加pods操作权限
// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the HffDemo object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *HffDemoReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	// TloggerODO(user): your logic here  增删查改逻辑
	// controllerMananger会监控 当时的集群状态与yaml 声明期望状态的差值，转而去执行Reconcile 里的 增删查改的操作。

	logger.Info("revice reconcile event", req.Namespace, req.Name)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
//默认情况下生成的controller只监听自定义资源
func (r *HffDemoReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hffappv1.HffDemo{}).
		Complete(r)
	// return ctrl.NewControllerManagedBy(mgr).
	// For(&hffappv1.HffDemo{}).Watches(&source.Kind{Type: &hffappv1.Pod{}}, &handler.EnqueueRequestForObject{}).
	// Complete(r)
}
