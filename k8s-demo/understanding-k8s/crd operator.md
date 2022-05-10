# CRD 
自定义资源的描述(Custom Resource Definition)

# Operator
Operator的工作原理，实际上是利用了Kubernetes的自定义API资源（CRD），来描述我们想要部署的“有状态应用”；然后在自定义控制器里，根据自定义API对象的变化，来完成具体的部署和运维工作。