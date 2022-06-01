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
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

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
	logger := r.Client.Log.WithValues("hffdemo", req.NamespacedName)
	// TODO(user): your logic here  增删查改逻辑
	// controllerMananger会监控 当时的集群状态与yaml 声明期望状态的差值，转而去执行Reconcile 里的 增删查改的操作。

	logger.Info("revice reconcile event", "name", req.Name)
	// 获取demo对象
	hffdemo := &hffappv1.HffDemo{}
	if err := r.Client.Get(ctx, req.NamespacedName, hffdemo); err != nil {
		if errors.IsNotFound(err) {
			log.Info("resource not found, skipping reconcile")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Unable to get resource")
		return ctrl.Result{}, err
	}
	// 如果处在删除中直接跳过
	if hffdemo.DeletionTimestamp != nil {
		logger.Info("hffdemo in deleting", "name", req.String())
		return ctrl.Result{}, nil
	}

	// 查找upload pod
	// 查找deployment
	uploadPod := &hffappv1.Pod{}

	// 用客户端工具查询
	err = r.Get(ctx, req.NamespacedName, uploadPod)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
//默认情况下生成的controller只监听自定义资源
func (r *HffDemoReconciler) SetupWithManager(mgr ctrl.Manager) error {
	ctx := context.Background()
	// return ctrl.NewControllerManagedBy(mgr).
	// 	For(&hffappv1.HffDemo{}).
	// 	Complete(r)

	// 创建controller
	c, err := controller.New("hffdemo-controller", mgr, controller.Options{
		Reconciler:              r,
		MaxConcurrentReconciles: 1, //controller运行的worker数
	})
	if err != nil {
		return err
	}
	//监听自定义资源
	var demoObj hffappv1.HffDemo
	if err := c.Watch(&source.Kind{Type: &demoObj}, &handler.EnqueueRequestForObject{}); err != nil {
		return err
	}
	//监听upload pod, 将owner namespace/name添加到队列
	if err := c.Watch(&source.Kind{Type: &v1.Pods{}}, &handler.EnqueueRequestForOwner{
		OwnerType:    demoObj,
		IsController: true,
	}); err != nil {
		return err
	}

	//3.调用FreshStatus方法刷新资源状态
	if err := r.Client.freshStatus(ctx, &demoObj); err != nil {
		log.Error(err, "Fail to refresh status subresource")
		return ctrl.Result{}, err
	}
	return nil
}

var (
	newPodNode      = ""
	containerName   = "upload"
	containerImange = "busybox"
	destVolume      = "tmp-upload"
	destPath        = "/tmp-upload"
)

// 新建uplaod pod
func createUploadPod(ctx context.Context, r *HffDemoReconciler, hffdemo *hffappv1.HffDemo) error {
	logger := r.Client.Log.WithValues("func", "createUploadPod")

	var newPodName = fmt.Sprintf("tmp-upload-pvc-%v", time.Now().Unix())

	log.Info(fmt.Sprintf("newPodName [%s]", newPodName))

	//创建upload pod
	newPod := &hffappv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      newPodName,
			Namespace: hffdemo.Namespace,
		},
		Spec: hffdemo.PodSpec{
			Containers: []hffdemo.Container{
				{
					Name:            containerName,
					Image:           containerImange,
					ImagePullPolicy: corev1.PullIfNotPresent,
					TTY:             true,
					VolumeMounts: []hffdemo.VolumeMount{
						{
							Name:      destVolume,
							MountPath: hffdemo.Spec.DestPath,
						},
					},
				},
			},
			Volumes: []hffdemo.Volume{
				{
					Name: destVolume,
					VolumeSource: hffdemo.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: hffdemo.Spec.PVCName,
						},
					},
				},
			},
		},
	}

	// 这一步非常关键！
	// 建立关联后，删除crd资源时就会将pod也删除掉
	logger.Info("set reference")
	if err := controllerutil.SetControllerReference(hffdemo, newPod, r.Scheme); err != nil {
		logger.Error(err, "SetControllerReference error")
		return err
	}

	// 创建pod
	logger.Info("start create upload pod")
	if err := r.Client.Create(ctx, newPod); err != nil {
		logger.Error(err, "create pod error")
		return err
	}

	logger.Info("create pod success")

	return nil
}

func updateStatus(ctx context.Context, r *HffDemoReconciler, hffdemo *hffappv1.HffDemo) error {
	logger := r.Client.Log.WithValues("func", "updateStatus")

	*(hffdemo.Status.Phase) = hffappv1.Running

	if err := r.Client.Update(ctx, hffdemo); err != nil {
		logger.Error(err, "update instance error")
		return err
	}

	return nil
}
