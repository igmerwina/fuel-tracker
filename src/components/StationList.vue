<template>
  <section class="stations-list">
    <div class="list-heading">
      <div>
        <p class="eyebrow">Daftar SPBU</p>
        <h2>Harga tersedia per liter</h2>
      </div>
      <span class="count">{{ stations.length }}</span>
    </div>

    <div v-if="stations.length === 0" class="empty-state">
      <i class="fa-solid fa-map-location-dot" aria-hidden="true"></i>
      <p>Tidak ada SPBU yang sesuai filter.</p>
    </div>

    <div v-else class="list-container">
      <article v-for="station in stations" :key="station.id" class="station-item">
        <div class="station-top">
          <span class="brand-badge" :class="station.brand.toLowerCase()">
            {{ station.brand }}
          </span>
          <span class="region">{{ station.region }}</span>
        </div>

        <h3>{{ station.name }}</h3>
        <p class="address">
          <i class="fa-solid fa-location-dot" aria-hidden="true"></i>
          {{ station.address }}
        </p>

        <div class="prices-info">
          <div
            v-for="price in visiblePrices(station)"
            :key="price.type"
            class="price-item"
            :class="{ highlighted: filters.fuelType === price.type }"
          >
            <span class="fuel-type">{{ getFuelTypeLabel(price.type) }}</span>
            <strong>Rp {{ price.price.toLocaleString('id-ID') }}</strong>
          </div>
        </div>

        <footer class="card-footer">
          <span>Update {{ formatUpdated(station.lastUpdated) }}</span>
          <div class="actions">
            <a
              class="maps-btn"
              :href="station.googleMapsUrl"
              target="_blank"
              rel="noopener noreferrer"
            >
              <i class="fa-solid fa-arrow-up-right-from-square" aria-hidden="true"></i>
              <span>Google Maps</span>
            </a>
            <button class="view-map-btn" type="button" @click="$emit('flyto', station)">
              <i class="fa-solid fa-location-crosshairs" aria-hidden="true"></i>
              <span>Lihat di Peta</span>
            </button>
          </div>
        </footer>
      </article>
    </div>
  </section>
</template>

<script setup lang="ts">
import type { FilterState, FuelPrice, FuelStation } from '../types/index';
import { fuelTypes } from '../data/stations';

const props = defineProps<{
  stations: FuelStation[];
  filters: FilterState;
}>();

defineEmits<{
  flyto: [station: FuelStation];
}>();

const getFuelTypeLabel = (typeId: string): string => {
  const fuelType = fuelTypes.find((fuel) => fuel.id === typeId);
  return fuelType?.label || typeId;
};

const visiblePrices = (station: FuelStation): FuelPrice[] => {
  if (!props.filters.fuelType) return station.prices;
  return station.prices.filter((price) => price.type === props.filters.fuelType);
};

const formatUpdated = (value: string): string => {
  return new Intl.DateTimeFormat('id-ID', {
    day: '2-digit',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit',
  }).format(new Date(value));
};
</script>

<style scoped>
.stations-list {
  display: flex;
  min-height: 0;
  flex: 1;
  flex-direction: column;
  overflow: hidden;
  border: 1px solid var(--line);
  border-radius: 8px;
  background: var(--panel);
  box-shadow: var(--small-shadow);
}

.list-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 15px 16px;
  border-bottom: 1px solid var(--line);
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
  font-size: 18px;
  line-height: 1.2;
}

.count {
  display: grid;
  width: 38px;
  height: 38px;
  place-items: center;
  border-radius: 8px;
  color: var(--brand-strong);
  background: #ffed9c;
  font-weight: 900;
}

.empty-state {
  display: grid;
  flex: 1;
  place-items: center;
  padding: 32px;
  color: var(--muted);
  text-align: center;
}

.empty-state i {
  margin-bottom: 12px;
  color: var(--brand);
  font-size: 30px;
}

.empty-state p {
  margin: 0;
}

.list-container {
  display: flex;
  flex: 1;
  min-height: 0;
  flex-direction: column;
  gap: 12px;
  overflow-y: auto;
  padding: 14px;
}

.station-item {
  padding: 14px;
  border: 1px solid var(--line);
  border-left: 5px solid var(--brand);
  border-radius: 8px;
  background: #ffffff;
}

.station-item:hover {
  box-shadow: 0 10px 24px rgba(15, 23, 42, 0.09);
}

.station-top,
.card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.brand-badge {
  display: inline-flex;
  min-height: 25px;
  align-items: center;
  padding: 0 9px;
  border-radius: 999px;
  color: #ffffff;
  font-size: 12px;
  font-weight: 900;
}

.brand-badge.pertamina {
  background: #d71e2b;
}

.brand-badge.shell {
  color: #1f2937;
  background: #f4c300;
}

.brand-badge.vivo {
  background: #1d4ed8;
}

.brand-badge.bp {
  background: #159447;
}

.region {
  color: var(--muted);
  font-size: 12px;
  font-weight: 800;
}

h3 {
  margin: 10px 0 6px;
  color: var(--ink);
  font-size: 17px;
  line-height: 1.2;
}

.address {
  display: flex;
  gap: 7px;
  margin: 0;
  color: var(--muted);
  font-size: 13px;
}

.prices-info {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8px;
  margin: 13px 0;
}

.price-item {
  display: grid;
  gap: 1px;
  min-width: 0;
  padding: 9px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  background: var(--panel-soft);
}

.price-item.highlighted {
  border-color: #f4c300;
  background: #fff8d9;
}

.fuel-type {
  overflow: hidden;
  color: var(--muted);
  font-size: 11px;
  font-weight: 900;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.price-item strong {
  color: var(--brand-strong);
  font-size: 14px;
  line-height: 1.2;
}

.card-footer {
  color: var(--muted);
  font-size: 12px;
  font-weight: 700;
}

.actions {
  display: inline-flex;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 8px;
}

.maps-btn,
.view-map-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 7px;
  min-height: 36px;
  padding: 0 12px;
  border-radius: 8px;
  color: #ffffff;
  background: var(--brand);
  cursor: pointer;
  font-weight: 900;
  text-decoration: none;
}

.view-map-btn:hover {
  background: var(--brand-strong);
}

.maps-btn {
  color: var(--brand-strong);
  background: #e7f0f7;
}

.maps-btn:hover {
  background: #d7e8f4;
}

@media (max-width: 1180px) {
  .prices-info {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 430px) {
  .prices-info {
    grid-template-columns: 1fr;
  }

  .card-footer {
    align-items: stretch;
    flex-direction: column;
  }

  .actions,
  .maps-btn,
  .view-map-btn {
    width: 100%;
  }
}
</style>
