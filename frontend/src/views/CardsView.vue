<script setup>
import { ref, onMounted } from "vue";
import { api } from "../api";

const cards = ref([]);
const uid = ref("");
const label = ref("");
const error = ref("");

async function loadCards() {
  try {
    cards.value = await api.listCards();
  } catch (e) {
    error.value = e.message;
  }
}

async function addCard() {
  error.value = "";
  try {
    await api.createCard(uid.value.trim().toUpperCase(), label.value.trim());
    uid.value = "";
    label.value = "";
    await loadCards();
  } catch (e) {
    error.value = e.message;
  }
}

async function removeCard(id) {
  try {
    await api.deleteCard(id);
    await loadCards();
  } catch (e) {
    error.value = e.message;
  }
}

onMounted(loadCards);
</script>

<template>
  <section>
    <h2>Kartu Terdaftar</h2>

    <form @submit.prevent="addCard" class="card-form">
      <input v-model="uid" placeholder="UID (contoh: AA:BB:CC:DD)" required />
      <input v-model="label" placeholder="Label (contoh: Kartu Admin)" required />
      <button type="submit">Tambah</button>
    </form>

    <p v-if="error" class="error">{{ error }}</p>

    <table>
      <thead>
        <tr>
          <th>UID</th>
          <th>Label</th>
          <th>Ditambahkan</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="card in cards" :key="card.id">
          <td>{{ card.uid }}</td>
          <td>{{ card.label }}</td>
          <td>{{ new Date(card.created_at).toLocaleString() }}</td>
          <td><button @click="removeCard(card.id)">Hapus</button></td>
        </tr>
        <tr v-if="cards.length === 0">
          <td colspan="4">Belum ada kartu terdaftar.</td>
        </tr>
      </tbody>
    </table>
  </section>
</template>
