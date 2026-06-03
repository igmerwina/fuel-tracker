<script setup lang="ts">
import type { FuelType, SortMode, StationResult } from '../types';
import { fuelTypes } from '../data/stations';

defineProps<{
  stations: StationResult[];
  activeFuel: FuelType;
  sortMode: SortMode;
}>();

const emit = defineEmits<{
  flyto: [station: StationResult];
  'toggle-favorite': [stationId: string];
}>();

const fuelLabel = (type: string) => fuelTypes.find((fuel) => fuel.id === type)?.label ?? type;

const visiblePriceTypes: FuelType[] = ['RON_90', 'RON_92', 'RON_95'];

const brandTone: Record<string, string> = {
  Pertamina: 'bg-red-50 text-red-700',
  Shell: 'bg-amber-50 text-amber-700',
  Vivo: 'bg-blue-50 text-blue-700',
  BP: 'bg-emerald-50 text-emerald-700',
};
</script>

<template>
  <aside
    class="fixed inset-x-0 bottom-0 z-[900] max-h-[48vh] overflow-hidden rounded-t-[24px] bg-white shadow-2xl shadow-slate-950/20 lg:static lg:z-auto lg:max-h-none lg:rounded-[12px] lg:shadow-xl lg:shadow-slate-200/80"
    aria-label="Station results"
  >
    <div class="mx-auto mt-2 h-1.5 w-12 rounded-full bg-slate-300 lg:hidden"></div>

    <div class="flex items-center justify-between gap-3 border-b border-slate-100 p-4">
      <div>
        <p class="text-xs font-black uppercase tracking-wide text-slate-500">
          {{ sortMode === 'nearest' ? 'Nearest stations' : 'Cheapest stations' }}
        </p>
        <h2 class="text-xl font-black tracking-tight text-slate-950">{{ stations.length }} results</h2>
      </div>
      <span class="rounded-full bg-blue-50 px-3 py-1 text-sm font-black text-blue-700">
        {{ fuelLabel(activeFuel) }}
      </span>
    </div>

    <div v-if="stations.length === 0" class="grid min-h-72 place-items-center p-8 text-center">
      <div>
        <i class="fa-solid fa-map-location-dot text-3xl text-slate-300" aria-hidden="true"></i>
        <p class="mt-3 font-bold text-slate-600">No stations match your search.</p>
      </div>
    </div>

    <div v-else class="max-h-[calc(48vh-76px)] space-y-3 overflow-y-auto p-3 lg:max-h-[calc(100vh-292px)]">
      <article
        v-for="station in stations"
        :key="station.id"
        class="rounded-[12px] bg-white p-4 shadow-md shadow-slate-200/80 ring-1 transition hover:shadow-lg"
        :class="station.isCheapest ? 'ring-emerald-300' : 'ring-slate-100'"
      >
        <div class="flex items-start justify-between gap-3">
          <div class="min-w-0">
            <div class="mb-2 flex flex-wrap items-center gap-2">
              <span class="rounded-full px-2.5 py-1 text-xs font-black" :class="brandTone[station.brand]">
                {{ station.brand }}
              </span>
              <span v-if="station.isCheapest" class="rounded-full bg-emerald-500 px-2.5 py-1 text-xs font-black text-white">
                Cheapest
              </span>
              <span class="rounded-full bg-slate-100 px-2.5 py-1 text-xs font-black text-slate-600">
                Open now
              </span>
            </div>
            <h3 class="truncate text-lg font-black leading-6 text-slate-950">{{ station.name }}</h3>
            <p class="mt-1 line-clamp-2 text-sm font-semibold text-slate-500">{{ station.address }}</p>
          </div>

          <button
            type="button"
            class="grid h-10 w-10 flex-none place-items-center rounded-full text-slate-400 hover:bg-slate-100 hover:text-rose-500 focus:outline-none focus:ring-2 focus:ring-blue-600"
            :aria-label="station.isFavorite ? 'Remove favorite' : 'Save favorite'"
            @click="emit('toggle-favorite', station.id)"
          >
            <i :class="station.isFavorite ? 'fa-solid fa-heart text-rose-500' : 'fa-regular fa-heart'" aria-hidden="true"></i>
          </button>
        </div>

        <div class="mt-3 flex items-center gap-3 text-sm font-bold text-slate-600">
          <span><i class="fa-solid fa-location-arrow mr-1 text-blue-600" aria-hidden="true"></i>{{ station.distanceKm.toFixed(1) }} km</span>
          <span>{{ station.travelMinutes }} min</span>
          <span v-if="station.savingsPerLiter > 0" class="text-emerald-700">
            Save Rp {{ station.savingsPerLiter.toLocaleString('id-ID') }}/L
          </span>
        </div>

        <div class="mt-4 grid grid-cols-3 gap-2">
          <div
            v-for="type in visiblePriceTypes"
            :key="type"
            class="rounded-xl p-3"
            :class="activeFuel === type ? 'bg-blue-600 text-white' : 'bg-slate-100 text-slate-700'"
          >
            <p class="text-[11px] font-black uppercase">{{ fuelLabel(type) }}</p>
            <p class="mt-1 text-sm font-black">
              Rp {{ station.prices.find((price) => price.type === type)?.price.toLocaleString('id-ID') ?? '-' }}
            </p>
          </div>
        </div>

        <div class="mt-4 grid grid-cols-2 gap-2">
          <a
            class="inline-flex min-h-11 items-center justify-center gap-2 rounded-xl bg-blue-600 px-4 text-sm font-black text-white shadow-lg shadow-blue-600/20 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-600 focus:ring-offset-2"
            :href="station.googleMapsUrl"
            target="_blank"
            rel="noopener noreferrer"
          >
            <i class="fa-solid fa-route" aria-hidden="true"></i>
            Directions
          </a>
          <button
            type="button"
            class="inline-flex min-h-11 items-center justify-center gap-2 rounded-xl bg-slate-100 px-4 text-sm font-black text-slate-700 hover:bg-slate-200 focus:outline-none focus:ring-2 focus:ring-blue-600 focus:ring-offset-2"
            @click="emit('flyto', station)"
          >
            <i class="fa-solid fa-circle-info" aria-hidden="true"></i>
            Details
          </button>
        </div>
      </article>
    </div>
  </aside>
</template>
