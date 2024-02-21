# Data

Historic Data
- Provides up to last 31 days of data
- Fetches data if data behind buffer, makes client wait until data fullfilled
  - Returns data to client as a in house Kline object

Live Data
- Gets number of Klines to follow (min 1, max 250)
  - Returns a channel publishes in house Kline objects
    - for a single symbol

Correlations
- Returns a symbol correlation
- Returns symbols above x avg_correlation

Kline
- Can be merged to create a pair