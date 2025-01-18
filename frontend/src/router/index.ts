import { createRouter, createWebHistory } from "vue-router";
import Home from "../views/Home.vue";
import CreateRoom from "../views/CreateRoom.vue";
import Room from "../views/Room.vue";
import JoinRoom from "../views/JoinRoom.vue";

const routes = [
  {
    path: "/",
    name: "home",
    component: Home,
  },
  {
    path: "/create-room",
    name: "create-room",
    component: CreateRoom,
  },
  {
    path: "/rooms/:roomId/join",
    name: "join-room",
    component: JoinRoom,
    props: true,
  },
  {
    path: "/rooms/:roomId",
    name: "room",
    component: Room,
    props: true,
  }
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
});

export default router;
