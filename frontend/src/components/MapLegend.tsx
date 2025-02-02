import React, { useEffect } from "react";
import { useMap } from "react-leaflet";
import L from "leaflet";

interface MapLegendProps {
  minTime: number;
  maxTime: number;
}

const MapLegend: React.FC<MapLegendProps> = ({ minTime, maxTime }) => {
  const map = useMap();

  useEffect(() => {
    const legend = L.control({ position: "bottomleft" });
    const legendColors = [
      { color: "#FF0000", percent: 0.14 },
      { color: "#FFA600", percent: 0.28 },
      { color: "#F9FF00", percent: 0.42 },
      { color: "#28FF00", percent: 0.56 },
      { color: "#00FFF8", percent: 0.7 },
      { color: "#008CFF", percent: 0.86 },
      { color: "#11007F", percent: 1.0 },
    ];

    legend.onAdd = () => {
      const div = L.DomUtil.create("div", "legend");
      div.style.backgroundColor = "white";
      div.style.padding = "10px";
      div.style.borderRadius = "4px";
      div.style.boxShadow = "0 2px 4px rgba(0,0,0,0.2)";

      let legendContent =
        '<h3 style="font-size: 14px; font-weight: 600; margin-bottom: 8px;">Travel Time</h3>';

      legendContent += `
          <div style="display: flex; align-items: center; margin-bottom: 4px;">
            <div style="width: 20px; height: 20px; background-color: #FF00FF; margin-right: 8px; border: 1px solid #ccc;"></div>
            <span style="font-size: 12px;">START</span>
          </div>
        `;

      legendColors.forEach((item, index) => {
        const startTime =
          index === 0
            ? minTime
            : minTime + (maxTime - minTime) * legendColors[index - 1].percent;
        const endTime = minTime + (maxTime - minTime) * item.percent;

        legendContent += `
          <div style="display: flex; align-items: center; margin-bottom: 4px;">
            <div style="width: 20px; height: 20px; background-color: ${item.color}; margin-right: 8px; border: 1px solid #ccc;"></div>
            <span style="font-size: 12px;">${Math.round(startTime / 60)} - ${Math.round(endTime / 60)} min</span>
          </div>
        `;
      });

      div.innerHTML = legendContent;
      return div;
    };

    legend.addTo(map);
    return () => legend.remove();
  }, [map, minTime, maxTime]);

  return null;
};

export default MapLegend;
