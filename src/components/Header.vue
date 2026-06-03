<script setup lang="ts">
defineProps<{
  stationCount: number;
  lowestPrice: number;
  lastUpdated: string;
}>();
</script>

<template>
  <header class="header">
    <div class="brand-block">
      <div class="mark" aria-hidden="true">
        <i class="fa-solid fa-gas-pump"></i>
      </div>
      <div>
        <p class="kicker">DKI Jakarta fuel intelligence</p>
        <h1>Jakarta BBM Tracker</h1>
      </div>
    </div>

    <div class="stats" aria-label="Ringkasan data">
      <div class="stat">
        <span class="label">SPBU tampil</span>
        <strong>{{ stationCount }}</strong>
      </div>
      <div class="stat">
        <span class="label">Harga terendah</span>
        <strong>Rp {{ lowestPrice.toLocaleString('id-ID') }}</strong>
      </div>
      <div class="stat">
        <span class="label">Update data</span>
        <strong>{{ lastUpdated }}</strong>
      </div>
    </div>
  </header>
</template>

<style scoped>
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
  padding: 18px 22px;
  color: #ffffff;
  background:
    linear-gradient(90deg, rgba(8, 58, 90, 0.96), rgba(11, 79, 122, 0.96)),
    linear-gradient(135deg, #ffcc32 0%, #ffffff 100%);
  border-bottom: 4px solid var(--accent);
}

.brand-block {
  display: flex;
  min-width: 0;
  align-items: center;
  gap: 14px;
}

.mark {
  display: grid;
  width: 46px;
  height: 46px;
  flex: 0 0 auto;
  place-items: center;
  border-radius: 8px;
  color: var(--brand-strong);
  background: var(--accent);
  box-shadow: inset 0 -2px 0 rgba(0, 0, 0, 0.14);
}

.mark i {
  font-size: 22px;
}

.kicker {
  margin: 0 0 2px;
  color: rgba(255, 255, 255, 0.72);
  font-size: 12px;
  font-weight: 800;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

h1 {
  margin: 0;
  font-size: 26px;
  line-height: 1.05;
  letter-spacing: 0;
}

.stats {
  display: grid;
  grid-template-columns: repeat(3, minmax(112px, 1fr));
  gap: 10px;
  min-width: 430px;
}

.stat {
  min-width: 0;
  padding: 10px 12px;
  border: 1px solid rgba(255, 255, 255, 0.18);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.1);
}

.label {
  display: block;
  color: rgba(255, 255, 255, 0.68);
  font-size: 11px;
  font-weight: 800;
  text-transform: uppercase;
}

strong {
  display: block;
  margin-top: 2px;
  overflow: hidden;
  color: #ffffff;
  font-size: 17px;
  font-weight: 800;
  text-overflow: ellipsis;
  white-space: nowrap;
}

@media (max-width: 920px) {
  .header {
    align-items: stretch;
    flex-direction: column;
  }

  .stats {
    min-width: 0;
    width: 100%;
  }
}

@media (max-width: 620px) {
  .header {
    padding: 15px 14px;
  }

  .stats {
    grid-template-columns: 1fr;
  }

  h1 {
    font-size: 23px;
  }
}
</style>
