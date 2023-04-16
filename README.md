# ChatgptApp
Intelligent customer service assistant

## Description

go serves as the webscoekt server, uses python to develop specific gpt application modules, and the two communicate through scoket

## Features

Modular development ideas, using python to quickly develop the required applications and integrate them online

- Supports Multi-user
- Support go and python process communication


## Project Structure
  
├── README.md  
├── SmartAssistant  // web socket server  
│&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;└── cmd // 启动的入口文件   
│&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;└── main.go         
├── README.md   
├── configs  // 维护配置文件信息  
│	&nbsp;&nbsp;&nbsp;&nbsp;└── settings //配置模块信息  
│&nbsp;&nbsp;&nbsp;&nbsp;	│── ipc_server.go //ipc 通信配置信息  
│&nbsp;&nbsp;&nbsp;&nbsp;	│── mysql.go //数据库配置信息  
│&nbsp;&nbsp;&nbsp;&nbsp;	│── open_ai.go //openap配置信息  
│&nbsp;&nbsp;&nbsp;&nbsp; └── appsettings.ini  
│── go.mod  
├── go.sum  
├── internal  // 核心实现逻辑    
│&nbsp;&nbsp;&nbsp;&nbsp;   ├── conf  
│&nbsp;&nbsp;&nbsp;&nbsp;   ├── data  // 业务数据访问，实体层、数据库迁移记录和对应业务的仓储实现  
│&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;   │  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;│── entity  
│&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;   │  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;│── migrations  
│   &nbsp;&nbsp;&nbsp;&nbsp;│  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;│── repo  
│   &nbsp;&nbsp;&nbsp;&nbsp;├── handle  // 处理服务端和用户对话信息  
│   &nbsp;&nbsp;&nbsp;&nbsp;└── service  // web socket服务端核心处理逻辑  
└── third_party  // api 依赖的第三方proto  
│&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;    ├── chatgpt //gpt相关接口封装  
│&nbsp;&nbsp;&nbsp;&nbsp;    ├── client //模拟websocket服务端  
│&nbsp;&nbsp;&nbsp;&nbsp;    └── ipc  //进程间通信核心代码  
│── utils   // 辅助工具包  
│    &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;├── event //事件  
│    &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;├── wire  //依赖注入管理  
│    &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;└── process_manager  //进程管理  
## Installation and Usage

In this section, provide instructions on how to install and use your project, such as:

### Installation

1. Install python3.0 and go1.20 environment
2. Clone the repository with the command `git clone git@github.com:Linwenqiang-git/ChatgptApp.git`.
3. Start the project with `go run main.go`.

### Usage

In this section, provide usage examples and instructions for using your project, such as:

```go 
cd cmd
go run main.go
cd /third_party/client
go run client.go
