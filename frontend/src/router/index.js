import { createRouter, createWebHashHistory } from "vue-router";
import { getToken } from "../api";
import LoginView from "../views/LoginView.vue";
import UnlockView from "../views/UnlockView.vue";
import LogsView from "../views/LogsView.vue";
import CardsView from "../views/CardsView.vue";

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: "/", redirect: "/unlock" },
    { path: "/login", name: "login", component: LoginView },
    { path: "/unlock", name: "unlock", component: UnlockView },
    { path: "/logs", name: "logs", component: LogsView },
    { path: "/cards", name: "cards", component: CardsView },
  ],
});

router.beforeEach((to) => {
  const isAuthenticated = !!getToken();

  if (to.name !== "login" && !isAuthenticated) {
    return { name: "login" };
  }

  if (to.name === "login" && isAuthenticated) {
    return { name: "unlock" };
  }
});

export default router;
