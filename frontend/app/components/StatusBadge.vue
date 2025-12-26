<template>
  <span :class="badgeClass">{{ label }}</span>
</template>

<script setup lang="ts">
const props = defineProps<{
  status: string
}>()

const labelMap: Record<string, string> = {
  awaiting_approval: 'Menunggu Persetujuan',
  approved: 'Disetujui',
  rejected: 'Ditolak',
  auto_approved: 'Auto-Approved',
  completed: 'Selesai'
}

const badgeClassMap: Record<string, string> = {
  awaiting_approval: 'badge-warning',
  approved: 'badge-success',
  rejected: 'badge-error',
  auto_approved: 'badge-info',
  completed: 'badge-neutral'
}

const label = computed(() => labelMap[props.status] || props.status)
const badgeClass = computed(
  () => `badge badge-sm ${badgeClassMap[props.status] || 'badge-ghost'}`
)
</script>
