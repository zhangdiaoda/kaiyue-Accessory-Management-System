-- 权限管理系统 - 初始化权限数据
-- 注意：AutoMigrate已创建表结构，此脚本仅插入初始数据

-- 插入18个权限
INSERT INTO `sys_permission` (`code`, `name`, `category`, `description`) VALUES
-- 配件管理 (6个)
('part:view', '查看配件', '配件管理', '查看配件列表和详情'),
('part:add', '添加配件', '配件管理', '添加新配件入库'),
('part:edit', '编辑配件', '配件管理', '修改配件信息'),
('part:delete', '删除配件', '配件管理', '删除配件记录'),
('part:import', '导入配件', '配件管理', '批量导入配件数据'),
('part:export', '导出配件', '配件管理', '导出配件数据'),

-- 领用管理 (4个)
('borrow:create', '创建领用', '领用管理', '创建领用记录'),
('borrow:view', '查看领用记录', '领用管理', '查看领用记录'),
('borrow:return', '处理归还', '领用管理', '处理配件归还'),
('borrow:dispose', '处置旧件', '领用管理', '处置旧件和废品'),

-- 员工管理 (6个)
('employee:view', '查看员工', '员工管理', '查看员工列表和详情'),
('employee:add', '添加员工', '员工管理', '添加新员工'),
('employee:edit', '编辑员工', '员工管理', '修改员工信息'),
('employee:delete', '删除员工', '员工管理', '删除员工记录'),
('employee:import', '导入员工', '员工管理', '批量导入员工数据'),
('employee:export', '导出员工', '员工管理', '导出员工数据'),

-- 报表管理 (2个)
('report:view', '查看报表', '报表管理', '查看各类报表'),
('report:export', '导出报表', '报表管理', '导出报表数据');

-- 插入仓库管理员默认权限 (8个基础权限)
INSERT INTO `sys_role_permission` (`role`, `permission_code`) VALUES
('WAREHOUSE_ADMIN', 'part:view'),
('WAREHOUSE_ADMIN', 'part:add'),
('WAREHOUSE_ADMIN', 'part:edit'),
('WAREHOUSE_ADMIN', 'borrow:create'),
('WAREHOUSE_ADMIN', 'borrow:view'),
('WAREHOUSE_ADMIN', 'borrow:return'),
('WAREHOUSE_ADMIN', 'employee:view'),
('WAREHOUSE_ADMIN', 'report:view');

-- 验证插入结果
SELECT '权限总数' as item, COUNT(*) as count FROM sys_permission
UNION ALL
SELECT '角色默认权限数', COUNT(*) FROM sys_role_permission;
