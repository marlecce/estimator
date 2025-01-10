<template>
  <div class="room-container min-h-screen bg-white flex flex-col items-center text-gray-800">
    <!-- Header -->
    <header class="w-full bg-blue-900 py-6 shadow-md">
      <h1 class="text-center text-3xl font-bold text-white">
        Room: <span class="text-blue-300">{{ roomName }}</span>
      </h1>
      <p class="text-center text-lg text-blue-100 mt-2">
        Estimation Type:
        <span class="font-semibold text-yellow-300">{{ estimationTypeLabel }}</span>
      </p>
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
          <p v-if="participant.hasEstimated" class="mt-2 text-sm text-green-600 font-semibold">
            Estimated
          </p>
          <button
            v-else-if="!isHost(participant.id)"
            @click="openEstimateDialog(participant.id)"
            class="mt-3 px-4 py-2 bg-blue-600 text-white text-sm font-semibold rounded-lg shadow hover:bg-blue-700"
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

  <!-- Dialog per la stima -->
  <div v-if="showEstimateDialog" class="fixed inset-0 bg-gray-900 bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white rounded-lg shadow-lg p-6 w-96">
      <h2 class="text-xl font-semibold text-blue-700 mb-4">
        Submit Estimate ({{ estimationTypeLabel }})
      </h2>
      <div class="mb-4">
        <label for="estimateValue" class="block text-sm font-medium text-gray-700">Value</label>
        <input
          id="estimateValue"
          v-model.number="estimateValue"
          type="number"
          class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
        />
      </div>
      <div class="flex justify-end gap-4">
        <button
          @click="sendEstimate()"
          class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
        >
          Submit
        </button>
        <button
          @click="closeEstimateDialog"
          class="px-4 py-2 bg-gray-300 rounded-lg hover:bg-gray-400"
        >
          Cancel
        </button>
      </div>
    </div>
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
      estimationType: "",
      participants: [],
      hostId: null,
      currentParticipantId: this.$route.query.participantId, 
      linkCopied: false,
      socket: null,
      showEstimateDialog: false, 
      estimateValue: null, 
      selectedParticipantId: null,
    };
  },
  computed: {
    estimationTypeLabel() {
      switch (this.estimationType) {
        case "hours":
          return "Hours";
        case "days":
          return "Days";
        case "story_points":
          return "Story Points";
        default:
          return "Unknown";
      }
    },
  },
  methods: {
    async fetchRoomDetails() {
      try {
        const response = await apiClient.get(`/rooms/${this.roomId}`);
        console.log(response.data)
        this.roomName = response.data.name;
        this.estimationType = response.data.estimationType;
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
            this.fetchRoomDetails();
            break;
          case "estimates_revealed":
            const participant = this.participants.find(p => p.id === message.participantId);
            if (participant) {
              participant.hasEstimated = true;
            }
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
    openEstimateDialog(participantId) {
      this.selectedParticipantId = participantId;
      this.estimateValue = null;
      this.showEstimateDialog = true;
    },
    closeEstimateDialog() {
      this.selectedParticipantId = null;
      this.estimateValue = null;
      this.showEstimateDialog = false;
    },
    async sendEstimate() {
      
      if (!this.estimateValue || this.estimateValue <= 0) {
        alert("Please enter a valid estimate value.");
        return;
      }

      const participantId = this.selectedParticipantId;

      if (!participantId || !this.estimationType) {
        console.error("Participant ID or Estimation Type is missing.");
        return;
      }

      try {
        await apiClient.post(`/rooms/${this.roomId}/estimate`, {
          participant_id: participantId,
          Value: this.estimateValue,
          estimation_type: this.estimationType,
        });

        const participant = this.participants.find(
          (p) => p.id === participantId
        );
        if (participant) {
          participant.hasEstimated = true;
        }

        this.closeEstimateDialog();
      } catch (error) {
        console.error("Failed to submit estimate:", error);
      }
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
header p {
  font-size: 1.125rem;
  margin-top: 0.5rem;
}
.text-green-600 {
  color: #16a34a;
}
.font-semibold {
  font-weight: 600;
}
.fixed {
  position: fixed;
}
.inset-0 {
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}
.bg-gray-900 {
  background-color: rgba(17, 24, 39, 0.9);
}
.opacity-50 {
  opacity: 0.5;
}
.z-50 {
  z-index: 50;
}
</style>
