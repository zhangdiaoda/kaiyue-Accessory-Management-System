package com.warehouse.service.impl;

import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.warehouse.common.BusinessException;
import com.warehouse.dto.LoginRequest;
import com.warehouse.entity.User;
import com.warehouse.mapper.UserMapper;
import com.warehouse.security.JwtTokenProvider;
import com.warehouse.service.AuthService;
import com.warehouse.vo.LoginResponse;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

/**
 * 认证服务实现类
 */
@Slf4j
@Service
public class AuthServiceImpl implements AuthService {

    @Autowired
    private UserMapper userMapper;

    @Autowired
    private PasswordEncoder passwordEncoder;

    @Autowired
    private JwtTokenProvider jwtTokenProvider;

    @Override
    public LoginResponse login(LoginRequest request) {
        // 查询用户
        User user = userMapper.selectOne(
                new LambdaQueryWrapper<User>()
                        .eq(User::getUsername, request.getUsername())
        );

        if (user == null) {
            throw new BusinessException("用户名或密码错误");
        }

        // 验证密码
        if (!passwordEncoder.matches(request.getPassword(), user.getPassword())) {
            throw new BusinessException("用户名或密码错误");
        }

        // 检查用户状态
        if (user.getStatus() == 0) {
            throw new BusinessException("账号已被禁用");
        }

        // 生成Token
        String token = jwtTokenProvider.generateToken(user.getUsername(), user.getRole());

        log.info("用户登录成功: {}", user.getUsername());

        return new LoginResponse(
                token,
                user.getUsername(),
                user.getRealName(),
                user.getRole(),
                user.getDepartment()
        );
    }

    @Override
    public LoginResponse getUserInfo(String username) {
        User user = userMapper.selectOne(
                new LambdaQueryWrapper<User>()
                        .eq(User::getUsername, username)
        );

        if (user == null) {
            throw new BusinessException("用户不存在");
        }

        return new LoginResponse(
                null,
                user.getUsername(),
                user.getRealName(),
                user.getRole(),
                user.getDepartment()
        );
    }
}
