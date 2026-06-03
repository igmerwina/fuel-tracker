<template>
  <section class="filter-bar">
    <div class="section-heading">
      <div>
        <p class="eyebrow">Filter pencarian</p>
        <h2>{{ resultCount }} SPBU cocok</h2>
      </div>
      <button class="reset-btn" type="button" title="Reset filter" @click="resetFilters">
        <i class="fa-solid fa-rotate-left" aria-hidden="true"></i>
        <span>Reset</span>
      </button>
    </div>

    <div class="filters-grid">
      <label class="filter-group" for="region">
        <span>Wilayah</span>
        <select id="region" v-model="filterState.region" @change="emitFilters">
          <option value="">Semua wilayah</option>
          <option v-for="reg in regions" :key="reg" :value="reg">
            {{ reg }}
          </option>
        </select>
      </label>

      <label class="filter-group" for="brand">
        <span>Merek SPBU</span>
        <select id="brand" v-model="filterState.brand" @change="emitFilters">
          <option value="">Semua merek</option>
          <option v-for="brand in brands" :key="brand" :value="brand">
            {{ brand }}
          </option>
        </select>
      </label>

      <label class="filter-group full" for="fuelType">
        <span>Jenis BBM</span>
        <select id="fuelType" v-model="filterState.fuelType" @change="emitFilters">
          <option value="">Semua jenis BBM</option>
          <option v-for="fuel in fuelTypes" :key="fuel.id" :value="fuel.id">
            {{ fuel.label }}
          </option>
        </select>
      </label>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import type { FilterState } from '../types/index';
import { brands, fuelTypes, regions } from '../data/stations';

defineProps<{
  resultCount: number;
}>();

const filterState = ref<FilterState>({
  region: '',
  brand: '',
  fuelType: '',
});

const emit = defineEmits<{
  'update-filters': [value: FilterState];
}>();

const emitFilters = () => {
  emit('update-filters', { ...filterState.value });
};

const resetFilters = () => {
  filterState.value = {
    region: '',
    brand: '',
    fuelType: '',
  };
  emitFilters();
};
</script>

<style scoped>
.filter-bar {
  flex: 0 0 auto;
  padding: 16px;
  border: 1px solid var(--line);
  border-radius: 8px;
  background: var(--panel);
  box-shadow: var(--small-shadow);
}

.section-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 14px;
}

.eyebrow {
  margin: 0;
  color: var(--muted);
  font-size: 11px;
  font-weight: 900;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

h2 {
  margin: 2px 0 0;
  color: var(--ink);
  font-size: 21px;
  line-height: 1.15;
}

.filters-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.filter-group {
  display: grid;
  gap: 6px;
  min-width: 0;
}

.filter-group.full {
  grid-column: 1 / -1;
}

.filter-group span {
  color: #334155;
  font-size: 12px;
  font-weight: 800;
}

select {
  width: 100%;
  min-height: 42px;
  padding: 0 36px 0 12px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  color: var(--ink);
  background: #ffffff;
  cursor: pointer;
  outline: none;
}

select:hover {
  border-color: #8aa4bd;
}

select:focus {
  border-color: var(--brand);
  box-shadow: 0 0 0 3px rgba(11, 79, 122, 0.13);
}

.reset-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  min-height: 38px;
  padding: 0 12px;
  border-radius: 8px;
  color: var(--brand-strong);
  background: #e7f0f7;
  cursor: pointer;
  font-weight: 800;
}

.reset-btn:hover {
  background: #d7e8f4;
}

@media (max-width: 430px) {
  .filters-grid {
    grid-template-columns: 1fr;
  }
}
</style>
