from kubernetes import client, config, watch
from conf import *

def k8s_node_watch():
    try:
        config.load_kube_config(config_file=K8S_CONFIG)
        api = client.CoreV1Api()
        nodes = api.list_node(timeout_seconds=3).items
        print(nodes)
        print("----------------------------------------")
        print("watch start...")
        node_watch = watch.Watch()
        for event in node_watch.stream(client.CoreV1Api().list_node):
            node_name = event["object"].metadata.name
            node_status = " READY"
            for one_con in event["object"].status.conditions:
                if one_con.type == "Ready" and one_con.status != "True":
                    node_status = " NOT_READY!!!!!"
            print(node_name + node_status)
    except Exception as e:  # pylint: disable=broad-except
        print(e)
