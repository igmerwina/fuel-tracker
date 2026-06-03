<script setup lang="ts">
import type { FuelType, SortMode, StationResult } from '../types';

defineProps<{
  stations: StationResult[];
  activeFuel: FuelType;
  sortMode: SortMode;
  selectedStationId: string;
}>();

const emit = defineEmits<{
  flyto: [station: StationResult];
  'toggle-favorite': [stationId: string];
  'start-navigation': [station: StationResult];
}>();

const price = (station: StationResult, type: FuelType) =>
  station.prices.find((item) => item.type === type)?.price.toLocaleString('id-ID') ?? '-';
</script>

<template>
  <aside
    class="fixed inset-x-0 bottom-0 z-[900] max-h-[42vh] overflow-hidden rounded-t-[20px] bg-white shadow-2xl shadow-slate-950/20 lg:static lg:max-h-none lg:rounded-[12px] lg:shadow-sm"
    aria-label="Hasil SPBU"
  >
    <div class="mx-auto mt-2 h-1.5 w-12 rounded-full bg-slate-300 lg:hidden"></div>

    <div class="flex h-14 items-center justify-between px-3">
      <div>
        <p class="text-[11px] font-black uppercase tracking-wide text-slate-500">
          {{ sortMode === 'nearest' ? 'Terdekat' : 'Harga terbaik' }}
        </p>
        <h2 class="text-base font-black text-slate-950">{{ stations.length }} SPBU</h2>
      </div>
    </div>

    <div class="max-h-[calc(42vh-64px)] overflow-y-auto px-2 pb-3 lg:max-h-[calc(100vh-150px)]">
      <button
        v-for="station in stations"
        :key="station.id"
        type="button"
        class="mb-2 grid min-h-[112px] w-full rounded-[12px] bg-white p-3 text-left shadow-sm ring-1 transition hover:shadow-md focus:outline-none focus:ring-2 focus:ring-blue-600"
        :class="selectedStationId === station.id ? 'ring-blue-500' : station.isCheapest ? 'ring-emerald-300' : 'ring-slate-100'"
        @click="emit('flyto', station)"
      >
        <div class="flex items-start justify-between gap-2">
          <div class="min-w-0">
            <div class="flex items-center gap-2">
              <h3 class="truncate text-sm font-black text-slate-950">{{ station.name }}</h3>
              <span v-if="station.isCheapest" class="shrink-0 rounded-full bg-emerald-500 px-2 py-0.5 text-[10px] font-black text-white">Terbaik</span>
            </div>
            <p class="mt-0.5 text-xs font-bold text-slate-500">
              {{ station.distanceKm.toFixed(1) }} km · Hemat Rp{{ station.savingsPerLiter.toLocaleString('id-ID') }}/L
            </p>
          </div>
          <button
            type="button"
            class="grid h-8 w-8 shrink-0 place-items-center rounded-full text-slate-400 hover:bg-slate-100 hover:text-rose-500"
            :aria-label="station.isFavorite ? 'Hapus favorit' : 'Simpan favorit'"
            @click.stop="emit('toggle-favorite', station.id)"
          >
            <i :class="station.isFavorite ? 'fa-solid fa-heart text-rose-500' : 'fa-regular fa-heart'" aria-hidden="true"></i>
          </button>
        </div>

        <div class="mt-2 grid grid-cols-3 gap-2 text-xs">
          <span><b>RON90</b> Rp{{ price(station, 'RON_90') }}</span>
          <span class="font-black text-blue-700"><b>RON92</b> Rp{{ price(station, 'RON_92') }}</span>
          <span><b>RON95</b> Rp{{ price(station, 'RON_95') }}</span>
        </div>

        <div class="mt-2 flex gap-2">
          <button
            type="button"
            class="h-8 flex-1 rounded-lg bg-blue-600 text-xs font-black text-white"
            @click.stop="emit('start-navigation', station)"
          >
            Rute
          </button>
          <button
            type="button"
            class="h-8 flex-1 rounded-lg bg-slate-100 text-xs font-black text-slate-700"
            @click.stop="emit('flyto', station)"
          >
            Detail
          </button>
        </div>
      </button>
    </div>
  </aside>
</template>
