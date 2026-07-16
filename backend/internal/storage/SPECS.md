# Structs and tables

- UserProfile table
-

## CRUD Operations

- Create UserProfile
- Read/Get UserProfile
- Update UserProfile
- Delete UserProfile (optional)

## Seeded users

Geological data generation

```Go
type CityCoords struct {
    City string
    Lat  float64
    Lon  float64
}

// Example for Finland
var targetCities = []CityCoords{
    {"Helsinki", 60.1699, 24.9384},
    {"Espoo",    60.2055, 24.6559},
    {"Tampere",  61.4978, 23.7610},
    {"Vantaa",   60.2934, 25.0378},
    {"Turku",    60.4518, 22.2666},
}

// To randomly assign city when generating 100 users:
baseCity := targetCities[rand.Intn(len(targetCities))]

// location diviation ~10-15 km
randomLat := baseCity.Lat + (rand.Float64() - 0.5) * 0.2
randomLon := baseCity.Lon + (rand.Float64() - 0.5) * 0.2
```
