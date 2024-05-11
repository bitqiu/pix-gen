# pix-gen

## 验证码图片生成

### URL

> GET http://localhost:8080/captcha?code={code}&width={width}&height={height}

### 参数

- `code`: 验证码内容
- `width` (可选): 验证码宽度，默认为 `120`
- `height` (可选): 验证码高度，默认为 `30`

示例请求：

> GET http://localhost:8080/captcha?code=abcdef&width=120&height=30

## 二维码图片生成

### URL

> GET http://localhost:8080/qrcode?text={text}&size={size}&level={level}&color={color}

### 参数

- `text`: 二维码内容
- `size` (可选): 二维码大小，默认为 `300`
- `level` (可选): 二维码容错率，默认为 `H`，可选 `L`, `M`, `Q`, `H`
- `color` (可选): 二维码颜色，默认为 `#549ecc`（16进制，不包含`#`号）

示例请求：

> GET http://localhost:8080/qrcode?text=helloworld&size=400&level=L&color=549ecc
