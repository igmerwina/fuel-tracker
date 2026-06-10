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
  <section class="flex min-h-10 items-center gap-1.5 overflow-x-auto rounded-xl bg-white px-2 py-1.5 shadow-sm shadow-slate-200/70 ring-1 ring-slate-100 md:min-h-12 md:gap-2 md:rounded-2xl md:px-2.5 md:py-2">
    <label class="inline-flex h-8 shrink-0 items-center gap-1.5 rounded-full bg-blue-600 px-2.5 text-xs font-black text-white shadow-sm shadow-blue-600/20 transition hover:bg-blue-700 md:h-9 md:gap-2 md:px-3 md:text-sm">
      <i class="fa-solid fa-gas-pump text-[10px] md:text-xs" aria-hidden="true"></i>
      <select
        :value="filters.fuelType"
        class="max-w-[100px] appearance-none bg-transparent outline-none md:max-w-[132px]"
        aria-label="Jenis BBM"
        @change="emit('update-filters', { fuelType: ($event.target as HTMLSelectElement).value as FuelType })"
      >
        <option v-for="fuel in fuelTypes" :key="fuel.id" :value="fuel.id">{{ fuelLabel(fuel.id) }}</option>
      </select>
    </label>

    <label class="inline-flex h-8 shrink-0 items-center gap-1.5 rounded-full bg-slate-50 px-2.5 text-xs font-black text-slate-700 ring-1 ring-slate-200 transition hover:bg-white hover:ring-slate-300 md:h-9 md:gap-2 md:px-3 md:text-sm">
      <i class="fa-solid fa-building-circle-check text-[10px] text-slate-500 md:text-xs" aria-hidden="true"></i>
      <select
        :value="filters.brand"
        class="max-w-[100px] appearance-none bg-transparent outline-none md:max-w-[128px]"
        aria-label="Merek"
        @change="emit('update-filters', { brand: ($event.target as HTMLSelectElement).value as FuelBrand | '' })"
      >
        <option value="">Semua merek</option>
        <option v-for="brand in brands" :key="brand" :value="brand">{{ brand }}</option>
      </select>
    </label>

    <label class="inline-flex h-8 shrink-0 items-center gap-1.5 rounded-full bg-slate-50 px-2.5 text-xs font-black text-slate-700 ring-1 ring-slate-200 transition hover:bg-white hover:ring-slate-300 md:h-9 md:gap-2 md:px-3 md:text-sm">
      <i class="fa-solid fa-location-dot text-[10px] text-slate-500 md:text-xs" aria-hidden="true"></i>
      <select
        :value="filters.region"
        class="max-w-[110px] appearance-none bg-transparent outline-none md:max-w-[150px]"
        aria-label="Wilayah"
        @change="emit('update-filters', { region: ($event.target as HTMLSelectElement).value as FilterState['region'] })"
      >
        <option value="">Semua wilayah</option>
        <option v-for="region in regions" :key="region" :value="region">{{ region }}</option>
      </select>
    </label>

    <button
      type="button"
      class="ml-auto inline-flex h-8 shrink-0 items-center gap-1.5 rounded-full px-2.5 text-xs font-black transition focus:outline-none focus:ring-2 focus:ring-blue-600 md:h-9 md:gap-2 md:px-3 md:text-sm"
      :class="sortMode === 'cheapest' ? 'bg-slate-950 text-white shadow-sm shadow-slate-950/20' : 'bg-slate-50 text-slate-700 ring-1 ring-slate-200 hover:bg-white hover:ring-slate-300'"
      @click="emit('update-sort', sortMode === 'cheapest' ? 'nearest' : 'cheapest')"
    >
      <i class="fa-solid fa-sliders" aria-hidden="true"></i>
      <span class="hidden md:inline">{{ sortMode === 'cheapest' ? 'Harga terendah' : 'Jarak terdekat' }}</span>
    </button>
  </section>
</template>
