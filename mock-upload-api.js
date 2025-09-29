// 模拟上传资料API接口
// 这个文件用于模拟后端上传资料功能，实际开发时需要替换为真实的后端接口

// 模拟上传资料接口（网盘链接形式）
export const uploadMaterial = async (courseId, materialData) => {
  // 模拟网络延迟
  await new Promise(resolve => setTimeout(resolve, 1000));
  
  // 模拟成功响应
  return {
    success: true,
    data: {
      id: Math.floor(Math.random() * 10000),
      courseId: courseId,
      materialName: materialData.materialName,
      type: materialData.type,
      service: materialData.service,
      link: materialData.link,
      accessPassword: materialData.accessPassword,
      description: materialData.description,
      uploaderId: 1, // 模拟当前用户ID
      uploaderName: "当前用户",
      uploadTime: new Date().toISOString(),
      downloadCount: 0
    }
  };
};

// 模拟获取课程资料列表接口
export const getCourseMaterials = async (courseId) => {
  // 模拟网络延迟
  await new Promise(resolve => setTimeout(resolve, 500));
  
  // 模拟资料列表数据（网盘链接形式）
  const mockMaterials = [
    {
      id: 1,
      courseId: courseId,
      materialName: "计算机网络课件第一章.pdf",
      type: "cloud",
      service: "百度网盘",
      link: "https://pan.baidu.com/s/1abc123def456",
      accessPassword: "1234",
      description: "计算机网络基础概念和协议介绍",
      uploaderId: 1,
      uploaderName: "张三",
      uploadTime: "2023-10-27T10:00:00Z",
      downloadCount: 156
    },
    {
      id: 2,
      courseId: courseId,
      materialName: "实验报告模板.docx",
      type: "cloud",
      service: "夸克网盘",
      link: "https://pan.quark.cn/s/def456ghi789",
      accessPassword: "",
      description: "网络实验报告标准模板",
      uploaderId: 2,
      uploaderName: "李四",
      uploadTime: "2023-10-26T15:30:00Z",
      downloadCount: 89
    },
    {
      id: 3,
      courseId: courseId,
      materialName: "考试重点整理.md",
      type: "cloud",
      service: "阿里云盘",
      link: "https://www.aliyundrive.com/s/ghi789jkl012",
      accessPassword: "exam2023",
      description: "期末考试重点知识点整理",
      uploaderId: 3,
      uploaderName: "王五",
      uploadTime: "2023-10-25T20:15:00Z",
      downloadCount: 234
    },
    {
      id: 4,
      courseId: courseId,
      materialName: "网络协议详解.pptx",
      type: "cloud",
      service: "OneDrive",
      link: "https://1drv.ms/p/s!abc123def456",
      accessPassword: "",
      description: "详细讲解各种网络协议的工作原理",
      uploaderId: 4,
      uploaderName: "赵六",
      uploadTime: "2023-10-24T09:20:00Z",
      downloadCount: 78
    },
    {
      id: 5,
      courseId: courseId,
      materialName: "期末冲刺思维导图",
      type: "note",
      service: "幕布",
      link: "https://mubu.com/doc/xyz123",
      accessPassword: "",
      description: "课程知识点全景思维导图，适合复习使用",
      uploaderId: 5,
      uploaderName: "小林",
      uploadTime: "2023-10-23T14:10:00Z",
      downloadCount: 312
    },
    {
      id: 6,
      courseId: courseId,
      materialName: "课堂笔记整理",
      type: "note",
      service: "OneNote",
      link: "https://onenote.com/some-note",
      accessPassword: "note2023",
      description: "课堂重点内容与老师补充信息，自带目录导航",
      uploaderId: 6,
      uploaderName: "小王",
      uploadTime: "2023-10-22T19:05:00Z",
      downloadCount: 147
    }
  ];
  
  return {
    success: true,
    data: {
      total: mockMaterials.length,
      materials: mockMaterials
    }
  };
};

// 模拟保存笔记接口
export const saveCourseNote = async (courseId, content) => {
  // 模拟网络延迟
  await new Promise(resolve => setTimeout(resolve, 300));
  
  return {
    success: true,
    data: {
      id: Math.floor(Math.random() * 10000),
      courseId: courseId,
      content: content,
      lastModified: new Date().toISOString()
    }
  };
};

// 模拟获取笔记列表接口
export const getCourseNotes = async (courseId) => {
  // 模拟网络延迟
  await new Promise(resolve => setTimeout(resolve, 200));
  
  // 模拟多用户笔记数据
  const mockNotes = [
    {
      id: 1,
      courseId: courseId,
      content: `# 计算机网络学习笔记

## 第一章 网络基础

### 1.1 网络协议
- **TCP/IP协议栈**：四层模型
- **OSI模型**：七层模型

### 1.2 重要概念
- 分组交换
- 电路交换
- 报文交换

## 第二章 物理层

### 2.1 传输介质
- 双绞线
- 光纤
- 无线传输

### 2.2 编码技术
- 曼彻斯特编码
- 差分曼彻斯特编码

## 学习心得

这门课程让我深入理解了网络的工作原理，特别是TCP/IP协议栈的设计思想。

> 重点：理解每一层的功能和协议的作用

## 待办事项
- [ ] 完成实验报告
- [ ] 复习第二章内容
- [ ] 准备期中考试`,
      authorId: 1,
      authorName: "张三",
      createdAt: "2023-10-27T10:00:00Z",
      isOwner: true
    },
    {
      id: 2,
      courseId: courseId,
      content: `# 网络协议重点整理

## TCP vs UDP

| 特性 | TCP | UDP |
|------|-----|-----|
| 连接性 | 面向连接 | 无连接 |
| 可靠性 | 可靠 | 不可靠 |
| 速度 | 较慢 | 较快 |

## 重要端口号
- HTTP: 80
- HTTPS: 443
- FTP: 21
- SSH: 22

## 考试重点
1. **三次握手**过程
2. **四次挥手**过程
3. **滑动窗口**机制`,
      authorId: 2,
      authorName: "李四",
      createdAt: "2023-10-26T15:30:00Z",
      isOwner: false
    },
    {
      id: 3,
      courseId: courseId,
      content: `# 实验报告模板

## 实验一：网络抓包分析

### 实验目的
学习使用Wireshark进行网络数据包分析

### 实验步骤
1. 安装Wireshark
2. 选择网络接口
3. 开始抓包
4. 分析HTTP请求

### 实验结果
成功捕获到HTTP请求包，分析了请求头和响应头的内容。

### 实验总结
通过本次实验，深入理解了HTTP协议的工作机制。`,
      authorId: 3,
      authorName: "王五",
      createdAt: "2023-10-25T20:15:00Z",
      isOwner: false
    }
  ];

  return {
    success: true,
    data: {
      notes: mockNotes
    }
  };
};
