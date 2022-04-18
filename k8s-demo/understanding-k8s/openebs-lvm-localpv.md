
https://github.com/openebs/lvm-localpv

# pv 
## 创建pv
pvname=${disk}1

sgdisk -n 1:2048 ${disk}
pvcreate ${pvname}
pvdisplay ${pvname}


## 创建vg
vgcreate ${vgname} ${pvname}
vgextend ${vgname} ${pvname}
pvdisplay ${pvname} | grep ${vgname}


# deploy
$ kubectl apply -f https://openebs.github.io/charts/lvm-operator.yaml


# 应用

```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: openebs-lvmpv
allowVolumeExpansion: true
parameters:
  storage: "lvm"
  volgroup: "lvmvg"
provisioner: local.csi.openebs.io
allowedTopologies:
- matchLabelExpressions:
  - key: kubernetes.io/hostname
    values:
      - lvmpv-node1
      - lvmpv-node2
```
volgroup：选择vg
allowedTopologies: 选择节点

# disk回收
pvname=${disk}1

vgremove ${vgname} 或 vgreduce ${vgname} ${pvname}

pvremove ${pvname}

pvdisplay ${pvname}

sgdisk --zap-all --clear --mbrtogpt ${disk}
