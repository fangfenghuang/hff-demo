# helm命令：
helm create myservice
helm install .
helm install myservice --name myservice --namespace hfftest
helm list
helm package .

helm ls --all hfftest
helm del --purge hfftest



# examples/
  Chart.yaml          # Yaml文件，用于描述Chart的基本信息，包括名称版本等
  LICENSE             # [可选] 协议
  README.md           # [可选] 当前Chart的介绍
  values.yaml         # Chart的默认配置文件
  requirements.yaml   # [可选] 用于存放当前Chart依赖的其它Chart的说明文件
  charts/             # [可选]: 该目录中放置当前Chart依赖的其它Chart
  templates/          # [可选]: 部署文件模版目录，模版使用的值来自values.yaml和由Tiller提供的值
  templates/NOTES.txt # [可选]: 放置Chart的使用指南

