package com.warehouse.controller;

import com.warehouse.common.Result;
import com.warehouse.dto.LoginRequest;
import com.warehouse.service.AuthService;
import com.warehouse.vo.LoginResponse;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.bind.annotation.*;

/**
 * 认证控制器
 */
@Tag(name = "认证管理")
@RestController
@RequestMapping("/api/auth")
public class AuthController {

    @Autowired
    private AuthService authService;

    @Operation(summary = "用户登录")
    @PostMapping("/login")
    public Result<LoginResponse> login(@Valid @RequestBody LoginRequest request) {
        LoginResponse response = authService.login(request);
        return Result.success(response);
    }

    @Operation(summary = "获取当前用户信息")
    @GetMapping("/userinfo")
    public Result<LoginResponse> getUserInfo() {
        String username = SecurityContextHolder.getContext().getAuthentication().getName();
        LoginResponse response = authService.getUserInfo(username);
        return Result.success(response);
    }

    @Operation(summary = "退出登录")
    @PostMapping("/logout")
    public Result<Void> logout() {
        // JWT是无状态的，前端删除Token即可
        // 如果需要，可以在Redis中维护Token黑名单
        return Result.success("退出登录成功", null);
    }
}
