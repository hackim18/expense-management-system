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
        <div class="hidden text-right text-xs lg:block">
          <p class="font-semibold leading-tight">
            {{ auth.user.value?.name || 'User' }}
          </p>
          <div class="badge badge-outline">{{ roleLabel }}</div>
        </div>
        <button class="btn btn-outline btn-sm" @click="handleLogout">
          Logout
        </button>
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

const roleLabel = computed(() =>
  auth.user.value?.role === 'manager' ? 'Manager' : 'Employee'
)

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
