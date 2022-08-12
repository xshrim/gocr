# 验证码识别工具

## 说明

本工具基于opencv图像处理库和tesseract光学字符识别引擎，采用golang开发，支持对简单的字符验证码图像进行识别。

工具对于仅包含数字、字母和干扰线的简单验证码识别率较高，对于字符形变和带其他字符的验证码识别效果可能较差。

由于工具使用的opencv和tesseract开发库不支持跨平台，且安装复杂（尤其是windows环境），因此我们直接在配置好基础环境的docker容器中编译运行。

## 编译

```
sh
docker build -t cmsent/ocr:latest .
```

## 运行

```
sh
# 测试
docker run --rm --name ocr -p 2333:2333 -it cmsent/ocr
# 生产
docker run  --restart=always --name ocr -p 2333:2333 -d cmsent/ocr
```

## 调用

支持file和base64两种调用方式，languages字段定义识别语系，whitelist字段定义可识别的字符，二者均是可选参数

- Base64

  ```sh
  curl -XPOST -s "http://23.234.239.243:2333/base64" -H "Content-Type: application/json" -d '{"languages":"eng", "whitelist":"0123456789abcdefghijkmnpqrstuvwxyzABCDEFGHIJKLMNPQRSTUVWXYZ", "data":"iVBORw0KGgoAAAANSUhEUgAAAGQAAAAoCAIAAACHGsgUAAABu0lEQVR42u3Z3U3DMBAHcE/CCGzA\nI4gNkLoBjMFrXhmFAViAFzoQWLIUItt3uS/bl8bWqarUJlJ/+vt0ccPXXOQVJsHEmlgnxwphiTWx\nboppGNYRmXpjJSMPTPfLi0esn9e7WMkovYdKdv+n98+1KN//fQxRKr4KqjkWhUmMtZU6NhadySRW\nFCylVNtk6TnoUkSvarcqUQb0rHZYValdLKivnwWL1eMNsJ4flqw+wjVV9QL800ZY8XdmtdIQsZBx\nQZ6sqFDylTUwWcml5NtWP6xtQZna1WwqlYWoDBdOuYWjS9V7VhWLvvvSsoonZQTFdyIUNFYwSVhJ\nhyuV7UGouKCQiAALTxbEF/AuDm1JEyxZQlnT6e5joHYbVrE0c4PJ4I5kZzAW1LlGTVjQFIqPprZS\nVCzNHuwsNQBL2bM8YFEOrQywqlKjsIiPMpkX8XhPi5XRyLz6Y2VePbDwwb0/FutgjyulwoJEuOEy\nl2qExZX6xyKeNAzB4j5g0w9OJVi7FqxwmWBxpdarmmPdzBL/x3U8rLfL98Tq5NVayuM21OfrRFie\nvZw2eJ9eE+v4WD69XM9Z3ry8D6WuvP4AN26jXn8WX7oAAAAASUVORK5CYII=","reqId":"ed00c4b4801b4563819ac6ab8a553d79"}'
  ```
- File

  ```sh
  curl -XPOST -s "http://127.0.0.1:2333/file" -F "file=@a.png" -F 'json={"languages":"eng","whitelist":"0123456789abcdefghijkmnpqrstuvwxyzABCDEFGHIJKLMNPQRSTUVWXYZ"}' -H "Content-Type: multipart/form-data"
  ```
