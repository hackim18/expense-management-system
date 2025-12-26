export default defineNuxtRouteMiddleware(() => {
  if (import.meta.server) {
    return
  }

  const auth = useAuth()
  auth.init()

  if (!auth.isAuthenticated.value) {
    return navigateTo('/login')
  }
})
