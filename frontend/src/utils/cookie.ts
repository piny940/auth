export const fromCookie = (key: string): string | undefined => {
  if (typeof document === 'undefined') return
  const match = document.cookie.match(new RegExp('(^| )' + key + '=([^;]+)'))
  if (match) return match[2]
  return
}
