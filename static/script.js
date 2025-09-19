// Start with some default center
const map = L.map('map').setView([0, 0], 2);
L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
  maxZoom: 19
}).addTo(map);

const now = Date.now();

fetch('/api/loc/all')
  .then(r => r.json())
  .then(points => {
    if (points.length === 0) return;

    // draw each segment with dynamic color
    for (let i = 0; i < points.length - 1; i++) {
      const p1 = points[i];
      const p2 = points[i + 1];

      const color = getEventColor(p1.timestamp);
      L.polyline(
        [
          [p1.lat, p1.lng],
          [p2.lat, p2.lng]
        ],
        { color, weight: 4 }
      ).addTo(map);
    }

    // center map on last point
    const last = points[points.length - 1];
    map.setView([last.lat, last.lng], 15);
  });

function getEventColor(eventTime, baseHue = 200) {
  const now = Date.now();
  const elapsed = now - eventTime;

  const halfHour = 60 * 60 * 1000;
  const oneDay = 24 * 60 * 60 * 1000;

  if (elapsed >= oneDay) {
    return `rgba(128,128,128,0)`; // fully invisible
  }

  if (elapsed <= halfHour) {
    // Phase 1: fade color → gray by lowering saturation
    const t = elapsed / halfHour; // 0 → 1
    const saturation = 100 - t * 100; // 100% → 0%
    return `hsl(${baseHue}, ${saturation}%, 50%)`;
  } else {
    // Phase 2: fade gray opacity over the rest of the day
    const t = (elapsed - halfHour) / (oneDay - halfHour); // 0 → 1
    const opacity = 1 - t; // 1 → 0
    return `hsla(0, 0%, 50%, ${opacity.toFixed(2)})`;
  }
}

