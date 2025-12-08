<font style="color:rgb(74, 74, 74);">KubeVela 作为一个声明式的应用交付控制平面，天然就可以以 GitOps 的方式进行使用，并且这样做会在 GitOps 的基础上为用户提供更多的益处和端到端的体验，包括：</font>

+ <font style="color:rgb(1, 1, 1);">应用交付工作流（CD 流水线）：KubeVela 支持在 GitOps 模式中描述过程式的应用交付，而不只是简单的声明终态；</font>
+ <font style="color:rgb(1, 1, 1);">处理部署过程中的各种依赖关系和拓扑结构；</font>
+ <font style="color:rgb(1, 1, 1);">在现有各种 GitOps 工具的语义之上提供统一的上层抽象，简化应用交付与管理过程；</font>
+ <font style="color:rgb(1, 1, 1);">统一进行云服务的声明、部署和服务绑定；</font>
+ <font style="color:rgb(1, 1, 1);">提供开箱即用的交付策略（金丝雀、蓝绿发布等）；</font>
+ <font style="color:rgb(1, 1, 1);">提供开箱即用的混合云/多云部署策略（放置规则、集群过滤规则等）；</font>
+ <font style="color:rgb(1, 1, 1);">在多环境交付中提供 Kustomize 风格的 Patch 来描述部署差异，而无需学习任何 Kustomize 本身的细节</font>

<font style="color:rgb(74, 74, 74);">GitOps 模式需要依赖 FluxCD 插件，所以在使用 GitOps 模式下交付应用之前需要先启用 FluxCD 插件。</font>

```bash
vela addon enable fluxcd
```

<font style="color:rgb(74, 74, 74);">GitOps 工作流分为 </font>**<font style="color:rgb(10, 10, 10);">CI</font>**<font style="color:rgb(74, 74, 74);"> 和 </font>**<font style="color:rgb(74, 74, 74);">CD</font>**<font style="color:rgb(74, 74, 74);"> 两个部分：</font>

+ **<font style="color:rgb(10, 10, 10);">CI</font>**<font style="color:rgb(1, 1, 1);">：持续集成对业务代码进行代码构建、构建镜像并推送至镜像仓库。目前有许多成熟的 CI 工具：如开源项目常用的 GitHub Action、Travis 等，以及企业中常用的 Jenkins、Tekton 等，KubeVela 围绕 GitOps 可以对接任意工具下的 CI 流程。</font>
+ **<font style="color:rgb(10, 10, 10);">CD</font>**<font style="color:rgb(1, 1, 1);">：持续部署会自动更新集群中的配置，如将镜像仓库中的最新镜像更新到集群中。目前主要有两种方案的 CD：</font>
    - **<font style="color:rgb(10, 10, 10);">Push-Based</font>**<font style="color:rgb(1, 1, 1);">：Push 模式的 CD 主要是通过配置 CI 流水线来完成的，这种方式需要将集群的访问秘钥共享给 CI，从而使得 CI 流水线能够通过命令将更改推送到集群中。前面我们讲解的 Jenkins 方式就属于该方案。</font>
    - **<font style="color:rgb(10, 10, 10);">Pull-Based</font>**<font style="color:rgb(1, 1, 1);">：Pull 模式的 CD 会在集群中监听仓库（代码仓库或者配置仓库）的变化，并且将这些变化同步到集群中。这种方式与 Push 模式相比，由集群主动拉取更新，从而避免了秘钥暴露的问题。前面课程中我们讲解的 Argo CD 与 Flux CD 就属于这种模式。</font>

<font style="color:rgb(74, 74, 74);">而交付面向的人员有以下两种：</font>

+ <font style="color:rgb(1, 1, 1);">面向平台管理员/运维人员的基础设施交付，用户可以通过直接更新仓库中的配置文件，从而更新集群中的基础设施配置，如系统的依赖软件、安全策略、存储、网络等基础设施配置。</font>
+ <font style="color:rgb(1, 1, 1);">面向终端开发者的交付，用户的代码一旦合并到应用代码仓库，就自动化触发集群中应用的更新，可以更高效的完成应用的迭代，与 KubeVela 的灰度发布、流量调拨、多集群部署等功能结合可以形成更为强大的自动化发布能力。</font>

## <font style="color:rgb(10, 10, 10);">面向平台管理员/运维人员的交付</font>
<font style="color:rgb(74, 74, 74);">如下图所示，对于平台管理员/运维人员而言，他们并不需要关心应用的代码，所以只需要准备一个 Git 配置仓库并部署 KubeVela 配置文件，后续对于应用及基础设施的配置变动，便可通过直接更新 Git 配置仓库来完成，使得每一次配置变更可追踪。</font>

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1700061786013-421cf655-ebcd-4a68-a198-8c1adf5b9308.png)

<font style="color:rgb(74, 74, 74);">这里我们将部署一个 MySQL 数据库作为项目的基础设施，同时部署一个业务应用，使用这个数据库。配置仓库的目录结构如下:</font>

+ **<font style="color:rgb(10, 10, 10);">clusters/</font>**<font style="color:rgb(1, 1, 1);"> </font><font style="color:rgb(1, 1, 1);">中包含集群中的 KubeVela GitOps 配置，用户需要将</font><font style="color:rgb(1, 1, 1);"> </font>**<font style="color:rgb(10, 10, 10);">clusters/</font>**<font style="color:rgb(1, 1, 1);"> </font><font style="color:rgb(1, 1, 1);">中的文件手动部署到集群中。这个是一次性的管控操作，执行完成后，KubeVela 便能自动监听配置仓库中的文件变动且自动更新集群中的配置。其中，</font>**<font style="color:rgb(10, 10, 10);">clusters/apps.yaml</font>**<font style="color:rgb(1, 1, 1);"> </font><font style="color:rgb(1, 1, 1);">将监听</font><font style="color:rgb(1, 1, 1);"> </font>**<font style="color:rgb(10, 10, 10);">apps/</font>**<font style="color:rgb(1, 1, 1);"> </font><font style="color:rgb(1, 1, 1);">下所有应用的变化，</font>**<font style="color:rgb(10, 10, 10);">clusters/infra.yaml</font>**<font style="color:rgb(1, 1, 1);"> </font><font style="color:rgb(1, 1, 1);">将监听</font><font style="color:rgb(1, 1, 1);"> </font>**<font style="color:rgb(10, 10, 10);">infrastructure/</font>**<font style="color:rgb(1, 1, 1);"> </font><font style="color:rgb(1, 1, 1);">下所有基础设施的变化。</font>
+ **<font style="color:rgb(10, 10, 10);">apps/</font>**<font style="color:rgb(1, 1, 1);"> </font><font style="color:rgb(1, 1, 1);">目录中包含业务应用的所有配置，在本例中为一个使用数据库的业务应用。</font>
+ **<font style="color:rgb(10, 10, 10);">infrastructure/</font>**<font style="color:rgb(1, 1, 1);"> </font><font style="color:rgb(1, 1, 1);">中包含一些基础设施相关的配置和策略，在本例中为 MySQL 数据库。</font>

```bash
├── apps
│   └── my-app.yaml
├── clusters
│   ├── apps.yaml
│   └── infra.yaml
└── infrastructure
    └── mysql.yaml
```

:::info
💡KubeVela 建议使用如上的目录结构管理你的 GitOps 仓库。**clusters/** 中存放相关的 KubeVela GitOps 配置并需要被手动部署到集群中，**apps/** 和 **infrastructure/** 中分别存放你的应用和基础设施配置。通过把应用和基础配置分开，能够更为合理的管理你的部署环境，隔离应用的变动影响。

:::

**<font style="color:rgb(10, 10, 10);">clusters/ 目录</font>**

<font style="color:rgb(74, 74, 74);">首先，我们来看下 clusters 目录，这也是 KubeVela 对接 GitOps 的初始化操作配置目录。</font>

<font style="color:rgb(74, 74, 74);">以</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">clusters/infra.yaml</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">为例：</font>

```yaml
apiVersion: core.oam.dev/v1beta1
kind: Application
metadata:
  name: infra
spec:
  components:
    - name: database-config
      type: kustomize
      properties:
        repoType: git
        # 将此处替换成你需要监听的 git 配置仓库地址
        url: https://github.com/cnych/KubeVela-GitOps-Infra-Demo
        # 如果是私有仓库，还需要关联 git secret
        # secretRef: git-secret
        # 自动拉取配置的时间间隔，由于基础设施的变动性较小，此处设置为十分钟
        pullInterval: 10m
        git:
          # 监听变动的分支
          branch: main
        # 监听变动的路径，指向仓库中 infrastructure 目录下的文件
        path: ./infrastructure
```

**<font style="color:rgb(10, 10, 10);">apps.yaml</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">与</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">infra.yaml</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">几乎保持一致，只不过监听的文件目录有所区别。在</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">apps.yaml</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">中，</font>**<font style="color:rgb(10, 10, 10);">properties.path</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">的值将改为</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">./apps</font>**<font style="color:rgb(74, 74, 74);">，表明监听</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">apps/</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">目录下的文件变动。</font>

<font style="color:rgb(74, 74, 74);">cluster 文件夹中的 GitOps 管控配置文件需要在初始化的时候一次性手动部署到集群中，在此之后 KubeVela 将自动监听</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">apps/</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">以及</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">infrastructure/</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">目录下的配置文件并定期更新同步。</font>

**<font style="color:rgb(10, 10, 10);">apps/ 目录</font>**

**<font style="color:rgb(10, 10, 10);">apps/</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">目录中存放着应用配置文件，这是一个配置了数据库信息以及 Ingress 的简单应用。该应用将连接到一个 MySQL 数据库，并简单地启动服务。在默认的服务路径下，会显示当前版本号。在</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">/db</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">路径下，会列出当前数据库中的信息。</font>

```yaml
apiVersion: core.oam.dev/v1beta1
kind: Application
metadata:
  name: my-app
  namespace: default
spec:
  components:
    - name: my-server
      type: webservice
      properties:
        image: cnych/kubevela-gitops-demo:main-76a34322-1697703461
        port: 8088
        env:
          - name: DB_HOST
            value: mysql-cluster-mysql.default.svc.cluster.local:3306
          - name: DB_PASSWORD
            valueFrom:
              secretKeyRef:
                name: mysql-secret
                key: ROOT_PASSWORD
      traits:
        - type: scaler
          properties:
            replicas: 1
        - type: gateway
          properties:
            class: nginx
            classInSpec: true
            domain: vela-gitops-demo.k8s.local
            http:
              /: 8088
            pathType: ImplementationSpecific
```

<font style="color:rgb(74, 74, 74);">这是一个使用了 KubeVela 内置组件类型</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">webservice</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">的应用，该应用绑定了</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">gateway</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">运维特征。通过在应用中声明运维能力的方式，只需一个文件，便能将底层的 Deployment、Service、Ingress 集合起来，从而更为便捷地管理应用。</font>

**<font style="color:rgb(10, 10, 10);">infrastructure/ 目录</font>**

**<font style="color:rgb(10, 10, 10);">infrastructure/</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">目录下存放一些基础设施的配置。此处我们使用 mysql controller 来部署了一个 MySQL 集群。</font>

```yaml
apiVersion: core.oam.dev/v1beta1
kind: Application
metadata:
  name: mysql
  namespace: default
spec:
  components:
    - name: mysql-secret
      type: k8s-objects # 需要添加一个包含 ROOT_PASSWORD 的 secret
      properties:
        objects:
          - apiVersion: v1
            kind: Secret
            metadata:
              name: mysql-secret
            type: Opaque
            stringData:
              ROOT_PASSWORD: root321
    - name: mysql-operator
      type: helm
      properties:
        repoType: helm
        url: https://helm-charts.bitpoke.io
        chart: mysql-operator
        version: 0.6.3
    - name: mysql-cluster
      type: raw
      dependsOn:
        - mysql-operator
        - mysql-secret
      properties:
        apiVersion: mysql.presslabs.org/v1alpha1
        kind: MysqlCluster
        metadata:
          name: mysql-cluster
        spec:
          replicas: 1
          secretName: mysql-secret
```

<font style="color:rgb(74, 74, 74);">在这个 MySQL 应用中，我们添加了 3 个 KubeVela 的组件，第一个是一个</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">k8s-objects</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">类型的组件，也就是直接应用 Kubernetes 资源对象，我们这里需要部署一个 Secret 对象；然后添加一个</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">helm</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">类型的组件，用来部署 MySQL 的 Operator。当 Operator 部署成功且正确运行后，最后我们将开始部署 MySQL 集群。</font>

**<font style="color:rgb(10, 10, 10);">部署 clusters/ 目录下的文件</font>**

<font style="color:rgb(74, 74, 74);">配置完以上文件并存放到 Git 配置仓库后，我们需要在集群中手动部署</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">clusters/</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">目录下的 KubeVela GitOps 配置文件。</font>

<font style="color:rgb(74, 74, 74);">首先，在集群中部署</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">clusters/infra.yaml</font>**<font style="color:rgb(74, 74, 74);">。可以看到它自动在集群中拉起了</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">infrastructure/</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">目录下的 MySQL 部署文件：</font>

```bash
$ kubectl apply -f clusters/infra.yaml
$ vela ls
APP             COMPONENT       TYPE            TRAITS          PHASE   HEALTHY STATUS                                                          CREATED-TIME
infra           database-config kustomize                       running healthy                                                                 2023-10-19 15:27:28 +0800 CST
mysql           mysql-operator  helm                            running healthy Fetch repository successfully, Create helm release              2023-10-19 15:27:31 +0800 CST
                                                                                successfully
└─              mysql-cluster   raw                             running healthy                                                                 2023-10-19 15:27:31 +0800 CST
```

<font style="color:rgb(74, 74, 74);">至此，我们通过部署 KubeVela GitOps 配置文件，自动在集群中拉起了数据库基础设施。</font>

```bash
$ kubectl get pods
NAME                                     READY   STATUS    RESTARTS         AGE
mysql-cluster-mysql-0                    4/4     Running   0                35m
mysql-operator-0                         2/2     Running   0                35m
```

<font style="color:rgb(74, 74, 74);">通过这种方式，我们可以方便地通过更新 Git 配置仓库中的文件，从而自动化更新集群中的配置。</font>

## <font style="color:rgb(10, 10, 10);">面向终端开发者的交付</font>
<font style="color:rgb(74, 74, 74);">对于终端开发者而言，在 KubeVela Git 配置仓库以外，还需要准备一个应用代码仓库。在用户更新了应用代码仓库中的代码后，需要配置一个 CI 来自动构建镜像并推送至镜像仓库中。KubeVela 会监听镜像仓库中的最新镜像，并自动更新配置仓库中的镜像配置，最后再更新集群中的应用配置。使用户可以达成在更新代码后，集群中的配置也自动更新的效果，代码仓库位于</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">https://github.com/cnych/KubeVela-GitOps-App-Demo</font>**<font style="color:rgb(74, 74, 74);">。</font>

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1700061786095-524bc92b-a4ae-4c4f-b5df-c948cabf85d6.png)

**<font style="color:rgb(10, 10, 10);">准备代码仓库</font>**

<font style="color:rgb(74, 74, 74);">准备一个代码仓库，里面包含一些源代码以及对应的</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">Dockerfile</font>**<font style="color:rgb(74, 74, 74);">。这些代码将连接到一个 MySQL 数据库，并简单地启动服务。在默认的服务路径下，会显示当前版本号。在</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">/db</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">路径下，会列出当前数据库中的信息，基本代码如下所示：</font>

```go
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    _, _ = fmt.Fprintf(w, "Version: %s\n", VERSION)
})
http.HandleFunc("/db", func(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("select * from userinfo;")
    if err != nil {
        _, _ = fmt.Fprintf(w, "Error: %v\n", err)
    }
    for rows.Next() {
        var username string
        var desc string
        err = rows.Scan(&username, &desc)
        if err != nil {
            _, _ = fmt.Fprintf(w, "Scan Error: %v\n", err)
        }
        _, _ = fmt.Fprintf(w, "User: %s \nDescription: %s\n\n", username, desc)
    }
})

if err := http.ListenAndServe(":8088", nil); err != nil {
    panic(err.Error())
}
```

<font style="color:rgb(74, 74, 74);">我们希望用户改动代码进行提交后，自动构建出最新的镜像并推送到镜像仓库。这一步 CI 可以通过前面我们讲解的 Jenkins 来实现，基本一致。</font>

<font style="color:rgb(74, 74, 74);">首先为代码仓库创建一个 Webhook，指向 Jenkins 的触发器地址：</font>

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1700061785982-6db1550b-f02a-4406-bee8-92e22915b93c.png)

<font style="color:rgb(74, 74, 74);">然后在 Jenkins 中创建一个名为</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">KubeVela-GitOps-App-Demo</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">的流水线：</font>

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1700061786015-3ef316f8-ce8b-44e9-b601-e8c52e9700cb.png)

<font style="color:rgb(74, 74, 74);">并勾选</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">GitHub hook trigger for GITScm polling</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">触发器。</font>

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1700061786055-61fd2297-3f54-4425-93df-303dc4585ea2.png)

<font style="color:rgb(136, 136, 136);">触发器</font>

<font style="color:rgb(74, 74, 74);">然后添加如下所示的流水线脚本：</font>

```go
void setBuildStatus(String message, String state) {
    step([
        $class: "GitHubCommitStatusSetter",
        reposSource: [$class: "ManuallyEnteredRepositorySource", url: "https://github.com/cnych/KubeVela-GitOps-App-Demo"],
        contextSource: [$class: "ManuallyEnteredCommitContextSource", context: "ci/jenkins/deploy-status"],
        errorHandlers: [[$class: "ChangingBuildStatusErrorHandler", result: "UNSTABLE"]],
        statusResultSource: [ $class: "ConditionalStatusResultSource", results: [[$class: "AnyBuildResult", message: message, state: state]] ]
    ]);
}
pipeline {
    agent {
        kubernetes {
            cloud 'Kubernetes'
            defaultContainer 'jnlp'
            yaml '''
            spec:
            serviceAccountName: jenkins
            containers:
            - name: golang
            image: golang:1.16-alpine3.15
            command:
            - cat
            tty: true
            - name: docker
            image: docker:latest
            command:
            - cat
            tty: true
            env:
            - name: DOCKER_HOST
            value: tcp://docker-dind:2375
            '''
        }
    }
    stages {
        stage('Prepare') {
            steps {
                script {
                    def checkout = git branch: 'main', url: 'https://github.com/cnych/KubeVela-GitOps-App-Demo.git'
                    env.GIT_COMMIT = checkout.GIT_COMMIT
                    env.GIT_BRANCH = checkout.GIT_BRANCH

                    def unixTime = (new Date().time.intdiv(1000))
                    def gitBranch = env.GIT_BRANCH.replace("origin/", "")
                    env.BUILD_ID = "${gitBranch}-${env.GIT_COMMIT.substring(0,8)}-${unixTime}"

                    echo "env.GIT_BRANCH=${env.GIT_BRANCH},env.GIT_COMMIT=${env.GIT_COMMIT}"
                    echo "env.BUILD_ID=${env.BUILD_ID}"

                    setBuildStatus("Deploy running", "PENDING");
                }
            }
        }
        stage('Test') {
            steps {
                container('golang') {
                    sh 'GOPROXY=https://goproxy.io CGO_ENABLED=0 GOCACHE=$(pwd)/.cache go test *.go'
                }
            }
        }
        stage('Build') {
            steps {
                withCredentials([[$class: 'UsernamePasswordMultiBinding',
                                credentialsId: 'docker-auth',
                                usernameVariable: 'DOCKER_USER',
                                passwordVariable: 'DOCKER_PASSWORD']]) {
                    container('docker') {
                        sh """
                        docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWORD}
                        docker build -t cnych/kubevela-gitops-demo:${env.BUILD_ID} .
                        docker push cnych/kubevela-gitops-demo:${env.BUILD_ID}
                        """
                    }
                }
            }
        }
    }
    post {
        success {
            setBuildStatus("Deploy success", "SUCCESS");
        }
        failure {
            setBuildStatus("Deploy failed", "FAILURE");
        }
    }
}
```

<font style="color:rgb(74, 74, 74);">构建后我们就可以将应用的镜像打包后推送到 Docker Hub 去。</font>

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1700061786345-512c4a6c-3310-471c-ad72-7fc98b66f369.png)

**<font style="color:rgb(10, 10, 10);">配置秘钥信息</font>**

<font style="color:rgb(74, 74, 74);">在新的镜像推送到镜像仓库后，KubeVela 会识别到新的镜像，并更新仓库及集群中的 Application 配置文件。因此，我们需要一个含有 Git 信息的 Secret，使 KubeVela 向 Git 仓库进行提交。部署如下文件，将其中的用户名和密码替换成你的 Git 用户名及密码（或 Token）：</font>

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: git-secret
type: kubernetes.io/basic-auth
stringData:
  username: <your username>
  password: <your password>
```

**<font style="color:rgb(10, 10, 10);">准备配置仓库</font>**

<font style="color:rgb(74, 74, 74);">配置仓库与之前面向运维人员的配置大同小异，只需要加上与镜像仓库相关的配置即可。</font>

<font style="color:rgb(74, 74, 74);">修改</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">clusters/</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">中的</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">apps.yaml</font>**<font style="color:rgb(74, 74, 74);">，该 GitOps 配置会监听仓库中</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">apps/</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">下的应用文件变动以及镜像仓库中的镜像更新：</font>

```yaml
# ...... 省略其他的
imageRepository:
  # 镜像地址
  image: <your image>
  # 如果这是一个私有的镜像仓库，可以通过 `kubectl create secret docker-registry` 创建对应的镜像秘钥并相关联
  secretRef: dockerhub-secret
  filterTags:
    # 可对镜像 tag 进行过滤
    pattern: "^main-[a-f0-9]+-(?P<ts>[0-9]+)"
    extract: "$ts"
  # 通过 policy 筛选出最新的镜像 Tag 并用于更新
  policy:
    numerical:
      order: asc
  # 追加提交信息
  commitMessage: "Image: {{range .Updated.Images}}{{println .}}{{end}}"
```

<font style="color:rgb(74, 74, 74);">修改</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">apps/my-app.yaml</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">中的</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">image</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">字段，在后面加上</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);"># {"$imagepolicy": "default:apps"}</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">的注释，KubeVela 会通过该注释去更新对应的镜像字段，</font>**<font style="color:rgb(10, 10, 10);">default:apps</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">是上面 GitOps 配置对应的命名空间和名称。</font>

```yaml
spec:
  components:
    - name: my-server
      type: webservice
      properties:
        image: cnych/kubevela-gitops-demo:main-9e8d2465-1697703645 # {"$imagepolicy": "default:apps"}
```

<font style="color:rgb(74, 74, 74);">将</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">clusters/</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">中包含镜像仓库配置的文件更新到集群中后，我们便可以通过修改代码来完成应用的更新。</font>

<font style="color:rgb(74, 74, 74);">部署</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">clusters/apps.yaml</font>**<font style="color:rgb(74, 74, 74);">：</font>

```bash
$ kubectl apply -f clusters/apps.yaml
$ vela ls
APP             COMPONENT       TYPE            TRAITS          PHASE           HEALTHY         STATUS                                                     CREATED-TIME
apps            apps            kustomize                       running         healthy                                                                    2023-10-19 16:31:49 +0800 CST
my-app          my-server       webservice      scaler,gateway  runningWorkflow unhealthy       Ready:0/1                                                  2023-10-19 16:32:11 +0800 CST
$ kubectl get pods
NAME                                     READY   STATUS    RESTARTS         AGE
my-server-6947fd65f9-84zhv               1/1     Running   0                2m
```

<font style="color:rgb(74, 74, 74);">这样我们就可以通过部署 KubeVela GitOps 配置文件，自动在集群中拉起应用了。我们可以通过</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">curl</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">应用的 Ingress 来验证结果是否正确，可以看到目前的版本是 0.1.5，并且成功地连接到了数据库：</font>

```bash
$ kubectl get ingress
NAME           CLASS   HOSTS                        ADDRESS   PORTS   AGE
my-server      nginx   vela-gitops-demo.k8s.local             80      115s
$ curl -H "Host:vela-gitops-demo.k8s.local" http://192.168.0.100
Version: 0.1.8
$ curl -H "Host:vela-gitops-demo.k8s.local" http://192.168.0.100/db
User: KubeVela
Description: It's a test user
```

**<font style="color:rgb(10, 10, 10);">修改代码</font>**

<font style="color:rgb(74, 74, 74);">将代码文件中的 Version 改为</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">0.2.0</font>**<font style="color:rgb(74, 74, 74);">，并修改数据库中的数据:</font>

```go
const VERSION = "0.2.0"

...

func InsertInitData(db *sql.DB) {
    stmt, err := db.Prepare(insertInitData)
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    _, err = stmt.Exec("KubeVela2", "It's another test user")
    if err != nil {
        panic(err)
    }
}
```

<font style="color:rgb(74, 74, 74);">提交该改动至代码仓库，正常我们配置的 CI 流水线就会自动开始构建镜像并推送至镜像仓库。</font>

<font style="color:rgb(74, 74, 74);">而 KubeVela 会通过监听镜像仓库，根据最新的镜像 Tag 来更新配置仓库中</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">apps/</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">下的应用</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">my-app</font>**<font style="color:rgb(74, 74, 74);">。</font>

<font style="color:rgb(74, 74, 74);">此时，可以看到配置仓库中有一条来自 kubevelabot 的提交，提交信息均带有</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">Update image automatically.</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">前缀。你也可以通过</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">{{range .Updated.Images}}{{println .}}{{end}}</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">在</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">commitMessage</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">字段中追加你所需要的信息。</font>

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1700061786510-0e1ed0ef-2b86-46f6-bf12-7015eac26a37.png)

<font style="color:rgb(74, 74, 74);">经过一段时间后，应用</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">my-app</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">就自动更新了。KubeVela 会通过你配置的</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">interval</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">时间间隔，来每隔一段时间分别从配置仓库及镜像仓库中获取最新信息：</font>

+ <font style="color:rgb(1, 1, 1);">当 Git 仓库中的配置文件被更新时，KubeVela 将根据最新的配置更新集群中的应用。</font>
+ <font style="color:rgb(1, 1, 1);">当镜像仓库中多了新的 Tag 时，KubeVela 将根据你配置的 policy 规则，筛选出最新的镜像 Tag，并更新到 Git 仓库中。而当代码仓库中的文件被更新后，KubeVela 将重复第一步，更新集群中的文件，从而达到了自动部署的效果。</font>

<font style="color:rgb(74, 74, 74);">通用我们可以通过 curl 对应的 Ingress 查看当前版本和数据库信息：</font>

```bash
$ kubectl get ingress
NAME           CLASS   HOSTS                        ADDRESS   PORTS   AGE
my-server      nginx   vela-gitops-demo.k8s.local             80      12m

$ curl -H "Host:vela-gitops-demo.k8s.local" http://<ingress-ip>
Version: 0.2.0

$ curl -H "Host:vela-gitops-demo.k8s.local" http://<ingress-ip>/db
User: KubeVela
Description: It's a test user

User: KubeVela2
Description: It's another test user
```

<font style="color:rgb(74, 74, 74);">版本已被成功更新！至此，我们完成了从变更代码，到自动部署至集群的全部操作。</font>

## <font style="color:rgb(10, 10, 10);">总结</font>
<font style="color:rgb(74, 74, 74);">在运维侧，如若需要更新基础设施（如数据库）的配置，或是应用的配置项，只需要修改配置仓库中的文件，KubeVela 将自动把配置同步到集群中，简化了部署流程。</font>

<font style="color:rgb(74, 74, 74);">在研发侧，用户修改代码仓库中的代码后，KubeVela 将自动更新配置仓库中的镜像，从而进行应用的版本更新。通过与 GitOps 的结合，KubeVela 加速了应用从开发到部署的整个流程。可能你会觉得这和 Flux CD 不是差不多吗？的确是这样的，KubeVela 的 GitOps 功能本身就是依赖 Flux CD 的，但是 KubeVela 的功能可远远不止于此，比如说上面我们的应用使用的 MySQL 数据我们是通过 MySQL Operator 来部署的，那如果我现在还换成云资源 RDS 呢？按照以前的方式方法，那么我们需要去云平台手动开通 RDS 或者使用 Terraform 来进行管理，但在 KubeVela 中我们完全可以帮助开发者集成、编排不同类型的云资源，涵盖混合多云环境，让你用统一地方式去使用不同厂商的云资源。同样的我们只需要在 GitOps 仓库中的配置文件 Application 对象中去添加云资源的管理配置即可，这样做到了一个对象管理多种资源的能力，这也是 KubeVela 的核心能力之一。</font>

<font style="color:rgb(74, 74, 74);">最后如果你觉得应用太多管理不太方便，那么我们还可以使用</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">vela top</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">命令获取平台的概览信息以及对应用程序的资源状态进行诊断。</font>

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1700061786514-cf348f88-12ad-49aa-ac43-74a8fb66b7c5.png)

参考文档：[https://kubevela.io/zh/docs/end-user/gitops/fluxcd/](https://kubevela.io/zh/docs/end-user/gitops/fluxcd/)

