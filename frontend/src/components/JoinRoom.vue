<template>
  <div>
    <h1>Join Room</h1>
    <p>Enter your name to join the room.</p>
    <input v-model="name" placeholder="Your Name" />
    <button @click="joinRoom">Join</button>
    <div v-if="participantId">
      <p>You've joined the room! Your participant ID is {{ participantId }}.</p>
    </div>
    <div v-if="errorMessage" style="color: red;">
      <p>{{ errorMessage }}</p>
    </div>

    <div v-if="participantId">
      <p>You have successfully joined the room!</p>
      <router-link :to="`/rooms/${roomId}`" class="btn btn-primary">Go to the Room</router-link>
    </div>
  </div>
</template>

<script>
import { useRouter } from "vue-router";
import apiClient from "../api-client";

export default {
  data() {
    return {
      name: "",
      errorMessage: "",
      participantId: null,
      ws: null, 
    };
  },
  computed: {
    roomId() {
      return this.$route.params.roomId; 
    }
  },
  methods: {
    async joinRoom() {
      const roomId = this.$route.params.roomId; 
      if (!this.name.trim()) {
        this.errorMessage = "Please enter your name";
        return;
      }

      try {
        const response = await apiClient.post(`/rooms/${roomId}/join`, {
          name: this.name,
        });

        const { participant } = response.data;
        this.errorMessage = '';

        // Redirect to room page with the participantId in query params
        this.$router.push({ 
          name: "room", 
          params: { roomId }, 
          query: { id: participant.id, name: participant.name }, 
        });
      } catch (error) {
        console.error("Error joining room:", error);
        this.errorMessage = "Failed to join the room.";
      }
    }
  },
};
</script>
