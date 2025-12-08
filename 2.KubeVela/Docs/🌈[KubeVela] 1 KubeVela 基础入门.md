KubeVela 是一个开箱即用的现代化应用交付与管理平台，它使得应用在面向混合云环境中的交付更简单、快捷，是**<u><font style="color:#DF2A3F;">开放应用模型（OAM）</font></u>**的一个实现，所以我们需要先了解下 OAM。

## 1 OAM 简介
OAM（Open Application Model） 是阿里巴巴和微软共同开源的**<font style="color:#DF2A3F;">云原生应用规范模型</font>**，OAM 的本质是根据软件设计的**<font style="color:#DF2A3F;">兴趣点分离</font>**原则对负责的 DevOps 流程的高度抽象和封装，一个以应用为中心的 Kubernetes API 分层，这种模型旨在**<font style="color:#DF2A3F;">定义云原生应用的标准</font>**。

![OAM](https://cdn.nlark.com/yuque/0/2023/png/2555283/1696949642357-5026ad35-7d68-4e37-935a-79c659a65757.png)

从 OAM 名称中可以看出，它是一个开放的应用模型：

+ 开放（Open）：支持异构的平台、容器运行时、调度系统、云供应商、硬件配置等，总之与底层无关
+ 应用（Application）：云原生应用
+ 模型（Model）：定义标准，以使其与底层平台无关

### 1.1 为什么我们需要 OAM 模型呢？
现阶段应用管理的主要面临两个挑战：

+ 对应用研发而言，Kubernetes 的 API 针对简单应用过于复杂，针对复杂应用却又难以上手；
+ 对应用运维而言，Kubernetes 的扩展能力难以管理；Kubernetes 原生的 API 没有对云资源全部涵盖。

总体而言，面临的主要挑战就是：**如何基于 Kubernetes 提供真正意义上的应用管理平台，让研发和运维只需关注到应用本身**。

比如下面是一个典型的 K8s 资源清单文件，该 yaml 文件已经是被简化过的，但实际上还是比较长。

![Kubernetes Yaml](https://cdn.nlark.com/yuque/0/2023/png/2555283/1696949642361-0b95d48b-4932-497a-a5a3-509502d512f5.png)

自上而下，我们可以大致把它们分为三块：

+ 一块是扩缩容、滚动升级相关的参数，这一块一般是运维的同学比较关心的；
+ 中间一块是镜像、端口、启动参数相关的，这一部分应该是开发的同学比较关心的；
+ 最后一块大家可能根本看不懂，当然大部分情况下也不太需要明白，一般来说这属于 K8s 平台层的同学需要关心的内容。

看到这样一个 yaml 文件，我们很容易想到，只要把里面的字段封装一下，把该暴露的暴露出来就好了。这个时候我们就可以去开发一个应该管理平台，并做一个漂亮的前端界面给用户，只暴露给用户 5 个左右的字段，这显然可以大大降低用户理解 Kubernetes 的心智负担，底层实现用类似模板的方式把这五个字段渲染成一个完整的 yaml 文件。

```yaml
image: :v0.34.0
args:
  - --logtostderr=true
ports:
  - containerPort: 8080
    name: http
    protocol: TCP
envs:
  - name: INNER-KEY
    value: app
volumes:
  - name: cache-volume
    emptyDir: {}
```

该方式针对简单无状态应用非常有效，精简 API 可以大大降低 Kubernetes 的门槛。但是当出现大规模业务后，就会遇到很多复杂的应用，这个时候就会发现该 PaaS 应用平台能力不够了。比如 ZK 多实例选主、主从切换这些逻辑，在这五个字段里就很难描述了。因为屏蔽大量字段的方式会限制基础设施本身的能力，而 Kubernetes 的能力是非常强大而灵活的，所以我们不可能为了简化而放弃掉 Kubernetes 本身强大的能力。

中间件的工程师跟我们说，我这有个 Zookeeper 该用哪种 Kubernetes 的工作负载接入呢？我们当然会想到可以让他们使用 Operator 了，于是他们就很懵逼的说到 Operator 是啥？

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1696949642391-3d5f1d9d-69e5-470c-84a1-17b01e9fd613.png)

### 1.2 Operator是啥
然后我们耐心的给他解释相关概念 **CRD**、**Controller**、**Informer**、**Reflector**、**Indexer** 这些就可以了，当然他们就更懵了，当然理论上也不需要理解。业务方更应该专注于他们的业务本身，当然我们就不得不帮他们一起写这个控制器了。为此我们需要一个统一的模型去解决研发对应用管理的诉求。

:::success
**除了研发侧的问题之外，在运维侧同样也会有很多挑战。**

:::

**<u>K8s 的 CRD Operator 机制非常灵活而强大，不光是复杂应用可以通过编写 CRD Operator 实现，运维能力当然也可以通过 Operator 来扩展，比如灰度发布、流量管理、弹性扩缩容等等。</u>**

比如有一个案例就是开发了一个 CronHPA 的 CRD，可以定时设置 HPA 的范围，但是应用运维却并不知道该 CRD 会跟原生的 HPA 会产生冲突，结果自然是引起了故障。这血的教训提醒我们要做事前检查，熟悉 K8s 的机制很容易让我们想到为每个 Operator 加上 Admission Webhook。这个 Admission Webhook 需要拿到这个应用绑定的所有运维能力以及应用本身的运行模式，然后做统一的校验。如果这些运维能力都是一方提供的还好，如果存在两方，甚至三方提供的扩展能力，我们就没有一个统一的方式去获知了。

如果再深入思考下就知道我们需要一个统一的模型来协商并管理这些复杂的扩展能力。

云原生应用有一个很大的特点，那就是它往往会依赖云上的资源，包括数据库、网络、负载均衡、缓存等一系列资源。

当我们交付应用的时候比如使用 Helm 进行打包，我们只能针对 K8s 原生 API，而如果我们还想启动 RDS 数据库，就比较困难了。如果不想去数据库的交互页面，想通过 K8s 的 API 来管理，那就又不得不去写一个 CRD 来定义了，然后通过 Operator 去调用实际云资源的 API。

这一整套交付物实际上就是一个应用的完整描述，即我们所说的“应用定义”。这种定义方式最终所有的配置还是会全部堆叠到一个 yaml 文件里，这跟前面说的 all-in-one 问题其实是一样的，而且，这些应用定义最终也都成为了黑盒，除了对应项目本身可以使用，其他系统基本无法复用了。

而且事实上很多公司和团队也在根据自身业务需要进行定义，比如 Pinterest 定义的应用规范如下所示：

:::color1
Pinterest CRD

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1696949642457-484fe4b6-df5b-4971-941a-33b99811ae3e.png)

:::

应用定义实际上是应用交付/分发不可或缺的部分，所以我们可以思考下是否可以定义足够开放的、可复用的应用模型呢？

一个应用定义需要容易上手，但又不失灵活性，更不能是一个黑盒。应用定义同样需要跟开源生态紧密结合，没有生态的应用定义注定是没有未来的，自然也很难持续的迭代和演进。

这也是为什么我们需要 OAM 的深层次的原因！！！

前面我们说的各种问题，归根结底在于 Kubernetes 的 All in One API 是为平台提供者设计的，我们不能像下图左侧显示的一样，让应用的研发、运维跟 Kubernetes 团队一样面对这一团 API。

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1696949642546-cc572816-4011-4dc6-ae26-0eca00bb4b4d.png)

一个合理的应用模型应该具有区分使用者角色的分层结构，同时将运维能力模块化的封装。让不同的角色使用不同的 API，如上图右侧部分。

### 1.3 OAM 模型定义
上面我们了解了为什么需要 OAM 模型，那么 OAM 模型到底是如何定义的呢？

在最新的 API 版本 v0.3.0 版本(core.oam.dev/v1beta1)中，OAM 定义了以下内容：

+ **ComponentDefiniton**：组件模型，OAM 中最基础的单元，应用程序中的每个微服务都可以被描述为一个组件，在实践中，一个简单的容器化工作负载、Helm Chart 或云数据库都可以定义为一个组件。
+ **WorkloadDefiniton**: 工作负载是一个特定组件定义的关键特征，由平台提供，以便用户可以检查平台并了解哪些工作负载类型可供使用。请注意，工作负载类型不允许最终用户创建新的（仅限平台提供商） 。
+ **TraitDefinition**: 为组件工作负载实例增加的运维特征，运维人员可以对组件的配置做出具体的决定。例如，向 WordPress Helm Chart 的工作负载注入 sidecar 容器的 sidecar trait。特征可以是适用于单个组件的分布式应用程序的任何配置，例如负载均衡策略、网络入口路由、断路器、速率限制、自动扩展策略、升级策略等，特征是运维人员的关注点。
+ **Application Scope**: 应用范围是通过提供不同形式的应用边界和相同组的行为，将组件组合成逻辑应用。应用范围可以决定组件是否可以被同时部署到同一应用范围类型的多个实例中。
+ **Application**: Application 定义了在部署应用程序后将被实例化的组件列表。

因此，一个应用程序是由一组具有一些运维特征的组件组成的，并且被限定在一个或多个应用程序边界中。

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1696949643290-d1d7a109-597f-44c9-87a2-bfb6cebb7353.png)

具体的模型定义规范可以查看 OAM Spec 文档了解更多，不过需要注意的是现在 KubeVela 的规范和 OAM 的规范并不是完全一样的。

## 2 KubeVela 简介
KubeVela 是 OAM 规范（实际上 OAM 规范会滞后于 KubeVela 中使用的规范）的一个实现，是一个开箱即用的现代化应用交付与管理平台，它使得应用在面向混合云环境中的交付更简单、快捷。使用 KubeVela 的软件开发团队，可以按需使用云原生能力构建应用，随着团队规模的发展、业务场景的变化扩展其功能，一次构建应用，随处运行。

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1696949643260-21151d4e-618f-477c-b244-a2b4780a9e63.png)

### 2.1 核心功能
KubeVela 具有以下几个核心功能：

+ 应用部署即代码（Deployment as Code），完整定义全交付流程使用 OAM 作为应用交付的顶层抽象，这种方式使你可以用声明式的方式描述应用交付全流程，自动化的集成 CI/CD 及 GitOps 体系，通过  CUE  轻松扩展或重新编写你的交付过程。
+ 天然支持企业级集成，安全、合规、可观测性一应俱全支持多集群认证和授权并与 K8s RBAC 集成，还可以从社区的插件中心找到一系列开箱即用的平台扩展，包括多种用户体系（LDAP 等）集成、多租户权限控制、安全校验和扫描、应用可观测性等大量企业级能力。
+ 面向多云多集群混合环境，丰富的应用交付和管理能力原生支持丰富的多集群/混合环境持续交付策略，包括金丝雀、蓝绿、多环境差异化配置等，同样也支持跨环境交付，这些交付策略为你的分布式交付流程提供了充足的效率和安全保证。
+ 轻量并且架构高度可扩展，满足企业不同场景的定制化需求KubeVela 最小的部署模式仅需 1 个 pod （0.5 核 1G 内存）就可以用于部署上千个应用。其微内核、高可扩展的架构可以轻松满足你的扩展和定制化需求，衔接企业内部的权限体系、微服务、流量治理、配置管理、可观测性等模块。不仅如此，社区还有一个正在快速增长的插件市场可供你选择和使用，你可以在这里贡献、复用社区丰富的功能模块。

### 2.2 关注点分离
关注点分离这个属于 KubeVela 的核心理念，它是 KubeVela 的设计哲学，也是 KubeVela 与众不同的地方。KubeVela 的用户天然分为两种角色，由公司的两个团队（或个人）承担。

+ **平台团队**由平台工程师完成，他们需要准备应用部署环境，维护稳定可靠的基础设施功能（如 mysql operator），并将基础设施能力作为 KubeVela 模块定义注册到集群中。他们需要具备丰富的基础设施经验。
+ **最终用户**最终用户即业务应用的开发者，使用平台的过程中首先选择部署环境，然后挑选能力模块，填写业务参数并组装成 KubeVela 应用。他们无需关心基础设施细节。

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1696949643292-1bc183c4-56a5-4f30-b6c9-b21860a1103f.png)

### 2.3 核心概念
KubeVela 遵循 OAM 规范通过一个 **Application** 的对象来声明一个微服务应用的完整交付流程，其中包含了待交付组件、关联的运维能力、交付流水线等内容。所有待交付组件、运维动作和流水线中的每一个步骤，都遵循 OAM 规范设计为独立的可插拔模块，允许用户按照自己的需求进行组合或者定制。

基本上 **Application**  对象是终端用户唯一需要了解的对象，它表达了一个微服务应用的部署计划。遵循 OAM 规范，一个应用部署计划（Application）由组件（Component）、运维特征（Trait）、部署工作流（Workflow）、应用执行策略（Policy）四部分组成，这些组件是平台构建者维护的可编程模块，这种抽象方式是高度可扩展、可定制的。

+ **组件（Component）**组件是构成微服务应用的基本单元。一个应用中可以包括多个组件，最佳的实践方案是一个应用中包括一个主组件（核心业务）和附属组件（强依赖或独享的中间件，运维组件等）。KubeVela 内置支持多种类型的组件交付，包括 Helm Chart、容器镜像、CUE 模块、Terraform 模块等。同时也允许平台管理员以 CUE 语言的形式定制其它任意类型的组件。
+ **运维特征（Trait）**运维特征是可以随时绑定给待部署组件的模块化、可拔插的运维能力，比如：副本数调整、数据持久化、设置网关策略、自动设置 DNS 解析等。用户可以从社区获取成熟的能力，也可以自行定义。
+ **工作流（Workflow）**工作流由多个步骤组成，允许用户自定义应用在某个环境的交付过程。典型的工作流步骤包括人工审核、数据传递、多集群发布、通知等。
+ **应用策略（Policy）**应用策略负责定义指定应用交付过程中的策略，比如多集群部署的差异化配置、安全组策略、防火墙规则等。

整体定义如下所示：

```yaml
apiVersion: core.oam.dev/v1beta1
kind: Application
metadata:
  name: <应用名称>
spec:
  components:
    - name: <组件名称1>
      type: <组件类型1>
      properties: <组件参数>
      traits:
        - type: <运维特征类型1>
          properties: <运维特征类型>
        - type: <运维特征类型2>
          properties: <运维特征类型>
    - name: <组件名称2>
      type: <组件类型2>
      properties: <组件参数>
  policies:
    - name: <应用策略名称>
      type: <应用策略类型>
      properties: <应用策略参数>
  workflow:
    - name: <工作流节点名称>
      type: <工作流节点类型>
      properties: <工作流节点参数>
```

无论待交付的组件是 Helm Chart 还是云数据库，目标基础设施是 Kubernetes 集群还是云平台，KubeVela 都通过 Application 这样一个统一的、上层的交付描述文件来同用户交互，不会泄露任何复杂的底层基础设施细节，真正做到让用户完全专注于应用研发和交付本身。

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1696949643908-578be5e7-d22b-4523-a6e1-6a63a579c9a3.png)

在实际使用时，用户通过上述 Application 对象来引用预置的组件、运维特征、应用策略、以及工作流节点模块，填写这些模块暴露的用户参数即可完成一次对应用交付的建模。

当然上面提到的几个类型背后都是由一组称为模块定义（Definitions）的可编程模块来提供具体功能。KubeVela 会像胶水一样基于 K8s API 定义基础设施定义的抽象并将不同的能力组合起来。

将定义的 OAM 模块和背后的 K8s CRD 控制器结合起来就可以形成 KubeVela 的 Addon 插件，社区已经有一个完善的且在不断扩大的插件市场，比如  **terraform** 插件提供了云资源的供给，**fluxcd** 插件提供了 GitOps 能力等等。我们可以自己根据需求开发插件，类似于 Helm 可以提供一个插件仓库来发现和分发插件。

### 2.4 KubeVela 架构
KubeVela 的整体架构如下所示：

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1696949644061-5722808a-29ad-40fd-a701-eb5011b7cc09.png)

**KubeVela 是一个的应用交付与管理控制平面**，它架在 Kubernetes 集群、云平台等基础设施之上，通过 OAM 来对组件、云服务、运维能力、交付工作流进行统一的编排和交付。KubeVela 这种与基础设施本身完全解耦的设计，很容易就能帮助你面向混合云/多云/多集群基础设施进行应用交付与管理。

而为了能够同任何 CI 流水线或者 GitOps 工具无缝集成，KubeVela 的 API 被设计为是声明式、完全以应用为中心的，它包括：

+ 帮助用户定义应用交付计划的 Application 对象
+ 帮助平台管理员通过 CUE 语言定义平台能力和抽象的 **X-Definition** 对象，比如 **ComponentDefinition**、**TraitDefinition** 等。

在具体实现上，KubeVela 依赖一个独立的 Kubernetes 集群来运行。具体来说，KubeVela 主要由如下几个部分组成:

+ **KubeVela 核心控制器**：为整个系统提供核心控制逻辑，完成诸如编排应用和工作流、修订版本快照、垃圾回收等等基础逻辑。
+ **Cluster Gateway 控制器**：提供统一的多集群访问接口和操作。
+ **插件体系**：注册和管理 KubeVela 的扩展功能，包括 CRD 控制器和相关模块定义。例如，下面列出了几个常用的插件：
    - **VelaUX** 插件是 KubeVela 的 Web UI。 此外，它在架构中更像是一个功能齐全的 “应用交付平台”，将业务逻辑耦合在起特定的 API 中，并为不了解 k8s 的业务开发者提供开箱即用的平台体验。
    - **Workflow** 插件是一个独立的工作流引擎，可以作为统一的 Pipeline 运行以部署多个应用程序或其他操作。与传统 Pipeline 相比，它主要使用 CUE 驱动基于 IaC 的 API，而不是每一步都运行容器（或 pod）。 它与 KubeVela 核心控制器的应用工作流使用相同的机制。
    - **Vela Prism** 插件是 KubeVela 的扩展 API 服务器，基于 Kubernetes Aggregated API 机制构建。它可以将诸如 Grafana 创建仪表盘等第三方服务 API 映射为 Kubernetes API，方便用户将第三方资源作为 Kubernetes 原生资源进行 IaC 化管理。
    - **Terraform** 插件允许用户使用 Terraform 通过 Kubernetes 自定义资源管理云资源。
    - 此外，KubeVela 有一个不断增长的插件市场，其中已经包含 50 多个用于集成的社区插件，包括 ArgoCD、FluxCD、Backstage、OpenKruise、Dapr、Crossplane、Terraform、OpenYurt 等等。
+ 如果你还没有任何 Kubernetes 集群，构建在 k3s 和 k3d 之上的 **VelaD** 工具可以帮助你一键启动所有这些东西。它将 KubeVela 与 Kubernetes 控制平面集成在一起，这对于构建开发/测试环境非常有帮助。

还有一个<font style="color:#DF2A3F;">非常重要的点是 </font>**<font style="color:#DF2A3F;">KubeVela 是可编程的</font>**。现实世界中的应用交付，往往是一个比较复杂的过程。哪怕是一个比较通用的交付流程，也会因为场景、环境、用户甚至团队的不同而千差万别。所以 KubeVela 从第一天起就采用了一种可编程式的方法来实现它的交付模型，这使得 KubeVela 可以以前所未有的灵活度适配到你的应用交付场景中。

## 3 KubeVela 安装
如果你没有 Kubernetes 环境，可以选择使用 VelaD 来独立安装 KubeVela。**<u><font style="color:#DF2A3F;">它是一个命令行工具，将 KubeVela 最小安装以及使用 VelaUX 的一切依赖打包为一个可执行文件，VelaD 会集成了 K3s 和 k3d 用于自动化管理 Kubernetes 集群。</font></u>**

我们这里当然选择基于先有的 Kubernetes 集群来安装 KubeVela。要求集群版本 **>= v1.19 && <= v1.26**。

### 3.1 安装 KubeVela CLI 命令行工具
首先需要安装 KubeVela 命令行工具，KubeVela CLI 提供了常用的集群和应用管理能力，直接使用下面的命令即可安装：

```bash
curl -fsSl https://kubevela.io/script/install.sh | bash
```

安装完成后，可以通过 **vela version** 命令查看版本信息：

```yaml
$ vela version
CLI Version: 1.9.6
Core Version:
GitRevision: git-9c57c098
GolangVersion: go1.19.12
```

然后我们可以使用如下命令来安装 KubeVela 控制平面：

```bash
vela install
```

安装完成后，会创建一个 **vela-system** 的命名空间，对应的 Pod 列表如下所示：

```bash
$ kubectl get pods -n vela-system
NAME                                                        READY   STATUS      RESTARTS   AGE
kubevela-cluster-gateway-b689d74dc-mtzrh                    1/1     Running     0          134m
kubevela-vela-core-85fd59d846-49q22                         1/1     Running     0          134m
kubevela-vela-core-admission-patch-8x9lv                    0/1     Completed   0          131m
kubevela-vela-core-cluster-gateway-tls-secret-patch-xjcw9   0/1     Completed   0          129m
```

### 3.2 使用 Helm 安装 KubeVela
当然如果你习惯使用 Helm，也可以通过如下 Helm 命令完成 VelaCore 的安装和升级：

```bash
$ helm repo add kubevela https://charts.kubevela.net/core
$ helm repo update
$ helm upgrade --install --create-namespace -n vela-system kubevela kubevela/vela-core --wait
```

上面的只是安装了 KubeVela 控制平面，我们一般情况下也会安装 VelaUX，它是 KubeVela 的 UI 控制台，可以通过浏览器访问它，当然你也可以不安装，这是可选的。

要安装也非常简单，只需要执行下面的命令启用 **velaux** 插件即可：

```bash
$ vela addon enable velaux
```

VelaUX 需要认证访问，默认的用户名是 **admin**，默认密码是 **VelaUX12345**。请务必在第一次登录之后重新设置和保管好你的新密码。

另外默认情况下，VelaUX 没有暴露任何端口。端口转发会以代理的方式允许你通过本地端口来访问 VelaUX 控制台。

```bash
vela port-forward addon-velaux -n vela-system
```

选择 **> local | velaux | velaux** 来启用端口转发。

VelaUX 控制台插件支持三种和 Kubernetes 服务一样的服务访问方式，它们是：**ClusterIP**、**NodePort** 以及 **LoadBalancer**，默认的服务访问方式为 ClusterIP。我们可以用下面的方式来改变 VelaUX 控制台的访问方式

```bash
vela addon enable velaux serviceType=LoadBalancer
# 或者
vela addon enable velaux serviceType=NodePort
```

一旦服务访问方式指定为 **LoadBalancer** 或者 **NodePort**，你可以通过执行 **vela status**来获取访问地址：

```bash
vela status addon-velaux -n vela-system --endpoint
```

期望得到的输出如下：

```bash
+----------------------------+----------------------+
|  REF(KIND/NAMESPACE/NAME)  |       ENDPOINT       |
+----------------------------+----------------------+
| Service/vela-system/velaux | http://<IP address> |
+----------------------------+----------------------+
```

如果你集群中拥有可用的 ingress 和域名，那么你可以按照下面的方式给你的 VelaUX 在部署过程中指定一个域名。

```bash
$ vela addon enable velaux domain=vela.k8s.local
Addon velaux enabled successfully.
Please access addon-velaux from the following endpoints:
+---------+---------------+-----------------------------------+--------------------------------+-------+
| CLUSTER |   COMPONENT   |     REF(KIND/NAMESPACE/NAME)      |            ENDPOINT            | INNER |
+---------+---------------+-----------------------------------+--------------------------------+-------+
| local   | velaux-server | Service/vela-system/velaux-server | velaux-server.vela-system:8000 | true  |
| local   | velaux-server | Ingress/vela-system/velaux-server | http://vela.k8s.local          | false |
+---------+---------------+-----------------------------------+--------------------------------+-------+
    To open the dashboard directly by port-forward:

    vela port-forward -n vela-system addon-velaux 8000:8000

    Please refer to https://kubevela.io/docs/reference/addons/velaux for more VelaUX addon installation and visiting method.
```

此外 VelaUX 支持 Kubernetes 和 MongoDB 作为其数据库。默认数据库为 Kubernetes，我们强烈建议你通过使用 MongoDB 来增强你的生产环境使用体验。

```bash
vela addon enable velaux dbType=mongodb dbURL=mongodb://<MONGODB_USER>:<MONGODB_PASSWORD>@<MONGODB_URL>
```

### 3.1 VelaUX
现在我们可以通过 `**http://vela.k8s.local**` 来访问 VelaUX 控制台了，第一次访问可以配置管理员账号信息：

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1696949644120-a7019bc4-a5d1-4a96-b99f-e242d29b63be.png)

VelaUX 是 KubeVela 的插件，它是一个企业可以开箱即用的云原生应用交付和管理平台。与此同时，也加入了一些企业使用中需要的概念。

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1696949644146-675a1f5c-940b-48c8-9ecc-befea943138e.png)

**项目（Project）**

项目作为在 KubeVela 平台组织人员和资源的业务承载，项目中可以设定成员、权限、应用和分配环境。在项目维度集成外部代码库、制品库，呈现完整 CI/CD Pipeline；集成外部需求管理平台，呈现项目需求管理；集成微服务治理，提供多环境业务联调和统一治理能力。项目提供了业务级的资源隔离能力。

默认情况下，VelaUX 会创建一个名为 **default** 的项目，你可以在 **项目管理** 中创建更多的项目。

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1696949644101-20c91d3d-e7c6-4f50-8659-08edb1267c80.png)

**环境（Environment）**

环境指通常意义的开发、测试、生产的环境业务描述，它可以包括多个交付目标。环境协调上层应用和底层基础设施的匹配关系，不同的环境对应管控集群的不同 Kubernetes Namespace。处在同一个环境中的应用才能具备内部互访和资源共享能力。

同样默认情况下，VelaUX 会创建一个名为 **default** 的环境，你可以在 **环境管理** 中创建更多的环境。

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1696949644682-a0e86670-6cdd-4dd6-83d7-0f810c8673a8.png)

应用可绑定多个环境进行发布，对于每一个环境可设置环境级部署差异。

**交付目标（Target）**

交付目标用于描述应用的相关资源实际部署的物理空间，对应 Kubernetes 集群或者云的区域（Region）和专有网络（VPC）。对于普通应用，组件渲染生成的资源会在交付目标指定的 Kubernetes 集群中创建（可以精确到指定集群的 Namespace）；对于云服务，资源创建时会根据交付目标中指定的云厂商的参数创建到对应的区域和专有网络中，然后将生成的云资源信息分发到交付目标指定的 Kubernetes 集群中。单个环境可关联多个交付目标，代表该环境需要多集群交付。单个交付目标只能对应一个环境。

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1696949644700-c80affca-1a69-41e2-be0d-8fa3d4721da3.png)

**应用（Application）**

应用是定义了一个微服务业务单元所包括的制品（二进制、Docker 镜像、Helm Chart...）或云服务的交付和管理需求，它由组件、运维特征、工作流、应用策略四部分组成，应用的生命周期操作包括：

+ **创建(Create)** 应用是创建元信息，并不会实际部署和运行资源。
+ **部署(Deploy)** 指执行指定的工作流， 将应用在某一个环境中完成实例化。
+ **回收(Recycle)** 删除应用部署到某一个环境的实例，回收其占用的资源。
+ **删除**应用会删除元数据，前提是应用实例已经完全被回收后才能删除。

VelaUX 应用中其他概念均与 KubeVela 控制器中的概念完全一致。

## 4 第一个 KubeVela 应用
上面我们已经安装好了 KubeVela，接下来我们就可以开始使用 KubeVela 来部署我们的第一个应用了。

下面我们定义了一个简单的 OAM 应用，它包括了一个无状态服务组件和运维特征，然后定义了三个部署策略和工作流步骤。此应用描述的含义是将一个服务部署到两个目标命名空间，并且在第一个目标部署完成后等待人工审核后部署到第二个目标，且在第二个目标时部署 2 个实例。

```yaml
# first-vela-app.yaml
apiVersion: core.oam.dev/v1beta1
kind: Application
metadata:
  name: first-vela-app
spec:
  components:
    - name: express-server
      type: webservice # webservice 是一个内置的组件类型
      properties: # 组件参数
        image: oamdev/hello-world
        ports:
          - port: 8000
            expose: true
      traits: # 组件运维特征
        - type: scaler
          properties:
            replicas: 1
  policies:
    - name: target-default
      type: topology # topology 是一个内置的应用策略类型，它可以将应用部署到多个目标
      properties:
        clusters: ["local"] # local 集群即 Kubevela 所在的集群
        namespace: "default"
    - name: target-prod
      type: topology
      properties:
        clusters: ["local"]
        namespace: "prod" # 此命名空间需要在应用部署前完成创建
    - name: deploy-ha
      type: override # override 是一个内置的应用策略类型，它可以覆盖组件的参数
      properties:
        components:
          - type: webservice
            traits:
              - type: scaler
                properties:
                  replicas: 2
  workflow: # 应用工作流
    steps:
      - name: deploy2default
        type: deploy # deploy 是一个内置的工作流类型，它可以将应用部署到指定的目标
        properties:
          policies: ["target-default"]
      - name: manual-approval
        type: suspend # suspend 是一个内置的工作流类型，它可以暂停工作流的执行
      - name: deploy2prod
        type: deploy
        properties:
          policies: ["target-prod", "deploy-ha"]
```

要先创建 **prod** 命名空间，可以使用 **vela env init** 命令，当然也可以直接使用 **kubectl create ns** 命令：

```bash
# 此命令用于在管控集群创建命名空间
vela env init prod --namespace prod
```

接下来就可以启动我们的第一个 KubeVela 应用了:

```bash
$ vela up -f first-vela-app.yaml
Applying an application in vela K8s object format...
✅ App has been deployed 🚀🚀🚀
    Port forward: vela port-forward first-vela-app -n prod
             SSH: vela exec first-vela-app -n prod
         Logging: vela logs first-vela-app -n prod
      App status: vela status first-vela-app -n prod
        Endpoint: vela status first-vela-app -n prod --endpoint
Application prod/first-vela-app applied.
```

**vela up** 命令会将上面定义的 Application 对象根据我们的描述翻译渲染成对应的 K8s 资源对象，部署完成后可以使用 vela 的相关命令来了解该应用的相关信息。

首先可以使用 **vela status** 命令来查看下应用的当前状态。由于上面应用定义的 Workflow 是先将应用部署到 local 集群的 default 命名空间中，然后进入第二个步骤的时候是一个 **suspend** 类型的工作流，所以正常情况下应用完成第一个目标部署后会进入暂停状态（左侧的 **workflowSuspending** 状态）。

```yaml
$ vela status first-vela-app -n prod
About:

  Name:         first-vela-app
  Namespace:    prod
  Created at:   2023-10-10 16:50:17 +0800 CST
  Status:       workflowSuspending

Workflow:

  mode: StepByStep-DAG
  finished: false
  Suspend: true
  Terminated: false
  Steps
  - id: kkotnerd76
    name: deploy2default
    type: deploy
    phase: succeeded
  - id: axtmf24jcx
    name: manual-approval
    type: suspend
    phase: suspending
    message: Suspended by field suspend

Services:

  - Name: express-server
    Cluster: local  Namespace: default
    Type: webservice
    Healthy Ready:1/1
    Traits:
      ✅ scaler
```

要继续工作流，则需要进行人工审核（左侧显示的第二个步骤），批准应用进入第二个目标部署，直接使用下面的命令即可：

```bash
vela workflow resume first-vela-app
```

当然在 VelaxUX 控制台中也可以看到应用的状态，也可以在控制台中直接进行人工审核操作。

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1696949644797-f73e1243-58e3-434d-8440-4bc731e4bc9a.png)

审批通过后会执行第三个步骤 **deploy2prod**，应用 **target-prod**、**deploy-ha** 这两个策略了。

经过上面的整个工作流过后，最终应用会在 default 命名空间下面创建一个 Pod，在 prod 命名空间下面创建两个副本的 Pod。

```plain
$ kubectl get pods -n prod
NAME                              READY   STATUS    RESTARTS   AGE
express-server-5447567596-jcpnh   1/1     Running   0          72s
express-server-5447567596-lgqdz   1/1     Running   0          72s
$ kubectl get pods
NAME                                     READY   STATUS    RESTARTS         AGE
express-server-5447567596-clbgb          1/1     Running   0                7m36s
```

在 VelaUX 控制台中也可以看到应用的状态：

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1696949644623-4b705257-90ff-40be-9fa2-3a592a8f91a8.png)

到这里就完成了我们的第一个 KubeVela 应用的部署流程。

