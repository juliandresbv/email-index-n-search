import axios from 'axios'
import type { AxiosInstance, Method } from 'axios'

const apiHost = import.meta.env.VITE_SERVER_API_HOST
const apiPort = import.meta.env.VITE_SERVER_API_PORT

const httpClient: AxiosInstance = axios.create({
  baseURL: `http://${apiHost}:${apiPort}`,
  headers: {
    'Content-Type': 'application/json'
  }
})

export const makeRequest = async <R>(
  method: Method,
  endpointResource: string,
  data: any
): Promise<R | null> => {
  try {
    const response = await httpClient.request({
      method,
      url: endpointResource,
      data
    })

    return response?.data
  } catch (error) {
    console.error(error)

    return null
  }
}
