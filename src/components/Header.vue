<script setup lang="ts">
defineProps<{
  query: string;
  stationCount: number;
  lastUpdated: string;
}>();

const emit = defineEmits<{
  'update-query': [value: string];
}>();

const handleInput = (event: Event) => {
  emit('update-query', (event.target as HTMLInputElement).value);
};
</script>

<template>
  <header class="sticky top-0 z-[1100] border-b border-slate-100 bg-white/95 shadow-sm shadow-slate-200/60 backdrop-blur">
    <div class="mx-auto grid h-[60px] max-w-[1800px] grid-cols-[auto_minmax(0,1fr)_auto] items-center gap-2 px-2 md:h-[72px] md:gap-3 md:px-4">
      <a class="flex min-w-0 items-center gap-2 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-600 focus:ring-offset-2 md:gap-3" href="#">
        <span class="grid h-9 w-9 place-items-center rounded-xl bg-blue-600 text-sm text-white shadow-lg shadow-blue-600/20 md:h-11 md:w-11 md:rounded-2xl md:text-base">
          <i class="fa-solid fa-gas-pump" aria-hidden="true"></i>
        </span>
        <span class="hidden min-w-0 sm:block">
          <span class="block truncate text-sm font-black leading-4 text-slate-950">FuelWatch Jakarta</span>
          <span class="block truncate text-xs font-semibold text-slate-500">Cek harga BBM hari ini</span>
        </span>
      </a>

      <label class="group grid h-10 grid-cols-[auto_minmax(0,1fr)] items-center gap-2 rounded-xl bg-slate-50 px-3 ring-1 ring-slate-100 transition focus-within:bg-white focus-within:ring-2 focus-within:ring-blue-600 md:h-12 md:gap-3 md:rounded-2xl md:px-4">
        <i class="fa-solid fa-magnifying-glass text-xs text-slate-400 transition group-focus-within:text-blue-600 md:text-sm" aria-hidden="true"></i>
        <input
          :value="query"
          type="search"
          class="min-w-0 bg-transparent text-sm font-semibold text-slate-950 outline-none placeholder:text-slate-500 md:text-base"
          placeholder="Cari SPBU..."
          aria-label="Cari SPBU atau jenis BBM"
          autocomplete="off"
          @input="handleInput"
        />
      </label>

      <div class="flex items-center justify-end gap-1.5 md:gap-2">
        <div class="hidden rounded-xl bg-slate-50 px-2 py-1.5 text-right ring-1 ring-slate-100 md:block md:rounded-2xl md:px-3 md:py-2">
          <span class="block text-[10px] font-bold uppercase tracking-wide text-slate-500 md:text-[11px]">SPBU</span>
          <strong class="block text-xs font-black text-slate-950 md:text-sm">{{ stationCount }}</strong>
        </div>
        <div class="rounded-xl bg-emerald-50 px-2 py-1.5 text-right ring-1 ring-emerald-100 md:rounded-2xl md:px-3 md:py-2">
          <span class="block text-[10px] font-bold uppercase tracking-wide text-emerald-700 md:text-[11px]">Diperbarui</span>
          <strong class="block whitespace-nowrap text-xs font-black text-emerald-900 md:text-sm">{{ lastUpdated }}</strong>
        </div>
      </div>
    </div>
  </header>
</template>
