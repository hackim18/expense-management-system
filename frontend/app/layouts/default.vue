<template>
  <div class="relative min-h-screen bg-base-200">
    <div class="pointer-events-none absolute inset-0 bg-grid opacity-40"></div>
    <div class="relative mx-auto max-w-6xl px-4 pb-12 sm:px-6 lg:px-8">
      <AppNav />
      <main class="pt-6">
        <slot />
      </main>
    </div>
  </div>
</template>
