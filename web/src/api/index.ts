import type { TikTokItem } from '@/types'
import { isAccessTokenEmpty } from '@/utils'
import {
  useInfiniteQuery,
  useMutation,
  useQuery,
  type UseMutationReturnType,
} from '@tanstack/vue-query'
import type { AxiosError } from 'axios'
import axios from 'axios'
import type { Ref } from 'vue'

export const BASE_URL = import.meta.env.VITE_API_BASE_URL

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

export const VIDEOS_PER_PAGE = 11

export const useSearchVideos = (searchTerm: Ref<string>) => {
  return useInfiniteQuery({
    queryKey: ['videos', 'search', searchTerm],
    initialPageParam: 1,
    enabled: false,
    queryFn: async ({ pageParam }) => {
      if (searchTerm.value.length < 3) return Promise.reject([])
      if (isAccessTokenEmpty()) return Promise.reject([])

      const response = await apiClient.get<BaseResponse<TikTokItem[]>>(
        '/video',
        {
          params: {
            keyword: searchTerm.value,
            page: pageParam,
            pageSize: VIDEOS_PER_PAGE,
          },
        },
      )

      if (response.data.data == null) {
        return Promise.reject([])
      }

      return response.data.data
    },
    getNextPageParam: (lastPage, pages) => {
      return lastPage.length > 0 ? pages.length + 1 : undefined
    },
  })
}

export const usePersonalizedVideos = () => {
  return useInfiniteQuery({
    queryKey: ['videos', 'personalized'],
    initialPageParam: 1,
    queryFn: async ({ pageParam }) => {
      if (isAccessTokenEmpty()) return Promise.reject([])

      const response = await apiClient.get<BaseResponse<TikTokItem[]>>(
        '/video/personalized',
        {
          params: {
            keyword: localStorage.getItem('keyword'),
            page: pageParam,
            pageSize: VIDEOS_PER_PAGE,
          },
        },
      )

      if (response.data.data == null) {
        return Promise.reject([])
      }

      return response.data.data
    },
    getNextPageParam: (lastPage, pages) => {
      return lastPage.length > 0 ? pages.length + 1 : undefined
    },
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

export const useCreateGuestUser = () => {
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

export const useGetUserProfileKeyword = () => {
  return useMutation({
    mutationFn: async () => {
      const response =
        await apiClient.get<BaseResponse<string>>(`/user/keyword`)

      if (response.data.data.length) {
        localStorage.setItem('keyword', response.data.data)
      } else {
        localStorage.setItem('keyword', 'trending')
      }
      return response.data.data
    },
  })
}

export const useDeleteUserProfileKeyword = () => {
  return useMutation({
    mutationFn: async (keywords: string) => {
      const response = await apiClient.post<BaseResponse<string>>(
        `/user/keyword`,
        { keywords },
      )
      if (response.data.data.length) {
        localStorage.setItem('keyword', response.data.data)
      } else {
        localStorage.setItem('keyword', 'trending')
      }
      return response.data.data
    },
  })
}

// export const useFetchVideoLink = (
//   videoUrl: string,
// ): UseQueryReturnType<string, AxiosError> => {
//   return useQuery({
//     queryKey: [VideoQueryKey, 'fetchLink', videoUrl],
//     queryFn: async () => {
//       const response = await apiClient.get<BaseResponse<string>>(
//         '/video/fetch-video',
//         {
//           params: { videoUrl },
//         },
//       )
//
//       // Parse the returned HTML response
//       const parser = new DOMParser()
//       const doc = parser.parseFromString(response.data.data, 'text/html')
//
//       // Find the 'a' tag inside the element with id 'button-download-ready'
//       const downloadLinkElement = doc.querySelector('#button-download-ready a')
//
//       if (!downloadLinkElement) {
//         throw new Error('Download link not found')
//       }
//
//       // Extract the href attribute from the found element
//       const videoLink = downloadLinkElement.getAttribute('href')
//
//       return videoLink // Return the extracted video link
//     },
//   })
// }

export const useFetchVideoLink = (videoUrl: string) => {
  return useQuery<string, AxiosError>({
    queryKey: ['fetchLink', videoUrl],
    queryFn: async () => {
      const requestBody = { language_id: 1, query: videoUrl }

      try {
        const response = await axios.post<string>(
          'https://cors-proxy.fringe.zone/https://ttsave.app/download',
          requestBody,
          { headers: { 'Content-Type': 'application/json' } },
        )

        const parser = new DOMParser()
        const doc = parser.parseFromString(response.data, 'text/html')

        const downloadLinkElement = doc.querySelector(
          '#button-download-ready a',
        )

        if (!downloadLinkElement) {
          throw new Error('Download link not found')
        }

        const videoLink = downloadLinkElement.getAttribute('href')

        if (!videoLink) {
          throw new Error('Video link is invalid')
        }

        return videoLink
      } catch (error) {
        throw new Error(
          error instanceof Error ? error.message : 'Unknown error',
        )
      }
    },
    enabled: !!videoUrl,
  })
}
