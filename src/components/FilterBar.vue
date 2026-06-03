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

const advancedOpen = ref(false);

const chipClass = (active: boolean) =>
  [
    'inline-flex min-h-10 items-center justify-center rounded-full px-4 text-sm font-black transition focus:outline-none focus:ring-2 focus:ring-blue-600 focus:ring-offset-2',
    active ? 'bg-blue-600 text-white shadow-lg shadow-blue-600/20' : 'bg-white text-slate-700 shadow-sm ring-1 ring-slate-200 hover:bg-slate-50',
  ].join(' ');

const fuelLabel = (fuel: FuelType) => {
  if (fuel === 'RON_92') return 'Pertamax';
  if (fuel === 'RON_90') return 'Pertalite';
  return fuelTypes.find((item) => item.id === fuel)?.label ?? fuel;
};
</script>

<template>
  <section class="rounded-[12px] bg-white p-3 shadow-lg shadow-slate-200/70">
    <div class="flex flex-wrap items-center gap-2">
      <button
        v-for="fuel in fuelTypes"
        :key="fuel.id"
        type="button"
        :class="chipClass(filters.fuelType === fuel.id)"
        @click="emit('update-filters', { fuelType: fuel.id })"
      >
        {{ fuelLabel(fuel.id) }}
      </button>

      <span class="mx-1 hidden h-8 w-px bg-slate-200 md:block"></span>

      <button
        v-for="brand in brands"
        :key="brand"
        type="button"
        :class="chipClass(filters.brand === brand)"
        @click="emit('update-filters', { brand: filters.brand === brand ? '' : (brand as FuelBrand) })"
      >
        {{ brand }}
      </button>

      <button
        type="button"
        :class="chipClass(filters.openNow)"
        @click="emit('update-filters', { openNow: !filters.openNow })"
      >
        <i class="fa-regular fa-clock mr-2" aria-hidden="true"></i>
        Open now
      </button>

      <button
        type="button"
        class="ml-auto inline-flex min-h-10 items-center gap-2 rounded-full bg-slate-100 px-4 text-sm font-black text-slate-700 hover:bg-slate-200 focus:outline-none focus:ring-2 focus:ring-blue-600 focus:ring-offset-2"
        @click="advancedOpen = !advancedOpen"
      >
        <i class="fa-solid fa-sliders" aria-hidden="true"></i>
        Filters
      </button>
    </div>

    <div class="mt-3 flex flex-wrap items-center gap-2">
      <button
        type="button"
        :class="chipClass(sortMode === 'cheapest')"
        @click="emit('update-sort', 'cheapest')"
      >
        Sort cheapest
      </button>
      <button
        type="button"
        :class="chipClass(sortMode === 'nearest')"
        @click="emit('update-sort', 'nearest')"
      >
        Sort nearest
      </button>
      <button
        type="button"
        class="inline-flex min-h-10 items-center justify-center rounded-full px-4 text-sm font-black text-slate-500 hover:bg-slate-100"
        @click="emit('update-filters', { query: '', region: '', brand: '', fuelType: 'RON_92', openNow: false })"
      >
        Reset
      </button>
    </div>

    <div v-if="advancedOpen" class="mt-3 grid gap-2 border-t border-slate-100 pt-3 md:grid-cols-2">
      <label class="grid gap-1 text-sm font-bold text-slate-600">
        Wilayah
        <select
          :value="filters.region"
          class="h-11 rounded-xl border-0 bg-slate-100 px-3 font-bold text-slate-950 outline-none ring-1 ring-slate-200 focus:ring-2 focus:ring-blue-600"
          @change="emit('update-filters', { region: ($event.target as HTMLSelectElement).value as FilterState['region'] })"
        >
          <option value="">Semua wilayah</option>
          <option v-for="region in regions" :key="region" :value="region">{{ region }}</option>
        </select>
      </label>
    </div>
  </section>
</template>
