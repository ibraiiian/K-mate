<template>
  <div class="flex items-center gap-3 p-3 bg-white rounded-xl border border-gray-100 shadow-sm hover:shadow-md transition-all group">
    <!-- Image -->
    <div class="h-12 w-12 rounded-lg bg-gray-50 border border-gray-100 flex items-center justify-center overflow-hidden shrink-0">
        <img 
            v-if="item.imagePath" 
            :src="`http://localhost:3000${item.imagePath}`" 
            class="h-full w-full object-cover"
            alt="Product"
            @error="$event.target.src = '/placeholder.png'" 
        >
        <i v-else class="text-gray-300"><Package class="w-5 h-5" /></i>
    </div>

    <!-- Info -->
    <div class="flex-1 min-w-0">
      <h4 class="font-bold text-gray-800 text-sm truncate">{{ item.productName }}</h4>
      <div class="flex items-center gap-2 mt-0.5">
         <span class="text-xs font-medium text-slate-500 bg-slate-100 px-1.5 py-0.5 rounded">{{ item.unitName }}</span>
         <span class="text-xs text-blue-600 font-bold">{{ formatRupiah(item.price) }}</span>
      </div>
    </div>

    <!-- Qty Controls -->
    <div class="flex items-center gap-2 bg-gray-50 rounded-lg p-1">
      <button 
        @click="$emit('update-qty', item.qty - 1)"
        class="w-6 h-6 flex items-center justify-center rounded-md bg-white border border-gray-200 text-gray-600 hover:bg-red-50 hover:text-red-500 hover:border-red-200 transition-colors shadow-sm disabled:opacity-50"
        :disabled="item.qty <= 1"
      >
        <Minus class="w-3 h-3" />
      </button>
      <span class="text-sm font-bold w-4 text-center text-gray-800">{{ item.qty }}</span>
      <button 
        @click="$emit('update-qty', item.qty + 1)"
        class="w-6 h-6 flex items-center justify-center rounded-md bg-white border border-gray-200 text-gray-600 hover:bg-green-50 hover:text-green-500 hover:border-green-200 transition-colors shadow-sm"
      >
        <Plus class="w-3 h-3" />
      </button>
    </div>

    <!-- Remove -->
    <button 
      @click="$emit('remove')"
      class="p-2 text-gray-400 hover:text-red-500 hover:bg-red-50 rounded-lg transition-colors"
    >
      <Trash2 class="w-5 h-5" />
    </button>
  </div>
</template>

<script setup>
const props = defineProps({
  item: Object
})

import { Package, Minus, Plus, Trash2 } from 'lucide-vue-next'

const emit = defineEmits(['update-qty', 'remove'])

const formatRupiah = (value) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(value || 0)
}
</script>
