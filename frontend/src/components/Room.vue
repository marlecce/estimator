<template>
  <div>
    <h1>Room ID: {{ roomId }}</h1>
    <div v-if="roomDetails">
      <h2>Room Name: {{ roomDetails.name }}</h2>
      <h3>Participants:</h3>
      <ul>
        <li v-for="participant in roomDetails.participants" :key="participant.id">
          {{ participant.name }}
        </li>
      </ul>
      <h3>Status: {{ roomDetails.status }}</h3>
    </div>
    <div v-else>
      <p>Loading room details...</p>
    </div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "Room",
  data() {
    return {
      roomId: this.$route.params.roomId, 
      roomDetails: null,
    };
  },
  async created() {
    try {
      const response = await axios.get(`/api/rooms/${this.roomId}`);
      this.roomDetails = response.data;
    } catch (error) {
      console.error("Failed to load room details:", error);
      alert("Could not load room details.");
    }
  },
};
</script>
