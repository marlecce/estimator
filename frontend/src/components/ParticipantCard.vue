<template>
  <div
    :class="[
      'participant-card',
      participant.hasEstimated ? 'bg-blue-100' : 'bg-gray-100',
    ]"
    class="flex flex-col items-center border border-gray-200 rounded-lg p-4 shadow-sm hover:shadow-lg transition-shadow"
  >
    <div class="avatar w-14 h-14 rounded-full bg-blue-500 text-white text-lg font-bold flex items-center justify-center">
      {{ participant.name.charAt(0).toUpperCase() }}
    </div>
    <p class="mt-2 text-lg font-medium">{{ participant.name }}</p>
    <p
      v-if="revealed && participant.revealedEstimate !== null"
      class="mt-2 text-lg font-semibold text-yellow-600"
    >
      Estimate: {{ participant.revealedEstimate }}
    </p>
    <p
      v-else-if="participant.hasEstimated && !revealed"
      class="mt-2 text-sm text-green-600 font-semibold"
    >
      Estimated
    </p>
    <button
      v-else-if="!isHost"
      @click="$emit('open-estimate-dialog', participant.id)"
      class="mt-3 px-4 py-2 bg-blue-600 text-white text-sm font-semibold rounded-lg shadow hover:bg-blue-700"
      :disabled="participant.hasEstimated"
    >
      Submit Estimate
    </button>
  </div>
</template>

<script>
export default {
  name: "ParticipantCard",
  props: {
    participant: Object,
    revealed: Boolean,
    isHost: Boolean,
  },
};
</script>
