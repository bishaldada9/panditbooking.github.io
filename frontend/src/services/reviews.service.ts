import api from './api'
import { Review } from '../types'
import { responseArray, responseData } from './response'

export const reviewsService = {
  async create(data: {
    booking_id: string
    pandit_id: string
    rating: number
    comment: string
  }): Promise<Review> {
    const response = await api.post('/reviews', data)
    return responseData<Review>(response)
  },

  async getPanditReviews(panditId: string): Promise<Review[]> {
    const response = await api.get(`/pandits/${panditId}/reviews`)
    return responseArray<Review>(response)
  },

  async moderate(id: string, data: { is_visible: boolean; admin_reply?: string }): Promise<Review> {
    const response = await api.post(`/admin/reviews/${id}/moderate`, data)
    return responseData<Review>(response)
  },
}
