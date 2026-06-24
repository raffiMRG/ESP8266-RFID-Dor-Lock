<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import { api, setToken } from "../api";

const router = useRouter();
const username = ref("");
const password = ref("");
const error = ref("");
const loading = ref(false);

async function submit() {
  error.value = "";
  loading.value = true;
  try {
    const { token } = await api.login(username.value, password.value);
    setToken(token);
    router.push("/unlock");
  } catch (e) {
    error.value = e.message;
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <section class="login">
    <h2>Login</h2>
    <form @submit.prevent="submit">
      <input v-model="username" placeholder="Username" autocomplete="username" required />
      <input
        v-model="password"
        type="password"
        placeholder="Password"
        autocomplete="current-password"
        required
      />
      <button type="submit" :disabled="loading">
        {{ loading ? "Memproses..." : "Login" }}
      </button>
    </form>
    <p v-if="error" class="error">{{ error }}</p>
  </section>
</template>
