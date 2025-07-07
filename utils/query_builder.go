package utils

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// QueryBuilder 查询构建器
type QueryBuilder struct {
	db *gorm.DB
}

// NewQueryBuilder 创建查询构建器
func NewQueryBuilder(db *gorm.DB) *QueryBuilder {
	return &QueryBuilder{db: db}
}

// Where 添加WHERE条件
func (qb *QueryBuilder) Where(query interface{}, args ...interface{}) *QueryBuilder {
	qb.db = qb.db.Where(query, args...)
	return qb
}

// Or 添加OR条件
func (qb *QueryBuilder) Or(query interface{}, args ...interface{}) *QueryBuilder {
	qb.db = qb.db.Or(query, args...)
	return qb
}

// Not 添加NOT条件
func (qb *QueryBuilder) Not(query interface{}, args ...interface{}) *QueryBuilder {
	qb.db = qb.db.Not(query, args...)
	return qb
}

// Order 添加排序
func (qb *QueryBuilder) Order(value interface{}) *QueryBuilder {
	qb.db = qb.db.Order(value)
	return qb
}

// Limit 添加限制
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.db = qb.db.Limit(limit)
	return qb
}

// Offset 添加偏移
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.db = qb.db.Offset(offset)
	return qb
}

// Select 选择字段
func (qb *QueryBuilder) Select(query interface{}, args ...interface{}) *QueryBuilder {
	qb.db = qb.db.Select(query, args...)
	return qb
}

// Preload 预加载关联
func (qb *QueryBuilder) Preload(query string, args ...interface{}) *QueryBuilder {
	qb.db = qb.db.Preload(query, args...)
	return qb
}

// Joins 添加JOIN
func (qb *QueryBuilder) Joins(query string, args ...interface{}) *QueryBuilder {
	qb.db = qb.db.Joins(query, args...)
	return qb
}

// Group 添加GROUP BY
func (qb *QueryBuilder) Group(name string) *QueryBuilder {
	qb.db = qb.db.Group(name)
	return qb
}

// Having 添加HAVING条件
func (qb *QueryBuilder) Having(query interface{}, args ...interface{}) *QueryBuilder {
	qb.db = qb.db.Having(query, args...)
	return qb
}

// Distinct 添加DISTINCT
func (qb *QueryBuilder) Distinct(args ...interface{}) *QueryBuilder {
	qb.db = qb.db.Distinct(args...)
	return qb
}

// GetDB 获取数据库实例
func (qb *QueryBuilder) GetDB() *gorm.DB {
	return qb.db
}

// Find 查找记录
func (qb *QueryBuilder) Find(dest interface{}) error {
	return qb.db.Find(dest).Error
}

// First 查找第一条记录
func (qb *QueryBuilder) First(dest interface{}) error {
	return qb.db.First(dest).Error
}

// Last 查找最后一条记录
func (qb *QueryBuilder) Last(dest interface{}) error {
	return qb.db.Last(dest).Error
}

// Take 获取一条记录
func (qb *QueryBuilder) Take(dest interface{}) error {
	return qb.db.Take(dest).Error
}

// Count 统计记录数
func (qb *QueryBuilder) Count(count *int64) error {
	return qb.db.Count(count).Error
}

// Pluck 获取单个字段
func (qb *QueryBuilder) Pluck(column string, dest interface{}) error {
	return qb.db.Pluck(column, dest).Error
}

// Delete 删除记录
func (qb *QueryBuilder) Delete(value interface{}) error {
	return qb.db.Delete(value).Error
}

// Update 更新记录
func (qb *QueryBuilder) Update(column string, value interface{}) error {
	return qb.db.Update(column, value).Error
}

// Updates 批量更新记录
func (qb *QueryBuilder) Updates(values interface{}) error {
	return qb.db.Updates(values).Error
}

// 常用查询条件构建方法

// WhereID 根据ID查询
func (qb *QueryBuilder) WhereID(id interface{}) *QueryBuilder {
	return qb.Where("id = ?", id)
}

// WhereIn 根据ID列表查询
func (qb *QueryBuilder) WhereIn(column string, values interface{}) *QueryBuilder {
	return qb.Where(fmt.Sprintf("%s IN ?", column), values)
}

// WhereNotIn 根据ID列表排除查询
func (qb *QueryBuilder) WhereNotIn(column string, values interface{}) *QueryBuilder {
	return qb.Where(fmt.Sprintf("%s NOT IN ?", column), values)
}

// WhereLike 模糊查询
func (qb *QueryBuilder) WhereLike(column, value string) *QueryBuilder {
	return qb.Where(fmt.Sprintf("%s LIKE ?", column), "%"+value+"%")
}

// WhereLeftLike 左模糊查询
func (qb *QueryBuilder) WhereLeftLike(column, value string) *QueryBuilder {
	return qb.Where(fmt.Sprintf("%s LIKE ?", column), "%"+value)
}

// WhereRightLike 右模糊查询
func (qb *QueryBuilder) WhereRightLike(column, value string) *QueryBuilder {
	return qb.Where(fmt.Sprintf("%s LIKE ?", column), value+"%")
}

// WhereBetween 范围查询
func (qb *QueryBuilder) WhereBetween(column string, start, end interface{}) *QueryBuilder {
	return qb.Where(fmt.Sprintf("%s BETWEEN ? AND ?", column), start, end)
}

// WhereDate 日期查询
func (qb *QueryBuilder) WhereDate(column string, date time.Time) *QueryBuilder {
	return qb.Where(fmt.Sprintf("DATE(%s) = ?", column), date.Format("2006-01-02"))
}

// WhereDateRange 日期范围查询
func (qb *QueryBuilder) WhereDateRange(column string, start, end time.Time) *QueryBuilder {
	return qb.Where(fmt.Sprintf("DATE(%s) BETWEEN ? AND ?", column), 
		start.Format("2006-01-02"), end.Format("2006-01-02"))
}

// WhereTimeRange 时间范围查询
func (qb *QueryBuilder) WhereTimeRange(column string, start, end time.Time) *QueryBuilder {
	return qb.Where(fmt.Sprintf("%s BETWEEN ? AND ?", column), start, end)
}

// WhereStatus 状态查询
func (qb *QueryBuilder) WhereStatus(status interface{}) *QueryBuilder {
	return qb.Where("status = ?", status)
}

// WhereUserID 用户ID查询
func (qb *QueryBuilder) WhereUserID(userID interface{}) *QueryBuilder {
	return qb.Where("user_id = ?", userID)
}

// WhereCreatedAt 创建时间查询
func (qb *QueryBuilder) WhereCreatedAt(start, end time.Time) *QueryBuilder {
	return qb.WhereTimeRange("created_at", start, end)
}

// WhereUpdatedAt 更新时间查询
func (qb *QueryBuilder) WhereUpdatedAt(start, end time.Time) *QueryBuilder {
	return qb.WhereTimeRange("updated_at", start, end)
}

// OrderByID 按ID排序
func (qb *QueryBuilder) OrderByID(desc bool) *QueryBuilder {
	if desc {
		return qb.Order("id DESC")
	}
	return qb.Order("id ASC")
}

// OrderByCreatedAt 按创建时间排序
func (qb *QueryBuilder) OrderByCreatedAt(desc bool) *QueryBuilder {
	if desc {
		return qb.Order("created_at DESC")
	}
	return qb.Order("created_at ASC")
}

// OrderByUpdatedAt 按更新时间排序
func (qb *QueryBuilder) OrderByUpdatedAt(desc bool) *QueryBuilder {
	if desc {
		return qb.Order("updated_at DESC")
	}
	return qb.Order("updated_at ASC")
}

// Paginate 分页查询
func (qb *QueryBuilder) Paginate(page, pageSize int) *QueryBuilder {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	return qb.Offset(offset).Limit(pageSize)
}

// 高级查询方法

// WhereCondition 根据条件map构建查询
func (qb *QueryBuilder) WhereCondition(conditions map[string]interface{}) *QueryBuilder {
	for key, value := range conditions {
		if value != nil {
			// 处理特殊操作符
			if strings.Contains(key, "__") {
				parts := strings.Split(key, "__")
				if len(parts) == 2 {
					column := parts[0]
					operator := parts[1]
					
					switch operator {
					case "like":
						qb.WhereLike(column, fmt.Sprintf("%v", value))
					case "in":
						qb.WhereIn(column, value)
					case "not_in":
						qb.WhereNotIn(column, value)
					case "gt":
						qb.Where(fmt.Sprintf("%s > ?", column), value)
					case "gte":
						qb.Where(fmt.Sprintf("%s >= ?", column), value)
					case "lt":
						qb.Where(fmt.Sprintf("%s < ?", column), value)
					case "lte":
						qb.Where(fmt.Sprintf("%s <= ?", column), value)
					case "between":
						if slice, ok := value.([]interface{}); ok && len(slice) == 2 {
							qb.WhereBetween(column, slice[0], slice[1])
						}
					default:
						qb.Where(fmt.Sprintf("%s = ?", column), value)
					}
				}
			} else {
				qb.Where(fmt.Sprintf("%s = ?", key), value)
			}
		}
	}
	return qb
}

// WhereSearch 搜索查询（支持多个字段）
func (qb *QueryBuilder) WhereSearch(search string, columns ...string) *QueryBuilder {
	if search == "" || len(columns) == 0 {
		return qb
	}
	
	var conditions []string
	var args []interface{}
	
	for _, column := range columns {
		conditions = append(conditions, fmt.Sprintf("%s LIKE ?", column))
		args = append(args, "%"+search+"%")
	}
	
	query := strings.Join(conditions, " OR ")
	return qb.Where(query, args...)
}

// WhereTimeFilter 时间过滤查询
func (qb *QueryBuilder) WhereTimeFilter(column string, filter string, value interface{}) *QueryBuilder {
	switch filter {
	case "today":
		return qb.Where(fmt.Sprintf("DATE(%s) = CURDATE()", column))
	case "yesterday":
		return qb.Where(fmt.Sprintf("DATE(%s) = DATE_SUB(CURDATE(), INTERVAL 1 DAY)", column))
	case "this_week":
		return qb.Where(fmt.Sprintf("YEARWEEK(%s) = YEARWEEK(NOW())", column))
	case "last_week":
		return qb.Where(fmt.Sprintf("YEARWEEK(%s) = YEARWEEK(DATE_SUB(NOW(), INTERVAL 1 WEEK))", column))
	case "this_month":
		return qb.Where(fmt.Sprintf("YEAR(%s) = YEAR(NOW()) AND MONTH(%s) = MONTH(NOW())", column, column))
	case "last_month":
		return qb.Where(fmt.Sprintf("YEAR(%s) = YEAR(DATE_SUB(NOW(), INTERVAL 1 MONTH)) AND MONTH(%s) = MONTH(DATE_SUB(NOW(), INTERVAL 1 MONTH))", column, column))
	case "this_year":
		return qb.Where(fmt.Sprintf("YEAR(%s) = YEAR(NOW())", column))
	case "last_year":
		return qb.Where(fmt.Sprintf("YEAR(%s) = YEAR(DATE_SUB(NOW(), INTERVAL 1 YEAR))", column))
	default:
		return qb
	}
}

// 统计查询方法

// CountByCondition 按条件统计
func (qb *QueryBuilder) CountByCondition(condition string, args ...interface{}) (int64, error) {
	var count int64
	err := qb.db.Where(condition, args...).Count(&count).Error
	return count, err
}

// Sum 求和
func (qb *QueryBuilder) Sum(column string) (float64, error) {
	var result float64
	err := qb.db.Select(fmt.Sprintf("SUM(%s) as total", column)).Scan(&result).Error
	return result, err
}

// Avg 平均值
func (qb *QueryBuilder) Avg(column string) (float64, error) {
	var result float64
	err := qb.db.Select(fmt.Sprintf("AVG(%s) as average", column)).Scan(&result).Error
	return result, err
}

// Max 最大值
func (qb *QueryBuilder) Max(column string) (interface{}, error) {
	var result interface{}
	err := qb.db.Select(fmt.Sprintf("MAX(%s) as maximum", column)).Scan(&result).Error
	return result, err
}

// Min 最小值
func (qb *QueryBuilder) Min(column string) (interface{}, error) {
	var result interface{}
	err := qb.db.Select(fmt.Sprintf("MIN(%s) as minimum", column)).Scan(&result).Error
	return result, err
} 