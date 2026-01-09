<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
    <!-- Backdrop -->
    <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" @click="$emit('close')"></div>
    
    <!-- Modal Content -->
    <div class="relative bg-white rounded-lg shadow-xl w-[650px] max-h-[90vh] overflow-hidden border border-gray-200 animate-fadeSlideIn">
      <!-- Header -->
      <div class="bg-red-50 px-6 py-4 border-b border-red-100 flex justify-between items-center">
        <div>
          <h3 class="font-bold text-red-800 text-lg flex items-center gap-2">
            <CornerUpLeft class="w-5 h-5" />
            Proses Retur Barang
          </h3>
          <p class="text-xs text-red-600 mt-0.5">{{ transaction?.ref_no }}</p>
        </div>
        <button @click="$emit('close')" class="text-red-400 hover:text-red-600 transition">
          <X class="w-6 h-6" />
        </button>
      </div>

      <!-- Body -->
      <div class="p-6 overflow-y-auto max-h-[60vh]">
        <!-- Info -->
        <div class="bg-amber-50 border border-amber-200 rounded-lg p-3 mb-6 flex items-start gap-3">
          <AlertCircle class="w-5 h-5 text-amber-600 mt-0.5" />
          <div>
            <p class="text-sm text-amber-800 font-medium">Perhatian</p>
            <p class="text-xs text-amber-700">Pilih barang yang akan diretur dan tentukan kondisinya. Stok akan dikembalikan otomatis jika kondisi "Layak Jual".</p>
          </div>
        </div>

        <!-- Items Table -->
        <div class="border border-gray-200 rounded-lg overflow-hidden mb-6">
          <table class="w-full text-sm">
            <thead class="bg-gray-50 text-slate-500 border-b border-gray-200">
              <tr>
                <th class="px-4 py-3 text-left font-medium">Barang</th>
                <th class="px-4 py-3 text-center font-medium w-20">Beli</th>
                <th class="px-4 py-3 text-center font-medium w-24">Qty Retur</th>
                <th class="px-4 py-3 text-center font-medium w-32">Kondisi</th>
                <th class="px-4 py-3 text-right font-medium w-28">Refund</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100">
              <tr v-for="(item, idx) in returnItems" :key="idx" class="hover:bg-gray-50">
                <td class="px-4 py-3">
                  <p class="text-slate-800 font-medium">{{ item.product_name }}</p>
                  <p class="text-xs text-slate-500 font-mono">{{ formatRupiah(item.unit_price) }}/pcs</p>
                </td>
                <td class="px-4 py-3 text-center text-slate-600 font-mono">{{ item.max_qty }}</td>
                <td class="px-4 py-3 text-center">
                  <input 
                    v-model.number="item.return_qty" 
                    type="number" 
                    min="0" 
                    :max="item.max_qty"
                    class="w-16 px-2 py-1.5 border border-gray-300 rounded-md text-center text-sm focus:outline-none focus:border-blue-600 focus:ring-1 focus:ring-blue-600"
                    @input="validateQty(idx)"
                  >
                </td>
                <td class="px-4 py-3 text-center">
                  <select 
                    v-model="item.condition"
                    class="px-2 py-1.5 border border-gray-300 rounded-md text-sm focus:outline-none focus:border-blue-600"
                    :disabled="item.return_qty === 0"
                  >
                    <option value="good">Layak Jual</option>
                    <option value="damaged">Rusak</option>
                  </select>
                </td>
                <td class="px-4 py-3 text-right font-mono text-sm" :class="item.return_qty > 0 ? 'text-green-600 font-medium' : 'text-slate-400'">
                  {{ formatRupiah(item.return_qty * item.unit_price) }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Total Refund -->
        <div class="bg-slate-800 text-white p-4 rounded-lg flex justify-between items-center">
          <div>
            <p class="text-xs text-slate-400 uppercase tracking-wide">Total Refund</p>
            <p class="text-xs text-slate-400 mt-1">{{ totalReturnItems }} item akan dikembalikan</p>
          </div>
          <p class="text-2xl font-bold font-mono text-green-400">{{ formatRupiah(totalRefund) }}</p>
        </div>
      </div>

      <!-- Footer -->
      <div class="bg-gray-50 px-6 py-4 border-t border-gray-200 flex justify-between items-center">
        <p class="text-xs text-slate-500 flex items-center gap-1">
          <Info class="w-4 h-4" />
          Pastikan data sudah benar sebelum memproses
        </p>
        <div class="flex gap-3">
          <button 
            @click="$emit('close')"
            class="px-4 py-2 text-sm font-medium bg-gray-200 hover:bg-gray-300 text-slate-700 rounded-md transition"
          >
            Batal
          </button>
          <button 
            @click="processReturn"
            :disabled="totalReturnItems === 0 || processing"
            class="px-4 py-2 text-sm font-medium bg-red-600 hover:bg-red-700 text-white rounded-md transition flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <CornerUpLeft class="w-4 h-4" :class="{'animate-spin': processing}" v-if="!processing" />
            <Loader2 class="w-4 h-4 animate-spin" v-else />
            {{ processing ? 'Memproses...' : 'Proses Retur' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

import api from '@/api/axiosInstance'
import { CornerUpLeft, X, AlertCircle, Info, Loader2 } from 'lucide-vue-next'

const props = defineProps({
  transaction: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['close', 'success'])

// State
const processing = ref(false)
const returnItems = ref([])

// Initialize return items from transaction
onMounted(() => {
  if (props.transaction?.items && props.transaction.items.length > 0) {
    returnItems.value = props.transaction.items.map(item => ({
      transaction_detail_id: item.id || item.detail_id,
      product_unit_id: item.product_unit_id,
      product_name: item.product_name || item.name,
      unit_price: parseFloat(item.unit_price || item.price || item.price_at_moment) || 0,
      max_qty: item.quantity || item.qty,
      return_qty: 0,
      condition: 'good'
    }))
  }
})

// Computed
const totalRefund = computed(() => {
  return returnItems.value.reduce((sum, item) => {
    return sum + (item.return_qty * item.unit_price)
  }, 0)
})

const totalReturnItems = computed(() => {
  return returnItems.value.reduce((sum, item) => sum + item.return_qty, 0)
})

// Methods
const validateQty = (idx) => {
  const item = returnItems.value[idx]
  if (item.return_qty < 0) item.return_qty = 0
  if (item.return_qty > item.max_qty) item.return_qty = item.max_qty
}

const formatRupiah = (value) => {
  const num = parseFloat(value) || 0
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(num)
}

const processReturn = async () => {
  // Filter items with return_qty > 0
  const itemsToReturn = returnItems.value
    .filter(item => item.return_qty > 0)
    .map(item => ({
      transaction_detail_id: item.transaction_detail_id,
      qty: item.return_qty,
      condition: item.condition
    }))

  if (itemsToReturn.length === 0) {
    alert('Pilih minimal 1 barang untuk diretur')
    return
  }

  processing.value = true

  try {
    // Get the correct numeric transaction ID
    let transactionId = props.transaction.transaction_id || props.transaction.original?.id
    
    // If id is string like "txn-2", extract the number
    if (!transactionId && typeof props.transaction.id === 'string' && props.transaction.id.startsWith('txn-')) {
      transactionId = parseInt(props.transaction.id.replace('txn-', ''))
    }
    
    if (!transactionId) {
      alert('ID transaksi tidak valid')
      return
    }

    const payload = {
      transaction_id: transactionId,
      items: itemsToReturn,
      total_refund: totalRefund.value
    }

    await api.post('/returns', payload)
    
    alert('Retur berhasil diproses!')
    emit('success')
    emit('close')
  } catch (error) {
    console.error('Error processing return:', error)
    const message = error.response?.data?.message || 'Gagal memproses retur'
    alert(message)
  } finally {
    processing.value = false
  }
}
</script>

<style scoped>
/* Fade slide in animation */
@keyframes fadeSlideIn {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-fadeSlideIn {
  animation: fadeSlideIn 0.2s ease-out;
}
</style>
