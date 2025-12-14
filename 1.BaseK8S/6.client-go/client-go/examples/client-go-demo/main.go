package main

import (
	"context"
	"fmt"
	"os"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"
)

func initClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	// 通过 config 来创建 clientset 客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		//panic(err.Error())
		return nil, fmt.Errorf("创建客户端失败：%v", err)
	}
	return clientset, nil
}

// 创建 Pod
func createPod(clientset *kubernetes.Clientset) error {
	// 需要找到对应 Resource 的 apiGroup 组
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-pod",
			Namespace: "default",
			Labels: map[string]string{
				"app":     "nginx",
				"env":     "test",
				"version": "v1",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "nginx",
					Image: "swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/nginx:1.29.3-linuxarm64",
					Ports: []corev1.ContainerPort{
						{
							Name:          "http",
							ContainerPort: 80,
							Protocol:      "TCP",
						},
					},
				},
			},
		},
	}
	// 相当于执行了
	// kubectl run nginx-pod --image=swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/nginx:1.29.3-linuxarm64 \
	// --labels app=nginx,env=test,version=v1
	ctx := context.Background()
	_, err := clientset.CoreV1().Pods("default").Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("创建 Pod 失败：%v", err)
	}

	fmt.Println("Pod 创建成功")
	return nil
}

// 获取 Pod
func getPod(clientset *kubernetes.Clientset) error {
	ctx := context.Background()

	pod, err := clientset.CoreV1().Pods("default").Get(ctx, "nginx-pod", metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("获取 Pod 失败: %v", err)
	}

	fmt.Printf("Pod 名称: %s\n", pod.Name)
	fmt.Printf("Pod 命名空间: %s\n", pod.Namespace)
	fmt.Printf("Pod 状态: %s\n", pod.Status.Phase)
	fmt.Printf("Pod IP: %s\n", pod.Status.PodIP)
	fmt.Printf("Pod 创建时间: %s\n", pod.CreationTimestamp)

	// 打印 Pod Labels
	if len(pod.Labels) > 0 {
		fmt.Println("Pod 标签:")
		for key, value := range pod.Labels {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	// 打印容器信息
	fmt.Println("容器信息: ")
	for _, container := range pod.Spec.Containers {
		fmt.Printf("  名称: %s, 镜像: %s\n", container.Name, container.Image)
	}
	return nil
}

// 更新 Pod（更新标签）
func updatePodLabel(clientset *kubernetes.Clientset) error {
	ctx := context.Background()

	// 使用 retry 机制处理并发更新
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// 获取最新的 Pod
		pod, err := clientset.CoreV1().Pods("default").Get(ctx, "nginx-pod", metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("获取 Pod 失败: %v", err)
		}

		// 修改标签
		if pod.Labels == nil {
			pod.Labels = make(map[string]string)
		}
		// 添加标签
		pod.Labels["update"] = "true"
		pod.Labels["time"] = fmt.Sprintf("%d", time.Now().Unix())

		// 更新 Pod
		_, updateErr := clientset.CoreV1().Pods("default").Update(ctx, pod, metav1.UpdateOptions{})
		return updateErr
	})

	if retryErr != nil {
		return fmt.Errorf("更新 Pod 失败: %v", retryErr)
	}

	return nil
}

// 补丁 Patch Pod
func patchPod(clientset *kubernetes.Clientset) error {
	patchData := []byte(`{"spec":{"containers":[{"name":"nginx","image":"swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/nginx:latest-linuxarm64"}]}}`)

	patchedPod, err := clientset.CoreV1().Pods("default").Patch(context.TODO(), "nginx-pod", types.StrategicMergePatchType, patchData, metav1.PatchOptions{})
	if err != nil {
		return fmt.Errorf("❎Patch 更新 Pod 失败: %v", err)
	} else {
		fmt.Println("✅Patch Pod 更新成功")
		fmt.Printf("Patch Pod 名称: %s, Patch Pod 镜像: %s\n", patchedPod.Name, patchedPod.Spec.Containers[0].Image)
	}

	return nil
}

// 删除 Pod
func deletePod(clientset *kubernetes.Clientset) error {
	ctx := context.Background()

	// 立刻删除
	deletePolicy := metav1.DeletePropagationForeground

	err := clientset.CoreV1().Pods("default").Delete(ctx, "nginx-pod", metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})

	if err != nil {
		return fmt.Errorf("删除 Pod 失败: %v", err)
	}

	fmt.Println("Pod 删除成功")
	return nil
}

// 列出 Pod
func listPods(clientset *kubernetes.Clientset) {
	namespace := "default"

	// 列出所有的 Pod
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("❎获取 Pod 列表失败: %v", err)
		return
	}

	fmt.Printf("✅Total Pods: %d\n", len(pods.Items))

	// 显示详细信息
	fmt.Println("\nDetails Pod List:")
	// 遍历 Pods 数组
	for i, pod := range pods.Items {
		fmt.Printf("Pod %d: %s\n", i+1, pod.Name)
		fmt.Printf("		Status: %s\n", pod.Status.Phase)
		fmt.Printf("		Node: %s\n", pod.Spec.NodeName)
		fmt.Printf("		Age: %v\n", time.Since(pod.CreationTimestamp.Time).Round(time.Second))

		for _, container := range pod.Spec.Containers {
			fmt.Printf("		Container: %s, Image: %s\n", container.Name, container.Image)
		}
		fmt.Println()
	}
}

// 监听 Pod
func watchPod(clientset *kubernetes.Clientset) {
	namespace := "default"

	fmt.Println("👁 开始进行监听 Pods")
	fmt.Println("Please Ctrl+C to stop")

	watcher, err := clientset.CoreV1().Pods(namespace).Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("监听 Pods 失败: %v", err)
		return
	}
	defer watcher.Stop()

	// 设置超时时间
	timeout := time.After(5 * time.Minute)

	for {
		select {
		case <-timeout:
			fmt.Println("⏰ 监听超时，已停止监听")
			return
			// 通过隧道的方式 Channel 来传递信息
		case event, ok := <-watcher.ResultChan():
			if !ok {
				fmt.Println("📮 监听已关闭")
				return
			}
			// 将 Object 转换为 Pod 对象
			pod, ok := event.Object.(*corev1.Pod)
			if !ok {
				continue
			}

			switch event.Type {
			case watch.Added:
				fmt.Printf("➕ Pod %s 已创建 (Phase: %s)\n", pod.Name, pod.Status.Phase)
			case watch.Modified:
				fmt.Printf("📋 Pod %s 已修改 (Phase: %s)\n", pod.Name, pod.Status.Phase)
			case watch.Deleted:
				fmt.Printf("➖ Pod %s 已删除\n", pod.Name)
			case watch.Bookmark:
				fmt.Printf("📋 Pod %s 已被书签\n", pod.Name)
			case watch.Error:
				fmt.Printf("❎ Pod %s 发生错误\n", event.Object)
			}
		}
	}
}

// 用户传递一个参数，根据参数决定执行什么操作
func main() {
	// 从参数中获取操作类型
	if len(os.Args) < 2 {
		panic("请指定操作类型(create/get/update/patch/delete/list/watch)")
		os.Exit(1)
	}
	operation := os.Args[1]
	clientset, err := initClient()
	if err != nil {
		panic(err)
	}
	switch operation {
	case "create":
		// 创建 Pod
		err = createPod(clientset)
	case "get":
		// 获取 Pod
		err = getPod(clientset)
	case "update":
		// 更新 Pod
		err = updatePodLabel(clientset)
	case "patch":
		//  patch Pod
		err = patchPod(clientset)
	case "delete":
		// 删除 Pod
		err = deletePod(clientset)
	case "list":
		// 列出所有 Pod
		listPods(clientset)
	case "watch":
		// 监听 Pod
		watchPod(clientset)
	default:
		fmt.Printf("未知操作类型: %s\n", operation)
		os.Exit(1)
	}
	if err != nil {
		fmt.Printf("操作失败: %v\n", err)
		os.Exit(1)
	}
}
