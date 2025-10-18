# NPTUI 使用指南

## 编译测试结果 ✅

### 编译成功
- **Go 版本**: 1.25.1
- **AMD64 二进制**: 3.7 MB (x86_64)
- **ARM64 二进制**: 3.6 MB (aarch64)

### DEB 包
- **AMD64 包**: `build/nptui-1.0.0-amd64.deb` (1.3 MB)
- **ARM64 包**: `build/nptui-1.0.0-arm64.deb` (1.2 MB)

### 安装状态
✅ 已成功安装到 `/usr/bin/nptui`

---

## 快速开始

### 1. 运行程序
```bash
sudo nptui
```

### 2. 界面导航

#### 主菜单选项
1. **Edit Network Interfaces** - 编辑网络接口
2. **Apply Configuration** - 应用配置
3. **Quit** - 退出

#### 快捷键
- `↑` / `↓` - 上下移动
- `Enter` - 选择/确认
- `Tab` - 切换表单字段
- `Esc` - 返回上级菜单
- `q` - 退出程序（主菜单）
- `b` - 返回（网卡列表）

### 3. 配置网络接口

#### DHCP 配置（自动获取 IP）
1. 选择 "Edit Network Interfaces"
2. 选择要配置的网卡（如 eth0）
3. 配置方式选择 "DHCP"
4. 点击 "Save"

#### 静态 IP 配置
1. 选择 "Edit Network Interfaces"
2. 选择要配置的网卡
3. 配置方式选择 "Static"
4. 填写以下信息：
   - **IP Address/CIDR**: 例如 `192.168.1.100/24`
   - **Gateway**: 例如 `192.168.1.1`
   - **DNS Server**: 例如 `8.8.8.8`
5. 点击 "Save"

### 4. 应用配置
配置保存后，选择 "Apply Configuration" 或运行：
```bash
sudo netplan apply
```

---

## 配置示例

### 示例 1: 家庭网络 DHCP
```
网卡: eth0
配置: DHCP
```

### 示例 2: 服务器静态 IP
```
网卡: eth0
配置: Static
IP: 192.168.1.100/24
网关: 192.168.1.1
DNS: 8.8.8.8
```

### 示例 3: 内网服务器
```
网卡: eth0
配置: Static
IP: 10.0.0.50/24
网关: 10.0.0.1
DNS: 10.0.0.1
```

---

## 配置文件位置

程序会生成/修改以下文件：
- `/etc/netplan/01-netcfg.yaml` - 主配置文件
- `/etc/netplan/*.yaml` - 其他 netplan 配置

### 查看生成的配置
```bash
cat /etc/netplan/01-netcfg.yaml
```

### 示例配置文件
```yaml
network:
  version: 2
  renderer: networkd
  ethernets:
    eth0:
      dhcp4: true
```

或静态 IP：
```yaml
network:
  version: 2
  renderer: networkd
  ethernets:
    eth0:
      dhcp4: false
      addresses:
        - 192.168.1.100/24
      gateway4: 192.168.1.1
      nameservers:
        addresses:
          - 8.8.8.8
```

---

## 故障排查

### 程序无法启动
**问题**: `This program must be run as root`  
**解决**: 使用 sudo 运行
```bash
sudo nptui
```

### 网络配置不生效
**问题**: 修改后网络未变化  
**解决**: 手动应用配置
```bash
sudo netplan apply
```

### 查看网络状态
```bash
# 查看所有网卡
ip addr

# 查看路由
ip route

# 测试连接
ping -c 4 8.8.8.8
```

### 配置错误恢复
如果配置错误导致网络断开：
```bash
# 恢复备份（如果有）
sudo cp /etc/netplan/01-netcfg.yaml.backup /etc/netplan/01-netcfg.yaml

# 或删除配置重新开始
sudo rm /etc/netplan/01-netcfg.yaml

# 应用配置
sudo netplan apply
```

---

## 卸载

```bash
sudo dpkg -r nptui
```

---

## 重新构建

### 编译
```bash
cd /root/nptui
make clean
make build
```

### 打包 DEB
```bash
# 打包所有架构
make deb-all

# 或单独打包
make deb-amd64  # AMD64
make deb-arm64  # ARM64
```

### 使用构建脚本
```bash
./build.sh
```

---

## 技术细节

### 依赖库
- `github.com/rivo/tview` - TUI 框架
- `github.com/gdamore/tcell/v2` - 终端库
- `gopkg.in/yaml.v3` - YAML 解析

### 系统要求
- Linux kernel 2.6+
- netplan.io
- Root 权限

### 架构支持
- ✅ AMD64 (x86_64)
- ✅ ARM64 (aarch64)
- ⚠️ 其他架构需要重新编译

---

## 开发信息

### 项目结构
```
nptui/
├── main.go           # 主入口
├── netplan/          # netplan 配置管理
├── ui/               # TUI 界面
├── debian/           # DEB 打包
└── Makefile          # 构建脚本
```

### 贡献
欢迎提交 Issue 和 PR！

---

## 常见问题

**Q: 支持 IPv6 吗？**  
A: 当前版本主要支持 IPv4，IPv6 支持计划在未来版本添加。

**Q: 可以配置多个 DNS 吗？**  
A: 当前版本支持一个 DNS，多 DNS 支持在开发中。

**Q: 如何备份配置？**  
A: 建议在修改前备份：
```bash
sudo cp /etc/netplan/01-netcfg.yaml /etc/netplan/01-netcfg.yaml.backup
```

**Q: 支持其他发行版吗？**  
A: 支持所有使用 netplan 的发行版（如 Ubuntu 17.10+）。

---

## 更新日志

### v1.0.0 (2025-10-18)
- ✨ 初始版本
- ✅ DHCP 和静态 IP 配置
- ✅ DNS 配置
- ✅ ARM64 和 AMD64 支持
- ✅ DEB 包支持

