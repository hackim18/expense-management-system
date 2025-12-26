<template>
  <section class="space-y-6 animate-rise">
    <div class="space-y-2">
      <div class="badge badge-outline">Manager</div>
      <h2 class="text-3xl font-bold text-balance">Antrian Persetujuan</h2>
      <p class="text-base-content/70 text-balance">
        Expense di atas 1 juta menunggu keputusanmu.
      </p>
    </div>

    <div v-if="!isManager" class="alert alert-warning">
      Halaman ini hanya untuk manager.
    </div>

    <div v-else>
      <div v-if="error" class="alert alert-error">
        {{ error }}
      </div>

      <div v-if="loading" class="card border border-base-200/80 bg-base-100/90 shadow-sm">
        <div class="card-body text-base-content/70">Memuat antrian...</div>
      </div>

      <div v-else>
        <div v-if="expenses.length" class="grid gap-4">
          <div
            v-for="expense in expenses"
            :key="expense.id"
            class="card border border-base-200/80 bg-base-100/90 shadow-soft"
          >
            <div class="card-body space-y-4">
              <div class="flex flex-col gap-3 md:flex-row md:items-start md:justify-between">
                <div class="space-y-1">
                  <h3 class="text-lg font-semibold">{{ expense.description }}</h3>
                  <p class="text-sm text-base-content/60">
                    {{ formatDate(expense.submitted_at) }}
                  </p>
                  <StatusBadge :status="expense.status" />
                </div>
                <div class="text-2xl font-semibold">
                  {{ expense.amount_idr_formatted || formatIdr(expense.amount_idr) }}
                </div>
              </div>

              <label class="form-control">
                <div class="label">
                  <span class="label-text">Catatan approval</span>
                </div>
                <textarea
                  v-model="notes[expense.id]"
                  class="textarea textarea-bordered w-full"
                  rows="2"
                  placeholder="Catatan approval (opsional)"
                />
              </label>
              <div class="flex flex-col gap-3 sm:flex-row">
                <button
                  class="btn btn-success flex-1"
                  :disabled="busyId === expense.id"
                  @click="approveExpense(expense.id)"
                >
                  Approve
                </button>
                <button
                  class="btn btn-error flex-1"
                  :disabled="busyId === expense.id"
                  @click="rejectExpense(expense.id)"
                >
                  Reject
                </button>
              </div>
            </div>
          </div>
        </div>

        <div
          v-else
          class="card border border-base-200/80 bg-base-100/90 shadow-sm"
        >
          <div class="card-body text-center">
            <p class="text-lg font-semibold">Antrian kosong.</p>
            <p class="text-sm text-base-content/70">
              Semua expense sudah diproses.
            </p>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth'
})

type Expense = {
  id: string
  amount_idr: number
  amount_idr_formatted?: string
  description: string
  status: string
  submitted_at: string
}

const { request } = useApi()
const { format: formatIdr } = useIdr()
const auth = useAuth()

const expenses = ref<Expense[]>([])
const loading = ref(false)
const error = ref('')
const notes = reactive<Record<string, string>>({})
const busyId = ref('')

const isManager = computed(() => auth.user.value?.role === 'manager')

const fetchQueue = async () => {
  if (!isManager.value) {
    return
  }

  loading.value = true
  error.value = ''
  try {
    const data = await request<Expense[]>('/api/expenses?status=pending')
    expenses.value = data || []
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Gagal memuat antrian'
  } finally {
    loading.value = false
  }
}

const approveExpense = async (id: string) => {
  busyId.value = id
  try {
    await request(`/api/expenses/${id}/approve`, {
      method: 'PUT',
      body: { notes: notes[id] || '' }
    })
    expenses.value = expenses.value.filter((item) => item.id !== id)
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Gagal approve'
  } finally {
    busyId.value = ''
  }
}

const rejectExpense = async (id: string) => {
  busyId.value = id
  try {
    await request(`/api/expenses/${id}/reject`, {
      method: 'PUT',
      body: { notes: notes[id] || '' }
    })
    expenses.value = expenses.value.filter((item) => item.id !== id)
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Gagal reject'
  } finally {
    busyId.value = ''
  }
}

const formatDate = (value: string) => {
  if (!value) return '-'
  return new Date(value).toLocaleString('id-ID', {
    dateStyle: 'medium',
    timeStyle: 'short'
  })
}

onMounted(fetchQueue)
</script>
