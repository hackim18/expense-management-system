<template>
  <div class="min-h-screen bg-base-200" />
</template>

<script setup lang="ts">
const auth = useAuth()

onMounted(() => {
  auth.init()
  if (auth.isAuthenticated.value) {
    navigateTo('/expenses')
  } else {
    navigateTo('/login')
  }
})
</script>
