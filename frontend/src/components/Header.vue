<template>
    <header class="w-full bg-blue-900 py-6 shadow-md">
      <h1 class="text-center text-3xl font-bold text-white">
        Room: <span class="text-blue-300">{{ roomName }}</span>
      </h1>
      <p class="text-center text-lg text-blue-100 mt-2">
        Estimation Type:
        <span class="font-semibold text-yellow-300">{{ estimationTypeLabel }}</span>
      </p>
    </header>
</template>
  
  <script>
  export default {
    name: "Header",
    props: {
      roomName: { type: String, required: true },
      estimationTypeLabel: { type: String, required: true },
    },
  };
  </script>

<style scoped>
header p {
  font-size: 1.125rem;
  margin-top: 0.5rem;
}
</style>
  