<template>
  <header class="pt-4">
    <div
      class="navbar rounded-box border border-base-200/80 bg-base-100/80 shadow-soft backdrop-blur"
    >
      <div class="navbar-start">
        <div class="dropdown">
          <label tabindex="0" class="btn btn-ghost lg:hidden">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="h-5 w-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M4 6h16M4 12h16M4 18h16"
              />
            </svg>
          </label>
          <ul
            tabindex="0"
            class="menu menu-sm dropdown-content mt-3 w-52 rounded-box bg-base-100 p-2 shadow"
          >
            <li v-for="link in visibleLinks" :key="link.to">
              <NuxtLink :to="link.to">{{ link.label }}</NuxtLink>
            </li>
          </ul>
        </div>
        <NuxtLink to="/expenses" class="btn btn-ghost text-xl">
          RupaExpenses
        </NuxtLink>
      </div>

      <div class="navbar-center hidden lg:flex">
        <ul class="menu menu-horizontal px-1">
          <li v-for="link in visibleLinks" :key="link.to">
            <NuxtLink :to="link.to" :class="isActive(link.to) ? 'active' : ''">
              {{ link.label }}
            </NuxtLink>
          </li>
        </ul>
      </div>

      <div class="navbar-end gap-2">
        <NuxtLink to="/expenses/new" class="btn btn-primary btn-sm hidden sm:inline-flex">
          Ajukan
        </NuxtLink>
        <div class="dropdown dropdown-end">
          <label tabindex="0" class="btn btn-ghost btn-sm gap-3 rounded-full px-3">
            <span
              class="flex h-9 w-9 items-center justify-center rounded-full border border-base-300/60 bg-base-200/70 text-base-content/70"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-4 w-4"
                viewBox="0 0 24 24"
                fill="currentColor"
              >
                <path
                  d="M12 12a4 4 0 1 0-4-4 4 4 0 0 0 4 4zm0 2c-3.87 0-7 2.24-7 5a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1c0-2.76-3.13-5-7-5z"
                />
              </svg>
            </span>
            <div class="hidden text-left sm:block">
              <p class="text-sm font-semibold leading-tight">{{ userName }}</p>
              <p class="text-xs text-base-content/60">{{ roleLabel }}</p>
            </div>
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="hidden h-4 w-4 opacity-60 sm:block"
              viewBox="0 0 20 20"
              fill="currentColor"
            >
              <path
                fill-rule="evenodd"
                d="M5.23 7.21a.75.75 0 0 1 1.06.02L10 11.168l3.71-3.938a.75.75 0 1 1 1.08 1.04l-4.24 4.5a.75.75 0 0 1-1.08 0l-4.24-4.5a.75.75 0 0 1 .02-1.06z"
                clip-rule="evenodd"
              />
            </svg>
          </label>
          <div
            tabindex="0"
            class="dropdown-content mt-3 w-56 rounded-box border border-base-200/80 bg-base-100 p-3 shadow-soft"
          >
            <div class="flex items-center gap-3">
              <span
                class="flex h-10 w-10 items-center justify-center rounded-full border border-base-300/60 bg-base-200/70 text-base-content/70"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  class="h-5 w-5"
                  viewBox="0 0 24 24"
                  fill="currentColor"
                >
                  <path
                    d="M12 12a4 4 0 1 0-4-4 4 4 0 0 0 4 4zm0 2c-3.87 0-7 2.24-7 5a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1c0-2.76-3.13-5-7-5z"
                  />
                </svg>
              </span>
              <div>
                <p class="text-sm font-semibold">{{ userName }}</p>
                <p class="text-xs text-base-content/60">{{ roleLabel }}</p>
              </div>
            </div>
            <div class="my-3 h-px bg-base-200"></div>
            <button class="btn btn-ghost btn-sm w-full justify-start" @click="handleLogout">
              Logout
            </button>
          </div>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
const route = useRoute()
const auth = useAuth()

const links = [
  { label: 'Dashboard', to: '/expenses' },
  { label: 'Ajukan', to: '/expenses/new' },
  { label: 'Antrian Approval', to: '/manager/approvals', managerOnly: true }
]

const visibleLinks = computed(() =>
  links.filter((link) => !link.managerOnly || auth.user.value?.role === 'manager')
)

const userName = computed(() => auth.user.value?.name || 'User')
const roleLabel = computed(() =>
  auth.user.value?.role === 'manager' ? 'Manajer' : 'Karyawan'
)
const initials = computed(() => {
  const parts = userName.value.trim().split(/\s+/).filter(Boolean)
  if (!parts.length) return 'U'
  if (parts.length === 1) return parts[0].slice(0, 2).toUpperCase()
  return `${parts[0][0]}${parts[parts.length - 1][0]}`.toUpperCase()
})

const isActive = (path: string) => {
  if (path === '/expenses') {
    return route.path === '/expenses'
  }
  return route.path.startsWith(path)
}

const handleLogout = () => {
  auth.logout()
  navigateTo('/login')
}
</script>
