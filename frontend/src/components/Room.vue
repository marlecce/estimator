<template>
  <div class="room-container min-h-screen bg-gray-800 flex flex-col items-center">
    <h1 class="text-4xl font-bold text-white mt-8 mb-4">Room ID: {{ roomId }}</h1>

    <!-- Tavolo: Assembla il layout del tavolo -->
    <div class="relative w-96 h-96 rounded-full bg-green-700 flex items-center justify-center">
      
      <!-- Tavolo: Circonda il tavolo con i partecipanti -->
      <div v-for="(participant, index) in roomDetails.participants" :key="participant.id" class="absolute">
        <div :class="`participant-card transform rotate-${index * 360 / roomDetails.participants.length} absolute w-24 h-24 rounded-full bg-blue-600 text-center text-white flex items-center justify-center`">
          <p>{{ participant.name }}</p>

          <!-- Il pulsante per stimare solo per il partecipante connesso -->
          <button 
            v-if="participant.id === currentParticipantId && !participant.hasEstimated"
            @click="makeEstimate(participant.id)"
            class="estimate-btn absolute bottom-1 w-full py-1 bg-blue-700 text-white rounded-full mt-2"
          >
            Estimate
          </button>

          <!-- Visualizzazione della carta, coperta fino a rivelazione -->
          <div v-if="participant.hasEstimated" class="covered-card bg-gray-500 rounded-full w-16 h-16 mx-auto mt-2 flex justify-center items-center">
            <p class="text-white text-xl">{{ participant.estimatedValue || '?' }}</p>
          </div>
        </div>
      </div>

      <!-- Tavolo centrale (opzionale, decorativo) -->
      <div class="absolute w-16 h-16 bg-yellow-500 rounded-full flex items-center justify-center text-white">
        <p class="text-sm">Center</p>
      </div>
    </div>

    <!-- Copia link stanza -->
    <div class="mt-6">
      <button @click="copyLink" class="px-4 py-2 bg-green-500 text-white rounded-full hover:bg-green-600 transition">Copy Room Link</button>
      <p v-if="linkCopied" class="mt-2 text-green-400">Link copied!</p>
    </div>

    <!-- Solo il creatore può rivelare le carte -->
    <div v-if="isCreator" class="mt-6">
      <button @click="revealCards" class="px-6 py-2 bg-red-500 text-white rounded-full hover:bg-red-600 transition">Reveal All Cards</button>
    </div>
  </div>
</template>

<script>
import apiClient from '../api-client';

export default {
  name: "Room",
  data() {
    return {
      roomId: this.$route.params.roomId,
      roomDetails: null,
      linkCopied: false,
      currentParticipantId: "participant-1", // ID del partecipante connesso (simulato)
      isCreator: false, // Controlla se l'utente è il creatore della stanza
    };
  },
  async created() {
    try {
      const response = await apiClient.get(`/rooms/${this.roomId}`);
      this.roomDetails = response.data;
      this.isCreator = this.roomDetails.creatorId === this.currentParticipantId; // Solo il creatore può rivelare le carte
    } catch (error) {
      console.error("Failed to load room details:", error);
      alert("Could not load room details.");
    }
  },
  methods: {
    makeEstimate(participantId) {
      // Logica per fare la stima, inviamo il valore al server
      console.log(`Participant ${participantId} making an estimate`);
      apiClient.post(`/rooms/${this.roomId}/estimate`, {
        participantId,
        value: 5 // La stima, puoi cambiarla in base alla logica dell'interfaccia
      }).then(response => {
        // Dopo aver inviato la stima, aggiorniamo lo stato
        const participant = this.roomDetails.participants.find(p => p.id === participantId);
        if (participant) {
          participant.hasEstimated = true;
          participant.estimatedValue = 5; // In un'applicazione vera, il valore sarebbe dinamico
        }
      });
    },
    copyLink() {
      const link = `${window.location.origin}/rooms/${this.roomId}`;
      navigator.clipboard.writeText(link).then(() => {
        this.linkCopied = true;
        setTimeout(() => {
          this.linkCopied = false;
        }, 2000);
      });
    },
    revealCards() {
      // Solo il creatore può rivelare le carte
      console.log("Revealing all cards");
      apiClient.post(`/rooms/${this.roomId}/reveal`).then(() => {
        // Dopo la rivelazione, le carte sono visibili
        this.roomDetails.participants.forEach(participant => {
          participant.hasEstimated = true;
        });
      });
    }
  }
};
</script>

<style scoped>
.room-container {
  padding: 20px;
}

.participant-card {
  position: absolute;
  width: 6rem;
  height: 6rem;
  text-align: center;
  border-radius: 50%;
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 1rem;
}

.estimate-btn {
  width: 100%;
  font-size: 0.875rem;
  transition: background-color 0.3s ease;
}

.estimate-btn:hover {
  background-color: #4c6b7b;
}

.covered-card {
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 1.25rem;
}

.table-container {
  position: relative;
  width: 24rem;
  height: 24rem;
  background-color: #4b8f29;
  border-radius: 50%;
  margin: 20px;
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.5);
}

.table-center {
  position: absolute;
  width: 4rem;
  height: 4rem;
  background-color: #f8b400;
  border-radius: 50%;
  display: flex;
  justify-content: center;
  align-items: center;
}

button {
  cursor: pointer;
  font-size: 1rem;
}

button:hover {
  opacity: 0.9;
}
</style>
