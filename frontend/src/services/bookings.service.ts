import api from './api'
import { Booking } from '../types'
import { responseArray, responseData } from './response'

export const bookingsService = {
  async create(data: {
    pandit_id: string
    ritual_id: string
    scheduled_date: string
    start_time: string
    end_time: string
    address: string
    special_notes?: string
  }): Promise<Booking> {
    const response = await api.post('/bookings', data)
    return responseData<Booking>(response)
  },

  async getMyBookings(): Promise<Booking[]> {
    const response = await api.get('/bookings')
    return responseArray<Booking>(response)
  },

  async getById(id: string): Promise<Booking> {
    const response = await api.get(`/bookings/${id}`)
    return responseData<Booking>(response)
  },

  async confirm(id: string): Promise<Booking> {
    const response = await api.put(`/bookings/${id}/confirm`)
    return responseData<Booking>(response)
  },

  async complete(id: string): Promise<Booking> {
    const response = await api.put(`/bookings/${id}/complete`)
    return responseData<Booking>(response)
  },

  async cancel(id: string, reason = 'Cancelled by user'): Promise<Booking> {
    const response = await api.put(`/bookings/${id}/cancel`, { reason })
    return responseData<Booking>(response)
  },

  async reject(id: string, reason = 'Rejected by pandit'): Promise<Booking> {
    const response = await api.put(`/bookings/${id}/reject`, { reason })
    return responseData<Booking>(response)
  },

  async getPanditBookings(): Promise<Booking[]> {
    const response = await api.get('/pandit/bookings')
    return responseArray<Booking>(response)
  },
}
