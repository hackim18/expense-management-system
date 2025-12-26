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
          Mulai kelola expense dengan rapi.
        </h1>
        <p class="max-w-xl text-base-content/70 text-balance">
          Buat akun employee untuk mengajukan expense, memantau status, dan
          berkolaborasi dengan manager.
        </p>
        <div
          class="stats stats-vertical border border-base-200/70 bg-base-100/80 shadow-soft backdrop-blur lg:stats-horizontal"
        >
          <div class="stat">
            <div class="stat-title">Transparan</div>
            <div class="stat-value text-primary">IDR</div>
            <div class="stat-desc">Format rupiah otomatis.</div>
          </div>
          <div class="stat">
            <div class="stat-title">Approval Cepat</div>
            <div class="stat-value text-secondary">1Jt+</div>
            <div class="stat-desc">Masuk antrian manager.</div>
          </div>
        </div>
      </div>

      <div
        class="card border border-base-200/80 bg-base-100/90 shadow-soft backdrop-blur"
      >
        <div class="card-body gap-6">
          <div class="space-y-2">
            <h2 class="card-title text-2xl">Daftar</h2>
            <p class="text-sm text-base-content/70">
              Isi data berikut untuk membuat akun baru.
            </p>
          </div>

          <form class="grid gap-4" @submit.prevent="handleRegister">
            <label class="form-control">
              <div class="label">
                <span class="label-text">Nama</span>
              </div>
              <input
                v-model="form.name"
                type="text"
                class="input input-bordered w-full"
                placeholder="Nama lengkap"
              />
            </label>
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
            <div v-if="success" class="alert alert-success text-sm">
              {{ success }}
            </div>

            <button class="btn btn-primary w-full" :disabled="loading">
              {{ loading ? 'Memproses...' : 'Buat Akun' }}
            </button>
          </form>

          <div
            class="flex flex-col gap-2 rounded-box bg-base-200/70 p-3 text-sm text-base-content/80 sm:flex-row sm:items-center sm:justify-between"
          >
            <span class="font-medium">Sudah punya akun?</span>
            <NuxtLink to="/login" class="btn btn-outline btn-sm">
              Masuk
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
  name: '',
  email: '',
  password: ''
})
const loading = ref(false)
const error = ref('')
const success = ref('')

onMounted(() => {
  auth.init()
  if (auth.isAuthenticated.value) {
    navigateTo('/expenses')
  }
})

const handleRegister = async () => {
  error.value = ''
  success.value = ''

  if (!form.name.trim() || !form.email.trim() || !form.password) {
    error.value = 'Mohon lengkapi semua field.'
    return
  }
  if (form.password.length < 8) {
    error.value = 'Password minimal 8 karakter.'
    return
  }

  loading.value = true
  try {
    await request('/api/auth/register', {
      method: 'POST',
      body: form,
      auth: false
    })
    success.value = 'Akun berhasil dibuat. Silakan login.'
    form.name = ''
    form.email = ''
    form.password = ''
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Gagal mendaftar'
  } finally {
    loading.value = false
  }
}
</script>
