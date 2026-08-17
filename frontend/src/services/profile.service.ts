import api from './api'
import { User } from '../types'

export const profileService = {
  async get(): Promise<User> {
    const response = await api.get('/profile')
    return response.data.data
  },

  async update(data: Partial<Pick<User, 'full_name' | 'phone'>>): Promise<User> {
    const response = await api.put('/profile', data)
    return response.data.data
  },
}
