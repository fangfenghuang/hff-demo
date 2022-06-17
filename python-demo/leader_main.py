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
                               onstopped_leading=on_stopped, release_on_stop=True)

    leaderelection.LeaderElection(config).run()
    print("Exited leader election")


"""
删除包： pip3 uninstall kubernetes
安装最新包： pip3 install kubernetes

修改前：
[root@tztest python-demo]# python3 leader_main.py
INFO:root:cca40786-84e8-4629-9372-cb62fd003fb8 is a follower
INFO:root:leader cca40786-84e8-4629-9372-cb62fd003fb8 has successfully acquired lease
INFO:root:cca40786-84e8-4629-9372-cb62fd003fb8 successfully acquired lease
I am leader now！！！！
INFO:root:Leader has entered renew loop and will try to update lease continuously
INFO:root:leader cca40786-84e8-4629-9372-cb62fd003fb8 has successfully acquired lease
INFO:root:leader cca40786-84e8-4629-9372-cb62fd003fb8 has successfully acquired lease
INFO:root:leader cca40786-84e8-4629-9372-cb62fd003fb8 has successfully acquired lease
INFO:root:leader cca40786-84e8-4629-9372-cb62fd003fb8 has successfully acquired lease
INFO:root:leader cca40786-84e8-4629-9372-cb62fd003fb8 has successfully acquired lease
INFO:root:leader cca40786-84e8-4629-9372-cb62fd003fb8 has successfully acquired lease
INFO:root:leader cca40786-84e8-4629-9372-cb62fd003fb8 has successfully acquired lease
INFO:root:leader cca40786-84e8-4629-9372-cb62fd003fb8 has successfully acquired lease
^CTraceback (most recent call last):
  File "leader_main.py", line 38, in <module>
    leaderelection.LeaderElection(config).run()
  File "/usr/local/python3/lib/python3.8/site-packages/kubernetes/leaderelection/leaderelection.py", line 64, in run
    self.renew_loop()
  File "/usr/local/python3/lib/python3.8/site-packages/kubernetes/leaderelection/leaderelection.py", line 101, in renew_loop
    time.sleep(retry_period)
KeyboardInterrupt
[root@tztest python-demo]#  python3 leader_main.py
INFO:root:9f2f1625-5a8f-491c-ab65-86d75a5eedea is a follower
INFO:root:yet to finish lease_duration, lease held by cca40786-84e8-4629-9372-cb62fd003fb8 and has not expired
INFO:root:yet to finish lease_duration, lease held by cca40786-84e8-4629-9372-cb62fd003fb8 and has not expired
INFO:root:yet to finish lease_duration, lease held by cca40786-84e8-4629-9372-cb62fd003fb8 and has not expired
^CTraceback (most recent call last):
  File "leader_main.py", line 38, in <module>
    leaderelection.LeaderElection(config).run()
  File "/usr/local/python3/lib/python3.8/site-packages/kubernetes/leaderelection/leaderelection.py", line 57, in run
    if self.acquire():
  File "/usr/local/python3/lib/python3.8/site-packages/kubernetes/leaderelection/leaderelection.py", line 80, in acquire
    time.sleep(retry_period)
KeyboardInterrupt

[root@tztest python-demo]# kubectl get configmaps python-demo -o yaml
apiVersion: v1
kind: ConfigMap
metadata:
  annotations:
    control-plane.alpha.kubernetes.io/leader: '{"holderIdentity": "cca40786-84e8-4629-9372-cb62fd003fb8",
      "leaseDurationSeconds": "30", "acquireTime": "2022-06-16 18:44:16.205808", "renewTime":
      "2022-06-16 18:44:51.338605"}'
  creationTimestamp: "2022-01-05T07:30:34Z"
  name: python-demo
  namespace: default
  resourceVersion: "85889315"
  selfLink: /api/v1/namespaces/default/configmaps/python-demo
  uid: ee023f69-4d91-475e-80b6-e4eac95e84bb




替换：
/usr/local/python3/lib/python3.8/site-packages/kubernetes/leaderelection


修改后：
[root@tztest python-demo]# python3 leader_main.py
INFO:root:5839e63b-77f4-4d28-92ab-bee011dab00a is a follower
INFO:root:leader 5839e63b-77f4-4d28-92ab-bee011dab00a has successfully acquired lease
INFO:root:5839e63b-77f4-4d28-92ab-bee011dab00a successfully acquired lease
I am leader now！！！！
INFO:root:Leader has entered renew loop and will try to update lease continuously
INFO:root:leader 5839e63b-77f4-4d28-92ab-bee011dab00a has successfully acquired lease
INFO:root:leader 5839e63b-77f4-4d28-92ab-bee011dab00a has successfully acquired lease
INFO:root:leader 5839e63b-77f4-4d28-92ab-bee011dab00a has successfully acquired lease
INFO:root:leader 5839e63b-77f4-4d28-92ab-bee011dab00a has successfully acquired lease
^CINFO:root:release lock: 5839e63b-77f4-4d28-92ab-bee011dab00a
stop leading now！！！
[root@tztest python-demo]# kubectl get configmaps python-demo -o yaml
apiVersion: v1
kind: ConfigMap
metadata:
  annotations:
    control-plane.alpha.kubernetes.io/leader: '{"holderIdentity": null, "leaseDurationSeconds":
      1, "acquireTime": "2022-06-16 18:54:35.323096", "renewTime": "2022-06-16 18:54:35.323096"}'
  creationTimestamp: "2022-01-05T07:30:34Z"
  name: python-demo
  namespace: default
  resourceVersion: "85891277"
  selfLink: /api/v1/namespaces/default/configmaps/python-demo
  uid: ee023f69-4d91-475e-80b6-e4eac95e84bb
[root@tztest python-demo]#
[root@tztest python-demo]#
[root@tztest python-demo]# kubectl get configmaps python-demo -o yaml
apiVersion: v1
kind: ConfigMap
metadata:
  annotations:
    control-plane.alpha.kubernetes.io/leader: '{"holderIdentity": "5839e63b-77f4-4d28-92ab-bee011dab00a",
      "leaseDurationSeconds": "30", "acquireTime": "2022-06-16 18:55:22.774642", "renewTime":
      "2022-06-16 18:55:32.822991"}'
  creationTimestamp: "2022-01-05T07:30:34Z"
  name: python-demo
  namespace: default
  resourceVersion: "85891471"
  selfLink: /api/v1/namespaces/default/configmaps/python-demo
  uid: ee023f69-4d91-475e-80b6-e4eac95e84bb




"""