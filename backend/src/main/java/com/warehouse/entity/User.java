package com.warehouse.entity;

import com.baomidou.mybatisplus.annotation.*;
import lombok.Data;

import java.io.Serializable;
import java.time.LocalDateTime;

/**
 * 管理员用户实体
 */
@Data
@TableName("sys_user")
public class User implements Serializable {

    private static final long serialVersionUID = 1L;

    /**
     * 用户ID
     */
    @TableId(value = "id", type = IdType.AUTO)
    private Long id;

    /**
     * 用户名
     */
    private String username;

    /**
     * 密码（BCrypt加密）
     */
    private String password;

    /**
     * 真实姓名
     */
    private String realName;

    /**
     * 角色：SUPER_ADMIN/WAREHOUSE_ADMIN
     */
    private String role;

    /**
     * 部门
     */
    private String department;

    /**
     * 手机号
     */
    private String phone;

    /**
     * 状态：0禁用/1启用
     */
    private Integer status;

    /**
     * 创建时间
     */
    @TableField(fill = FieldFill.INSERT)
    private LocalDateTime createdAt;

    /**
     * 更新时间
     */
    @TableField(fill = FieldFill.INSERT_UPDATE)
    private LocalDateTime updatedAt;
}
