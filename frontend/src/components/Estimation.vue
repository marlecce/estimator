<template>
  <div>
    <h1>Estimation</h1>
    <p>Select your estimate for the room.</p>
    <input v-model="estimate" placeholder="Your Estimate" />
    <button @click="sendEstimate">Submit Estimate</button>
  </div>
</template>

<script>
import { useRouter } from "vue-router";
import apiClient from "../api-client";

export default {
  data() {
    return {
      estimate: "",
    };
  },
  methods: {
    async sendEstimate() {
      const { roomId } = this.$route.params;
      
      try {
        await apiClient.post(`/rooms/${roomId}/estimation`, {
          estimate: this.estimate,
        });
        this.$router.push({ name: "room", params: { roomId } });
      } catch (error) {
        console.error("Error sending estimate:", error);
      }
    },
  },
};
</script>
