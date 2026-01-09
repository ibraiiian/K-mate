import axios from 'axios'

// Base URL - sesuaikan dengan port backend
const API_BASE_URL = 'http://localhost:3000/api/v1'

// Create axios instance
const api = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
    timeout: 10000, // 10 seconds
})

// Request interceptor - Auto inject token
api.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('token')
        if (token) {
            config.headers.Authorization = `Bearer ${token}`
        }
        return config
    },
    (error) => {
        return Promise.reject(error)
    }
)

// Response interceptor - Handle errors
api.interceptors.response.use(
    (response) => {
        return response
    },
    (error) => {
        // Handle 401 Unauthorized
        if (error.response && error.response.status === 401) {
            // Token expired or invalid
            localStorage.removeItem('token')
            localStorage.removeItem('user')
            window.location.href = '/login'
        }
        return Promise.reject(error)
    }
)

export default api
