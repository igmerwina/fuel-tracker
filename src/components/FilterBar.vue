<script setup lang="ts">
import type { FilterState, FuelBrand, FuelType, SortMode } from '../types';
import { brands, fuelTypes, regions } from '../data/stations';

const props = defineProps<{
  filters: FilterState;
  sortMode: SortMode;
}>();

const emit = defineEmits<{
  'update-filters': [value: Partial<FilterState>];
  'update-sort': [value: SortMode];
}>();

const fuelLabel = (fuel: FuelType) => {
  if (fuel === 'RON_92') return 'Pertamax';
  if (fuel === 'RON_90') return 'Pertalite';
  return fuelTypes.find((item) => item.id === fuel)?.label ?? fuel;
};
</script>

<template>
  <section class="flex min-h-12 items-center gap-2 overflow-x-auto rounded-2xl bg-white px-2.5 py-2 shadow-sm shadow-slate-200/70 ring-1 ring-slate-100">
    <label class="inline-flex h-9 shrink-0 items-center gap-2 rounded-full bg-blue-600 px-3 text-sm font-black text-white shadow-sm shadow-blue-600/20 transition hover:bg-blue-700">
      <i class="fa-solid fa-gas-pump text-xs" aria-hidden="true"></i>
      <select
        :value="filters.fuelType"
        class="max-w-[132px] appearance-none bg-transparent outline-none"
        aria-label="Jenis BBM"
        @change="emit('update-filters', { fuelType: ($event.target as HTMLSelectElement).value as FuelType })"
      >
        <option v-for="fuel in fuelTypes" :key="fuel.id" :value="fuel.id">{{ fuelLabel(fuel.id) }}</option>
      </select>
    </label>

    <label class="inline-flex h-9 shrink-0 items-center gap-2 rounded-full bg-slate-50 px-3 text-sm font-black text-slate-700 ring-1 ring-slate-200 transition hover:bg-white hover:ring-slate-300">
      <i class="fa-solid fa-building-circle-check text-xs text-slate-500" aria-hidden="true"></i>
      <select
        :value="filters.brand"
        class="max-w-[128px] appearance-none bg-transparent outline-none"
        aria-label="Merek"
        @change="emit('update-filters', { brand: ($event.target as HTMLSelectElement).value as FuelBrand | '' })"
      >
        <option value="">Semua merek</option>
        <option v-for="brand in brands" :key="brand" :value="brand">{{ brand }}</option>
      </select>
    </label>

    <label class="inline-flex h-9 shrink-0 items-center gap-2 rounded-full bg-slate-50 px-3 text-sm font-black text-slate-700 ring-1 ring-slate-200 transition hover:bg-white hover:ring-slate-300">
      <i class="fa-solid fa-location-dot text-xs text-slate-500" aria-hidden="true"></i>
      <select
        :value="filters.region"
        class="max-w-[150px] appearance-none bg-transparent outline-none"
        aria-label="Wilayah"
        @change="emit('update-filters', { region: ($event.target as HTMLSelectElement).value as FilterState['region'] })"
      >
        <option value="">Semua wilayah</option>
        <option v-for="region in regions" :key="region" :value="region">{{ region }}</option>
      </select>
    </label>

    <button
      type="button"
      class="inline-flex h-9 shrink-0 items-center gap-2 rounded-full px-3 text-sm font-black ring-1 transition focus:outline-none focus:ring-2 focus:ring-blue-600"
      :class="filters.openNow ? 'bg-emerald-500 text-white ring-emerald-500 shadow-sm shadow-emerald-500/20' : 'bg-slate-50 text-slate-700 ring-slate-200 hover:bg-white hover:ring-slate-300'"
      @click="emit('update-filters', { openNow: !filters.openNow })"
    >
      <i class="fa-regular fa-clock" aria-hidden="true"></i>
      Buka sekarang
    </button>

    <button
      type="button"
      class="ml-auto inline-flex h-9 shrink-0 items-center gap-2 rounded-full px-3 text-sm font-black transition focus:outline-none focus:ring-2 focus:ring-blue-600"
      :class="sortMode === 'cheapest' ? 'bg-slate-950 text-white shadow-sm shadow-slate-950/20' : 'bg-slate-50 text-slate-700 ring-1 ring-slate-200 hover:bg-white hover:ring-slate-300'"
      @click="emit('update-sort', sortMode === 'cheapest' ? 'nearest' : 'cheapest')"
    >
      <i class="fa-solid fa-sliders" aria-hidden="true"></i>
      {{ sortMode === 'cheapest' ? 'Termurah dulu' : 'Terdekat dulu' }}
    </button>
  </section>
</template>
