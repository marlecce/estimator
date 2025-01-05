<template>
  <div class="room-container min-h-screen bg-white flex flex-col items-center text-gray-800">
    <!-- Header -->
    <header class="w-full bg-blue-900 py-6 shadow-md">
      <h1 class="text-center text-3xl font-bold text-white">
        Room: <span class="text-blue-300">{{ roomName }}</span>
      </h1>
    </header>

    <!-- Participants Section -->
    <main class="flex flex-col items-center w-full max-w-5xl p-6">
      <h2 class="text-xl font-semibold mb-4 text-blue-700">Participants</h2>
      <div class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-6 w-full">
        <div
          v-for="participant in participants"
          :key="participant.id"
          :class="['participant-card', participant.hasEstimated ? 'bg-blue-100' : 'bg-gray-100']"
          class="flex flex-col items-center border border-gray-200 rounded-lg p-4 shadow-sm hover:shadow-lg transition-shadow"
        >
          <div
            class="avatar w-14 h-14 rounded-full bg-blue-500 text-white text-lg font-bold flex items-center justify-center"
          >
            {{ participant.name.charAt(0).toUpperCase() }}
          </div>
          <p class="mt-2 text-lg font-medium">{{ participant.name }}</p>
          <button
            v-if="!isHost(participant.id)"
            @click="sendEstimate(participant.id)"
            class="mt-3 px-4 py-2 bg-blue-600 text-white text-sm font-semibold rounded-lg shadow hover:bg-blue-700 disabled:bg-gray-400"
            :disabled="participant.hasEstimated"
          >
            Submit Estimate
          </button>
        </div>
      </div>
    </main>

    <!-- Host Actions -->
    <footer class="mt-10 w-full bg-gray-100 py-6 flex flex-col items-center">
      <div v-if="isHost(currentParticipantId)" class="flex gap-4">
        <button
          @click="revealEstimates"
          class="px-5 py-3 bg-yellow-500 text-sm font-semibold text-white rounded-lg shadow hover:bg-yellow-600"
        >
          Reveal Estimates
        </button>
      </div>
      <button
        @click="copyLink"
        class="mt-4 px-5 py-3 bg-blue-600 text-white text-sm font-semibold rounded-lg shadow hover:bg-blue-700"
      >
        Copy Room Link
      </button>
      <p v-if="linkCopied" class="mt-3 text-sm text-green-600 font-medium">Room link copied!</p>
    </footer>
  </div>
</template>

<script>
import apiClient from "../api-client";

export default {
  name: "Room",
  data() {
    return {
      roomId: this.$route.params.roomId,
      roomName: "",
      participants: [],
      hostId: null,
      currentParticipantId: this.$route.query.participantId, 
      linkCopied: false,
      socket: null,
    };
  },
  methods: {
    async fetchRoomDetails() {
      try {
        const response = await apiClient.get(`/rooms/${this.roomId}`);
        this.roomName = response.data.name;
        this.participants = response.data.participants || [];
        this.hostId = response.data.hostId;

      } catch (error) {
        console.error("Failed to load room details:", error);
      }
    },
    isHost(participantId) {
      return participantId === this.hostId;
    },
    setupWebSocket() {
      let userId = null;  
      let isJoined = false; 

      this.socket = new WebSocket("ws://localhost:8181/ws");
      this.socket.onopen = () => {
        console.log("WebSocket connection opened");

        const setupMessage = {
            type: 'setup_connection',
            participantId: this.currentParticipantId
        };

        this.socket.send(JSON.stringify(setupMessage));
      };

      this.socket.onmessage = (event) => {
        const message = JSON.parse(event.data);

        if (message.type === "user_id_assigned") {
          this.userId = message.userId;
          console.log("Received user ID:", this.userId);

          if (!isJoined && this.userId) {
            this.socket.send(
              JSON.stringify({
                type: "participant_joined",
                roomId: this.roomId,
                participantId: userId, 
              })
            );
            isJoined = true; 
          }
      } else {
        switch (message.type) {
          case "participant_joined":
          case "participant_left":
          case "estimate_submitted":
          case "estimates_revealed":
            this.fetchRoomDetails();
            break;
          default:
            console.log("Unknown message type:", message.type);
        }
      }
    };

      this.socket.onerror = (err) => {
        console.error("WebSocket error:", err);
      };

      this.socket.onclose = () => {
        console.log("Disconnected from WebSocket");
      };
    },
    sendEstimate(participantId) {
      this.socket.send(
        JSON.stringify({
          type: "estimate_submitted",
          roomId: this.roomId,
          participantId,
        })
      );
    },
    revealEstimates() {
      this.socket.send(
        JSON.stringify({
          type: "reveal_estimates",
          roomId: this.roomId,
        })
      );
    },
    copyLink() {
      const link = `${window.location.origin}/rooms/${this.roomId}/join`;
      navigator.clipboard.writeText(link).then(() => {
        this.linkCopied = true;
        setTimeout(() => (this.linkCopied = false), 2000);
      });
    },
  },
  async created() {
    await this.fetchRoomDetails();
    this.setupWebSocket();
  },
  beforeDestroy() {
    if (this.socket) {
      this.socket.close();
    }
  },
};
</script>

<style scoped>
.room-container {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}
.participant-card {
  text-align: center;
  transition: transform 0.2s ease-in-out;
}
.participant-card:hover {
  transform: translateY(-5px);
}
</style>
