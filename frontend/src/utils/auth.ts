/**
 * Token 管理工具
 */

const TOKEN_KEY = 'access_token'
const REFRESH_TOKEN_KEY = 'refresh_token'
const USERNAME_KEY = 'username'
const IS_ADMIN_KEY = 'is_admin'
const CAPTCHA_API_KEY = 'captcha_api'

export interface TokenInfo {
  accessToken: string
  refreshToken: string
  username: string
  isAdmin: boolean
  captchaApi: string
}

/**
 * 保存 Token 信息到 localStorage
 */
export function saveTokenInfo(tokenInfo: TokenInfo) {
  localStorage.setItem(TOKEN_KEY, tokenInfo.accessToken)
  localStorage.setItem(REFRESH_TOKEN_KEY, tokenInfo.refreshToken)
  localStorage.setItem(USERNAME_KEY, tokenInfo.username)
  localStorage.setItem(IS_ADMIN_KEY, tokenInfo.isAdmin.toString())
  localStorage.setItem(CAPTCHA_API_KEY, tokenInfo.captchaApi || '')
}

/**
 * 获取 Access Token
 */
export function getAccessToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

/**
 * 获取 Refresh Token
 */
export function getRefreshToken(): string | null {
  return localStorage.getItem(REFRESH_TOKEN_KEY)
}

/**
 * 更新 Access Token
 */
export function updateAccessToken(accessToken: string) {
  localStorage.setItem(TOKEN_KEY, accessToken)
}

/**
 * 获取用户名
 */
export function getUsername(): string | null {
  return localStorage.getItem(USERNAME_KEY)
}

/**
 * 获取是否管理员
 */
export function isAdmin(): boolean {
  return localStorage.getItem(IS_ADMIN_KEY) === 'true'
}

/**
 * 获取验证码 API
 */
export function getCaptchaApi(): string | null {
  return localStorage.getItem(CAPTCHA_API_KEY)
}

/**
 * 更新验证码 API
 */
export function updateCaptchaApi(captchaApi: string) {
  localStorage.setItem(CAPTCHA_API_KEY, captchaApi)
}

/**
 * 检查是否已登录
 */
export function isAuthenticated(): boolean {
  return !!getAccessToken() && !!getRefreshToken()
}

/**
 * 清除所有 Token 信息
 */
export function clearTokenInfo() {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(REFRESH_TOKEN_KEY)
  localStorage.removeItem(USERNAME_KEY)
  localStorage.removeItem(IS_ADMIN_KEY)
  localStorage.removeItem(CAPTCHA_API_KEY)
}

/**
 * 获取完整的 Token 信息
 */
export function getTokenInfo(): TokenInfo | null {
  const accessToken = getAccessToken()
  const refreshToken = getRefreshToken()
  const username = getUsername()
  const adminStatus = isAdmin()
  const captchaApi = getCaptchaApi()

  if (!accessToken || !refreshToken || !username) {
    return null
  }

  return {
    accessToken,
    refreshToken,
    username,
    isAdmin: adminStatus,
    captchaApi: captchaApi || ''
  }
}
