import { defineStore } from 'pinia'
import api from '@/api/axiosInstance'

export const useAuthStore = defineStore('auth', {
    state: () => ({
        user: JSON.parse(localStorage.getItem('user')) || null,
        token: localStorage.getItem('token') || null,
        isLoading: false,
        error: null,
    }),

    getters: {
        isAuthenticated: (state) => !!state.token,
        getUser: (state) => state.user,
        getRole: (state) => state.user?.role || null,
        isOwner: (state) => state.user?.role === 'owner',
        isAdmin: (state) => state.user?.role === 'admin' || state.user?.role === 'owner',
    },

    actions: {
        async login(username, password) {
            this.isLoading = true
            this.error = null

            try {
                const response = await api.post('/auth/login', {
                    username,
                    password,
                })

                if (response.data.status === 'success') {
                    const { token, user } = response.data.data

                    // Save to state
                    this.token = token
                    this.user = user

                    // Save to localStorage
                    localStorage.setItem('token', token)
                    localStorage.setItem('user', JSON.stringify(user))

                    return { success: true, user }
                } else {
                    throw new Error(response.data.message)
                }
            } catch (error) {
                const message = error.response?.data?.message || error.message || 'Login gagal'
                this.error = message
                return { success: false, message }
            } finally {
                this.isLoading = false
            }
        },

        async register(fullName, username, password, role) {
            this.isLoading = true
            this.error = null

            try {
                const response = await api.post('/auth/register', {
                    full_name: fullName,
                    username,
                    password,
                    role,
                })

                if (response.data.status === 'success') {
                    return { success: true, data: response.data.data }
                } else {
                    throw new Error(response.data.message)
                }
            } catch (error) {
                const message = error.response?.data?.message || error.message || 'Register gagal'
                this.error = message
                return { success: false, message }
            } finally {
                this.isLoading = false
            }
        },

        logout() {
            // Clear state
            this.user = null
            this.token = null
            this.error = null

            // Clear localStorage
            localStorage.removeItem('token')
            localStorage.removeItem('user')
        },

        // Check if token is still valid
        async checkAuth() {
            if (!this.token) return false

            try {
                const response = await api.get('/me')
                if (response.data.status === 'success') {
                    return true
                }
                return false
            } catch {
                this.logout()
                return false
            }
        },
    },
})
