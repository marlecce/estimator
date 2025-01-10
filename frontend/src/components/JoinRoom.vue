<template>
  <div class="join-container min-h-screen bg-gray-100 flex flex-col">
    <!-- Header -->
    <header class="bg-blue-900 py-4 shadow-md">
      <div class="container mx-auto text-center">
        <h1 class="text-2xl font-bold text-white">Join Room</h1>
        <p class="text-lg text-blue-100 mt-1">
          Enter your name to join the room <span class="text-yellow-300">#{{ roomId }}</span>
        </p>
      </div>
    </header>

    <!-- Form -->
    <main class="flex-1 flex items-center justify-center">
      <div class="w-full max-w-md bg-white rounded-lg shadow-md p-6">
        <form @submit.prevent="joinRoom" class="space-y-4">
          <div>
            <label for="name" class="block text-sm font-medium text-gray-700">
              Your Name
            </label>
            <input
              id="name"
              v-model="name"
              type="text"
              placeholder="Enter your name"
              class="w-full mt-1 p-3 border border-gray-300 rounded-md text-gray-800 focus:ring-blue-500 focus:border-blue-500"
            />
          </div>

          <div v-if="errorMessage" class="text-red-600 text-sm">
            {{ errorMessage }}
          </div>

          <button
            type="submit"
            class="w-full px-5 py-3 bg-blue-600 text-white text-sm font-semibold rounded-md shadow hover:bg-blue-700"
          >
            Join
          </button>
        </form>
      </div>
    </main>

    <!-- Success Message -->
    <div v-if="participantId" class="mt-6 text-center">
      <p class="text-lg text-green-600 font-semibold">You've successfully joined the room!</p>
      <p class="mt-2 text-sm text-gray-600">
        Your participant ID is <span class="font-semibold">{{ participantId }}</span>.
      </p>
      <router-link
        :to="`/rooms/${roomId}`"
        class="mt-4 inline-block px-5 py-3 bg-blue-600 text-white text-sm font-semibold rounded-lg shadow hover:bg-blue-700"
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

        const { participant } = response.data;
        this.participantId = participant.id;
        this.errorMessage = "";

        // Redirect to room page
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
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

header {
  margin-bottom: 0;
}

input {
  box-shadow: 0px 2px 4px rgba(0, 0, 0, 0.1);
}
</style>
