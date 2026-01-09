<template>
  <div class="h-screen flex overflow-hidden bg-[#F8FAFC]">
    <!-- Sidebar -->
    <Sidebar :isOpen="sidebarOpen" @close="sidebarOpen = false" @logout="handleLogout" />

    <!-- Main Content Wrapper -->
    <div class="flex-1 flex flex-col min-w-0 bg-[#F8FAFC] relative">
      <!-- Navbar -->
      <Navbar @toggle-sidebar="sidebarOpen = !sidebarOpen" />

      <!-- Main Content Area -->
      <main class="flex-1 overflow-y-auto p-6 md:p-8 custom-scroll relative">
        <div class="max-w-7xl mx-auto space-y-8">
          <!-- Slot for Dynamic Content -->
          <slot />
        </div>
      </main>
    </div>

    <!-- AI Chat Widget -->
    <ChatWidget />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/authStore'
import Sidebar from './Sidebar.vue'
import Navbar from './Navbar.vue'
import ChatWidget from '../common/ChatWidget.vue'

const router = useRouter()
const authStore = useAuthStore()
const sidebarOpen = ref(false)

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}
</script>

<style>
/* Global Styles for Layout */
.custom-scroll::-webkit-scrollbar { width: 6px; }
.custom-scroll::-webkit-scrollbar-track { background: transparent; }
.custom-scroll::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 10px; }

/* Font Family Global Override */
body {
  font-family: 'Poppins', sans-serif;
}
</style>
