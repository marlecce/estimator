import { createRouter, createWebHistory } from "vue-router";
import Home from "../components/Home.vue";
import CreateRoom from "../components/CreateRoom.vue";
import Room from "../components/Room.vue";
import Estimation from "../components/Estimation.vue";
import JoinRoom from "../components/JoinRoom.vue";

/*
const routes = [
  { path: "/", name: "Home", component: Home },
  { path: "/rooms/create", name: "CreateRoom", component: CreateRoom },
  { path: "/rooms/:roomId", name: "Room", component: Room },
  { path: "/rooms/:roomId/join", name: "JoinRoom", component: JoinRoom },
  { path: "/rooms/estimate", name: "Estimate", component: Estimation },
];
*/

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
  },
  {
    path: "/rooms/:roomId/estimation",
    name: "estimation",
    component: Estimation,
    props: true,
  },
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
});

export default router;
