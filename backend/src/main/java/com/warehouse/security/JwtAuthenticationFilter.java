package com.warehouse.security;

import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.web.authentication.WebAuthenticationDetailsSource;
import org.springframework.stereotype.Component;
import org.springframework.util.StringUtils;
import org.springframework.web.filter.OncePerRequestFilter;

import java.io.IOException;
import java.util.Collections;

/**
 * JWT认证过滤器
 */
@Slf4j
@Component
public class JwtAuthenticationFilter extends OncePerRequestFilter {

    @Autowired
    private JwtTokenProvider jwtTokenProvider;

    @Override
    protected void doFilterInternal(HttpServletRequest request,
                                    HttpServletResponse response,
                                    FilterChain filterChain) throws ServletException, IOException {
        try {
            String jwt = getJwtFromRequest(request);

            if (StringUtils.hasText(jwt) && jwtTokenProvider.validateToken(jwt)) {
                String username = jwtTokenProvider.getUsernameFromToken(jwt);
                String role = jwtTokenProvider.getRoleFromToken(jwt);

                // 创建认证对象
                UsernamePasswordAuthenticationToken authentication =
                        new UsernamePasswordAuthenticationToken(
                                username,
                                null,
                                Collections.singletonList(new SimpleGrantedAuthority(role))
                        );
                authentication.setDetails(new WebAuthenticationDetailsSource().buildDetails(request));

                // 设置到SecurityContext
                SecurityContextHolder.getContext().setAuthentication(authentication);
                
                log.debug("设置用户认证: {}, 角色: {}", username, role);
            }
        } catch (Exception ex) {
            log.error("无法设置用户认证", ex);
        }

        filterChain.doFilter(request, response);
    }

    /**
     * 从请求头中获取JWT Token
     */
    private String getJwtFromRequest(HttpServletRequest request) {
        String bearerToken = request.getHeader("Authorization");
        if (StringUtils.hasText(bearerToken) && bearerToken.startsWith("Bearer ")) {
            return bearerToken.substring(7);
        }
        return null;
    }
}
