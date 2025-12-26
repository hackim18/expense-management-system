<template>
  <section class="mx-auto w-full max-w-6xl animate-rise">
    <div class="grid items-center gap-10 lg:grid-cols-[1.1fr_0.9fr]">
      <div class="space-y-7 lg:pr-8">
        <div class="flex flex-wrap items-center gap-2">
          <div class="badge badge-outline">Portal</div>
          <div class="badge badge-ghost">Expense Studio</div>
        </div>
        <h1
          class="max-w-2xl text-4xl font-bold leading-tight md:text-5xl lg:text-6xl text-balance"
        >
          Kendalikan pengeluaran, tetap gesit.
        </h1>
        <p class="max-w-xl text-base-content/70 text-balance">
          Masuk untuk mengajukan pengeluaran, memantau status, dan mempercepat
          persetujuan di seluruh tim.
        </p>
        <div
          class="stats stats-vertical border border-base-200/70 bg-base-100/80 shadow-soft backdrop-blur lg:stats-horizontal"
        >
          <div class="stat">
            <div class="stat-title">IDR Native</div>
            <div class="stat-value text-primary">Rp</div>
            <div class="stat-desc">Format rupiah otomatis.</div>
          </div>
          <div class="stat">
            <div class="stat-title">Persetujuan Cepat</div>
            <div class="stat-value text-secondary">1Jt+</div>
            <div class="stat-desc">Masuk antrean manajer.</div>
          </div>
        </div>
      </div>

      <div
        class="card border border-base-200/80 bg-base-100/90 shadow-soft backdrop-blur"
      >
        <div class="card-body gap-6">
          <div class="space-y-2">
            <h2 class="card-title text-2xl">Masuk</h2>
            <p class="text-sm text-base-content/70">
            Gunakan akun karyawan atau manajer yang sudah disediakan.
            </p>
          </div>

          <form class="grid gap-4" @submit.prevent="handleLogin">
            <label class="form-control">
              <div class="label">
                <span class="label-text">Email</span>
              </div>
              <input
                v-model="form.email"
                type="email"
                class="input input-bordered w-full"
                placeholder="nama@perusahaan.com"
              />
            </label>
            <label class="form-control">
              <div class="label">
                <span class="label-text">Password</span>
              </div>
              <input
                v-model="form.password"
                type="password"
                class="input input-bordered w-full"
                placeholder="Minimal 8 karakter"
              />
            </label>

            <div v-if="error" class="alert alert-error text-sm">
              {{ error }}
            </div>

            <button class="btn btn-primary w-full" :disabled="loading">
              {{ loading ? 'Memproses...' : 'Masuk' }}
            </button>
          </form>

          <div class="divider">Demo login</div>
          <div class="rounded-box bg-base-200/70 p-3 text-xs text-base-content/70">
            <p>Karyawan: john@mail.com / 12345678</p>
            <p>Manajer: manager@mail.com / 12345678</p>
          </div>

          <div
            class="flex flex-col gap-2 rounded-box bg-base-200/70 p-3 text-sm text-base-content/80 sm:flex-row sm:items-center sm:justify-between"
          >
            <span class="font-medium">Belum punya akun?</span>
            <NuxtLink to="/register" class="btn btn-outline btn-sm">
              Daftar
            </NuxtLink>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'clean'
})

const auth = useAuth()
const { request } = useApi()

const form = reactive({
  email: '',
  password: ''
})
const loading = ref(false)
const error = ref('')

onMounted(() => {
  auth.init()
  if (auth.isAuthenticated.value) {
    navigateTo('/expenses')
  }
})

const handleLogin = async () => {
  error.value = ''
  loading.value = true
  try {
    const data = await request<{
      id: string
      name: string
      email: string
      role: string
      access_token: string
    }>('/api/auth/login', {
      method: 'POST',
      body: form,
      auth: false
    })

    auth.setAuth(data)
    navigateTo('/expenses')
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Gagal login'
  } finally {
    loading.value = false
  }
}
</script>
