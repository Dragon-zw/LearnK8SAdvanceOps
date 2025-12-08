<font style="color:rgb(74, 74, 74);">KubeVela 打通了应用与基础设施之间的交付管控的壁垒，相较于原生的 Kubernetes 对象，KubeVela 的 Application 更好地简化抽象了开发者需要关心的配置，将复杂的基础设施能力及编排细节留给了平台工程师。而 KubeVela 的 </font>**<font style="color:rgb(10, 10, 10);">apiserver</font>**<font style="color:rgb(74, 74, 74);"> 则是进一步为开发者提供了使用 HTTP Request 直接操纵 Application 的途径，使得开发者即使没有 Kubernetes 的使用经验与集群访问权限也可以轻松部署自己的应用。</font>

<font style="color:rgb(74, 74, 74);">接下来我们就以 Jenkins 为基础，结合 KubeVela 来实现一个简单的应用持续交付的流程。</font>

<font style="color:rgb(74, 74, 74);">要实现一个简单的应用持续交付，我们需要做如下几件事情：</font>

+ <font style="color:rgb(1, 1, 1);">需要一个 git 仓库来存放应用程序代码、测试代码，以及描述 KubeVela Application 的 YAML 文件。</font>
+ <font style="color:rgb(1, 1, 1);">需要一个持续集成的工具帮你自动化完成程序代码的测试，并打包成镜像上传到仓库中。</font>
+ <font style="color:rgb(1, 1, 1);">需要在 Kubernetes 集群上安装 KubeVela 并启用 apiserver 功能。</font>

<font style="color:rgb(74, 74, 74);">我们这里的演示 Demo 采用 Github 作为 git 仓库，Jenkins 作为 CI 工具，DockerHub 作为镜像仓库。应用程序以一个简单的 Golang HTTP Server 为例，整个持续交付的流程如下。</font>

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1700062008189-b4531f65-8df6-4677-adb1-2bfec943058c.png)

<font style="color:rgb(74, 74, 74);">从整个流程可以看出开发者只需要关心应用的开发并使用 Git 进行代码版本的维护，即可自动走完测试流程并部署应用到 Kubernetes 集群中。</font>

<font style="color:rgb(74, 74, 74);">关于 Jenkins 在 Kubernetes 集群中的安装配置前面我们已经介绍过了，这里我们就不再赘述。</font>

#### <font style="color:black;">应用配置</font>
<font style="color:rgb(74, 74, 74);">这里我们采用了 Github 作为代码仓库，仓库地址为 https://github.com/cnych/KubeVela-demo-CICD-app，当然也可以根据各自的需求与喜好，使用其他代码仓库，如 Gitlab。为了 Jenkins 能够获取到 GitHub 中的更新，并将流水线的运行状态反馈回 GitHub，需要在 GitHub 中完成以下两步操作。</font>

1. <font style="color:rgb(1, 1, 1);">配置</font><font style="color:rgb(1, 1, 1);"> </font>**<font style="color:rgb(10, 10, 10);">Personal Access Token</font>**<font style="color:rgb(1, 1, 1);">。注意将</font><font style="color:rgb(1, 1, 1);"> </font>**<font style="color:rgb(10, 10, 10);">repo:status</font>**<font style="color:rgb(1, 1, 1);"> </font><font style="color:rgb(1, 1, 1);">勾选，以获得向 GitHub 推送 Commit 状态的权限，将生成的 Token 复制下来，下面会用到。</font>

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1700062008164-69abc5b7-2652-49dc-9d6e-f02ca5eba32e.png)

1. <font style="color:rgb(1, 1, 1);">然后在 Jenkins 的 Credential 中加入</font><font style="color:rgb(1, 1, 1);"> </font>**<font style="color:rgb(10, 10, 10);">Secret Text</font>**<font style="color:rgb(1, 1, 1);"> </font><font style="color:rgb(1, 1, 1);">类型的 Credential 并将上述的 GitHub 的 Personal Access Token 填入。</font>

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1700062008192-a271634e-6c92-458a-9f37-ea1b366645d9.png)<font style="color:rgb(136, 136, 136);"></font>

2. <font style="color:rgb(1, 1, 1);">接下来到 Jenkins 的</font><font style="color:rgb(1, 1, 1);"> </font>**<font style="color:rgb(10, 10, 10);">Dashboard > Manage Jenkins > Configure System > GitHub</font>**<font style="color:rgb(1, 1, 1);"> </font><font style="color:rgb(1, 1, 1);">中点击</font><font style="color:rgb(1, 1, 1);"> </font>**<font style="color:rgb(10, 10, 10);">Add GitHub Server</font>**<font style="color:rgb(1, 1, 1);"> </font><font style="color:rgb(1, 1, 1);">并将刚才创建的 Credential 填入。完成后可以点击 Test connection 来验证配置是否正确。</font>

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1700062008189-23ee0a46-4568-406a-a09b-a9ffdea5e0a9.png)

3. <font style="color:rgb(1, 1, 1);">由于我们这里的 Jenkins 位于本地环境，要让 GitHub 通过 Webhook 来触发 Jenkins，我们需要提供一个可访问的地址，这里我们可以使用</font><font style="color:rgb(1, 1, 1);"> </font>**<font style="color:rgb(10, 10, 10);">ngrok</font>**<font style="color:rgb(1, 1, 1);"> </font><font style="color:rgb(1, 1, 1);">来实现，首先前往 https://dashboard.ngrok.com 注册一个账号，将</font><font style="color:rgb(1, 1, 1);"> </font>**<font style="color:rgb(10, 10, 10);">Authtoken</font>**<font style="color:rgb(1, 1, 1);"> </font><font style="color:rgb(1, 1, 1);">和</font><font style="color:rgb(1, 1, 1);"> </font>**<font style="color:rgb(10, 10, 10);">APIKEY</font>**<font style="color:rgb(1, 1, 1);"> </font><font style="color:rgb(1, 1, 1);">记录下来。</font>

```bash
export NGROK_AUTHTOKEN=<your-ngrok-authtoken>
export NGROK_API_KEY=<your-ngrok-apikey>
```

<font style="color:rgb(74, 74, 74);">然后我们可以在本地 Kubernetes 集群中安装 ngrok ingress controller：</font>

```bash
helm repo add ngrok https://ngrok.github.io/kubernetes-ingress-controller
# 使用下面命令安装 ngrok ingress controller
helm install ngrok-ingress-controller ngrok/kubernetes-ingress-controller \
--namespace ngrok-ingress-controller \
--create-namespace \
--set credentials.apiKey=$NGROK_API_KEY \
--set credentials.authtoken=$NGROK_AUTHTOKEN
```

<font style="color:rgb(74, 74, 74);">安装完成后为 Jenkins 创建一个 ngrok 的 ingress 路由：</font>

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
 name: jenkins-ngrok
 namespace: kube-ops
spec:
 ingressClassName: ngrok
 rules:
   - host: prompt-adjusted-sculpin.ngrok-free.app
     http:
       paths:
         - backend:
             service:
               name: jenkins
               port:
                 name: web
           path: /
           pathType: Prefix
```

<font style="color:rgb(74, 74, 74);">上面的 host 域名是 ngrok 为我们分配的，你可以在 ngrok 的控制台中手动创建，应用上面的 ingress 对象后我们就可以通过 ngrok 为我们分配的域名来访问 Jenkins 了。</font>

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1700062008164-17c4332e-ba64-4f5d-b0d2-e29e1c736643.png)

4. <font style="color:rgb(1, 1, 1);">接下来我们就可以在 GitHub 的代码仓库的设定里添加 Webhook，将 Jenkins 的地址对应的 Webhook 地址填入</font><font style="color:rgb(1, 1, 1);"> </font>**<font style="color:rgb(10, 10, 10);"><ngrok domain>/github-webhook/</font>**<font style="color:rgb(1, 1, 1);">，这样该代码仓库的所有 Push 事件推送到 Jenkins 中。</font>

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1700062008631-448139a3-5587-4afd-9125-6fae64bf1ffc.png)

#### <font style="color:black;">编写应用</font>
<font style="color:rgb(74, 74, 74);">我们这里采用的应用是一个基于 Golang 语言编写的简单的 HTTP Server。在代码中，声明了一个名叫</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">VERSION</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">的常量，并在访问该服务时打印出来。同时还附带一个简单的测试，用来校验</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">VERSION</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">的格式是否符合标准。</font>

```go
// main.go
package main

import (
    "fmt"
    "net/http"
)

const VERSION = "0.1.0-v1alpha1"

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        _, _ = fmt.Fprintf(w, "Version: %s\n", VERSION)
    })
    if err := http.ListenAndServe(":8088", nil); err != nil {
        println(err.Error())
    }
}
```

<font style="color:rgb(74, 74, 74);">测试代码如下所示：</font>

```go
// main_test.go

package main

import (
    "regexp"
    "testing"
)

const verRegex string = `^v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?` +
`(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?` +
`(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?$`

func TestVersion(t *testing.T) {
    if ok, _ := regexp.MatchString(verRegex, VERSION); !ok {
        t.Fatalf("invalid version: %s", VERSION)
    }
}
```

<font style="color:rgb(74, 74, 74);">在应用交付时需要将 Golang 服务打包成镜像并以 KubeVela Application 的形式发布到 Kubernetes 集群中，因此在代码仓库中还包含</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">Dockerfile</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">文件，用来描述镜像的打包方式。</font>

```go
# Dockerfile
FROM golang:1.13-rc-alpine3.10 as builder
WORKDIR /app
COPY main.go .
RUN go build -o kubevela-demo-cicd-app main.go

FROM alpine:3.10
WORKDIR /app
COPY --from=builder /app/kubevela-demo-cicd-app /app/kubevela-demo-cicd-app
ENTRYPOINT ./kubevela-demo-cicd-app
EXPOSE 8088
```

#### <font style="color:black;">配置 CI 流水线</font>
<font style="color:rgb(74, 74, 74);">在这里我们将包含两条流水线，一条是用来进行测试的流水线 (对应用代码运行测试) ，一条是交付流水线 (将应用代码打包上传镜像仓库，同时更新目标环境中的应用，实现自动更新) 。</font>

**<font style="color:rgb(10, 10, 10);">测试流水线</font>**

<font style="color:rgb(74, 74, 74);">在 Jenkins 中创建一条新的名为</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">KubeVela-demo-CICD-app-test</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">的流水线：</font>

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1700062008652-06b18ee3-b684-48a5-87a6-be69555c2e3d.png)

<font style="color:rgb(74, 74, 74);">然后配置构建触发器为</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">GitHub hook trigger for GITScm polling</font>**<font style="color:rgb(74, 74, 74);">:</font>

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1700062008746-986d470a-23d7-4517-b907-615de0b16ea7.png)

<font style="color:rgb(74, 74, 74);">在这条流水线中，首先是采用了 golang 的镜像作为执行环境，方便后续运行测试。然后将分支配置为 GitHub 仓库中的 dev 分支，代表该条流水线被 Push 事件触发后会拉取 dev 分支上的内容并执行测试，测试结束后将流水线的状态回写至 GitHub 中。这里我们使用的是基于 Kubernetes 的动态 Slave Agent，因此在流水线中需要配置 Kubernetes 的相关信息，包括 Kubernetes 的地址、Service Account 等。</font>

```go
void setBuildStatus(String message, String state) {
    step([
        $class: "GitHubCommitStatusSetter",
        reposSource: [$class: "ManuallyEnteredRepositorySource", url: "https://github.com/cnych/KubeVela-demo-CICD-app"],
        contextSource: [$class: "ManuallyEnteredCommitContextSource", context: "ci/jenkins/test-status"],
        errorHandlers: [[$class: "ChangingBuildStatusErrorHandler", result: "UNSTABLE"]],
        statusResultSource: [ $class: "ConditionalStatusResultSource", results: [[$class: "AnyBuildResult", message: message, state: state]] ]
    ]);
}

pipeline {
    agent {
        kubernetes {
            cloud 'Kubernetes'
            containerTemplate {
                name 'golang'
                image 'golang:1.13-rc-alpine3.10'
                command 'cat'
                ttyEnabled true
            }
            serviceAccount 'jenkins'
        }
    }

    stages {
        stage('Prepare') {
            steps {
                script {
                    def checkout = git branch: 'dev', url: 'https://github.com/cnych/KubeVela-demo-CICD-app.git'
                    env.GIT_COMMIT = checkout.GIT_COMMIT
                    env.GIT_BRANCH = checkout.GIT_BRANCH
                    echo "env.GIT_BRANCH=${env.GIT_BRANCH},env.GIT_COMMIT=${env.GIT_COMMIT}"
                }
                setBuildStatus("Test running", "PENDING");
            }
        }
        stage('Test') {
            steps {
                container('golang') {
                    sh 'CGO_ENABLED=0 GOCACHE=$(pwd)/.cache go test *.go'
                }
            }
        }
    }

    post {
        success {
            setBuildStatus("Test success", "SUCCESS");
        }
        failure {
            setBuildStatus("Test failed", "FAILURE");
        }
    }
}
```

<font style="color:rgb(74, 74, 74);">我们可以使用上面的代码来执行流水线：</font>

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1700062008748-7c8410ea-309c-43e1-a504-c76f9e774f5d.png)

**<font style="color:rgb(10, 10, 10);">部署流水线</font>**

<font style="color:rgb(74, 74, 74);">类似测试流水线创建一个名为</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">KubeVela-demo-CICD-app-deploy</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">的部署流水线，首先将代码仓库中的分支拉取下来，区别是这里采用 prod 分支。然后使用 Docker 进行镜像构建并推送至远端镜像仓库。构建成功后，再将 Application 对应的 YAML 文件转换为 JSON 文件并注入</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">GIT_COMMIT</font>**<font style="color:rgb(74, 74, 74);">，最后向 KubeVela apiserver 发送请求进行创建或更新。</font>

<font style="color:rgb(74, 74, 74);">首先我们需要通过 VelaUX 来创建一个应用，这里我们创建一个名为</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">kubevela-demo-app</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">的应用，包含一个名为</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">kubevela-demo-app-web</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">的组件，组件类型为</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">webservice</font>**<font style="color:rgb(74, 74, 74);">，并将组件的镜像设置为</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">cnych/kubevela-demo-cicd-app</font>**<font style="color:rgb(74, 74, 74);">，如下图所示：</font>

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1700062008766-d1a5795f-4bfb-450c-bbe3-4ac2b5500c00.png)

<font style="color:rgb(74, 74, 74);">在应用面板上，我们可以找到一个默认的触发器，点击</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">手动触发</font>**<font style="color:rgb(74, 74, 74);">，我们可以看到</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">Webhook URL</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">和</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">Curl Command</font>**<font style="color:rgb(74, 74, 74);">，我们可以在 Jenkins 的流水线中使用任意一个。</font>

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1700062009094-e8028cfd-c836-4f52-a4f8-84b1903164e2.png)

<font style="color:rgb(74, 74, 74);">Webhook URL 是这个触发器的触发地址，在 Curl Command 里，还提供了手动 Curl 该触发器的请求示例。我们来详细解析一下请求体：</font>

```go
{
    // 必填，此次触发的更新信息
    "upgrade": {
        // Key 为应用的名称
        "<application-name>": {
            // 需要更新的值，这里的内容会被 Patch 更新到应用上
            "image": "<image-name>"
        }
    },
    // 可选，此次触发携带的代码信息
    "codeInfo": {
        "commit": "<commit-id>",
        "branch": "<branch>",
        "user": "<user>"
    }
}
```

**<font style="color:rgb(10, 10, 10);">upgrade</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">下是本次触发要携带的更新信息，在应用名下，是需要被 Patch 更新的值。默认推荐的是更新镜像 image，也可以扩展这里的字段来更新应用的其他属性。</font>**<font style="color:rgb(10, 10, 10);">codeInfo</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">中是代码信息，可以选择性地携带，比如 commit ID、分支、提交者等，一般这些值可以通过在 CI 系统中使用变量替换来指定。</font>

<font style="color:rgb(74, 74, 74);">然后我们可以是部署流水线中使用上面的触发器来部署应用，的代码如下所示：</font>

```go
void setBuildStatus(String message, String state) {
    step([
        $class: "GitHubCommitStatusSetter",
        reposSource: [$class: "ManuallyEnteredRepositorySource", url: "https://github.com/cnych/KubeVela-demo-CICD-app"],
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
            image: golang:1.13-rc-alpine3.10
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
                    def checkout = git branch: 'prod', url: 'https://github.com/cnych/KubeVela-demo-CICD-app.git'
                    env.GIT_COMMIT = checkout.GIT_COMMIT
                    env.GIT_BRANCH = checkout.GIT_BRANCH
                    echo "env.GIT_BRANCH=${env.GIT_BRANCH},env.GIT_COMMIT=${env.GIT_COMMIT}"
                    setBuildStatus("Deploy running", "PENDING");
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
                        docker build -t cnych/kubevela-demo-cicd-app .
                        docker push cnych/kubevela-demo-cicd-app
                        """
                    }
                }
            }
        }
        stage('Deploy') {
            steps {
                sh '''#!/bin/bash
                set -ex
                curl -X POST -H 'content-type: application/json' --url http://vela.k8s.local/api/v1/webhook/x0i7t8jdsz2uvime -d '{"action":"execute","upgrade":{"kubevela-demo-app":{"image":"cnych/kubevela-demo-cicd-app"}},"codeInfo":{"commit":"","branch":"","user":""}}'
                '''
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

**<font style="color:rgb(10, 10, 10);">测试效果</font>**

<font style="color:rgb(74, 74, 74);">在完成上述的配置流程后，持续交付的流程便已经搭建完成。我们可以来检验一下它的效果。</font>

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1700062009168-a64305ef-b744-46a9-852f-ed3b49aebcfd.png)

<font style="color:rgb(74, 74, 74);">我们首先将</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">main.go</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">中的</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">VERSION</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">字段修改为</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">Bad Version Number</font>**<font style="color:rgb(74, 74, 74);">，即</font>

```plain
const VERSION = "Bad Version Number"
```

<font style="color:rgb(74, 74, 74);">然后提交该修改至</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">dev</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">分支，我们可以看到 Jenkins 上的测试流水线被触发运行，失败后将该状态回写给 GitHub。</font>

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1700062009200-94085117-c875-42ed-9558-1c9b37a716d4.png)

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1700062009331-cae70590-7433-457a-abc9-f4acb7eec391.png)

<font style="color:rgb(74, 74, 74);">我们重新将</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">VERSION</font>**<font style="color:rgb(74, 74, 74);"> </font><font style="color:rgb(74, 74, 74);">修改为</font><font style="color:rgb(74, 74, 74);"> </font>**<font style="color:rgb(10, 10, 10);">0.1.1</font>**<font style="color:rgb(74, 74, 74);">，然后再次提交。可以看到这一次测试流水线成功完成执行，并在 GitHub 对应的 Commit 上看到了成功的标志。</font>

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1700062009231-518ef81e-364a-4485-b9e6-1443073c4155.png)

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1700062009516-378e9d78-6d23-44b9-beeb-646197c0d374.png)

<font style="color:rgb(74, 74, 74);">接下来我们在 GitHub 上提交 Pull Request 尝试将 dev 分支上的更新合并至 prod 分支上。</font>

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1700062009590-e0e8a1e1-4342-4c64-9a17-d09f6015caaa.png)

<font style="color:rgb(74, 74, 74);">可以看到在 Jenkins 的部署流水线成功运行结束后，GitHub 上 prod 分支最新的 Commit 也显示了成功的标志。</font>

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1700062009614-ac31f45d-8d70-4a5b-890f-f42a46adfe47.png)

![](https://cdn.nlark.com/yuque/0/2023/png/2555283/1700062009666-6713bbb1-7e1c-4ab9-a25c-4ba6a5554392.png)

<font style="color:rgb(74, 74, 74);">我们的应用已经成功部署了，当前 Deployment 的副本数是 3，并且还有一个 Ingress 对象，这时我们可以访问 Ingress 所配置的域名，成功显示了当前的版本号。</font>

```bash
$ vela ls
APP                     COMPONENT               TYPE            TRAITS          PHASE   HEALTHY STATUS          CREATED-TIME
kubevela-demo-app       kubevela-demo-app       webservice      scaler,gateway  running healthy Ready:3/3       2023-10-14 19:11:59 +0800 CST
$ kubectl get pods
NAME                                     READY   STATUS    RESTARTS       AGE
kubevela-demo-app-675896596f-87kxl       1/1     Running   0              9m39s
kubevela-demo-app-675896596f-q5pvz       1/1     Running   0              9m39s
kubevela-demo-app-675896596f-v895m       1/1     Running   0              44m
$ kubectl get ingress
NAME                CLASS   HOSTS                              ADDRESS   PORTS   AGE
kubevela-demo-app   nginx   kubevela-demo-cicd-app.k8s.local             80      10m
$ curl -H "Host: kubevela-demo-cicd-app.k8s.local" http://<ingress controller address>
Version: 0.1.1
```

<font style="color:rgb(74, 74, 74);">如果想实现金丝雀发布，则可以使用上节的 kruise rollout 来实现，至此，我们便已经成功实现了一整套持续交付流程。在这个流程中，应用的开发者借助 KubeVela + Jenkins 的能力，可以轻松完成应用的迭代更新、集成测试、自动发布与滚动升级，而整个流程在各个环节也可以按照开发者的喜好和条件选择不同的工具，比如使用 Gitlab 替代 GitHub，或是使用 TravisCI 替代 Jenkins。</font>

> 参考文档：[https://kubevela.io/docs/tutorials/jenkins/](https://kubevela.io/docs/tutorials/jenkins/)
>

