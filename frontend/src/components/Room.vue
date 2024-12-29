<template>
  <div class="room-container min-h-screen bg-gray-800 flex flex-col items-center">
    <h1 class="text-4xl font-bold text-white mt-8 mb-4">Room ID: {{ roomId }}</h1>

    <div class="relative w-96 h-96 rounded-full bg-green-700 flex items-center justify-center">
      <div
        v-for="(participant, index) in participants"
        :key="participant.id"
        class="absolute"
        :style="getParticipantPositionStyle(index)"
      >
        <div class="participant-card w-24 h-24 rounded-full bg-blue-600 text-center text-white flex items-center justify-center">
          <p>{{ participant.name }}</p>
        </div>
      </div>
    </div>

    <button @click="copyLink" class="px-4 py-2 bg-green-500 text-white rounded-full hover:bg-green-600 mt-6">
      Copy Room Link
    </button>
    <p v-if="linkCopied" class="mt-2 text-green-400">Link copied!</p>
  </div>
</template>

<script>
import apiClient from '../api-client';
import { io } from 'socket.io-client';

export default {
  name: 'Room',
  data() {
    return {
      roomId: this.$route.params.roomId,
      participants: [],
      linkCopied: false,
      socket: null,
    };
  },
  methods: {
    async fetchRoomDetails() {
      try {
        const response = await apiClient.get(`/rooms/${this.roomId}`);
        this.participants = response.data.participants || [];
      } catch (error) {
        console.error('Failed to load room details:', error);
      }
    },
    setupWebSocket() {
      const participantId = this.$route.query.participantId;
      if (!participantId) {
        console.error("Missing participantId for WebSocket connection");
        return;
      }

      this.socket = io("http://localhost:8181", {
        query: { participantId },
      });

      this.socket.on('connect', () => {
        console.log('Connected to WebSocket');
        this.socket.emit("join_room", { roomId: this.roomId, participantId });
      });

      this.socket.on('participant_joined', (newParticipant) => {
        console.log("New participant joined the room:", newParticipant);
        this.participants.push(newParticipant);
      });

      this.socket.on('participant_left', (participantId) => {
        this.participants = this.participants.filter((p) => p.id !== participantId);
      });

      this.socket.on('roomStateUpdated', (updatedParticipants) => {
        this.participants = updatedParticipants;
      });

      this.socket.on("disconnect", () => {
        console.log("Disconnected from WebSocket");
      });
    },
    copyLink() {
      const link = `${window.location.origin}/rooms/${this.roomId}`;
      navigator.clipboard.writeText(link).then(() => {
        this.linkCopied = true;
        setTimeout(() => (this.linkCopied = false), 2000);
      });
    },
    getParticipantPositionStyle(index) {
      const angle = (index * 360) / this.participants.length;
      const radius = 120; // Adjust radius based on table size
      return {
        transform: `translate(-50%, -50%) rotate(${angle}deg) translate(${radius}px) rotate(-${angle}deg)`,
      };
    },
  },
  async created() {
    await this.fetchRoomDetails();
    this.setupWebSocket();
  },
  beforeDestroy() {
    if (this.socket) {
      this.socket.disconnect();
    }
  },
};
</script>

<style scoped>
.room-container {
  padding: 20px;
}
.participant-card {
  position: absolute;
  text-align: center;
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 1rem;
}
</style>
