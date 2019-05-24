# 周报维护方法文档

周报系统的运作主要与三个仓库有关：

- [dyweb/weekly][]，这是周报静态网站的仓库，这里存放了所有往期的周报。
- [dyweb/dy-bot][]，这是周报 GitHub bot 代码仓库。
- [dyweb/dy-weekly-generator][]，这是根据 [dyweb/weekly][] 的 Issue 生成周报内容的命令行工具所在的仓库。

周报的具体工作流程如下：

1. （需要人工干预）在每周的周一，维护者关于带有 working 标签的周报 Issue。
1. （自动化）GitHub bot dy-bot 通过 webhook 得到这一 Issue 的关闭事件，自动地把 working 标签删掉，并且创建出新的一期周报 Issue，然后给新的 Issue 打上 working 标签。
1. （需要人工干预）维护者利用 [dyweb/dy-weekly-generator][]，生成周报的 Markdown 文件，并且放在 [dyweb/weekly][] 对应的目录下。
1. （需要人工干预）维护者利用 [dyweb/weekly][] 中的脚本，生成新的静态页面，并推送到 GitHub 上

## 示例

以某个 Issue [dyweb/weekly#108](https://github.com/dyweb/weekly/issues/108) 为例，我们来讲解一下整体工作流程。

首先可以看到，在 Issue 被关闭后，GitHub user [@gaocegege-bot](https://github.com/gaocegege-bot) 把 Issue 的 working 标签去掉了。[@gaocegege-bot](https://github.com/gaocegege-bot) 背后运行的就是 [dyweb/dy-bot][]，它在 [https://github.com/dyweb/weekly/settings/hooks](https://github.com/dyweb/weekly/settings/hooks) 中配置了一个 hook。这一 hook 只针对 Issue 的 close 事件返回 200OK，其他请求都会返回 500，response body 为 `unknown event type %event_type%`。因此如果在 hook 页面看到不少 event 的返回都是 500，请不要惊慌，这是 feature。

接下来，在关闭后，我们可以利用 [dyweb/dy-weekly-generator][] 生成周报文档：

```bash
$ cd $PATH_TO_WEEKLY_DIR
$ weekly-gen --repo dyweb/weekly --issue $CURRENT_ISSUE_NUM > ./$YEAR/$YEAR-$MONTH-$DAY-weekly.md
# 在本例中，命令为
$ weekly-gen --repo dyweb/weekly --issue 108 > ./2019/2019-05-13-weekly.md
```

:tada: 感谢 [@codeworm96][] 为我们实现了周报生成器。

随后，新的周报就会被生成在 `$PATH_TO_WEEKLY_DIR/2019/2019-05-13-weekly.md` 中。

最后，我们可以利用 [dyweb/weekly][] 的构建脚本来构建真正的静态周报页面：

```bash
$ scripts/build.sh
```

随后就可以直接将其加入到 Git 版本管理中。

```bash
$ git add .
$ git commit -m "weekly: Add 2019-05-13"
$ git push
```

## 构建环境

### [dyweb/dy-bot][]

```bash
$ make install
$ ./dy-bot -o dyweb -r weekly -w $PATH_TO_WEEKLY_DIR -l :8123 -t $GITHUB_TOKEN
```

bot 是运行在服务器上的，但因为 bot 在实现的时候，想支持所有过程的自动化，因此需要在服务器上指定一个路径 `-w $PATH_TO_WEEKLY_DIR`。尽管现在这个是多余的。目前 dy-bot 由于某种 bug，不能支持自动地利用 [dyweb/dy-weekly-generator][] 生成 markdown 并且将其加入到 weekly 的版本管理中，所以仍需手动进行。

### [dyweb/dy-weekly-generator][]

在 [https://github.com/dyweb/dy-weekly-generator/releases/tag/v0.3.3](https://github.com/dyweb/dy-weekly-generator/releases/tag/v0.3.3) 下载对应版本的二进制，加入到系统路径变量中即可。

### [dyweb/weekly][]

首先安装 ruby 相关依赖，然后：

```bash
$ ./scrips/install-dep.sh
```

[dyweb/weekly]: https://github.com/dyweb/weekly
[dyweb/dy-bot]: https://github.com/dyweb/dy-bot
[dyweb/dy-weekly-generator]: https://github.com/dyweb/dy-weekly-generator
