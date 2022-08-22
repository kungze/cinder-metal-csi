# rabbitmq-standalone

在 K8s 环境中部署 [rabbitmq-standalone](https://github.com/kungze/rabbitmq-standalone) 

该组件提供一种不依赖于 pvc 的单节点部署方式

## 快速部署

```
helm repo add kungze https://kungze.github.io/cinder-metal-csi
helm install rabbitmq-standalone kungze rabbitmq-standalone
```


## Parameters

### Cluster Paramters

| Name                    | Form title            | Description                                  | Value           |
| ----------------------- | --------------------- | -------------------------------------------- | --------------- |
| `cluster_domain_suffix` | Cluster Domain Suffix | The doamin suffix of the current k8s cluster | `cluster.local` |


### Image Parameters

| Name             | Form title        | Description                                     | Value                   |
| ---------------- | ----------------- | ----------------------------------------------- | ----------------------- |
| `imageRegistry`  | Image Registry    | The registry address of openstack kolla image   | `registry.aliyuncs.com` |
| `imageNamespace` | Image Namespace   | The registry namespace of openstack kolla image | `kolla-helm`            |
| `openstackTag`   | Openstack version | The openstack version                           | `yoga`                  |
| `pullPolicy`     | Pull Policy       | The image pull policy                           | `IfNotPresent`          |


### Deployment Parameters

| Name                   | Form title              | Description                                                              | Value    |
| ---------------------- | ----------------------- | ------------------------------------------------------------------------ | -------- |
| `replicaCount`         |                         | Number of cinder replicas to deploy (requires ReadWriteMany PVC support) | `1`      |
| `serviceAccountName`   |                         | ServiceAccount name                                                      | `cinder` |
| `enableLivenessProbe`  | Enable Liveness Probe   | Whether or not enable liveness probe                                     | `true`   |
| `enableReadinessProbe` | Enable Readliness Probe | Whether or not enable readiness probe                                    | `true`   |


### rabbitmq Config parameters

| Name             | Form title        | Description            | Value                  |
| ---------------- | ----------------- | ---------------------- | ---------------------- |
| `rabbitmq.username` | RabbitMQ username   | RabbitMQ username | `openstack`            |
| `rabbitmq.password` | RabbitMQ password   | RabbitMQ password | `TXO9pjYDoX`           |
