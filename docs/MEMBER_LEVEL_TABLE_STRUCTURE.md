# Member Level 表结构文档

## 表名
`member_level`

## 表描述
用户等级配置表 - 存储用户等级配置信息，包括等级、名称、logo、返现比例、单数字额等

## 最新表结构

```sql
CREATE TABLE `member_level` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `level` int NOT NULL COMMENT '等级数值',
  `name` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '等级名称',
  `logo` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '等级logo',
  `remark` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注',
  `cashback_ratio` decimal(5,2) DEFAULT '0.00' COMMENT '返现比例（百分比）',
  `single_amount` int DEFAULT 1 COMMENT '单数字额',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_level` (`level`),
  KEY `idx_member_level_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

## 字段说明

| 字段名 | 类型 | 是否为空 | 默认值 | 说明 |
|--------|------|----------|--------|------|
| id | bigint unsigned | NOT NULL | AUTO_INCREMENT | 主键ID |
| level | int | NOT NULL | - | 等级数值 |
| name | varchar(20) | NOT NULL | - | 等级名称 |
| logo | varchar(255) | NULL | NULL | 等级logo |
| remark | varchar(255) | NULL | NULL | 备注 |
| cashback_ratio | decimal(5,2) | NULL | 0.00 | 返现比例（百分比） |
| single_amount | int | NULL | 1 | 单数字额 |
| created_at | datetime(3) | NULL | NULL | 创建时间 |
| updated_at | datetime(3) | NULL | NULL | 更新时间 |
| deleted_at | datetime(3) | NULL | NULL | 软删除时间 |

## 索引说明

| 索引名 | 类型 | 字段 | 说明 |
|--------|------|------|------|
| PRIMARY | PRIMARY KEY | id | 主键索引 |
| uniq_level | UNIQUE | level | 等级唯一索引 |
| idx_member_level_deleted_at | INDEX | deleted_at | 软删除时间索引 |

## 变更记录

### 最新变更（当前版本）
- **移除字段**：
  - `play_count` (bigint) - 每日游戏次数
  - `upgrade_amount` (decimal(18,2)) - 升级要求（充值金额）
- **新增字段**：
  - `single_amount` (int) - 单数字额，默认值为1
- **字段类型调整**：
  - `level` 从 bigint 改为 int
  - `id` 从 uint 改为 uint64

### 字段变更详情
1. **删除字段**：
   ```sql
   ALTER TABLE `member_level` 
   DROP COLUMN `play_count`,
   DROP COLUMN `upgrade_amount`;
   ```

2. **新增字段**：
   ```sql
   ALTER TABLE `member_level` 
   ADD COLUMN `single_amount` int DEFAULT 1 COMMENT '单数字额' AFTER `cashback_ratio`;
   ```

3. **修改字段类型**：
   ```sql
   ALTER TABLE `member_level` 
   MODIFY COLUMN `level` int NOT NULL COMMENT '等级数值';
   ```

## Go模型对应

```go
type MemberLevel struct {
    ID            uint64         `gorm:"primarykey" json:"id"`
    Level         int            `gorm:"not null;uniqueIndex:uniq_level;comment:等级数值" json:"level"`
    Name          string         `gorm:"size:20;not null;comment:等级名称" json:"name"`
    Logo          string         `gorm:"size:255;comment:等级logo" json:"logo"`
    Remark        string         `gorm:"size:255;comment:备注" json:"remark"`
    CashbackRatio float64        `gorm:"type:decimal(5,2);default:0;comment:返现比例（百分比）" json:"cashback_ratio"`
    SingleAmount  int            `gorm:"default:1;comment:单数字额" json:"single_amount"`
    CreatedAt     time.Time      `gorm:"type:datetime(3);autoCreateTime" json:"created_at"`
    UpdatedAt     time.Time      `gorm:"type:datetime(3);autoUpdateTime" json:"updated_at"`
    DeletedAt     gorm.DeletedAt `gorm:"type:datetime(3);index;comment:软删除时间" json:"-"`
}
```

## 使用说明

1. **等级数值**：用于表示用户等级的数字，必须唯一
2. **等级名称**：等级的显示名称，最大20个字符
3. **等级logo**：等级对应的图标或logo图片URL
4. **返现比例**：该等级用户的返现比例，以百分比形式存储（如：5.00表示5%）
5. **单数字额**：该等级用户单次操作的金额限制，整数类型，默认值为1
6. **软删除**：使用gorm.DeletedAt实现软删除功能

## 注意事项

1. 等级数值必须唯一，不能重复
2. 返现比例范围为0-100%
3. 单数字额必须为正整数，最小值为1
4. 支持软删除，删除的数据不会物理删除，只是标记deleted_at时间
5. 表使用utf8mb4字符集，支持emoji等特殊字符
6. 使用gorm.DeletedAt类型实现标准的软删除功能 