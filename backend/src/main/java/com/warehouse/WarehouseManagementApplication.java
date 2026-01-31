package com.warehouse;

import org.mybatis.spring.annotation.MapperScan;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.scheduling.annotation.EnableScheduling;

/**
 * 仓储管理系统启动类
 */
@SpringBootApplication
@MapperScan("com.warehouse.mapper")
@EnableScheduling
public class WarehouseManagementApplication {

    public static void main(String[] args) {
        SpringApplication.run(WarehouseManagementApplication.class, args);
        System.out.println("\n======================================");
        System.out.println("仓储管理系统启动成功！");
        System.out.println("API文档地址: http://localhost:8080/doc.html");
        System.out.println("======================================\n");
    }
}
