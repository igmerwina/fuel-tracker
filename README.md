# Jakarta BBM Tracker

A modern Vue.js 3 fuel price tracker application for Jakarta, Indonesia. Track fuel prices from major Indonesian fuel brands (Pertamina, Shell, Vivo, BP) across different regions of Jakarta.

## Features

✨ **Key Features:**
- 📍 Browse fuel stations across 5 Jakarta regions (North, South, East, West, Central)
- ⛽ Track prices for Regular, Premium, and Diesel fuel
- 🏢 Filter by brand (Pertamina, Shell, Vivo, BP)
- 📊 Sort by price, brand, or region
- 📱 Fully responsive mobile-first design
- 🎨 Modern gradient UI with smooth interactions
- 🔍 Detailed station information in modal popups

## Tech Stack

- **Vue.js 3** - Progressive JavaScript framework
- **TypeScript** - Type-safe JavaScript
- **Vite** - Lightning-fast build tool
- **Scoped CSS** - Component-scoped styling

## Project Structure

```
src/
├── components/
│   ├── Header.vue        # Main header with branding
│   ├── FilterBar.vue     # Region, brand, and sort filters
│   ├── StatsPanel.vue    # Statistics overview
│   ├── StationCard.vue   # Individual station card with modal
│   └── StationList.vue   # Grid of station cards
├── data/
│   └── stations.ts       # Dummy fuel station and region data
├── types/
│   └── index.ts          # TypeScript interfaces
├── App.vue               # Main application component
├── main.ts               # Application entry point
└── style.css             # Global styles
```

## Getting Started

### Prerequisites
- Node.js 16+
- npm or yarn

### Installation

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

## Development

The development server runs at `http://localhost:5173/`

### Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run type-check` - Run TypeScript type checking (if configured)

## Data

**Note**: All fuel price data is dummy data for demonstration purposes. This application uses hardcoded mock data stored in `src/data/stations.ts`.

### Sample Data Includes:
- 25 fuel stations across 5 Jakarta regions
- 4 fuel brands with authentic branding colors
- 3 fuel types with different pricing
- Station locations with coordinates
- Last updated timestamps

## Styling

The application uses a modern color scheme:
- **Primary Gradient**: Purple (#667eea to #764ba2)
- **Brand Colors**:
  - Pertamina: Navy Blue (#003d7a)
  - Shell: Red (#e81c1c)
  - Vivo: Orange (#ffa500)
  - BP: Green (#00a651)

## Components

### Header
Displays the application title and branding with gradient background.

### FilterBar
Provides filtering and sorting controls:
- Filter by region
- Filter by brand
- Sort by price, brand, or region
- Reset filters button

### StatsPanel
Shows overview statistics:
- Total stations
- Number of regions
- Cheapest fuel price
- Most expensive fuel price

### StationCard
Displays individual fuel station:
- Station name and brand badge
- Address and region
- Three fuel price types
- Click to view detailed information in modal

### StationList
Grid view of all fuel stations with filtering and sorting applied.

## Browser Support

- Chrome/Edge (latest)
- Firefox (latest)
- Safari (latest)
- Mobile browsers

## License

MIT

## Author

Created as a demonstration fuel tracker application inspired by FuelWatch WA.

