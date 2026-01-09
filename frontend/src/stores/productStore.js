import api from '@/api/axiosInstance'
import { defineStore } from 'pinia'

export const useProductStore = defineStore('product', {
    state: () => ({
        products: [],
        loading: false,
        error: null,
        totalProducts: 0,
        currentPage: 1,
        limit: 10
    }),

    actions: {
        async fetchProducts(page = 1, limit = 10, search = '') {
            this.loading = true
            this.error = null
            try {
                const response = await api.get('/products', {
                    params: { page, limit, search }
                })

                if (response.data.status === 'success') {
                    const data = response.data.data
                    if (Array.isArray(data)) {
                        this.products = data
                        this.totalProducts = data.length
                    } else {
                        this.products = data.products
                        this.totalProducts = data.total
                        this.currentPage = data.page
                        this.limit = data.limit
                    }
                }
            } catch (err) {
                console.error('Fetch products error:', err)
                this.error = 'Gagal memuat data produk'
            } finally {
                this.loading = false
            }
        },

        async addProduct(formDataPayload) {
            this.loading = true
            this.error = null
            try {
                const response = await api.post('/products', formDataPayload, {
                    headers: {
                        'Content-Type': 'multipart/form-data'
                    }
                })

                if (response.data.status === 'success') {
                    await this.fetchProducts(this.currentPage, this.limit)
                    return { success: true, message: 'Produk berhasil ditambahkan' }
                }
            } catch (err) {
                console.error('Add product error:', err)
                return {
                    success: false,
                    message: err.response?.data?.message || 'Gagal menambahkan produk'
                }
            } finally {
                this.loading = false
            }
        },

        async deleteProduct(id) {
            console.log('deleteProduct called with id:', id)

            // Langsung delete tanpa konfirmasi untuk testing

            try {
                console.log('Sending DELETE request to /products/' + id)
                await api.delete(`/products/${id}`)
                this.products = this.products.filter(p => p.id !== id)
                console.log('Delete successful')
                return { success: true }
            } catch (err) {
                console.error('Delete error:', err)
                return {
                    success: false,
                    message: err.response?.data?.message || 'Gagal menghapus produk'
                }
            }
        },

        async updateProduct(id, data, imageFile = null) {
            this.loading = true
            this.error = null
            try {
                let response

                if (imageFile) {
                    // If there's an image, send as FormData
                    const formData = new FormData()
                    formData.append('image', imageFile)
                    formData.append('data', JSON.stringify(data))

                    response = await api.put(`/products/${id}`, formData, {
                        headers: { 'Content-Type': 'multipart/form-data' }
                    })
                } else {
                    // No image, send as JSON
                    response = await api.put(`/products/${id}`, data)
                }

                if (response.data.status === 'success') {
                    await this.fetchProducts(this.currentPage, this.limit)
                    return { success: true, message: 'Produk berhasil diupdate' }
                }
            } catch (err) {
                console.error('Update product error:', err)
                return {
                    success: false,
                    message: err.response?.data?.message || 'Gagal mengupdate produk'
                }
            } finally {
                this.loading = false
            }
        }
    }
})
