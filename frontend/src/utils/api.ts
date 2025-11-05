/**
 * API 请求封装
 * 支持自动携带 Token、自动刷新 Token、统一错误处理
 */

import {
  getAccessToken,
  getRefreshToken,
  updateAccessToken,
  clearTokenInfo,
} from "./auth";
import { ElMessage } from "element-plus";

// API 基础路径
const API_BASE_URL = "/api";

// 是否正在刷新 token
let isRefreshing = false;
// 刷新 token 时的等待队列
let refreshSubscribers: Array<(token: string) => void> = [];

/**
 * 添加到刷新队列
 */
function subscribeTokenRefresh(callback: (token: string) => void) {
  refreshSubscribers.push(callback);
}

/**
 * 刷新成功后通知所有等待的请求
 */
function onRefreshed(token: string) {
  refreshSubscribers.forEach((callback) => callback(token));
  refreshSubscribers = [];
}

/**
 * 刷新 Access Token
 */
async function refreshAccessToken(): Promise<string | null> {
  const refreshToken = getRefreshToken();
  if (!refreshToken) {
    return null;
  }

  try {
    const response = await fetch(`${API_BASE_URL}/auth/refresh`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        refresh_token: refreshToken,
      }),
    });

    const data = await response.json();
    if (data.message === "success" && data.data?.access_token) {
      const newAccessToken = data.data.access_token;
      updateAccessToken(newAccessToken);
      return newAccessToken;
    }

    return null;
  } catch (error) {
    console.error("刷新 Token 失败:", error);
    return null;
  }
}

/**
 * 统一的请求方法
 */
export async function request<T = any>(
  url: string,
  options: RequestInit = {},
): Promise<T> {
  // 构建完整 URL
  const fullUrl = url.startsWith("http") ? url : `${API_BASE_URL}${url}`;

  // 准备请求头
  const headers: HeadersInit = {
    "Content-Type": "application/json",
    ...(options.headers || {}),
  };

  // 添加 Access Token
  const token = getAccessToken();
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  // 发起请求
  try {
    const response = await fetch(fullUrl, {
      ...options,
      headers,
    });

    // 处理 401 未授权
    if (response.status === 401) {
      // 如果正在刷新 token，等待刷新完成
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          subscribeTokenRefresh(async (newToken: string) => {
            // 使用新 token 重试请求
            headers["Authorization"] = `Bearer ${newToken}`;
            try {
              const retryResponse = await fetch(fullUrl, {
                ...options,
                headers,
              });
              const retryData = await retryResponse.json();
              resolve(retryData);
            } catch (error) {
              reject(error);
            }
          });
        });
      }

      // 开始刷新 token
      isRefreshing = true;
      const newToken = await refreshAccessToken();
      isRefreshing = false;

      if (newToken) {
        // 通知所有等待的请求
        onRefreshed(newToken);

        // 使用新 token 重试当前请求
        headers["Authorization"] = `Bearer ${newToken}`;
        const retryResponse = await fetch(fullUrl, {
          ...options,
          headers,
        });
        return await retryResponse.json();
      } else {
        // 刷新失败，清除登录信息并跳转到登录页
        clearTokenInfo();
        ElMessage.error("登录已过期，请重新登录");
        window.location.href = "/login";
        throw new Error("Token 过期且刷新失败");
      }
    }

    // 处理其他 HTTP 错误
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }

    // 解析响应
    const data = await response.json();
    return data;
  } catch (error) {
    console.error("API 请求错误:", error);
    throw error;
  }
}

/**
 * GET 请求
 */
export async function get<T = any>(
  url: string,
  params?: Record<string, any>,
): Promise<T> {
  let fullUrl = url;
  if (params) {
    const queryString = new URLSearchParams(params).toString();
    fullUrl = `${url}?${queryString}`;
  }
  return request<T>(fullUrl, { method: "GET" });
}

/**
 * POST 请求
 */
export async function post<T = any>(url: string, data?: any): Promise<T> {
  return request<T>(url, {
    method: "POST",
    body: data ? JSON.stringify(data) : undefined,
  });
}

/**
 * PUT 请求
 */
export async function put<T = any>(url: string, data?: any): Promise<T> {
  return request<T>(url, {
    method: "PUT",
    body: data ? JSON.stringify(data) : undefined,
  });
}

/**
 * PATCH 请求
 */
export async function patch<T = any>(url: string, data?: any): Promise<T> {
  return request<T>(url, {
    method: "PATCH",
    body: data ? JSON.stringify(data) : undefined,
  });
}

/**
 * DELETE 请求
 */
export async function del<T = any>(url: string, data?: any): Promise<T> {
  return request<T>(url, {
    method: "DELETE",
    body: data ? JSON.stringify(data) : undefined,
  });
}

/**
 * 登录接口（无需 Token）
 */
export async function login(username: string, password: string) {
  const response = await fetch(`${API_BASE_URL}/auth/login`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ username, password }),
  });
  return await response.json();
}

/**
 * 登出接口
 */
export async function logout() {
  const refreshToken = getRefreshToken();
  if (!refreshToken) {
    clearTokenInfo();
    return;
  }

  try {
    await post("/auth/logout", {
      refresh_token: refreshToken,
    });
  } catch (error) {
    console.error("登出失败:", error);
  } finally {
    clearTokenInfo();
  }
}

/**
 * 注册接口（无需 Token）
 */
export async function register(username: string, password: string) {
  const response = await fetch(`${API_BASE_URL}/user/register`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ username, password }),
  });
  return await response.json();
}

// ==================== 捡漏模式相关接口 ====================

/**
 * 创建捡漏任务
 */
export async function createBargainTask(data: {
  account_id_1: number;
  account_id_2: number;
  venue_site_id: string;
  reservation_date: string;
  site_name?: string;
  reservation_time?: string;
  scan_interval: number;
  deadline?: string; // 预约截止时间（可选，格式：YYYY-MM-DD HH:mm:ss）
}) {
  return post("/bargain/create", data);
}

/**
 * 获取捡漏任务列表
 */
export async function getBargainTasks() {
  return get("/bargain/list");
}

/**
 * 获取捡漏任务详情
 */
export async function getBargainTaskDetail(taskId: string) {
  return get(`/bargain/${taskId}`);
}

/**
 * 取消捡漏任务
 */
export async function cancelBargainTask(taskId: string) {
  return del(`/bargain/${taskId}`);
}

/**
 * 获取捡漏任务日志
 */
export async function getBargainTaskLogs(taskId: string) {
  return get(`/bargain/${taskId}/logs`);
}

/**
 * 获取所有捡漏任务（管理员）
 */
export async function getAllBargainTasks() {
  return get("/bargain/all");
}
