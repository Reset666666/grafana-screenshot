# 📸 Grafana Screenshot Automation (Go + Renderer)

> 自动调用 **Grafana Image Renderer** 生成仪表盘截图  
> 适用于日报、监控周报或定时生成图表快照,中文 PDF、定时任务、企业微信机器人推送

---

## 🚀 功能特性

✅ 自动请求 Grafana 渲染接口生成图像  
✅ 支持中文字体（无乱码）  
✅ 自动生成文件名（含时间戳）  
✅ 支持时间范围动态计算  
✅ 保留顶部栏、隐藏左边栏  
✅ 可配置图片尺寸、时区  
✅ 兼容自编译版 `grafana-image-renderer`
✅ Cron 定时任务自动执行
✅ DevMode 调试模式（立即运行）
✅ 企业微信机器人消息推送
✅ 自动清理临时截图


---

### 安装 grafana-image-renderer 服务
安装参考官方文档：
https://github.com/grafana/grafana-image-renderer?tab=readme-ov-file#remote-rendering-service-installation
https://grafana.com/grafana/plugins/grafana-image-renderer/


---
### 运行
1. 克隆代码
   ```bash
    git clone
    cd grafana-screenshot
    ```
2. 修改配置文件 `config_bak.yaml` 为 `config.yaml`，并根据实际情况修改配置项
3. 编译运行