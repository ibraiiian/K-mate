import { useReceiptStore } from '@/stores/receiptStore'
import { nextTick } from 'vue'

/**
 * Composable for handling thermal receipt printing
 * 
 * Usage:
 * const { printTransaction, printFromData } = usePrint()
 * printTransaction(transactionData)
 */
export function usePrint() {
    const receiptStore = useReceiptStore()

    /**
     * Print a transaction receipt
     * @param {Object} transactionData - Transaction data from checkout or history
     */
    const printTransaction = async (transactionData) => {
        if (!transactionData) {
            console.error('No transaction data provided')
            return false
        }

        // Set the receipt data
        receiptStore.setReceiptData(transactionData)

        // Wait for DOM to update
        await nextTick()

        // Trigger print
        receiptStore.printReceipt()

        return true
    }

    /**
     * Prepare receipt data without printing (for preview)
     * @param {Object} transactionData - Transaction data
     */
    const prepareReceipt = (transactionData) => {
        receiptStore.setReceiptData(transactionData)
    }

    /**
     * Clear the receipt from display
     */
    const clearReceipt = () => {
        receiptStore.clearReceipt()
    }

    /**
     * Print receipt from existing stored data
     */
    const printStoredReceipt = () => {
        receiptStore.printReceipt()
    }

    return {
        printTransaction,
        prepareReceipt,
        clearReceipt,
        printStoredReceipt
    }
}
