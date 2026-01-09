import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useReceiptStore = defineStore('receipt', () => {
    // State
    const showReceipt = ref(false)
    const receiptData = ref(null)

    // Actions
    const setReceiptData = (data) => {
        receiptData.value = {
            invoiceNo: data.invoiceNo || data.invoice_no || '',
            date: data.date || data.transaction_date || new Date().toISOString(),
            cashierName: data.cashierName || data.cashier_name || 'Admin',
            customerName: data.customerName || data.customer_name || 'Pelanggan Umum',
            paymentMethod: data.paymentMethod || data.payment_method || 'cash',
            items: data.items || [],
            totalItems: data.totalItems || data.total_items || 0,
            totalPrice: data.totalPrice || data.total_amount || 0,
            cashGiven: data.cashGiven || data.cash_given || 0,
            change: data.change || 0
        }
        showReceipt.value = true
    }

    const clearReceipt = () => {
        receiptData.value = null
        showReceipt.value = false
    }

    const printReceipt = () => {
        if (!receiptData.value) return

        // Show the receipt element
        showReceipt.value = true

        // Wait for DOM to update, then print
        setTimeout(() => {
            window.print()
        }, 100)
    }

    return {
        showReceipt,
        receiptData,
        setReceiptData,
        clearReceipt,
        printReceipt
    }
})
