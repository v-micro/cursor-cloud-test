# 星链智科官网（静态演示）

本目录为虚构科技公司「星链智科」的单页中文官网，使用纯 HTML、CSS 与 JavaScript，无后端依赖。

## 本地预览

### 方式一：Python 内置服务器（推荐）

在项目仓库根目录执行：

```bash
cd website
python3 -m http.server 8080
```

浏览器打开 [http://localhost:8080](http://localhost:8080) 即可。

### 方式二：直接打开文件

部分浏览器在 `file://` 协议下对字体或模块加载有限制，建议使用上面的本地 HTTP 服务。

### 方式三：Node.js（如已安装 npx）

```bash
cd website
npx --yes serve -p 8080
```

## 目录结构

- `index.html` — 页面结构与各区块内容
- `styles.css` — 样式与响应式布局
- `main.js` — 导航菜单、表单校验与页脚年份

## 说明

联系表单仅在前端做校验与成功提示，不会向服务器发送数据。
