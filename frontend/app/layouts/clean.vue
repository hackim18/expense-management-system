<template>
  <div class="relative min-h-screen overflow-hidden bg-base-200">
    <div class="pointer-events-none absolute inset-0 bg-grid opacity-60"></div>
    <div
      class="pointer-events-none absolute -top-24 right-0 h-64 w-64 rounded-full bg-primary/15 blur-3xl"
    ></div>
    <div
      class="pointer-events-none absolute bottom-0 left-0 h-72 w-72 rounded-full bg-secondary/20 blur-3xl"
    ></div>
    <main class="relative mx-auto flex min-h-screen max-w-6xl items-center px-6 py-14">
      <slot />
    </main>
  </div>
</template>
