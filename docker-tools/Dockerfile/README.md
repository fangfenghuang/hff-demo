# fio-iperf
```yaml
---
apiVersion: v1
kind: Pod
metadata:
  name: test-runc
spec:
  containers:
  - name: pod1
    image: hff/hff-fio-iperf:v0.1
    imagePullPolicy: IfNotPresent
    tty: true
    command: ["/bin/sh","-c","netserver -p 4444 -4; iperf3 -s -i 1;"] 
    ports:
    - name: netperf-port
      containerPort: 4444
    - name: iperf-port
      containerPort: 5210
    volumeMounts:
    - mountPath: /test
      name: test-volume
    resources:
      limits:
        memory: "1Gi"
        cpu: "1"
  volumes:
  - name: test-volume
    hostPath:
      path: /hff/test/runc
      type: Directory
```

