<template>
  <div>
    <!-- Mobile Overlay -->
    <div v-if="isOpen" @click="$emit('close')" class="fixed inset-0 bg-black/50 z-20 md:hidden"></div>

    <!-- Sidebar -->
    <aside 
      class="fixed inset-y-0 left-0 w-72 bg-white border-r border-gray-200 flex flex-col z-30 font-sans transform transition-transform duration-200 ease-in-out md:translate-x-0 md:static md:h-full"
      :class="isOpen ? 'translate-x-0' : '-translate-x-full'"
    >
    <div class="h-20 flex items-center px-8 border-b border-gray-100">
      <div class="flex items-center gap-3">
        <div class="bg-blue-600 text-white p-2 rounded-lg flex items-center justify-center shadow-sm">
          <Home class="w-5 h-5" />
        </div>
        <div>
          <h1 class="text-xl font-bold text-slate-800 tracking-tight leading-none">K-Mate</h1>
          <p class="text-[10px] font-medium text-slate-500 uppercase tracking-wide mt-1">Sistem Manajemen</p>
        </div>
      </div>
    </div>
    
    <nav class="flex-1 overflow-y-auto py-6 space-y-1">
      <div class="px-6 mb-2 text-xs font-bold text-slate-400 uppercase tracking-wider">Menu Utama</div>
      
      <router-link to="/dashboard" class="sidebar-link flex items-center px-4 py-3 text-sm" active-class="active">
        <LayoutDashboard class="w-5 h-5 mr-3" /> Dashboard
      </router-link>
      
      <router-link to="/inventory" class="sidebar-link flex items-center px-4 py-3 text-sm" active-class="active">
        <Package class="w-5 h-5 mr-3" /> Inventory
      </router-link>
      
      <router-link to="/pos" class="sidebar-link flex items-center px-4 py-3 text-sm" active-class="active">
        <ShoppingCart class="w-5 h-5 mr-3" /> Kasir (POS)
      </router-link>
      
      <router-link to="/transactions" class="sidebar-link flex items-center px-4 py-3 text-sm" active-class="active">
        <Receipt class="w-5 h-5 mr-3" /> Transaksi
      </router-link>

      <div class="px-6 mb-2 mt-6 text-xs font-bold text-slate-400 uppercase tracking-wider">Admin</div>
      
      <router-link to="/reports" class="sidebar-link flex items-center px-4 py-3 text-sm" active-class="active">
        <Archive class="w-5 h-5 mr-3" /> Laporan Keuangan
      </router-link>

      <router-link v-if="user?.role === 'owner'" to="/audit-logs" class="sidebar-link flex items-center px-4 py-3 text-sm" active-class="active">
        <ClipboardList class="w-5 h-5 mr-3" /> Audit Log
      </router-link>
    </nav>
    
    <div class="p-4 border-t border-gray-100">
      <div class="flex items-center gap-3 px-2">
        <div class="h-10 w-10 rounded-full bg-slate-800 text-white flex items-center justify-center font-bold">
          {{ userInitials }}
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-bold text-slate-800 truncate">{{ user?.full_name || 'User' }}</p>
          <p class="text-xs text-slate-500 truncate capitalize">{{ user?.role || 'Staff' }}</p>
        </div>
        <button @click="$emit('logout')" class="text-slate-400 hover:text-red-500 transition">
          <LogOut class="w-5 h-5" />
        </button>
      </div>
    </div>
  </aside>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useAuthStore } from '@/stores/authStore'
import { 
  LayoutDashboard, 
  ShoppingCart, 
  Package, 
  Receipt, 
  Archive, 
  LogOut, 
  Home, 
  ClipboardList 
} from 'lucide-vue-next'

defineProps({
  isOpen: Boolean
})

const authStore = useAuthStore()
const user = computed(() => authStore.user)

const userInitials = computed(() => {
  if (user.value?.full_name) {
    return user.value.full_name.split(' ').map(n => n[0]).join('').substring(0, 2).toUpperCase()
  }
  return 'U'
})
</script>

<style scoped>
.sidebar-link {
    border-radius: 8px;
    margin: 4px 16px;
    transition: all 0.2s ease;
    font-weight: 500;
    color: #64748B;
    text-decoration: none;
}
.sidebar-link.active {
    background-color: #EFF6FF; 
    color: #2563EB; 
    font-weight: 600;
}
.sidebar-link:hover:not(.active) {
    background-color: #F1F5F9;
    color: #1E293B;
}
</style>
