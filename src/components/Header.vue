<script setup lang="ts">
defineProps<{
  query: string;
  lastUpdated: string;
  isRefreshing: boolean;
}>();

const emit = defineEmits<{
  'update-query': [value: string];
  refresh: [];
}>();
</script>

<template>
  <header class="flex h-14 shrink-0 items-center gap-3 bg-petrol px-3 text-white md:h-16 md:gap-5 md:px-5">
    <a href="#" class="flex shrink-0 items-center gap-2 rounded-md">
      <i class="fa-solid fa-gas-pump text-lg text-lampu" aria-hidden="true"></i>
      <span class="hidden text-[15px] font-bold tracking-tight sm:inline">FuelWatch <span class="font-medium text-white/65">Jakarta</span></span>
    </a>

    <label class="flex h-10 min-w-0 flex-1 items-center gap-2 rounded-lg bg-white/12 px-3 transition focus-within:bg-white focus-within:text-aspal focus-within:ring-2 focus-within:ring-lampu md:max-w-md">
      <i class="fa-solid fa-magnifying-glass text-sm opacity-70" aria-hidden="true"></i>
      <input
        :value="query"
        type="search"
        class="min-w-0 flex-1 bg-transparent text-sm outline-none placeholder:text-current placeholder:opacity-60"
        placeholder="Cari SPBU, jalan, atau Pertamax"
        aria-label="Cari SPBU, jalan, atau jenis BBM"
        autocomplete="off"
        @input="emit('update-query', ($event.target as HTMLInputElement).value)"
      />
    </label>

    <div class="ml-auto flex shrink-0 items-center gap-3 text-sm">
      <span class="hidden text-white/65 lg:inline">Harga diperbarui {{ lastUpdated }}</span>
      <button
        type="button"
        class="inline-flex h-10 items-center gap-2 rounded-lg border border-white/25 px-3 font-semibold transition hover:bg-white/10 disabled:cursor-wait disabled:opacity-60"
        :disabled="isRefreshing"
        :aria-label="isRefreshing ? 'Memperbarui harga' : 'Perbarui harga'"
        @click="emit('refresh')"
      >
        <i class="fa-solid fa-rotate-right text-xs" :class="{ 'animate-spin': isRefreshing }" aria-hidden="true"></i>
        <span class="hidden sm:inline">{{ isRefreshing ? 'Memperbarui' : 'Perbarui harga' }}</span>
      </button>
    </div>
  </header>
</template>
