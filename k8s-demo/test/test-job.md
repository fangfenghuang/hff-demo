apiVersion: batch/v1
kind: Job
metadata:
  name: systemd-unit
spec:
  completions: 10  # Job 执行的次数
  parallelism: 2 # 并行执行的 Pod 个数
  backoffLimit: 4 # Job 重试的次数
  template:
    spec:
      containers:
      - name: systemd-unit
        image: systemd-unit：latest
        command: ["systemd-unit",  "--url=http://127.0.0.1/package/test.tar"]
      restartPolicy: Never


