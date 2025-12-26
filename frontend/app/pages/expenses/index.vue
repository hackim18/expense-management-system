<template>
  <section class="space-y-6 animate-rise">
    <div class="flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
      <div class="space-y-2">
        <div class="badge badge-outline">Dashboard</div>
        <h2 class="text-3xl font-bold text-balance">Ringkasan Expense</h2>
        <p class="text-base-content/70 text-balance">Pantau pengajuanmu, filter status, dan cek detail terbaru.</p>
      </div>
      <NuxtLink to="/expenses/new" class="btn btn-primary"> Ajukan Expense </NuxtLink>
    </div>

    <div class="card border border-base-200/80 bg-base-100/80 shadow-sm">
      <div class="card-body gap-3">
        <div class="text-sm font-semibold uppercase tracking-wide text-base-content/60">Filter status</div>
        <div class="flex flex-wrap gap-2">
          <button v-for="option in statusOptions" :key="option.value" class="btn btn-sm" :class="selectedStatus === option.value ? 'btn-primary' : 'btn-outline'" @click="changeStatus(option.value)">
            {{ option.label }}
          </button>
        </div>
      </div>
    </div>

    <div v-if="error" class="alert alert-error">
      {{ error }}
    </div>

    <div v-if="loading" class="card border border-base-200/80 bg-base-100/90 shadow-sm">
      <div class="card-body text-base-content/70">Memuat data...</div>
    </div>

    <div v-else class="space-y-4">
      <!-- Kalau ADA data -->
      <template v-if="expenses.length">
        <div class="hidden overflow-x-auto rounded-box border border-base-200/80 bg-base-100/90 shadow-sm lg:block">
          <table class="table table-zebra">
            <thead>
              <tr>
                <th>Deskripsi</th>
                <th>Tanggal</th>
                <th>Jumlah</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="expense in expenses" :key="expense.id">
                <td class="max-w-xs">
                  <div class="font-semibold">{{ expense.description }}</div>
                </td>
                <td class="text-sm text-base-content/60">
                  {{ formatDate(expense.submitted_at) }}
                </td>
                <td class="font-semibold">
                  {{ expense.amount_idr_formatted || formatIdr(expense.amount_idr) }}
                </td>
                <td>
                  <StatusBadge :status="expense.status" />
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="grid gap-4 lg:hidden">
          <div v-for="expense in expenses" :key="expense.id" class="card border border-base-200/80 bg-base-100/90 shadow-soft">
            <div class="card-body">
              <div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
                <div class="space-y-1">
                  <h3 class="text-lg font-semibold">{{ expense.description }}</h3>
                  <p class="text-sm text-base-content/60">
                    {{ formatDate(expense.submitted_at) }}
                  </p>
                </div>
                <div class="flex flex-col items-start gap-2 md:items-end">
                  <p class="text-2xl font-semibold">
                    {{ expense.amount_idr_formatted || formatIdr(expense.amount_idr) }}
                  </p>
                  <StatusBadge :status="expense.status" />
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>

      <!-- Kalau TIDAK ADA data: tampilkan HANYA ini -->
      <div v-else class="card border border-base-200/80 bg-base-100/90 shadow-sm">
        <div class="card-body text-center">
          <p class="text-lg font-semibold">Belum ada pengajuan.</p>
          <p class="text-sm text-base-content/70">Mulai ajukan expense pertama kamu hari ini.</p>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: "auth",
});

type Expense = {
  id: string;
  amount_idr: number;
  amount_idr_formatted?: string;
  description: string;
  status: string;
  submitted_at: string;
};

const { request } = useApi();
const { format: formatIdr } = useIdr();

const expenses = ref<Expense[]>([]);
const loading = ref(false);
const error = ref("");
const selectedStatus = ref("");

const statusOptions = [
  { label: "Semua", value: "" },
  { label: "Pending", value: "pending" },
  { label: "Approved", value: "approved" },
  { label: "Rejected", value: "rejected" },
  { label: "Auto-approved", value: "auto-approved" },
  { label: "Completed", value: "completed" },
];

const fetchExpenses = async () => {
  loading.value = true;
  error.value = "";
  try {
    const query = selectedStatus.value ? `?status=${encodeURIComponent(selectedStatus.value)}` : "";
    const data = await request<Expense[]>(`/api/expenses${query}`);
    expenses.value = data || [];
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Gagal memuat data";
  } finally {
    loading.value = false;
  }
};

const changeStatus = (value: string) => {
  selectedStatus.value = value;
  fetchExpenses();
};

const formatDate = (value: string) => {
  if (!value) return "-";
  return new Date(value).toLocaleString("id-ID", {
    dateStyle: "medium",
    timeStyle: "short",
  });
};

onMounted(fetchExpenses);
</script>
