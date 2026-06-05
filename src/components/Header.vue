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
    <div class="mx-auto grid h-[72px] max-w-[1800px] grid-cols-[auto_minmax(0,720px)_auto] items-center gap-3 px-3 md:px-4">
      <a class="flex min-w-0 items-center gap-3 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-600 focus:ring-offset-2" href="#">
        <span class="grid h-11 w-11 place-items-center rounded-2xl bg-blue-600 text-white shadow-lg shadow-blue-600/20">
          <i class="fa-solid fa-gas-pump" aria-hidden="true"></i>
        </span>
        <span class="hidden min-w-0 sm:block">
          <span class="block truncate text-sm font-black leading-4 text-slate-950">FuelWatch Jakarta</span>
          <span class="block truncate text-xs font-semibold text-slate-500">Cek harga BBM hari ini</span>
        </span>
      </a>

      <label class="group grid h-12 grid-cols-[auto_minmax(0,1fr)] items-center gap-3 rounded-2xl bg-slate-50 px-4 ring-1 ring-slate-100 transition focus-within:bg-white focus-within:ring-2 focus-within:ring-blue-600">
        <i class="fa-solid fa-magnifying-glass text-slate-400 transition group-focus-within:text-blue-600" aria-hidden="true"></i>
        <input
          :value="query"
          type="search"
          class="min-w-0 bg-transparent text-base font-semibold text-slate-950 outline-none placeholder:text-slate-500"
          placeholder="Cari area, merek, SPBU, atau jenis BBM..."
          aria-label="Cari SPBU atau jenis BBM"
          autocomplete="off"
          @input="handleInput"
        />
      </label>

      <div class="flex items-center justify-end gap-2">
        <div class="hidden rounded-2xl bg-slate-50 px-3 py-2 text-right ring-1 ring-slate-100 md:block">
          <span class="block text-[11px] font-bold uppercase tracking-wide text-slate-500">SPBU</span>
          <strong class="block text-sm font-black text-slate-950">{{ stationCount }}</strong>
        </div>
        <div class="rounded-2xl bg-emerald-50 px-3 py-2 text-right ring-1 ring-emerald-100">
          <span class="block text-[11px] font-bold uppercase tracking-wide text-emerald-700">Diperbarui</span>
          <strong class="block whitespace-nowrap text-sm font-black text-emerald-900">{{ lastUpdated }}</strong>
        </div>
      </div>
    </div>
  </header>
</template>
