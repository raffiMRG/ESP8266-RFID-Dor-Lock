<script setup>
import { RouterLink, RouterView, useRoute, useRouter } from "vue-router";
import { clearToken, getToken } from "./api";

const route = useRoute();
const router = useRouter();

function logout() {
  clearToken();
  router.push("/login");
}
</script>

<template>
  <header>
    <h1>RFID Door Lock</h1>
    <nav v-if="route.name !== 'login'">
      <RouterLink to="/unlock">Buka Pintu</RouterLink>
      <RouterLink to="/logs">Log Akses</RouterLink>
      <RouterLink to="/cards">Kartu</RouterLink>
      <a v-if="getToken()" href="#" @click.prevent="logout">Logout</a>
    </nav>
  </header>
  <main>
    <RouterView />
  </main>
</template>
