import api from './api'
import { Payment } from '../types'
import { responseArray, responseData } from './response'

export const paymentsService = {
  async initiate(data: {
    booking_id: string
    gateway: 'esewa' | 'khalti' | 'cash'
    amount: number
  }): Promise<Payment> {
    const response = await api.post('/payments/initiate', data)
    return responseData<Payment>(response)
  },

  async verify(data: { payment_id: string; gateway_ref_id?: string }): Promise<Payment> {
    const response = await api.post('/payments/verify', {
      transaction_id: data.payment_id,
      gateway_ref_id: data.gateway_ref_id,
      status: 'completed',
    })
    return responseData<Payment>(response)
  },

  async getById(id: string): Promise<Payment> {
    const response = await api.get(`/payments/${id}`)
    return responseData<Payment>(response)
  },

  async getByBooking(bookingId: string): Promise<Payment> {
    const response = await api.get(`/payments/booking/${bookingId}`)
    return responseData<Payment>(response)
  },

  async getAll(): Promise<Payment[]> {
    const response = await api.get('/admin/payments')
    return responseArray<Payment>(response)
  },

  async refund(data: { payment_id: string; amount?: number; reason: string }): Promise<any> {
    const response = await api.post('/admin/payments/refund', data)
    return responseData<any>(response)
  },
}
