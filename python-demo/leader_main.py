# encoding: utf-8

# 参考：python\kubernetes\base\leaderelection\example.py
import sys
import uuid
from gevent import pywsgi
from kubernetes import client, config
from kubernetes.leaderelection import leaderelection
from kubernetes.leaderelection.resourcelock.configmaplock import ConfigMapLock
from kubernetes.leaderelection import electionconfig
from conf import *
from app import app


config.load_kube_config(config_file=K8S_CONFIG)

candidate_id = uuid.uuid4()

lock_name = "python-demo"

lock_namespace = "default"


def on_started():
    print("I am leader now！！！！")
    # server = pywsgi.WSGIServer(('0.0.0.0', SERVICE_PORT), app)
    # server.serve_forever()

def on_stopped():
    print("stop leading now！！！")
    sys.exit(1)

if __name__ == '__main__':
    config = electionconfig.Config(ConfigMapLock(lock_name, lock_namespace, candidate_id), lease_duration=30,
                               renew_deadline=15, retry_period=5, onstarted_leading=on_started,
                               onstopped_leading=on_stopped, release_on_cancel=True)

    leaderelection.LeaderElection(config).run()
    print("Exited leader election")


"""
删除包： pip3 uninstall kubernetes
安装最新包： pip3 uninstall kubernetes

修改前：
[root@tztest python-demo]# python3 leader_main.py
INFO:root:d4ae0d02-3e30-490a-b154-990a54f6378e is a follower
INFO:root:leader d4ae0d02-3e30-490a-b154-990a54f6378e has successfully acquired lease
INFO:root:d4ae0d02-3e30-490a-b154-990a54f6378e successfully acquired lease
I am leader now！！！！
INFO:root:Leader has entered renew loop and will try to update lease continuously
INFO:root:leader d4ae0d02-3e30-490a-b154-990a54f6378e has successfully acquired lease
INFO:root:leader d4ae0d02-3e30-490a-b154-990a54f6378e has successfully acquired lease
INFO:root:leader d4ae0d02-3e30-490a-b154-990a54f6378e has successfully acquired lease
^CTraceback (most recent call last):
  File "leader_main.py", line 38, in <module>
    leaderelection.LeaderElection(config).run()
  File "/usr/local/python3/lib/python3.8/site-packages/kubernetes/leaderelection/leaderelection.py", line 64, in run
    self.renew_loop()
  File "/usr/local/python3/lib/python3.8/site-packages/kubernetes/leaderelection/leaderelection.py", line 101, in renew_loop
    time.sleep(retry_period)
KeyboardInterrupt

[root@tztest python-demo]# python3 leader_main.py
INFO:root:7027ffcc-e5ae-468b-ac9d-834d4276dc2e is a follower
INFO:root:yet to finish lease_duration, lease held by d4ae0d02-3e30-490a-b154-990a54f6378e and has not expired
INFO:root:yet to finish lease_duration, lease held by d4ae0d02-3e30-490a-b154-990a54f6378e and has not expired
INFO:root:yet to finish lease_duration, lease held by d4ae0d02-3e30-490a-b154-990a54f6378e and has not expired




修改后：
[root@tztest python-demo]# python3 leader_main.py
INFO:root:44d32a22-e320-4f0d-a95c-42adb3ba443f is a follower
INFO:root:yet to finish lease_duration, lease held by 5f33765a-38ff-44ef-9d44-14d05584d003 and has not expired
INFO:root:yet to finish lease_duration, lease held by 5f33765a-38ff-44ef-9d44-14d05584d003 and has not expired
INFO:root:yet to finish lease_duration, lease held by 5f33765a-38ff-44ef-9d44-14d05584d003 and has not expired
INFO:root:yet to finish lease_duration, lease held by 5f33765a-38ff-44ef-9d44-14d05584d003 and has not expired
INFO:root:yet to finish lease_duration, lease held by 5f33765a-38ff-44ef-9d44-14d05584d003 and has not expired
INFO:root:yet to finish lease_duration, lease held by 5f33765a-38ff-44ef-9d44-14d05584d003 and has not expired
INFO:root:leader 44d32a22-e320-4f0d-a95c-42adb3ba443f has successfully acquired lease
INFO:root:44d32a22-e320-4f0d-a95c-42adb3ba443f successfully acquired lease
I am leader now！！！！
INFO:root:Leader has entered renew loop and will try to update lease continuously
INFO:root:leader 44d32a22-e320-4f0d-a95c-42adb3ba443f has successfully acquired lease
INFO:root:leader 44d32a22-e320-4f0d-a95c-42adb3ba443f has successfully acquired lease
INFO:root:leader 44d32a22-e320-4f0d-a95c-42adb3ba443f has successfully acquired lease
INFO:root:leader 44d32a22-e320-4f0d-a95c-42adb3ba443f has successfully acquired lease
INFO:root:leader 44d32a22-e320-4f0d-a95c-42adb3ba443f has successfully acquired lease
^Cstop leading now！！！


[root@tztest kbuser]# kubectl get configmaps python-demo -o yaml
apiVersion: v1
kind: ConfigMap
metadata:
  annotations:
    control-plane.alpha.kubernetes.io/leader: '{"holderIdentity": "22f33846-5d5b-45f4-8d4b-b2213765e2f0",
      "leaseDurationSeconds": "30", "acquireTime": "2022-03-15 09:52:47.173846", "renewTime":
      "2022-03-15 09:53:07.256312"}'
  creationTimestamp: "2022-01-05T07:30:34Z"
  name: python-demo
  namespace: default
  resourceVersion: "58651326"
  selfLink: /api/v1/namespaces/default/configmaps/python-demo
  uid: ee023f69-4d91-475e-80b6-e4eac95e84bb




"""