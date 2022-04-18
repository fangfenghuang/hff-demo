
https://github.com/openebs/lvm-localpv

创建pv

创建vg

$ kubectl apply -f https://openebs.github.io/charts/lvm-operator.yaml

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

