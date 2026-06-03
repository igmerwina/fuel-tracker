<script setup lang="ts">
import type { StationResult } from '../types';

defineProps<{
  fuelName: string;
  cheapestStation?: StationResult;
  nearestStation?: StationResult;
  lastUpdated: string;
}>();
</script>

<template>
  <section class="grid grid-cols-1 gap-3 md:grid-cols-3">
    <article class="rounded-[12px] bg-white p-4 shadow-lg shadow-slate-200/70">
      <div class="flex items-center justify-between gap-3">
        <p class="text-sm font-bold text-slate-500">Cheapest {{ fuelName }}</p>
        <span class="grid h-10 w-10 place-items-center rounded-xl bg-emerald-50 text-emerald-600">
          <i class="fa-solid fa-fire-flame-curved" aria-hidden="true"></i>
        </span>
      </div>
      <p class="mt-3 text-3xl font-black tracking-tight text-slate-950">
        Rp {{ cheapestStation?.selectedPrice.toLocaleString('id-ID') ?? '-' }}
        <span class="text-base font-bold text-slate-500">/L</span>
      </p>
      <p class="mt-1 truncate text-sm font-semibold text-slate-500">
        {{ cheapestStation?.name ?? 'Belum ada data' }}
      </p>
    </article>

    <article class="rounded-[12px] bg-white p-4 shadow-lg shadow-slate-200/70">
      <div class="flex items-center justify-between gap-3">
        <p class="text-sm font-bold text-slate-500">Nearest Station</p>
        <span class="grid h-10 w-10 place-items-center rounded-xl bg-blue-50 text-blue-600">
          <i class="fa-solid fa-location-dot" aria-hidden="true"></i>
        </span>
      </div>
      <p class="mt-3 text-3xl font-black tracking-tight text-slate-950">
        {{ nearestStation?.distanceKm.toFixed(1) ?? '-' }} km
        <span class="text-base font-bold text-slate-500">· {{ nearestStation?.travelMinutes ?? '-' }} min</span>
      </p>
      <p class="mt-1 truncate text-sm font-semibold text-slate-500">
        {{ nearestStation?.name ?? 'Belum ada data' }}
      </p>
    </article>

    <article class="rounded-[12px] bg-white p-4 shadow-lg shadow-slate-200/70">
      <div class="flex items-center justify-between gap-3">
        <p class="text-sm font-bold text-slate-500">Last Updated</p>
        <span class="grid h-10 w-10 place-items-center rounded-xl bg-amber-50 text-amber-600">
          <i class="fa-solid fa-clock" aria-hidden="true"></i>
        </span>
      </div>
      <p class="mt-3 text-3xl font-black tracking-tight text-slate-950">{{ lastUpdated }}</p>
      <p class="mt-1 text-sm font-semibold text-slate-500">Harga otomatis memakai cache backend saat tersedia</p>
    </article>
  </section>
</template>
