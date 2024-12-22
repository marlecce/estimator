import { createRouter, createWebHistory } from "vue-router";
import Home from "../components/Home.vue";
import CreateRoom from "../components/CreateRoom.vue";
import Room from "../components/Room.vue";
import Participant from "../components/Participant.vue";
import Estimation from "../components/Estimation.vue";

const routes = [
  { path: "/", name: "Home", component: Home },
  { path: "/rooms/create", name: "CreateRoom", component: CreateRoom },
  { path: "/rooms/:roomId", name: "Room", component: Room },
  { path: "/rooms/join", name: "JoinRoom", component: Participant },
  { path: "/rooms/estimate", name: "Estimate", component: Estimation },
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
});

export default router;
