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

    points.forEach(p => {
      const color = getEventColor(p.timestamp);
      L.circle([p.lat, p.lon], {
        color: color,
        fillColor: color,
        fillOpacity: 1,
        radius: 100
      }).addTo(map)
        .bindPopup(`Time: ${p.timestamp}<br>Speed: ${p.speed}`);
    });

    // get last point and center map
    const last = points[points.length - 1];
    map.setView([last.lat, last.lon], 15);
  });
function getEventColor(eventTime, baseHue = 200) {
  const now = Date.now();
  const elapsed = now - eventTime;

  const fifteenMinutes = 15 * 60 * 1000;
  const oneDay = 24 * 60 * 60 * 1000;

  if (elapsed >= oneDay) {
    return `rgba(128,128,128,0)`; // fully invisible
  }

  if (elapsed <= fifteenMinutes) {
    // Phase 1: fade color → gray by lowering saturation
    const t = elapsed / fifteenMinutes; // 0 → 1
    const saturation = 100 - t * 100; // 100% → 0%
    return `hsl(${baseHue}, ${saturation}%, 50%)`;
  } else {
    // Phase 2: fade gray opacity over the rest of the day
    const t = (elapsed - fifteenMinutes) / (oneDay - fifteenMinutes); // 0 → 1
    const opacity = 1 - t; // 1 → 0
    return `hsla(0, 0%, 50%, ${opacity.toFixed(2)})`;
  }
}

