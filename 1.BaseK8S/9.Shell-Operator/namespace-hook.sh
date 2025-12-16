#!/bin/bash
if [[ $1 == "--config" ]]; then
# configuration
# 创建 NameSpace 的时候会去触发 shell-operator 钩子函数
cat <<EOF
configVersion: v1
kubernetes:
- apiVersion: v1
  kind: Namespace
  executeHookOnEvent: ["Added"]
EOF
else
  # BINDING_CONTEXT_PATH 事件创建时候的上下文
  # 在创建命名空间的时候，会打印命名空间的描述信息
  binding=$(jq -r .[0].binding $BINDING_CONTEXT_PATH)
  type=$(jq -r .[0].type $BINDING_CONTEXT_PATH)
  watchEvent=$(jq -r .[0].watchEvent $BINDING_CONTEXT_PATH)
  createdNamespace=$(jq -r .[0].object.metadata.name $BINDING_CONTEXT_PATH)
  echo "binding '${binding}', type '${type}',watchEvent '${watchEvent}'"
  echo "namespace '${createdNamespace}' added"
fi