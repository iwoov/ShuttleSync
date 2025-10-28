import {
  createRouter,
  createWebHistory,
  type RouteRecordRaw,
} from "vue-router";
import { ElMessage } from "element-plus";
import Login from "@/views/Login.vue";
import MyReservation from "@/views/MyReservation.vue";
import VenueReservation from "@/views/VenueReservation.vue";
import UserManagement from "@/views/UserManagement.vue";
import MyAccount from "@/views/MyAccount.vue";
import AllReservations from "@/views/AllReservations.vue";
import { isAuthenticated, isAdmin as checkIsAdmin } from "@/utils/auth";

const routes: RouteRecordRaw[] = [
  {
    path: "/",
    redirect: "/login",
  },
  {
    path: "/login",
    name: "Login",
    component: Login,
  },
  {
    path: "/my-reservation",
    name: "MyReservation",
    component: MyReservation,
    meta: { requiresAuth: true },
  },
  {
    path: "/venue-reservation",
    name: "VenueReservation",
    component: VenueReservation,
    meta: { requiresAuth: true },
  },
  {
    path: "/user-management",
    name: "UserManagement",
    component: UserManagement,
    meta: { requiresAuth: true, requiresAdmin: true },
  },
  {
    path: "/all-reservations",
    name: "AllReservations",
    component: AllReservations,
    meta: { requiresAuth: true, requiresAdmin: true },
  },
  {
    path: "/my-account",
    name: "MyAccount",
    component: MyAccount,
    meta: { requiresAuth: true },
  },
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
});

router.beforeEach((to: any, from: any, next: any) => {
  // 检查是否需要认证
  if (to.meta.requiresAuth) {
    if (!isAuthenticated()) {
      ElMessage.warning("请先登录");
      next("/login");
      return;
    }
  }

  // 检查是否需要管理员权限
  if (to.meta.requiresAdmin) {
    if (!checkIsAdmin()) {
      ElMessage.error("需要管理员权限");
      next("/my-reservation");
      return;
    }
  }

  next();
});

export default router;
