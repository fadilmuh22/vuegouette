import type { TikTokItem } from '@/types'
import {
  useMutation,
  useQuery,
  type UseMutationReturnType,
  type UseQueryReturnType,
} from '@tanstack/vue-query'
import type { AxiosError } from 'axios'
import axios from 'axios'
import type { Ref } from 'vue'

export const BASE_URL = 'http://localhost:1323/api'

export const VideoQueryKey = 'video'

export type BaseResponse<T> = {
  success: boolean
  message: string
  data: T
}

export const apiClient = axios.create({
  baseURL: BASE_URL,
})

apiClient.interceptors.request.use(config => {
  const jsonToken = localStorage.getItem('access_token')
  try {
    const token =
      jsonToken && jsonToken !== 'undefined' ? jsonToken.toString() : null

    if (token) {
      config.headers['Authorization'] = token
    }
  } catch (e) {
    console.error(e)
  }
  return config
})

export const useSearchVideoQuery = (
  searchText: Ref<string>,
): UseQueryReturnType<TikTokItem[], AxiosError> => {
  return useQuery({
    queryKey: [VideoQueryKey, 'search', searchText],
    enabled: false,
    queryFn: async () => {
      if (searchText.value.length < 3) return Promise.resolve([])

      const response = await apiClient.get<BaseResponse<TikTokItem[]>>(
        `/video?keyword=${searchText.value}`,
      )
      return response.data.data
    },
  })
}

export const usePresonalizedVideoQuery = (): UseQueryReturnType<
  TikTokItem[],
  AxiosError
> => {
  return useQuery({
    queryKey: [VideoQueryKey, 'personalized'],
    queryFn: async () => {
      const response =
        await apiClient.get<BaseResponse<TikTokItem[]>>(`/video/personalized`)
      return response.data.data
    },
    staleTime: 0,
  })
}

export const useUpdaUserProfile = (): UseMutationReturnType<
  BaseResponse<never>,
  AxiosError,
  TikTokItem,
  unknown
> => {
  return useMutation({
    mutationFn: async (tiktokItem: TikTokItem) => {
      const response = await apiClient.put<BaseResponse<never>>(
        `/user/profile`,
        tiktokItem,
      )
      return response.data
    },
  })
}

export const useCreateGuestUser = (): UseMutationReturnType<
  BaseResponse<never>,
  AxiosError,
  unknown,
  unknown
> => {
  return useMutation({
    mutationFn: async () => {
      const accessToken = localStorage.getItem('access_token')
      if (accessToken !== null && accessToken !== '') {
        return Promise.resolve({
          success: false,
          message: 'Already logged in',
        } as BaseResponse<never>)
      }

      const response = await apiClient.post<BaseResponse<never>>(`/user/guest`)

      if (response.data.success) {
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        localStorage.setItem('access_token', (response.data.data as any).token)
      }

      return response.data
    },
  })
}

export const useFetchVideoLink = (
  videoUrl: string,
): UseQueryReturnType<string, AxiosError> => {
  return useQuery({
    queryKey: [VideoQueryKey, 'fetchLink', videoUrl],
    queryFn: async () => {
      const response = await apiClient.get<BaseResponse<string>>(
        '/video/fetch-video',
        {
          params: { videoUrl },
        },
      )

      // Parse the returned HTML response
      const parser = new DOMParser()
      const doc = parser.parseFromString(response.data.data, 'text/html')

      // Find the 'a' tag inside the element with id 'button-download-ready'
      const downloadLinkElement = doc.querySelector('#button-download-ready a')

      if (!downloadLinkElement) {
        throw new Error('Download link not found')
      }

      // Extract the href attribute from the found element
      const videoLink = downloadLinkElement.getAttribute('href')

      return videoLink // Return the extracted video link
    },
  })
}
