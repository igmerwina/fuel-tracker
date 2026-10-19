<script setup lang="ts">
import type { FilterState, FuelBrand, SortMode } from '../types';
import { brands, fuelTypes, regions } from '../data/stations';

defineProps<{
  filters: FilterState;
  sortMode: SortMode;
}>();

const emit = defineEmits<{
  'update-filters': [value: Partial<FilterState>];
  'update-sort': [value: SortMode];
}>();

const selectClass =
  'h-9 shrink-0 rounded-lg border border-garis bg-white pl-3 pr-8 text-sm font-medium outline-none transition hover:border-aspal-soft focus-visible:ring-2 focus-visible:ring-teal';
</script>

<template>
  <div class="flex shrink-0 items-center gap-2 overflow-x-auto border-b border-garis bg-white px-3 py-2 md:px-5">
    <div class="flex shrink-0 rounded-lg bg-trotoar p-1" role="radiogroup" aria-label="Jenis BBM">
      <button
        v-for="fuel in fuelTypes"
        :key="fuel.id"
        type="button"
        role="radio"
        :aria-checked="filters.fuelType === fuel.id"
        class="h-7 shrink-0 rounded-md px-3 text-sm whitespace-nowrap transition"
        :class="filters.fuelType === fuel.id ? 'bg-teal font-semibold text-white' : 'text-aspal-soft hover:text-aspal'"
        :title="fuel.common"
        @click="emit('update-filters', { fuelType: fuel.id })"
      >
        {{ fuel.label }}
      </button>
    </div>

    <select
      :value="filters.brand"
      :class="selectClass"
      aria-label="Merek"
      @change="emit('update-filters', { brand: ($event.target as HTMLSelectElement).value as FuelBrand | '' })"
    >
      <option value="">Semua merek</option>
      <option v-for="brand in brands" :key="brand" :value="brand">{{ brand }}</option>
    </select>

    <select
      :value="filters.region"
      :class="selectClass"
      aria-label="Wilayah"
      @change="emit('update-filters', { region: ($event.target as HTMLSelectElement).value as FilterState['region'] })"
    >
      <option value="">Seluruh Jakarta</option>
      <option v-for="region in regions" :key="region" :value="region">{{ region }}</option>
    </select>

    <div class="ml-auto flex shrink-0 rounded-lg bg-trotoar p-1" role="radiogroup" aria-label="Urutkan">
      <button
        v-for="mode in (['cheapest', 'nearest'] as const)"
        :key="mode"
        type="button"
        role="radio"
        :aria-checked="sortMode === mode"
        class="h-7 rounded-md px-3 text-sm whitespace-nowrap transition"
        :class="sortMode === mode ? 'bg-white font-semibold text-petrol shadow-sm' : 'text-aspal-soft hover:text-aspal'"
        @click="emit('update-sort', mode)"
      >
        {{ mode === 'cheapest' ? 'Termurah' : 'Terdekat' }}
      </button>
    </div>
  </div>
</template>
