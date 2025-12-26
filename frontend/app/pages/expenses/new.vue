<template>
  <section class="mx-auto max-w-5xl space-y-6 animate-rise">
    <div class="space-y-2">
      <div class="badge badge-outline">Ajukan Expense</div>
      <h2 class="text-3xl font-bold text-balance">Form Pengajuan</h2>
      <p class="text-base-content/70 text-balance">
        Biaya di atas Rp 1.000.000 wajib approval manager.
      </p>
    </div>

    <div class="card border border-base-200/80 bg-base-100/90 shadow-soft">
      <div class="card-body gap-6">
        <form class="grid gap-5" @submit.prevent="handleSubmit">
          <label class="form-control">
            <div class="label">
              <span class="label-text">Jumlah (IDR)</span>
            </div>
            <div class="join w-full">
              <span class="btn join-item btn-ghost pointer-events-none">Rp</span>
              <input
                v-model="amountInput"
                type="text"
                inputmode="numeric"
                class="input input-bordered join-item w-full"
                placeholder="1.500.000"
                @input="handleAmountInput"
              />
            </div>
            <div class="label">
              <span class="label-text-alt text-base-content/60">
                Min 10.000, max 50.000.000
              </span>
            </div>
          </label>

          <label class="form-control">
            <div class="label">
              <span class="label-text">Deskripsi</span>
            </div>
            <textarea
              v-model="form.description"
              class="textarea textarea-bordered w-full"
              rows="4"
              placeholder="Contoh: Meeting dengan klien di Plaza Indonesia"
            />
          </label>

          <div class="grid gap-4 lg:grid-cols-2">
            <label class="form-control">
              <div class="label">
                <span class="label-text">Upload Receipt</span>
              </div>
              <input
                type="file"
                class="file-input file-input-bordered w-full"
                @change="handleReceiptFile"
              />
              <div v-if="receiptFileName" class="label">
                <span class="label-text-alt text-base-content/60">
                  File dipilih: {{ receiptFileName }}
                </span>
              </div>
            </label>
            <label class="form-control">
              <div class="label">
                <span class="label-text">URL Receipt</span>
              </div>
              <input
                v-model="form.receipt_url"
                type="text"
                class="input input-bordered w-full"
                placeholder="https://example.com/receipt.jpg"
              />
            </label>
          </div>

          <div v-if="requiresApproval" class="alert alert-warning">
            Pengajuan ini memerlukan approval manager sebelum diproses.
          </div>
          <div v-else-if="amountValue > 0" class="alert alert-success">
            Pengajuan di bawah 1 juta akan auto-approved dan langsung diproses.
          </div>

          <div v-if="error" class="alert alert-error">
            {{ error }}
          </div>
          <div v-if="success" class="alert alert-success">
            {{ success }}
          </div>

          <div class="flex flex-col gap-3 sm:flex-row">
            <button class="btn btn-primary flex-1" :disabled="loading">
              {{ loading ? 'Mengirim...' : 'Kirim Pengajuan' }}
            </button>
            <NuxtLink to="/expenses" class="btn btn-outline flex-1">
              Kembali
            </NuxtLink>
          </div>
        </form>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth'
})

const { request } = useApi()
const { formatInput } = useIdr()

const amountInput = ref('')
const amountValue = ref(0)
const receiptFileName = ref('')
const error = ref('')
const success = ref('')
const loading = ref(false)

const form = reactive({
  amount_idr: 0,
  description: '',
  receipt_url: ''
})

const requiresApproval = computed(() => amountValue.value >= 1000000)

const handleAmountInput = () => {
  const { formatted, value } = formatInput(amountInput.value)
  amountInput.value = formatted
  amountValue.value = value
  form.amount_idr = value
}

const handleReceiptFile = (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  receiptFileName.value = file?.name || ''

  if (file && !form.receipt_url) {
    form.receipt_url = `mock://${file.name}`
  }
}

const handleSubmit = async () => {
  error.value = ''
  success.value = ''
  loading.value = true

  try {
    await request('/api/expenses', {
      method: 'POST',
      body: form
    })
    success.value = 'Pengajuan berhasil dikirim.'
    amountInput.value = ''
    amountValue.value = 0
    receiptFileName.value = ''
    form.amount_idr = 0
    form.description = ''
    form.receipt_url = ''
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Gagal mengirim'
  } finally {
    loading.value = false
  }
}
</script>
