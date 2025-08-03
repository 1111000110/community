# Post 帖子模块

## 概述

Post模块是社区平台的核心内容模块，提供完整的帖子管理功能。用户可以发布、编辑、删除帖子，进行点赞、收藏、评论等社交互动。该模块支持多媒体内容、分类管理、标签系统和强大的搜索功能。

## 模块架构

```
post/
└── api/                    # HTTP API服务
    ├── post.api           # API接口定义文件
    ├── post.go            # 服务启动入口
    ├── etc/               # 配置文件目录
    │   └── post.yaml      # 服务配置文件
    └── internal/          # 内部实现
        ├── config/        # 配置结构定义
        ├── handler/       # HTTP请求处理器
        ├── logic/         # 业务逻辑层
        ├── middleware/    # 中间件
        ├── svc/          # 服务上下文
        └── types/        # 类型定义
```

## 核心功能

### 📝 帖子管理
- **创建帖子**：支持富文本内容、图片上传、分类标签
- **编辑帖子**：修改标题、内容、分类等信息
- **删除帖子**：软删除机制，保留数据完整性
- **帖子状态**：草稿、已发布、已删除、已隐藏

### 🔍 内容发现
- **帖子列表**：支持分页、筛选、搜索
- **分类浏览**：按技术、生活、娱乐等分类查看
- **标签系统**：多标签支持，便于内容组织
- **排序方式**：时间、热度、点赞数排序

### 💬 社交互动
- **点赞系统**：帖子和评论点赞功能
- **收藏功能**：个人收藏夹管理
- **评论系统**：多级评论、评论回复
- **分享功能**：社交媒体分享统计

### 📊 数据统计
- **浏览统计**：帖子浏览次数记录
- **互动统计**：点赞、评论、分享、收藏数据
- **用户行为**：个人互动状态跟踪
- **热度计算**：综合互动数据的热度排名

## 数据模型

### 帖子基本信息 (PostBase)
```go
type PostBase struct {
    PostId      int64  // 帖子唯一标识ID
    UserId      int64  // 发布者用户ID
    Title       string // 帖子标题
    Content     string // 帖子内容
    Images      string // 图片URL列表，JSON格式
    Category    string // 帖子分类
    Tags        string // 标签列表，逗号分隔
    Status      int32  // 帖子状态
    CreateTime  int64  // 创建时间戳
    UpdateTime  int64  // 更新时间戳
}
```

### 帖子统计信息 (PostStats)
```go
type PostStats struct {
    PostId       int64 // 帖子ID
    LikeCount    int64 // 点赞数
    CommentCount int64 // 评论数
    ShareCount   int64 // 分享数
    ViewCount    int64 // 浏览数
    CollectCount int64 // 收藏数
}
```

### 评论信息 (Comment)
```go
type Comment struct {
    CommentId    int64    // 评论ID
    PostId       int64    // 所属帖子ID
    UserId       int64    // 评论者用户ID
    ParentId     int64    // 父评论ID，0表示顶级评论
    Content      string   // 评论内容
    LikeCount    int64    // 评论点赞数
    CreateTime   int64    // 创建时间戳
    Author       UserBase // 评论者信息
    IsLiked      bool     // 当前用户是否已点赞
}
```

## API接口

### 🔐 认证接口组
需要JWT认证的接口，用于用户的帖子管理和互动操作。

#### 帖子管理
- `POST /post/create` - 创建帖子
- `POST /post/delete` - 删除帖子
- `POST /post/update` - 更新帖子
- `POST /post/detail` - 获取帖子详情
- `POST /post/list` - 获取帖子列表

#### 社交互动
- `POST /post/like` - 点赞/取消点赞帖子
- `POST /post/collect` - 收藏/取消收藏帖子
- `POST /comment/create` - 创建评论
- `POST /comment/list` - 获取评论列表
- `POST /comment/like` - 点赞/取消点赞评论

### 🌐 公开接口组
无需认证的公开接口，供游客浏览内容。

- `POST /public/post/list` - 公开帖子列表
- `POST /public/post/detail` - 公开帖子详情
- `POST /public/comment/list` - 公开评论列表

## 使用示例

### 创建帖子
```http
POST /post/create
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "title": "Go语言微服务架构实践",
  "content": "本文介绍如何使用go-zero框架构建微服务...",
  "images": "[\"https://example.com/image1.jpg\", \"https://example.com/image2.jpg\"]",
  "category": "tech",
  "tags": "Go,微服务,go-zero",
  "status": 1
}
```

### 获取帖子列表
```http
POST /post/list
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "page": 1,
  "page_size": 10,
  "category": "tech",
  "keyword": "Go语言",
  "sort_by": "hot"
}
```

### 点赞帖子
```http
POST /post/like
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "post_id": 12345,
  "action": 1
}
```

### 创建评论
```http
POST /comment/create
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "post_id": 12345,
  "parent_id": 0,
  "content": "很棒的文章，学到了很多！"
}
```

## 业务规则

### 帖子状态管理
- **草稿 (0)**：用户可以保存未完成的帖子
- **已发布 (1)**：正常显示在列表中的帖子
- **已删除 (2)**：软删除，不在列表中显示但保留数据
- **已隐藏 (3)**：管理员隐藏的帖子

### 权限控制
- 用户只能编辑和删除自己的帖子
- 管理员可以隐藏任何帖子
- 游客可以浏览公开内容但不能互动
- 登录用户可以进行所有互动操作

### 内容限制
- 帖子标题：1-100字符
- 帖子内容：1-10000字符
- 图片数量：最多9张
- 标签数量：最多5个
- 评论内容：1-500字符

## 技术特性

### 🚀 性能优化
- **分页查询**：避免大数据量查询
- **索引优化**：关键字段建立索引
- **缓存策略**：热门内容Redis缓存
- **图片处理**：支持多尺寸图片存储

### 🛡️ 安全特性
- **内容过滤**：敏感词过滤机制
- **防刷机制**：操作频率限制
- **权限验证**：严格的用户权限控制
- **数据验证**：输入参数严格验证

### 📊 可观测性
- **操作日志**：详细的用户操作记录
- **性能监控**：接口响应时间统计
- **错误追踪**：异常信息收集分析
- **业务指标**：内容发布和互动统计

## 扩展功能

### 🎯 推荐算法
- **个性化推荐**：基于用户行为的内容推荐
- **热门内容**：基于互动数据的热度计算
- **相关推荐**：基于标签和分类的相关内容
- **新用户引导**：为新用户推荐优质内容

### 🔔 通知系统
- **点赞通知**：帖子被点赞时通知作者
- **评论通知**：帖子被评论时通知作者
- **回复通知**：评论被回复时通知评论者
- **关注动态**：关注用户的新帖子通知

### 📈 数据分析
- **内容分析**：热门话题和趋势分析
- **用户画像**：基于发布和互动行为的用户分析
- **运营数据**：日活、内容发布量等运营指标
- **质量评估**：内容质量评分和推荐权重

## 部署配置

### 服务配置示例
```yaml
Name: post.api
Host: 0.0.0.0
Port: 8890
Timeout: 5s
MaxBytes: 10485760  # 10MB支持图片上传

# JWT配置
Auth:
  AccessSecret: "your_jwt_secret_key"
  AccessExpire: 86400

# 数据库配置
DataSource: "user:password@tcp(localhost:3306)/community?charset=utf8mb4&parseTime=true"

# Redis配置
RedisConf:
  Host: localhost:6379
  Type: node
  Pass: ""
  DB: 2

# 文件存储配置
FileStorage:
  Type: "local"  # local/oss/s3
  Path: "/uploads"
  MaxSize: 5242880  # 5MB
  AllowedTypes: ["jpg", "jpeg", "png", "gif"]
```

### 数据库设计

#### 帖子表 (posts)
```sql
CREATE TABLE posts (
    post_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    images JSON,
    category VARCHAR(50),
    tags VARCHAR(200),
    status TINYINT DEFAULT 1,
    create_time BIGINT NOT NULL,
    update_time BIGINT NOT NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_category (category),
    INDEX idx_status (status),
    INDEX idx_create_time (create_time)
);
```

#### 帖子统计表 (post_stats)
```sql
CREATE TABLE post_stats (
    post_id BIGINT PRIMARY KEY,
    like_count BIGINT DEFAULT 0,
    comment_count BIGINT DEFAULT 0,
    share_count BIGINT DEFAULT 0,
    view_count BIGINT DEFAULT 0,
    collect_count BIGINT DEFAULT 0,
    FOREIGN KEY (post_id) REFERENCES posts(post_id)
);
```

#### 评论表 (comments)
```sql
CREATE TABLE comments (
    comment_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    post_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    parent_id BIGINT DEFAULT 0,
    content VARCHAR(500) NOT NULL,
    like_count BIGINT DEFAULT 0,
    create_time BIGINT NOT NULL,
    INDEX idx_post_id (post_id),
    INDEX idx_user_id (user_id),
    INDEX idx_parent_id (parent_id),
    FOREIGN KEY (post_id) REFERENCES posts(post_id)
);
```

## 开发指南

### 添加新功能
1. 在 `post.api` 中定义新的接口和数据结构
2. 运行 `goctl api go -api post.api -dir .` 生成代码
3. 在 `internal/logic/` 中实现业务逻辑
4. 在 `internal/handler/` 中处理HTTP请求
5. 更新数据库模型和配置文件

### 代码规范
```go
// 业务逻辑示例
func (l *PostCreateLogic) PostCreate(req *types.PostCreateReq) (*types.PostCreateResp, error) {
    // 1. 参数验证
    if err := l.validateCreateReq(req); err != nil {
        return nil, err
    }
    
    // 2. 业务处理
    post := &model.Post{
        UserId:     l.ctx.Value("userId").(int64),
        Title:      req.Title,
        Content:    req.Content,
        Images:     req.Images,
        Category:   req.Category,
        Tags:       req.Tags,
        Status:     req.Status,
        CreateTime: time.Now().UnixMilli(),
        UpdateTime: time.Now().UnixMilli(),
    }
    
    // 3. 数据存储
    result, err := l.svcCtx.PostModel.Insert(l.ctx, post)
    if err != nil {
        return nil, err
    }
    
    // 4. 返回结果
    return &types.PostCreateResp{
        PostId: result.InsertedID.(int64),
    }, nil
}
```

## 测试

### 单元测试
```bash
# 运行所有测试
go test ./...

# 运行特定包测试
go test ./internal/logic

# 生成测试覆盖率报告
go test -cover ./...
```

### 接口测试
```bash
# 测试创建帖子
curl -X POST http://localhost:8890/post/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -d '{
    "title": "测试帖子",
    "content": "这是一个测试帖子的内容",
    "category": "test",
    "status": 1
  }'
```

## 监控告警

### 关键指标
- **帖子发布量**：每日新增帖子数量
- **用户活跃度**：点赞、评论、收藏等互动指标
- **内容质量**：平均点赞数、评论数等质量指标
- **系统性能**：接口响应时间、错误率等

### 告警规则
- 帖子发布失败率超过5%告警
- 接口响应时间超过2秒告警
- 数据库连接数超过80%告警
- 存储空间使用率超过85%告警

## 注意事项

1. **内容安全**：
   - 实现敏感词过滤机制
   - 建立内容审核流程
   - 用户举报和处理机制

2. **性能优化**：
   - 合理使用数据库索引
   - 实现有效的缓存策略
   - 避免N+1查询问题

3. **数据一致性**：
   - 统计数据的实时更新
   - 分布式事务处理
   - 数据备份和恢复策略

4. **用户体验**：
   - 响应式设计支持
   - 图片懒加载优化
   - 离线内容缓存

## 联系方式

- 作者：张璇
- 邮箱：xatuzx2025@163.com
- 版本：1.0