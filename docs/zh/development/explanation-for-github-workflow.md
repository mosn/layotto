本文档解释了Layotto的github工作流的每个组件。

注意：Layotto的github工作流的配置文件在[这里](https://github.com/mosn/layotto/tree/main/.github/workflows)

## 1. 定时任务
### stale bot 
![img_1.png](../../img/development/workflow/img_1.png)

我们将会[关闭超时的Issues](https://github.com/marketplace/actions/close-stale-issues)

如果一个issue或PR在过去30天内没有最近的活动，它将自动标记为陈旧的。它将在7天后关闭，除非它被特殊 label 标记 (pinned，security，good first issue或help wanted) 或有其他活动发生。

**如果这个社区任务issue被这个自动程序关闭，那么这个任务将会被分配给其他人。**

合并于 https://github.com/mosn/layotto/pull/246

## 2. 持续集成/持续交付
![img.png](../../img/development/workflow/img.png)
### 2.1. Chore
#### <1> cla bot	

检查贡献者是否签署了贡献许可协议

#### TODO: 在修改proto文件时自动生成新的 API 文档

目前[我们必须手动完成](https://mosn.io/layotto/#/en/api_reference/how_to_generate_api_doc)。

生成的文件在[这里](https://github.com/mosn/layotto/blob/main/docs/en/api_reference/runtime_v1.md)或者[这里](https://github.com/mosn/layotto/blob/main/docs/en/api_reference/appcallback_v1.md)

### 2.2. Test	
#### <5> 执行单元测试
#### <5> 检查你是否完成了`go fmt`
#### <2><3> 确保单元测试的覆盖率不会下降

具体请查看 https://docs.codecov.com/docs/commit-status#branches

#### TODO: 集成测试


### 2.3. Lint	
#### <4> License checker		
我们使用https://github.com/marketplace/actions/license-eye

合并于https://github.com/mosn/layotto/pull/247

##### 如何为所有文件自动添加许可证头文件

在Layotto目录下执行:

```shell
docker run -it --rm -v $(pwd):/github/workspace apache/skywalking-eyes header fix
```

它将递归地为代码文件添加许可头。

##### 如何配置许可证检查器忽略指定类型的文件

忽略检查列表在 `.licenserc.yaml`中，你可以添加新的类型进去。

##### 有关此工具的更多详细信息
请查看 https://github.com/marketplace/actions/license-eye#docker-image 来获取详细信息。

#### PR title lint	
为了规范PR标题，我们添加了这个检查操作。你可以在https://github.com/thehanimo/pr-title-checker获取更多详细信息。

##### 配置pr标题检查
对于一个 pull request 如果标题在 `prefixes` 或 `regexp` 中，则检查通过。否则，一个标签`title needs formatting`将被添加到那个pull request中。

`regexpFlags` 意味着正则表达式, 例如 : `i`(Case-insensitive search) `g`(Global search) .

```
"CHECKS": {
"prefixes": ["fix: ", "feat: ","doc: "], 
"regexp": "docs\\(v[0-9]\\): ",
"regexpFlags": "i",
"ignoreLabels" : ["dont-check-PRs-with-this-label", "meta"]
}
```
~~#### TODO: PR body lint?~~

#### TODO: Code style lint
举个例子，找出 `go xxx()` 没有 `recover`

我们可以使用 go lint, 参考MOSN的配置

####  ~~- Commit message lint~~ (reverted)
具体请查看 https://github.com/mosn/layotto/issues/243