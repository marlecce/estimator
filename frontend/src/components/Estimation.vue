<template>
  <div class="join-container min-h-screen bg-gradient-to-r from-blue-50 to-blue-100 flex flex-col">
    <!-- Header -->
    <header class="bg-gradient-to-r from-blue-800 to-blue-900 py-6 shadow-lg">
      <div class="container mx-auto text-center">
        <h1 class="text-3xl font-extrabold text-white">Join Room</h1>
        <p class="text-lg text-blue-200 mt-2">
          Enter your name to join the room <span class="text-yellow-300">#{{ roomId }}</span>
        </p>
      </div>
    </header>

    <!-- Form -->
    <main class="flex-1 flex items-center justify-center px-4">
      <div class="w-full max-w-md bg-white rounded-lg shadow-xl p-8">
        <h2 class="text-2xl font-semibold text-gray-800 mb-6 text-center">Welcome to the Room!</h2>
        <form @submit.prevent="joinRoom" class="space-y-6">
          <div>
            <label for="name" class="block text-sm font-medium text-gray-700">
              <span class="flex items-center">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-500 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5.121 17.804A4 4 0 117.11 9.09a8.028 8.028 0 0110.223 0 4 4 0 11-1.988 8.715M15 13a3 3 0 10-6 0 3 3 0 006 0z" />
                </svg>
                Your Name
              </span>
            </label>
            <input
              id="name"
              v-model="name"
              type="text"
              placeholder="Enter your name"
              class="w-full mt-2 p-4 border border-gray-300 rounded-lg shadow-sm text-gray-700 focus:ring-blue-500 focus:border-blue-500 transition duration-150"
            />
          </div>

          <div v-if="errorMessage" class="text-red-500 text-sm font-medium text-center">
            {{ errorMessage }}
          </div>

          <button
            type="submit"
            class="w-full px-6 py-3 bg-gradient-to-r from-blue-600 to-blue-500 text-white text-lg font-semibold rounded-lg shadow-md hover:shadow-lg transform hover:scale-105 transition duration-200 flex items-center justify-center space-x-2"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19V6l-7 7m8 8l7-7m0 0H3" />
            </svg>
            <span>Join</span>
          </button>
        </form>
      </div>
    </main>

    <!-- Success Message -->
    <div v-if="participantId" class="mt-8 text-center">
      <p class="text-lg text-green-600 font-semibold">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 inline-block mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m0 0a9 9 0 11-9 9 9 9 0 019-9z" />
        </svg>
        You've successfully joined the room!
      </p>
      <p class="mt-2 text-sm text-gray-600">
        Your participant ID is <span class="font-semibold">{{ participantId }}</span>.
      </p>
      <router-link
        :to="`/rooms/${roomId}`"
        class="mt-4 inline-block px-6 py-3 bg-gradient-to-r from-blue-600 to-blue-500 text-white text-sm font-semibold rounded-lg shadow hover:bg-blue-700 transition duration-200"
      >
        Go to the Room
      </router-link>
    </div>
  </div>
</template>

<script>
import apiClient from "../api-client";

export default {
  data() {
    return {
      name: "",
      errorMessage: "",
      participantId: null,
    };
  },
  computed: {
    roomId() {
      return this.$route.params.roomId;
    },
  },
  methods: {
    async joinRoom() {
      if (!this.name.trim()) {
        this.errorMessage = "Please enter your name.";
        return;
      }

      try {
        const response = await apiClient.post(`/rooms/${this.roomId}/join`, {
          name: this.name,
        });

        const participant = response.data;
        this.participantId = participant.id;
        this.errorMessage = "";

        this.$router.push({
          name: "room",
          params: { roomId: this.roomId },
          query: { participantId: participant.id },
        });
      } catch (error) {
        console.error("Error joining room:", error);
        this.errorMessage = "Failed to join the room. Please try again.";
      }
    },
  },
};
</script>

<style scoped>
.join-container {
  font-family: 'Inter', sans-serif;
}

header {
  margin-bottom: 0;
}

button:hover {
  box-shadow: 0px 4px 12px rgba(0, 0, 0, 0.2);
}

input:focus {
  box-shadow: 0px 0px 5px rgba(0, 0, 255, 0.2);
}
</style>
