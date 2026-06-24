<script setup>
import { ref } from "vue";
import { api } from "../api";

const status = ref("");
const loading = ref(false);

async function unlock() {
  status.value = "";
  loading.value = true;
  try {
    await api.unlock();
    status.value = "Perintah buka pintu terkirim.";
  } catch (e) {
    status.value = e.message;
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <section class="unlock">
    <h2>Buka Pintu</h2>
    <p>Alternatif membuka pintu tanpa kartu RFID.</p>
    <button class="unlock-btn" :disabled="loading" @click="unlock">
      {{ loading ? "Mengirim..." : "Buka Pintu" }}
    </button>
    <p v-if="status">{{ status }}</p>
  </section>
</template>
