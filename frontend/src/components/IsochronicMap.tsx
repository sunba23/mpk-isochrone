import React, { useState, useEffect } from 'react';
import { MapContainer, TileLayer, Marker, Popup } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import L from 'leaflet';
import { fetchStops } from '../services/apiService';
import { StopData } from '../types/stops';

import markerIcon2x from 'leaflet/dist/images/marker-icon-2x.png';
import markerIcon from 'leaflet/dist/images/marker-icon.png';
import markerShadow from 'leaflet/dist/images/marker-shadow.png';

delete L.Icon.Default.prototype._getIconUrl;
L.Icon.Default.mergeOptions({
  iconRetinaUrl: markerIcon2x,
  iconUrl: markerIcon,
  shadowUrl: markerShadow,
});

const IsochronicMap: React.FC = () => {
  const [stops, setStops] = useState<StopData[]>([]);

  useEffect(() => {
    const loadStops = async () => {
      try {
        const fetchedStops = await fetchStops();
        setStops(fetchedStops);
      } catch (error) {
        console.error('Error loading stops:', error);
      }
    };

    loadStops();
  }, []);

  return (
    <MapContainer 
      center={stops.length ? [stops[0].stop_location.latitude, stops[0].stop_location.longitude] : [0, 0]} 
      zoom={13} 
      style={{ height: '100vh', width: '100%' }}
    >
      <TileLayer
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
      />
      {stops.map(stop => (
        <Marker 
          key={stop.id}
          position={[stop.stop_location.latitude, stop.stop_location.longitude]}
        >
          <Popup>
            {stop.stop_name}<br />
            Code: {stop.stop_code}
          </Popup>
        </Marker>
      ))}
    </MapContainer>
  );
};

export default IsochronicMap;
