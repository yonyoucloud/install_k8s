import axios from 'axios'

// create an axios instance
const baseURL = window.CONFIG && window.CONFIG.apiHost || 'http://127.0.0.1:8081/'
const request = axios.create({
  baseURL: baseURL,
  withCredentials: true, // 跨域请求时发送 cookies
  timeout: 10000 // request timeout
})

// request interceptor
request.interceptors.request.use(
  config => {
    config.headers['jweToken'] = 'ABCD'
    return config
  },
  error => {
    // Do something with request error
    // console.log(888, error) // for debug
    Promise.reject(error)
  }
)

// response interceptor
request.interceptors.response.use(
  response => {
    return response.data
  },
  error => {
    console.log('err: ' + error) // for debug
    if (axios.isCancel(error)) {
      throw new Error('request has been cancelled')
    }
    return Promise.reject(error)
  }
)

export default request