const baseURL = window.CONFIG && window.CONFIG.apiHost || 'http://127.0.0.1:8081/'

export function stream(url) {
  return new EventSource(baseURL + url)
}