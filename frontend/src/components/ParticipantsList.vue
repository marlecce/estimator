<template>
      <h2 class="text-xl font-semibold mb-4 text-blue-700">Participants</h2>
      <div class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-6 w-full">
        <ParticipantCard
          v-for="participant in participants"
          :key="participant.id"
          :participant="participant"
          :revealed="revealed"
          :is-host="isHost(participant.id)"
          @open-estimate-dialog="openEstimateDialog"
        />
      </div>
  </template>
  
  <script>
  import ParticipantCard from './ParticipantCard.vue';
  
  export default {
    name: "ParticipantsList",
    components: { ParticipantCard },
    props: {
      participants: Array,
      revealed: Boolean,
      isHost: Function,
    },
    methods: {
      openEstimateDialog(participantId) {
        this.$emit('open-estimate-dialog', participantId);
      },
    },
  };
  </script>
  