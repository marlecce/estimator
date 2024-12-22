<template>
  <div>
    <h1>Create Room</h1>
    <input v-model="roomName" placeholder="Room Name" />
    <button @click="createRoom">Create</button>
    <div v-if="roomLink">
      <p>Room Created! <router-link :to="roomLink">Go to Room</router-link></p>
    </div>
  </div>
</template>

<script>
import apiClient from '../api-client';

export default {
  name: "CreateRoom",
  data() {
    return {
      roomName: "",
      roomLink: null,
    };
  },
  methods: {
    async createRoom() {
      try {
        const response = await apiClient.post('/rooms', { name: this.roomName });
        const roomId = response.data.room_id;
        this.roomLink = `/rooms/${roomId}`;
      } catch (error) {
        console.error("Failed to create room:", error);
        alert("An error occurred while creating the room.");
      }
    },
  },
};
</script>
