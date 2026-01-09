import { defineStore } from 'pinia'
import api from '@/api/axiosInstance'

export const useCartStore = defineStore('cart', {
    state: () => ({
        cart: [],
        customerName: '',
        paymentMethod: 'cash', // Default payment method
        loading: false,
        error: null,
        successMessage: null,
        lastTransaction: null // Store last successful transaction for receipt
    }),

    getters: {
        totalItems: (state) => {
            return state.cart.reduce((total, item) => total + item.qty, 0)
        },
        totalPrice: (state) => {
            return state.cart.reduce((total, item) => total + (item.price * item.qty), 0)
        }
    },

    actions: {
        // Add item to cart
        addToCart(product, selectedUnit, qty = 1) {
            // Check if item with same product ID AND same unit ID already exists
            const existingItemIndex = this.cart.findIndex(item =>
                item.productId === product.id && item.unitId === selectedUnit.id
            )

            if (existingItemIndex !== -1) {
                // Update qty if exists
                this.cart[existingItemIndex].qty += qty
            } else {
                // Add new item
                this.cart.push({
                    productId: product.id,
                    productName: product.name,
                    productCode: product.product_code,
                    unitId: selectedUnit.id,
                    unitName: selectedUnit.unit_name,
                    price: parseFloat(selectedUnit.sell_price), // Ensure price is number
                    originalPrice: parseFloat(selectedUnit.sell_price),
                    qty: qty,
                    imagePath: product.image_path
                })
            }
        },

        // Remove item from cart
        removeFromCart(index) {
            this.cart.splice(index, 1)
        },

        // Update item quantity
        updateQty(index, qty) {
            if (qty > 0) {
                this.cart[index].qty = qty
            }
        },

        // Clear cart
        clearCart() {
            this.cart = []
            this.customerName = ''
            this.paymentMethod = 'cash'
            this.error = null
            this.successMessage = null
        },

        // Checkout
        async checkout() {
            if (this.cart.length === 0) {
                this.error = 'Keranjang belanja kosong'
                return false
            }

            this.loading = true
            this.error = null
            this.successMessage = null

            try {
                // Format payload matching backend requirement
                const payload = {
                    customer_name: this.customerName || 'Guest',
                    payment_method: this.paymentMethod,
                    items: this.cart.map(item => ({
                        product_id: item.productId,
                        unit_id: item.unitId,
                        qty: item.qty,
                        price: item.price
                    }))
                }

                const response = await api.post('/pos/checkout', payload)

                if (response.data.status === 'success') {
                    // Calculate values BEFORE clearing
                    const savedItems = this.cart.map(item => ({
                        productId: item.productId,
                        productName: item.productName,
                        productCode: item.productCode,
                        unitId: item.unitId,
                        unitName: item.unitName,
                        price: item.price,
                        qty: item.qty
                    }))
                    const savedTotalItems = savedItems.reduce((total, item) => total + item.qty, 0)
                    const savedTotalPrice = savedItems.reduce((total, item) => total + (item.price * item.qty), 0)
                    const savedCustomerName = this.customerName || 'Guest'
                    const savedPaymentMethod = this.paymentMethod

                    // Save transaction data for receipt
                    this.lastTransaction = {
                        invoiceNo: response.data.data?.invoice_no || 'INV-' + Date.now(),
                        date: new Date().toLocaleString('id-ID'),
                        customerName: savedCustomerName,
                        paymentMethod: savedPaymentMethod,
                        items: savedItems,
                        totalItems: savedTotalItems,
                        totalPrice: savedTotalPrice,
                        cashGiven: 0, // Will be updated from POSView
                        change: 0
                    }

                    this.clearCart()
                    this.successMessage = 'Transaksi berhasil!'
                    return true
                }
            } catch (err) {
                console.error('Checkout error:', err)
                this.error = err.response?.data?.message || 'Gagal memproses transaksi'
                return false
            } finally {
                this.loading = false
            }
        }
    }
})
