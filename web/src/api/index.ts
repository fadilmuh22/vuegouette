import type { TikTokItem } from '@/types'
import {
  useMutation,
  useQuery,
  type UseMutationReturnType,
  type UseQueryReturnType,
} from '@tanstack/vue-query'
import type { AxiosError } from 'axios'
import axios from 'axios'

export const BASE_URL = 'http://localhost:1323/api'
export const VideoQueryKey = 'video'

export type BaseResponse<T> = {
  success: boolean
  message: string
  data: T
}

export const useSearchVideoQuery = (
  searchText: string,
): UseQueryReturnType<TikTokItem[], AxiosError> => {
  return useQuery({
    queryKey: [VideoQueryKey, searchText],
    enabled: false,
    queryFn: async () => {
      const response = await axios.get<BaseResponse<TikTokItem[]>>(
        `${BASE_URL}/video?=${searchText}?lang=en`,
      )
      return response.data.data
    },
  })
}

export const useUpdaUserProfile = (): UseMutationReturnType<
  never,
  AxiosError,
  TikTokItem,
  never
> => {
  return useMutation({
    mutationFn: async () => {
      const response = await axios.put<BaseResponse<never>>(
        `${BASE_URL}/profile`,
      )
      return response.data.data
    },
  })
}
