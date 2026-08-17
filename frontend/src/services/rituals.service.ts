import api from './api'
import { Ritual, RitualCategory } from '../types'
import { responseArray, responseData } from './response'

export const ritualsService = {
  async getAll(): Promise<Ritual[]> {
    const response = await api.get('/rituals')
    return responseArray<Ritual>(response)
  },

  async getById(id: string): Promise<Ritual> {
    const response = await api.get(`/rituals/${id}`)
    return responseData<Ritual>(response)
  },

  async getCategories(): Promise<RitualCategory[]> {
    const response = await api.get('/categories')
    return responseArray<RitualCategory>(response)
  },
}
