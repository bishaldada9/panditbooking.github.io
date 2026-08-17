import api from './api'
import { AuthResponse, User } from '../types'

export const authService = {
  async register(data: { email: string; password: string; full_name: string; phone: string }) {
    const response = await api.post<{ success: boolean; data: AuthResponse }>('/auth/register', data)
    return response.data.data
  },

  async login(data: { email: string; password: string; device_id?: string }) {
    const response = await api.post<{ success: boolean; data: AuthResponse }>('/auth/login', data)
    return response.data.data
  },

  async logout(refreshToken: string) {
    await api.post('/auth/logout', { refresh_token: refreshToken })
  },

  async refreshToken(refreshToken: string) {
    const response = await api.post('/auth/refresh', { refresh_token: refreshToken })
    return response.data.data
  },

  async forgotPassword(email: string) {
    await api.post('/auth/forgot-password', { email })
  },

  async resetPassword(token: string, newPassword: string) {
    await api.post('/auth/reset-password', { token, new_password: newPassword })
  },

  async changePassword(oldPassword: string, newPassword: string) {
    await api.post('/auth/change-password', { old_password: oldPassword, new_password: newPassword })
  },

  async getProfile(): Promise<User> {
    const response = await api.get('/profile')
    return response.data.data
  },

  async setupMFA() {
    const response = await api.post('/auth/mfa/setup')
    return response.data.data
  },

  async verifyMFA(code: string) {
    await api.post('/auth/mfa/verify', { code })
  },

  async disableMFA() {
    await api.post('/auth/mfa/disable')
  },
}
