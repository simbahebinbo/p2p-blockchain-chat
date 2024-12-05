启动一个节点 nodeA

```shell
./chat-blockchain -s 9999 -n nodeA
```

启动另一个节点 nodeB

```shell
./chat-blockchain -s 8888 -k 9999 -n nodeB
```

能够单向从 nodeB 向 nodeA 发消息
