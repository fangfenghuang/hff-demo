package cmd

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-demo/pkg/k8sclient"

	"go-demo/conf"

	"github.com/google/uuid"
	"k8s.io/klog"

	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var LeaseLockID string

var LeaderelectionCmd = &cobra.Command{
	Use:   "go-demo",
	Short: "goland test",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		LeaseLockID = uuid.New().String()
		initializaton()

		run := func(ctx context.Context) {
			// 添加运行逻辑代码
			klog.InitFlags(nil)
			flag.Parse()
			defer klog.Flush()
			klog.Infoln("service start...")
			main() //****************************************//
			//router.InitRouter()   //http server
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-ch
			klog.Warning("Received termination, signaling shutdown")
			cancel()
		}()

		// 指定锁的资源对象，这里使用了Lease资源，还支持configmap，endpoint，或者multilock(即多种配合使用)
		lock := &resourcelock.LeaseLock{
			LeaseMeta: metav1.ObjectMeta{
				Name:      conf.LeaseLockName,
				Namespace: conf.LeaseLockNamespace,
			},
			Client: k8sclient.K8sClientSet.CoordinationV1(),
			LockConfig: resourcelock.ResourceLockConfig{
				Identity: LeaseLockID,
			},
		}

		leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
			Lock:            lock,
			ReleaseOnCancel: true,
			LeaseDuration:   30 * time.Second, //租约时间
			RenewDeadline:   15 * time.Second, //更新租约的
			RetryPeriod:     5 * time.Second,  //非leader节点重试时间
			Callbacks: leaderelection.LeaderCallbacks{
				OnStartedLeading: func(ctx context.Context) {
					//变为leader执行的业务代码
					run(ctx)
				},
				OnStoppedLeading: func() {
					// 进程退出
					klog.Infof("leader lost: %s", LeaseLockID)
					os.Exit(0)
				},
				OnNewLeader: func(identity string) {
					//当产生新的leader后执行的方法
					if identity == LeaseLockID {
						klog.Infof("i am leader now: %s", identity)
						return
					}
					klog.Infof("new leader elected: %s, wait...", identity)
				},
			},
		})

	},
}

var RootCmd = &cobra.Command{
	Use:   "go-demo",
	Short: "golang test",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		klog.InitFlags(nil)
		flag.Parse()
		defer klog.Flush()

		klog.Infoln("service start...")
		initializaton()
		main() //****************************************//
		//router.InitRouter() //http server
	},
}

// add init func
func initializaton() {
	k8sclient.InitClientSet()
	//calico.InitGlobalNetworkPolicy()   //初始化网络策略
}

func Execute() {
	//if err := LeaderelectionCmd.Execute(); err != nil {
	if err := RootCmd.Execute(); err != nil {
		klog.Errorln(os.Stderr, err)
		os.Exit(1)
	}
}

// add test func
func main() {
	klog.Infoln("test...")
	k8sclient.Test()
}
