<!-- Use this file to provide workspace-specific custom instructions to Copilot. For more details, visit https://code.visualstudio.com/docs/copilot/copilot-customization#_use-a-githubcopilotinstructionsmd-file -->

# Jakarta BBM Tracker - Copilot Instructions

## Project Overview
This is a Vue.js 3 fuel price tracker application for Jakarta, Indonesia. The app displays real-time (dummy) fuel prices from major Indonesian fuel brands (Pertamina, Shell, Vivo, BP) across different regions of Jakarta.

## Tech Stack
- **Framework**: Vue.js 3 with TypeScript
- **Build Tool**: Vite
- **Styling**: Scoped CSS with responsive design
- **Data**: Dummy data in TypeScript

## Key Features to Maintain
1. **Responsive UI**: Mobile-first design approach
2. **Component Structure**: Reusable Vue components with proper props and emits
3. **Type Safety**: TypeScript interfaces for all data structures
4. **Filtering & Sorting**: Dynamic filtering by region, brand, and sorting by price
5. **Modal Details**: Click on station cards to view detailed pricing information

## File Organization
- `src/components/`: Vue components (Header, FilterBar, StationCard, StationList, StatsPanel)
- `src/data/`: Dummy data and constants (regions, fuel stations)
- `src/types/`: TypeScript type definitions
- `src/App.vue`: Main application component

## Development Guidelines
- Keep components small and focused
- Use TypeScript interfaces for type safety
- Maintain consistent styling with the purple gradient theme (#667eea to #764ba2)
- Ensure all modals use Teleport to body
- Test responsive behavior on mobile, tablet, and desktop views

## Color Scheme
- Primary: #667eea (purple)
- Secondary: #764ba2
- Brand Colors:
  - Pertamina: #003d7a
  - Shell: #e81c1c
  - Vivo: #ffa500
  - BP: #00a651
