[TOC]
# kubevirt常用命令


## 虚拟机资源
```
kubectl get vm -n kubevirt
kubectl get vmi -n kubevirt
kubectl get dv -n kubevirt
kubectl get pvc -n kubevirt
kubectl get pod -n kubevirt
```

## 创建虚拟机
```bash
Creating a virtual machine
$ kubectl apply -f  vm.yaml
[root@node1 kubevirt]# kubectl get vm --all-namespaces
NAMESPACE   NAME     AGE   VOLUME
default     testvm   55s

After deployment you can manage VMs using the usual verbs:
$ kubectl describe vm testvm
To start a VM you can use, this will create a VM instance (VMI)
$ ./virtctl start testvm
The interested reader can now optionally inspect the instance
$ kubectl describe vmi testvm


To shut the VM down again:
$ ./virtctl stop testvm

kubectl can be used too:

Start the virtual machine:
kubectl patch virtualmachine myvm --type merge -p \
    '{"spec":{"running":true}}'

Stop the virtual machine:
kubectl patch virtualmachine myvm --type merge -p \
    '{"spec":{"running":false}}'
kubectl patch vm testvm --type merge -p     '{"spec":{"running":false}}'


To delete
$ kubectl delete vm testvm# To create your own
$ kubectl apply -f $YOUR_VM_SPEC
```

## 进入虚拟机：
```bash
virtctl console testvm
Successfully connected to testvm console. The escape sequence is ^]

login as 'cirros' user. default password: 'gocubsgo'. use 'sudo' for root.
testvm login:cirros
Password:gocubsgo
login as 'cirros' user. default password: 'gocubsgo'. use 'sudo' for root.
testvm login:
 d 
```
## 重启虚拟机：

```bash
Restart the virtual machine (you delete the instance!):
kubectl delete virtualmachineinstance myvm
To restart a VirtualMachine named myvm using virtctl:
$ virtctl restart myvm
virtctl restart myvm --force --grace-period=0
Note: Force restart can cause data corruption, and should be used in cases of kernal panic or VirtualMachine being unresponsive to normal restarts.
```

## 删除、关闭虚拟机：
```
Stop the virtual machine instance:
    kubectl patch virtualmachine myvm --type merge -p \
        '{"spec":{"running":false}}'

    # Restart the virtual machine (you delete the instance!):
    kubectl delete virtualmachineinstance myvm

    # Implicit cascade delete (first deletes the virtual machine and then the virtual machine)
    kubectl delete virtualmachine myvm

    # Explicit cascade delete (first deletes the virtual machine and then the virtual machine)
    kubectl delete virtualmachine myvm --cascade=true

    # Orphan delete (The running virtual machine is only detached, not deleted)
    # Recreating the virtual machine would lead to the adoption of the virtual machine instance
    kubectl delete virtualmachine myvm --cascade=false
```
