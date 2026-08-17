import api from './api'
import { Notification } from '../types'
import { responseArray, responseData } from './response'

export const notificationsService = {
  async getAll(): Promise<Notification[]> {
    const response = await api.get('/notifications')
    return responseArray<Notification>(response)
  },

  async markRead(id: string): Promise<void> {
    await api.put(`/notifications/${id}/read`)
  },

  async markAllRead(): Promise<void> {
    await api.put('/notifications/read-all')
  },

  async getUnreadCount(): Promise<number> {
    const response = await api.get('/notifications/unread-count')
    return responseData<{ count: number }>(response)?.count ?? 0
  },
}
