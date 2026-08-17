import api from './api'
import { Pandit } from '../types'
import { responseArray, responseData } from './response'

export const panditsService = {
  async getAll(): Promise<Pandit[]> {
    const response = await api.get('/pandits')
    return responseArray<Pandit>(response)
  },

  async getById(id: string): Promise<Pandit> {
    const response = await api.get(`/pandits/${id}`)
    return responseData<Pandit>(response)
  },

  async getAvailability(id: string): Promise<any[]> {
    const response = await api.get(`/pandits/${id}/availability`)
    return responseArray<any>(response)
  },

  async register(data: {
    bio: string
    experience_years: number
    specialization: string
    languages: string[]
    base_price: number
    service_area: string
  }): Promise<any> {
    const response = await api.post('/pandit/register', data)
    return responseData<Pandit>(response)
  },

  async getProfile(): Promise<Pandit> {
    const response = await api.get('/pandit/profile')
    return responseData<Pandit>(response)
  },

  async updateAvailability(data: { date: string; start_time: string; end_time: string; is_booked?: boolean }): Promise<any> {
    const response = await api.post('/pandit/availability', data)
    return responseData<any>(response)
  },
}
