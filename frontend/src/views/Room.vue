<template>
  <div class="room-container min-h-screen bg-white flex flex-col items-center text-gray-800">
    <!-- Header -->
    <Header :roomName="roomName" :estimationTypeLabel="estimationTypeLabel" />

    <div v-if="revealed" class="mt-6 p-4 bg-green-100 text-green-800 rounded-lg shadow">
      <p class="text-lg font-semibold">Estimates have been revealed!</p>
    </div>

    <!-- Participants Section -->
    <main class="flex flex-col items-center w-full max-w-5xl p-6">
      <ParticipantsList
        :participants="participants"
        :revealed="revealed"
        :isHost="isHost"
        @openEstimateDialog="openEstimateDialog"
      />
    </main>

    <!-- Host Actions -->
    <footer class="mt-10 w-full bg-gray-100 py-6 flex flex-col items-center">
      <div v-if="isHost(currentParticipantId) && !revealed" class="flex gap-4">
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

  <!-- Dialog for the estimation -->
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
import Header from "../components/Header.vue";
import ParticipantsList from "../components/ParticipantsList.vue";

export default {
  name: "Room",
  components: {
    Header,
    ParticipantsList
  },
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
      revealed: false
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
        this.roomName = response.data.name;
        this.estimationType = response.data.estimation_type;
        this.hostId = response.data.host_id;
        this.revealed = response.data.revealed

        const estimates = response.data.estimates || [];
        this.participants = (response.data.participants || []).map((participant) => {
          const estimate = estimates.find(
            (e) => e.participant_id === participant.id
          );
          return {
            ...participant,
            revealedEstimate: estimate ? estimate.value : null,
          };
        });

        console.log(response.data)
        console.log(this.participants)

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

        switch (message.type) {
          case "user_id_assigned":
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
            break;
          case "participant_joined":
          case "participant_left":
          case "estimate_submitted":
          case "estimates_revealed":
            this.fetchRoomDetails();
            break;
         
          default:
            console.log("Unknown message type:", message.type);
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
    async revealEstimates() {
      try {
        await apiClient.post(`/rooms/${this.roomId}/reveal`, {
          participant_id: this.currentParticipantId
        });
      } catch (error) {
        console.error("Failed to reveal the estimations:", error);
      }
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
button:disabled {
  background-color: #d1d5db; /* Gray-400 */
  cursor: not-allowed;
}
.text-green-800 {
  color: #065f46;
}
.bg-green-100 {
  background-color: #d1fae5;
}
</style>
