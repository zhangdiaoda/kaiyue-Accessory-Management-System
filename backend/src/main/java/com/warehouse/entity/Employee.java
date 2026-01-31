package com.warehouse.entity;

import com.baomidou.mybatisplus.annotation.*;
import lombok.Data;

import java.io.Serializable;
import java.time.LocalDateTime;

/**
 * 员工实体
 */
@Data
@TableName("employee")
public class Employee implements Serializable {

    private static final long serialVersionUID = 1L;

    @TableId(value = "id", type = IdType.AUTO)
    private Long id;

    /**
     * 员工工号
     */
    private String employeeNo;

    /**
     * 姓名
     */
    private String name;

    /**
     * 部门
     */
    private String department;

    /**
     * 岗位
     */
    private String position;

    /**
     * 手机号
     */
    private String phone;

    /**
     * 状态：0离职/1在职
     */
    private Integer status;

    @TableField(fill = FieldFill.INSERT)
    private LocalDateTime createdAt;

    @TableField(fill = FieldFill.INSERT_UPDATE)
    private LocalDateTime updatedAt;
}
