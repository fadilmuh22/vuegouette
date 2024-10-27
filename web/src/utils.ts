export const isAccessTokenEmpty = () => {
  const accessToken = localStorage.getItem('access_token')

  return accessToken === '' || accessToken === undefined || accessToken === null
}
