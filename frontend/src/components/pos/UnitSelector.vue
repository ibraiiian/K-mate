<template>
  <div v-if="show" class="fixed inset-0 z-50 flex items-center justify-center p-4">
    <!-- Backdrop -->
    <div 
      class="absolute inset-0 bg-black/50 backdrop-blur-sm transition-opacity"
      @click="$emit('close')"
    ></div>

    <!-- Modal Content -->
    <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden transform transition-all scale-100">
      <!-- Header -->
      <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between bg-gray-50/50">
        <h3 class="text-lg font-bold text-gray-800">Pilih Satuan</h3>
        <button 
          @click="$emit('close')"
          class="text-gray-400 hover:text-gray-600 transition-colors p-1 rounded-full hover:bg-gray-100"
        >
          <X class="w-6 h-6" />
        </button>
      </div>

      <!-- Product Info -->
      <div class="px-6 py-4 flex items-center gap-4 bg-blue-50/30">
        <div class="h-16 w-16 rounded-lg bg-gray-100 border border-gray-200 flex items-center justify-center overflow-hidden shrink-0">
           <img 
            v-if="product?.image_path" 
            :src="`http://localhost:3000${product.image_path}`" 
            class="h-full w-full object-cover"
            alt="Product"
            @error="$event.target.src = '/placeholder.png'" 
           >
           <div v-else class="text-gray-400"><Package class="w-6 h-6" /></div>
        </div>
        <div>
          <h4 class="font-bold text-gray-800 text-lg leading-tight">{{ product?.name }}</h4>
          <p class="text-gray-500 text-sm mt-1">Stok: <span class="font-medium text-blue-600">{{ product?.stock_qty }} {{ product?.base_unit }}</span></p>
        </div>
      </div>

      <!-- Units List -->
      <div class="p-6 space-y-3">
        <p class="text-sm font-medium text-gray-500 mb-2">Tersedia dalam satuan:</p>
        
        <button
          v-for="unit in product?.units"
          :key="unit.id"
          @click="selectUnit(unit)"
          class="w-full flex items-center justify-between p-4 rounded-xl border border-gray-200 hover:border-blue-500 hover:bg-blue-50/50 hover:shadow-md transition-all group text-left"
        >
          <div>
            <span class="block font-bold text-gray-800 group-hover:text-blue-700">{{ unit.unit_name }}</span>
            <span class="text-xs text-gray-400">Konversi: 1 {{ unit.unit_name }} = {{ unit.conversion_factor }} {{ product?.base_unit }}</span>
          </div>
          <div class="text-right">
            <span class="block font-bold text-blue-600 text-lg">{{ formatRupiah(unit.sell_price) }}</span>
          </div>
        </button>

        <div v-if="!product?.units || product?.units.length === 0" class="text-center py-8 text-gray-400">
            <AlertTriangle class="w-8 h-8 mb-2 mx-auto" />
            <p>Tidak ada satuan yang tersedia</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { X, Package, AlertTriangle } from 'lucide-vue-next'

const props = defineProps({
  show: Boolean,
  product: Object
})

const emit = defineEmits(['close', 'select'])

const selectUnit = (unit) => {
  emit('select', unit)
}

const formatRupiah = (value) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(value || 0)
}
</script>
