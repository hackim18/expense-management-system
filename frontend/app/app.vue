<template>
  <div data-theme="rupate">
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
  </div>
</template>

<script setup lang="ts">
const auth = useAuth()

onMounted(() => {
  auth.init()
})
</script>
