# Review 评论审核服务

## 概述

Review服务是社区平台的评论和内容审核管理系统，提供完整的评论发布、审核、点赞点踩等功能。支持多级评论、内容审核、缓存优化等特性，为社区内容互动提供安全可靠的服务。

## 主要功能

### 💬 评论管理
- **发布评论**：支持多级评论结构（一级评论、二级评论等）
- **评论审核**：集成文本检测服务进行内容审核
- **评论状态**：正常、审核中、违规删除等多种状态管理
- **评论统计**：点赞数、点踩数、子评论数量统计

### 👍 互动功能
- **点赞点踩**：用户对评论进行点赞或点踩操作
- **互动记录**：记录用户的点赞点踩历史
- **状态切换**：支持点赞到点踩的无缝切换

### 🔍 查询功能
- **评论列表**：根据BaseId获取评论列表（支持分页）
- **评论详情**：根据ReviewId获取特定评论信息
- **根评论查询**：根据RootId获取整个评论树结构

### 🗑️ 管理功能
- **删除评论**：递归删除评论及其子评论
- **批量删除**：根据HeadId批量删除相关评论
- **状态更新**：更新评论审核状态

## 技术架构

- **数据存储**：MongoDB存储评论数据和用户互动数据
- **缓存策略**：Redis Sorted Set缓存评论列表，支持权重排序
- **消息队列**：Kafka集成文本检测服务进行内容审核
- **缓存穿透防护**：设置短时间空缓存防止缓存穿透

## 缓存机制

### 评论列表缓存
- **数据结构**：Redis Sorted Set，Key为BaseId
- **权重计算**：(点赞数 * -100) + 时间戳分钟数
- **过期策略**：1分钟自动过期
- **缓存更新**：数据变更时主动删除缓存

### 性能优化
- **分页查询**：支持offset和limit的精确分页
- **批量操作**：支持批量查询和更新操作
- **异步处理**：评论审核采用异步处理机制

## 数据模型

### 评论数据表 (MongoDB)
```javascript
{
  reviewId: Number,      // 评论ID（主键）
  userId: Number,        // 用户ID
  baseId: Number,        // 被评论对象ID
  rootId: Number,        // 根评论ID
  path: String,          // 评论路径（快速删除用）
  text: String,          // 评论文本
  status: Number,        // 状态（正常/审核中/违规）
  likeCount: Number,     // 点赞总数
  dislikeCount: Number,  // 点踩总数
  level: Number,         // 评论级别
  subCommentCount: Number,     // 子评论数量
  rootSubCommentCount: Number, // 根节点子评论总数
  createDate: Number,    // 创建时间戳
  updateDate: Number     // 更新时间戳
}
```

### 用户互动表 (MongoDB)
```javascript
{
  operationId: String,   // 操作ID（userId+reviewId）
  userId: Number,        // 用户ID
  reviewId: Number,      // 评论ID
  status: Number,        // 互动状态（点赞/点踩）
  createDate: Number,    // 创建时间戳
  updateDate: Number     // 更新时间戳
}
```

## 审核流程

1. **评论发布**：用户提交评论 → 并发写入数据库（状态：审核中）
2. **内容审核**：发送到Kafka队列 → 文本检测服务审核
3. **审核结果**：审核通过 → 更新状态为正常；审核失败 → 更新状态为违规
4. **缓存清理**：审核完成后清理相关Redis缓存

## API接口

- `ReviewCreate` - 创建评论
- `ReviewDelete` - 删除评论
- `ReviewDeleteByHeadId` - 根据头ID删除评论
- `ReviewGetByBaseId` - 根据基础ID获取评论列表
- `ReviewGetByRootId` - 根据根ID获取评论树
