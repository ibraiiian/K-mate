<template>
  <!-- This component is hidden in normal view, only shown when printing -->
  <div v-if="receiptStore.showReceipt" id="thermal-receipt" class="thermal-receipt">
    <!-- Store Header -->
    <div class="receipt-header">
      <h1 class="store-name">{{ storeInfo.name }}</h1>
      <p class="store-address">{{ storeInfo.address }}</p>
      <p class="store-phone">Telp: {{ storeInfo.phone }}</p>
    </div>

    <div class="separator double"></div>

    <!-- Transaction Info -->
    <div class="transaction-info">
      <div class="info-row">
        <span>No</span>
        <span>: {{ receipt.invoiceNo }}</span>
      </div>
      <div class="info-row">
        <span>Tgl</span>
        <span>: {{ formatDate(receipt.date) }}</span>
      </div>
      <div class="info-row">
        <span>Kasir</span>
        <span>: {{ receipt.cashierName || 'Admin' }}</span>
      </div>
      <div class="info-row" v-if="receipt.customerName">
        <span>Plg</span>
        <span>: {{ receipt.customerName }}</span>
      </div>
    </div>

    <div class="separator"></div>

    <!-- Items List -->
    <div class="items-list">
      <div v-for="(item, index) in receipt.items" :key="index" class="item">
        <div class="item-name">{{ truncate(item.productName, 24) }}</div>
        <div class="item-detail">
          <span>{{ item.qty }}x {{ formatNumber(item.price) }}</span>
          <span>{{ formatNumber(item.price * item.qty) }}</span>
        </div>
      </div>
    </div>

    <div class="separator"></div>

    <!-- Totals -->
    <div class="totals">
      <div class="total-row">
        <span>TOTAL</span>
        <span class="total-amount">{{ formatRupiah(receipt.totalPrice) }}</span>
      </div>
      <div class="total-row" v-if="receipt.cashGiven">
        <span>TUNAI</span>
        <span>{{ formatRupiah(receipt.cashGiven) }}</span>
      </div>
      <div class="total-row" v-if="receipt.change !== undefined">
        <span>KEMBALI</span>
        <span>{{ formatRupiah(receipt.change) }}</span>
      </div>
    </div>

    <div class="separator double"></div>

    <!-- Footer -->
    <div class="receipt-footer">
      <p>Terima kasih atas kunjungan Anda!</p>
      <p class="footer-note">Barang yg sudah dibeli tidak dapat</p>
      <p class="footer-note">dikembalikan kecuali ada perjanjian.</p>
    </div>

    <!-- Print Time -->
    <div class="print-time">
      <p>{{ getCurrentTime() }}</p>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useReceiptStore } from '@/stores/receiptStore'

const receiptStore = useReceiptStore()

// Store Information (configurable)
const storeInfo = {
  name: 'TB KARIB',
  address: 'Jl. Mt Haryono No. 41, Cirebon',
  phone: '021-1234567'
}

// Receipt data from store
const receipt = computed(() => receiptStore.receiptData || {})

// Helpers
const formatRupiah = (value) => {
  if (!value) return 'Rp 0'
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(value)
}

const formatNumber = (value) => {
  if (!value) return '0'
  return new Intl.NumberFormat('id-ID').format(value)
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('id-ID', {
    day: '2-digit',
    month: '2-digit',
    year: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const truncate = (str, len) => {
  if (!str) return ''
  return str.length > len ? str.slice(0, len) + '..' : str
}

const getCurrentTime = () => {
  return new Date().toLocaleString('id-ID')
}
</script>

<style>
/* ========================================
   THERMAL RECEIPT STYLES
   For 58mm (32 char) & 80mm (48 char) printers
======================================== */

.thermal-receipt {
  display: none; /* Hidden by default */
}

/* Print-only Styles */
@media print {
  /* Hide everything except receipt */
  body * {
    visibility: hidden !important;
  }
  
  #thermal-receipt,
  #thermal-receipt * {
    visibility: visible !important;
  }

  #thermal-receipt {
    display: block !important;
    position: absolute !important;
    left: 0 !important;
    top: 0 !important;
    width: 58mm !important; /* 58mm thermal paper */
    padding: 2mm !important;
    margin: 0 !important;
    background: white !important;
    font-family: 'Courier New', Courier, monospace !important;
    font-size: 10px !important;
    line-height: 1.3 !important;
    color: #000 !important;
  }

  /* For 80mm printers - uncomment this and comment the 58mm block above */
  /* 
  #thermal-receipt {
    width: 80mm !important;
    font-size: 12px !important;
  }
  */

  /* Header */
  .receipt-header {
    text-align: center;
    margin-bottom: 2mm;
  }

  .store-name {
    font-size: 14px !important;
    font-weight: bold;
    margin: 0;
    letter-spacing: 1px;
  }

  .store-address,
  .store-phone {
    font-size: 9px !important;
    margin: 1px 0;
  }

  /* Separator Lines */
  .separator {
    border-top: 1px dashed #000;
    margin: 2mm 0;
  }

  .separator.double {
    border-top: 1px solid #000;
    border-bottom: 1px solid #000;
    height: 2px;
    margin: 3mm 0;
  }

  /* Transaction Info */
  .transaction-info {
    margin-bottom: 2mm;
  }

  .info-row {
    display: flex;
    justify-content: flex-start;
    font-size: 9px !important;
  }

  .info-row span:first-child {
    width: 35px;
  }

  /* Items List */
  .items-list {
    margin: 2mm 0;
  }

  .item {
    margin-bottom: 1mm;
  }

  .item-name {
    font-weight: bold;
    font-size: 9px !important;
  }

  .item-detail {
    display: flex;
    justify-content: space-between;
    padding-left: 2mm;
    font-size: 9px !important;
  }

  /* Totals */
  .totals {
    margin: 2mm 0;
  }

  .total-row {
    display: flex;
    justify-content: space-between;
    font-size: 10px !important;
  }

  .total-row:first-child {
    font-weight: bold;
    font-size: 12px !important;
  }

  .total-amount {
    font-weight: bold;
  }

  /* Footer */
  .receipt-footer {
    text-align: center;
    margin-top: 3mm;
    font-size: 8px !important;
  }

  .footer-note {
    font-size: 7px !important;
    margin: 1px 0;
  }

  .print-time {
    text-align: center;
    margin-top: 2mm;
    font-size: 7px !important;
    color: #666;
  }

  /* Remove page margins */
  @page {
    margin: 0;
    size: 58mm auto; /* Width fixed, height auto */
  }
}
</style>
