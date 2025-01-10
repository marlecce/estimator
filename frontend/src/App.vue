<template>
  <div id="app">
    <header class="bg-blue-600 text-white shadow-md">
      <div class="max-w-7xl mx-auto px-4 py-3 flex items-center justify-between">
        <!-- Logo with Icon -->
        <div class="flex items-center space-x-3">
          <div class="w-10 h-10 bg-white text-blue-600 rounded-full flex items-center justify-center font-bold shadow">
            <span>🛠️</span>
          </div>
          <h1 class="text-2xl font-bold tracking-wide">Estimator</h1>
        </div>

        <!-- Navigation -->
        <nav class="hidden md:flex space-x-4">
          <a
            href="#"
            class="text-white hover:underline"
          >
            Home
          </a>
          <a
            href="#features"
            class="text-white hover:underline"
          >
            Features
          </a>
          <a
            href="#contact"
            class="text-white hover:underline"
          >
            Contact
          </a>
        </nav>
      </div>
    </header>
    <router-view />
  </div>
</template>

<script>
export default {
  name: "App",
};
</script>

<style scoped>
/* Icon-specific customization */
.icon-container {
  background: linear-gradient(135deg, #ffafbd, #ffc3a0);
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}
</style>
