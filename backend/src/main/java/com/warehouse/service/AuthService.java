package com.warehouse.service;

import com.warehouse.dto.LoginRequest;
import com.warehouse.vo.LoginResponse;

/**
 * 认证服务接口
 */
public interface AuthService {

    /**
     * 用户登录
     */
    LoginResponse login(LoginRequest request);

    /**
     * 获取用户信息
     */
    LoginResponse getUserInfo(String username);
}
