<template>
  <section class="space-y-6 animate-rise">
    <div class="flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
      <div class="space-y-2">
        <div class="badge badge-outline">Dashboard</div>
        <h2 class="text-3xl font-bold text-balance">Ringkasan Pengeluaran</h2>
        <p class="text-base-content/70 text-balance">Pantau pengajuanmu, filter status, dan cek detail terbaru.</p>
      </div>
      <NuxtLink to="/expenses/new" class="btn btn-primary"> Ajukan Pengeluaran </NuxtLink>
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

        <div v-if="paging" class="flex flex-col gap-3 rounded-box border border-base-200/80 bg-base-100/90 px-4 py-3 shadow-sm md:flex-row md:items-center md:justify-between">
          <div class="text-sm text-base-content/70">
            Menampilkan
            <span class="font-semibold text-base-content">{{ rangeStart }}</span>
            -
            <span class="font-semibold text-base-content">{{ rangeEnd }}</span>
            dari
            <span class="font-semibold text-base-content">{{ totalItem }}</span>
            pengajuan
          </div>
          <div class="flex flex-wrap items-center gap-2">
            <button class="btn btn-sm btn-outline" :disabled="!paging.has_previous" @click="goToPage(page - 1)">
              Sebelumnya
            </button>
            <div class="join">
              <button
                v-for="pageNumber in visiblePages"
                :key="pageNumber"
                class="btn btn-sm join-item"
                :class="pageNumber === page ? 'btn-primary' : 'btn-outline'"
                @click="goToPage(pageNumber)"
              >
                {{ pageNumber }}
              </button>
            </div>
            <button class="btn btn-sm btn-outline" :disabled="!paging.has_next" @click="goToPage(page + 1)">
              Berikutnya
            </button>
            <select class="select select-sm select-bordered" :value="pageSize" @change="changePageSize($event)">
              <option v-for="size in pageSizeOptions" :key="size" :value="size">
                {{ size }} / halaman
              </option>
            </select>
          </div>
        </div>
      </template>

      <div v-else class="card border border-base-200/80 bg-base-100/90 shadow-sm">
        <div class="card-body text-center">
          <p class="text-lg font-semibold">Tidak ada data untuk filter ini.</p>
          <p class="text-sm text-base-content/70">Coba pilih status lain atau buat pengajuan baru.</p>
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

type PageMetadata = {
  current_page: number;
  page_size: number;
  total_item: number;
  total_page: number;
  has_next: boolean;
  has_previous: boolean;
};

const { requestWithMeta } = useApi();
const { format: formatIdr } = useIdr();

const expenses = ref<Expense[]>([]);
const paging = ref<PageMetadata | null>(null);
const loading = ref(false);
const error = ref("");
const selectedStatus = ref("");
const page = ref(1);
const pageSize = ref(5);
const pageSizeOptions = [5, 10, 20, 50];

const statusOptions = [
  { label: "Semua", value: "" },
  { label: "Pending", value: "awaiting_approval" },
  { label: "Approved", value: "approved" },
  { label: "Rejected", value: "rejected" },
  { label: "Auto-approved", value: "auto_approved" },
  { label: "Completed", value: "completed" },
];

const fetchExpenses = async () => {
  loading.value = true;
  error.value = "";
  try {
    const query = new URLSearchParams({
      page: String(page.value),
      size: String(pageSize.value),
    });
    if (selectedStatus.value) {
      query.set("status", selectedStatus.value);
    }
    const payload = await requestWithMeta<Expense[], PageMetadata>(`/api/expenses?${query.toString()}`);
    expenses.value = payload.data || [];
    paging.value = payload.paging || null;
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Gagal memuat data";
    paging.value = null;
  } finally {
    loading.value = false;
  }
};

const changeStatus = (value: string) => {
  selectedStatus.value = value;
  page.value = 1;
  fetchExpenses();
};

const changePageSize = (event: Event) => {
  const target = event.target as HTMLSelectElement | null;
  const nextSize = target ? Number(target.value) : pageSize.value;
  if (!Number.isNaN(nextSize) && nextSize > 0) {
    pageSize.value = nextSize;
    page.value = 1;
    fetchExpenses();
  }
};

const goToPage = (nextPage: number) => {
  const total = paging.value?.total_page ?? 1;
  if (nextPage < 1 || nextPage > Math.max(1, total)) {
    return;
  }
  page.value = nextPage;
  fetchExpenses();
};

const totalItem = computed(() => paging.value?.total_item ?? expenses.value.length);
const rangeStart = computed(() => (expenses.value.length ? (page.value - 1) * pageSize.value + 1 : 0));
const rangeEnd = computed(() => (expenses.value.length ? (page.value - 1) * pageSize.value + expenses.value.length : 0));
const totalPages = computed(() => Math.max(1, paging.value?.total_page ?? 1));
const visiblePages = computed(() => {
  const total = totalPages.value;
  const current = page.value;
  const delta = 2;
  const start = Math.max(1, current - delta);
  const end = Math.min(total, current + delta);
  const pages: number[] = [];
  for (let i = start; i <= end; i += 1) {
    pages.push(i);
  }
  return pages;
});

const formatDate = (value: string) => {
  if (!value) return "-";
  return new Date(value).toLocaleString("id-ID", {
    dateStyle: "medium",
    timeStyle: "short",
  });
};

onMounted(fetchExpenses);
</script>
