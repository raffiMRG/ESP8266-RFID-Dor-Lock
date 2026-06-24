<script setup>
import { ref, onMounted, onUnmounted } from "vue";
import { api } from "../api";

const logs = ref([]);
const error = ref("");
let pollHandle = null;

async function loadLogs() {
  try {
    logs.value = await api.listLogs();
    error.value = "";
  } catch (e) {
    error.value = e.message;
  }
}

onMounted(() => {
  loadLogs();
  pollHandle = setInterval(loadLogs, 5000);
});

onUnmounted(() => clearInterval(pollHandle));
</script>

<template>
  <section>
    <h2>Log Akses</h2>
    <p v-if="error" class="error">{{ error }}</p>

    <table>
      <thead>
        <tr>
          <th>Waktu</th>
          <th>UID</th>
          <th>Label</th>
          <th>Status</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="log in logs" :key="log.id">
          <td>{{ new Date(log.created_at).toLocaleString() }}</td>
          <td>{{ log.uid }}</td>
          <td>{{ log.label || "-" }}</td>
          <td>
            <span :class="log.granted ? 'granted' : 'denied'">
              {{ log.granted ? "Diterima" : "Ditolak" }}
            </span>
          </td>
        </tr>
        <tr v-if="logs.length === 0">
          <td colspan="4">Belum ada log akses.</td>
        </tr>
      </tbody>
    </table>
  </section>
</template>
