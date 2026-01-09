<template>
  <header class="bg-white border-b border-gray-200 h-20 px-8 flex items-center justify-between sticky top-0 z-10 font-sans">
    <div class="flex items-center gap-4">
      <button @click="$emit('toggleSidebar')" class="md:hidden text-slate-500 text-2xl">
        <Menu class="w-6 h-6" />
      </button>
      <div class="hidden md:flex flex-col">
        <span class="text-xs font-bold text-slate-400 uppercase">Hari Ini</span>
        <div class="flex items-center font-bold text-slate-700 text-sm">
          <Calendar class="w-4 h-4 mr-2 text-blue-600" />
          <span>{{ currentDate }}</span>
        </div>
      </div>
    </div>
    <div class="flex items-center gap-4">
      <button class="relative p-2 text-slate-400 hover:text-blue-600 transition">
        <Bell class="w-6 h-6" />
        <span class="absolute top-2 right-2 h-2 w-2 bg-red-500 rounded-full border border-white"></span>
      </button>
      <router-link to="/pos" class="bg-blue-600 text-white rounded-lg flex items-center gap-2 px-6 py-2.5 text-sm font-bold shadow-sm hover:bg-blue-700 transition active:translate-y-0.5">
        <PlusCircle class="w-5 h-5" />
        <span>Transaksi Baru</span>
      </router-link>
    </div>
  </header>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Menu, Calendar, Bell, PlusCircle } from 'lucide-vue-next'

const currentDate = ref('')

onMounted(() => {
  const options = { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' }
  currentDate.value = new Date().toLocaleDateString('id-ID', options)
})
</script>

<style scoped>
/* Optional specific overrides if needed */
</style>
