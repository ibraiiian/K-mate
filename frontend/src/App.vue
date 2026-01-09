<template>
  <div>
    <router-view />
    <!-- Global Receipt Component for Printing -->
    <ThermalReceipt />
  </div>
</template>

<script setup>
// App.vue - Main wrapper component
import ThermalReceipt from '@/components/receipt/ThermalReceipt.vue'
</script>

<style>
/* Global styles handled by Tailwind */
</style>

