<script setup lang="ts">
import { ref } from 'vue';
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

const moreOpen = ref(false);

const fuelLabel = (fuel: FuelType) => {
  if (fuel === 'RON_92') return 'Pertamax';
  if (fuel === 'RON_90') return 'Pertalite';
  return fuelTypes.find((item) => item.id === fuel)?.label ?? fuel;
};
</script>

<template>
  <section class="flex items-center gap-2 overflow-x-auto rounded-[12px] bg-white/95 px-3 py-2 shadow-sm shadow-slate-200/70">
    <select
      :value="filters.fuelType"
      class="h-10 rounded-full bg-blue-600 px-4 text-sm font-black text-white outline-none focus:ring-2 focus:ring-blue-600 focus:ring-offset-2"
      aria-label="Jenis BBM"
      @change="emit('update-filters', { fuelType: ($event.target as HTMLSelectElement).value as FuelType })"
    >
      <option v-for="fuel in fuelTypes" :key="fuel.id" :value="fuel.id">{{ fuelLabel(fuel.id) }}</option>
    </select>

    <select
      :value="filters.brand"
      class="h-10 rounded-full bg-white px-4 text-sm font-black text-slate-700 shadow-sm ring-1 ring-slate-200 outline-none focus:ring-2 focus:ring-blue-600"
      aria-label="Merek"
      @change="emit('update-filters', { brand: ($event.target as HTMLSelectElement).value as FuelBrand | '' })"
    >
      <option value="">Semua merek</option>
      <option v-for="brand in brands" :key="brand" :value="brand">{{ brand }}</option>
    </select>

    <button
      type="button"
      class="inline-flex h-10 shrink-0 items-center gap-2 rounded-full px-4 text-sm font-black shadow-sm ring-1 ring-slate-200 focus:outline-none focus:ring-2 focus:ring-blue-600"
      :class="filters.openNow ? 'bg-emerald-500 text-white ring-emerald-500' : 'bg-white text-slate-700'"
      @click="emit('update-filters', { openNow: !filters.openNow })"
    >
      <i class="fa-regular fa-clock" aria-hidden="true"></i>
      Buka sekarang
    </button>

    <button
      type="button"
      class="inline-flex h-10 shrink-0 items-center gap-2 rounded-full bg-white px-4 text-sm font-black text-slate-700 shadow-sm ring-1 ring-slate-200 focus:outline-none focus:ring-2 focus:ring-blue-600"
      @click="moreOpen = !moreOpen"
    >
      <i class="fa-solid fa-sliders" aria-hidden="true"></i>
      Filter lainnya
    </button>

    <button
      type="button"
      class="ml-auto hidden h-10 shrink-0 rounded-full px-4 text-sm font-black md:inline-flex md:items-center"
      :class="sortMode === 'cheapest' ? 'bg-slate-950 text-white' : 'bg-white text-slate-700 shadow-sm ring-1 ring-slate-200'"
      @click="emit('update-sort', sortMode === 'cheapest' ? 'nearest' : 'cheapest')"
    >
      {{ sortMode === 'cheapest' ? 'Termurah dulu' : 'Terdekat dulu' }}
    </button>
  </section>

  <section v-if="moreOpen" class="grid gap-2 rounded-[12px] bg-white p-3 shadow-sm md:grid-cols-3">
    <label class="grid gap-1 text-sm font-bold text-slate-600">
      Wilayah
      <select
        :value="filters.region"
        class="h-10 rounded-xl bg-slate-100 px-3 font-bold text-slate-950 outline-none focus:ring-2 focus:ring-blue-600"
        @change="emit('update-filters', { region: ($event.target as HTMLSelectElement).value as FilterState['region'] })"
      >
        <option value="">Semua wilayah</option>
        <option v-for="region in regions" :key="region" :value="region">{{ region }}</option>
      </select>
    </label>
    <button class="h-10 self-end rounded-xl bg-slate-100 px-4 text-sm font-black text-slate-700" @click="emit('update-sort', 'cheapest')">Urutkan termurah</button>
    <button class="h-10 self-end rounded-xl bg-slate-100 px-4 text-sm font-black text-slate-700" @click="emit('update-sort', 'nearest')">Urutkan terdekat</button>
  </section>
</template>
