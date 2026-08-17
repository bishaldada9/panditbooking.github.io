export function responseData<T>(response: { data?: { data?: T } }): T {
  return response.data?.data as T
}

export function responseArray<T>(response: { data?: { data?: T[] | null } }): T[] {
  const data = response.data?.data
  return Array.isArray(data) ? data : []
}
