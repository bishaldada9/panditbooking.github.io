import api from './api'
import { DashboardMetrics, User, AuditLog } from '../types'
import { responseArray, responseData } from './response'

export const adminService = {
  async getDashboard(): Promise<DashboardMetrics> {
    const response = await api.get('/admin/dashboard')
    return responseData<DashboardMetrics>(response)
  },

  async getUsers(): Promise<User[]> {
    const response = await api.get('/admin/users')
    return responseArray<User>(response)
  },

  async suspendUser(id: string, reason: string): Promise<void> {
    await api.put(`/admin/users/${id}/suspend`, { reason })
  },

  async activateUser(id: string): Promise<void> {
    await api.put(`/admin/users/${id}/activate`)
  },

  async verifyPandit(id: string, status: 'pending' | 'approved' | 'rejected', notes?: string): Promise<void> {
    await api.put(`/admin/pandits/${id}/verify`, { status, notes })
  },

  async getPayments(): Promise<any[]> {
    const response = await api.get('/admin/payments')
    return responseArray<any>(response)
  },

  async getAuditLogs(): Promise<AuditLog[]> {
    const response = await api.get('/admin/audit-logs')
    return responseArray<AuditLog>(response)
  },

  async refund(data: { payment_id: string; amount?: number; reason: string }): Promise<any> {
    const response = await api.post('/admin/payments/refund', data)
    return responseData<any>(response)
  },

  async createCategory(data: { name: string; description: string; icon: string }): Promise<any> {
    const response = await api.post('/admin/categories', data)
    return responseData<any>(response)
  },

  async createRitual(data: {
    category_id: string
    name: string
    description: string
    duration: string
    base_price: number
    required_items: string
    procedure: string
  }): Promise<any> {
    const response = await api.post('/admin/rituals', data)
    return responseData<any>(response)
  },

  async getPandits(): Promise<any[]> {
    const response = await api.get('/admin/pandits')
    return responseArray<any>(response)
  },
}
