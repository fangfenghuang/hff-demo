
http://liubin.org/kata-dev-book/src/your-first-pr.html

# 如何找到 issue
- Kata Containers 的 issue 页面查找自己刚兴趣的去 fix
- 自己使用中发现的 bug
- 自己需要的功能
- 学习、阅读源代码时候发现的问题

# Signed-off-by 信息
```bash
runtime: use concrete KataAgentConfig instead of interface type

Kata Containers 2.0 only have one type of agent, so there is no
need to use interface as config's type

Fixes: #1610

Signed-off-by: bin <bin@hyper.sh>
```
- 标题:这部分由 subsystem: commit summary 组成，根据具体的 pr 内容填写
- body 部分:和上面的标题隔开一行，记录具体的修改内容，用于对标题进行补充说明。可以是问题背景、如何修改、注意事项、以及其他参考资料等。
- fix 的 issue 编号:这部分以 Fixes: 开始，后面跟 issue 编号，如果有多个 issue， issue 之间用逗号分隔。
- Signed-off-by：这部分通过 commit -s 自动填写。

# 合并条件：

所有必须的 CI 测试通过
有 2 个或 2 个以上的项目维护者的 approve 。

# 代码修改
不管是 reivew 人员指出的错误，还是在 CI 测试中发现的错误，都需要开发者在本地重新修改代码。

如果只有一个提交，可以直接使用 git commit --amend 和 git push -f 即可。如果是一个 pr 多个提交，则可能需要找到相应的提交进行 --amend 处理，本文中不对这样的例子进行说明。

关于提交粒度：对于比较大的 pr ，建议以单个功能点、单个顶层文件（src/runtime, src/agent, docs）分开提交，同时需要确保每个 commit 都能编译通过。

# 一些小技巧
很多项目维护者并不是专注于某一 issue 或 pr ，所以必要的时候，可能需要 pr 提交者主动“催”一下，可以在 pr 的 comment 栏里使用 @ 来提醒你想要联系的项目维护者。
