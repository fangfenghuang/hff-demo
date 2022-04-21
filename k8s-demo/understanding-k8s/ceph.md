# 挂载点查询
```bash
root@kubectl:/# kubectl get storageclass csi-iscsi-sc -o yaml
allowVolumeExpansion: true
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  annotations:
    storageclass.kubernetes.io/is-default-class: "true"
  creationTimestamp: "2020-12-17T02:08:10Z"
  name: csi-iscsi-sc
  resourceVersion: "1"
  selfLink: /apis/storage.k8s.io/v1/storageclasses/csi-iscsi-sc
  uid: 4ec4a55c-e2c5-434e-8413-8dc1eefba14a
mountOptions:
- _netdev
parameters:
  accessPaths: iscsiap
  csi.storage.k8s.io/provisioner-secret-name: csi-iscsi-secret
  csi.storage.k8s.io/provisioner-secret-namespace: xsky
  fstype: xfs
  pool: data-pool
  xmsServers: 10.10.35.111
provisioner: iscsi.csi.xsky.com
reclaimPolicy: Delete
volumeBindingMode: Immediate


---
root@kubectl:/# kubectl get pvc -n spnspp-prod www-web-0 -o yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  annotations:
    pv.kubernetes.io/bind-completed: "yes"
    pv.kubernetes.io/bound-by-controller: "yes"
    volume.beta.kubernetes.io/storage-provisioner: iscsi.csi.xsky.com
  creationTimestamp: "2021-06-23T07:12:07Z"
  finalizers:
  - kubernetes.io/pvc-protection
  labels:
    app: nginx
  name: www-web-0
  namespace: spnspp-prod
  resourceVersion: "174741770"
  selfLink: /api/v1/namespaces/spnspp-prod/persistentvolumeclaims/www-web-0
  uid: d5a09c7a-0dc7-4f9d-89dd-2d4cbc52783c
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
  storageClassName: csi-iscsi-sc
  volumeMode: Filesystem
  volumeName: pvc-d5a09c7a-0dc7-4f9d-89dd-2d4cbc52783c
status:
  accessModes:
  - ReadWriteOnce
  capacity:
    storage: 1Gi
  phase: Bound
---
root@kubectl:/# kubectl get pv pvc-d5a09c7a-0dc7-4f9d-89dd-2d4cbc52783c -o yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  annotations:
    pv.kubernetes.io/provisioned-by: iscsi.csi.xsky.com
  creationTimestamp: "2021-06-23T07:12:10Z"
  finalizers:
  - kubernetes.io/pv-protection
  - external-attacher/iscsi-csi-xsky-com
  name: pvc-d5a09c7a-0dc7-4f9d-89dd-2d4cbc52783c
  resourceVersion: "174741786"
  selfLink: /api/v1/persistentvolumes/pvc-d5a09c7a-0dc7-4f9d-89dd-2d4cbc52783c
  uid: 97d67c4a-03f8-4e64-8806-0b110320aa57
spec:
  accessModes:
  - ReadWriteOnce
  capacity:
    storage: 1Gi
  claimRef:
    apiVersion: v1
    kind: PersistentVolumeClaim
    name: www-web-0
    namespace: spnspp-prod
    resourceVersion: "174741728"
    uid: d5a09c7a-0dc7-4f9d-89dd-2d4cbc52783c
  csi:
    driver: iscsi.csi.xsky.com
    fsType: xfs
    volumeAttributes:
      accessPaths: iscsiap
      fstype: xfs
      pool: data-pool
      storage.kubernetes.io/csiProvisionerIdentity: 1624431607122-8081-iscsi.csi.xsky.com
      xmsServers: 10.10.35.4,10.10.35.5,10.10.35.6
    volumeHandle: csi-iscsi-pvc-d5a09c7a-0dc7-4f9d-89dd-2d4cbc52783c
  mountOptions:
  - _netdev
  persistentVolumeReclaimPolicy: Delete
  storageClassName: csi-iscsi-sc
  volumeMode: Filesystem
status:
  phase: Bound






```


# rbd vs cephfs
rbd不支持ReadWriteMany
在同一个节点中,是可以支持多个pod共享rbd image运行的

# xsky


# 